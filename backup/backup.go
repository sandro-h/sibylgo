package backup

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/sandro-h/sibylgo/util"
	log "github.com/sirupsen/logrus"
)

const sibylCommitAuthor = "sibylgo@example.com"
const dailyBackupPrefix = "Daily backup for "

var getNow = func() time.Time {
	return time.Now()
}

// Save creates a new backup of the todo file
func Save(files *util.FileConfig, message string) (*Backup, error) {
	todoDir := filepath.Dir(files.TodoFile)
	if !isRepoInitiated(todoDir) {
		err := initRepo(todoDir)
		if err != nil {
			return nil, err
		}
	}

	commit, err := commit(todoDir, message, sibylCommitAuthor, files.TodoFile, files.TrashFile)
	if err != nil {
		return nil, err
	}

	backup := toBackup(commit)
	return backup, nil
}

// Restore restores the todoFile to the passed backup and creates a new backup for this restored state.
// It does not delete any of the intermediate backups that were reverted, so it's still possible to restore
// a different state.
func Restore(files *util.FileConfig, restoreTo *Backup) (*Backup, error) {
	todoDir := filepath.Dir(files.TodoFile)
	if !isRepoInitiated(todoDir) {
		return nil, fmt.Errorf("no backups set up for %s", files.TodoDir)
	}

	restoreMessage := fmt.Sprintf("Restore backup %s '%s'", restoreTo.Identifier, restoreTo.Message)
	revertCommit, err := revertToCommit(todoDir, restoreTo.Identifier, restoreMessage, sibylCommitAuthor)
	if err != nil {
		return nil, err
	}

	restoreBackup := toBackup(revertCommit)
	return restoreBackup, nil
}

// CheckAndMakeDailyBackup creates a daily backup of the todofile if there isn't one already for today.
func CheckAndMakeDailyBackup(files *util.FileConfig) (*Backup, error) {
	newestDailyCommitDate, err := findNewestDailyCommitTimestamp(files.TodoDir)
	if err != nil {
		return nil, err
	}

	today := util.SetToStartOfDay(getNow())
	if !today.After(newestDailyCommitDate) {
		// Already have a daily backup for today
		return nil, nil
	}

	log.Infof("Creating daily backup for %s\n", today.Format("2006-01-02"))
	return Save(files, fmt.Sprintf("%s%s", dailyBackupPrefix, today.Format("02.01.2006")))
}

// SyncToRemote pushes the backup to the remote location defined in the backup configuration.
func SyncToRemote(backupCfg *util.Config, files *util.FileConfig) error {
	log.Infof("Pushing backup to remote %s\n", backupCfg.GetStringOrFail("remote_url"))

	return push(files.TodoDir,
		backupCfg.GetString("remote_url", ""),
		backupCfg.GetString("remote_user", ""),
		backupCfg.GetString("remote_password", ""))
}

// findNewestDailyCommitTimestamp returns the timestamp of the newest daily backup commit,
// or the base epoch time if there is no daily backup commit yet.
func findNewestDailyCommitTimestamp(todoDir string) (time.Time, error) {
	if !isRepoInitiated(todoDir) {
		return time.Unix(0, 0), nil
	}

	newestDailyCommit, err := findNewestCommit(todoDir, func(c *commitEntry) bool {
		return c.AuthorEmail == sibylCommitAuthor &&
			strings.HasPrefix(c.Message, dailyBackupPrefix)
	})
	if err != nil {
		return time.Unix(0, 0), nil
	}
	if newestDailyCommit != nil {
		return newestDailyCommit.Timestamp, nil
	}
	return time.Unix(0, 0), nil
}

// ListBackups lists all backups saved for the todoFile. They are ordered from newest to oldest backup.
func ListBackups(todoDir string) ([]*Backup, error) {
	commits, err := findCommits(todoDir, func(c *commitEntry) bool {
		return c.AuthorEmail == sibylCommitAuthor
	})
	if err != nil {
		return nil, err
	}

	var backups []*Backup
	for _, c := range commits {
		backups = append(backups, toBackup(c))
	}
	return backups, nil
}

// Backup denotes a specific backup of the todofile. It doesn't contain the content, but
// acts as a reference for restoring.
type Backup struct {
	Identifier string
	Timestamp  time.Time
	Message    string
}

func toBackup(c *commitEntry) *Backup {
	return &Backup{
		Identifier: c.Hash,
		Timestamp:  c.Timestamp,
		Message:    c.Message,
	}
}
