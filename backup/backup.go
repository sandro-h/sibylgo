package backup

import (
	"fmt"
	"path/filepath"
	"time"
)

// TODO: provide a daily backup goroutine, make it check last daily backup commit so that it works if sibyl turned off for a while

const sibylCommitAuthor = "sibylgo@example.com"

// Save creates a new backup of the todo file
func Save(todoFile string, message string) (*Backup, error) {
	todoDir := filepath.Dir(todoFile)
	if !isRepoInitiated(todoDir) {
		err := initRepo(todoDir)
		if err != nil {
			return nil, err
		}
	}

	commit, err := commit(todoDir, message, sibylCommitAuthor, todoFile)
	if err != nil {
		return nil, err
	}

	backup := toBackup(commit)
	return backup, nil
}

// Restore restores the todoFile to the passed backup and creates a new backup for this restored state.
// It does not delete any of the intermediate backups that were reverted, so it's still possible to restore
// a different state.
func Restore(todoFile string, restoreTo *Backup) (*Backup, error) {
	todoDir := filepath.Dir(todoFile)
	if !isRepoInitiated(todoDir) {
		return nil, fmt.Errorf("No backups set up for %s", todoFile)
	}

	restoreMessage := fmt.Sprintf("Restore backup %s '%s'", restoreTo.Identifier, restoreTo.Message)
	revertCommit, err := revertToCommit(todoDir, restoreTo.Identifier, restoreMessage, sibylCommitAuthor)
	if err != nil {
		return nil, err
	}

	restoreBackup := toBackup(revertCommit)
	return restoreBackup, nil
}

// ListBackups lists all backups saved for the todoFile. They are ordered from newest to oldest backup.
func ListBackups(todoFile string) ([]*Backup, error) {
	todoDir := filepath.Dir(todoFile)
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
