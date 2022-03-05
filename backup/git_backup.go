package backup

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sandro-h/sibylgo/util"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// IsRepoInitiated returns true if the passed folder is a git repository.
func isRepoInitiated(repoPath string) bool {
	_, err := git.PlainOpen(repoPath)
	return err == nil
}

// InitRepo initiates a new non-bare Git repo in the passed folder. Fails if there already is a git repository.
func initRepo(repoPath string) error {
	_, err := git.PlainInit(repoPath, false)
	return err
}

// EnableGitEncryption configures the backup git repository to use filters to automatically
// encrypt any committed backups, using sibylgo as the provider for encryption and decryption.
func EnableGitEncryption(repoPath string, executable string) error {
	if !isRepoInitiated(repoPath) {
		err := initRepo(repoPath)
		if err != nil {
			return err
		}
	}

	gitAttributes := "*.txt filter=sibylgo_filter diff=sibylgo_filter"
	err := util.WriteFile(filepath.Join(repoPath, ".gitattributes"), gitAttributes)
	if err != nil {
		return err
	}

	gitConfigFile := filepath.Join(repoPath, ".git", "config")
	existingGitConfig := ""
	if util.Exists(gitConfigFile) {
		existingGitConfig, err = util.ReadFile(gitConfigFile)
		if err != nil {
			return err
		}
	}

	if !strings.Contains(existingGitConfig, "[filter \"sibylgo_filter\"]") {
		normalizedExec := strings.ReplaceAll(executable, "\\", "/")

		updatedGitConfig := existingGitConfig + fmt.Sprintf(`
[filter "sibylgo_filter"]
	clean = %s --encrypt
	smudge = %s --decrypt

[diff "sibylgo_filter"]
	textconv = %s --decrypt
`, normalizedExec, normalizedExec, normalizedExec)

		err = util.WriteFile(gitConfigFile, updatedGitConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

// Commit stages and commits the passed files in the passed folder.
// Also commits if none of the passed files changed.
func commit(repoPath string, message string, authorEmail string, files ...string) (*commitEntry, error) {

	// Use git CLI here to support filters for encryption. go-git does not support this.
	for _, f := range files {
		if util.Exists(f) {
			rel, _ := filepath.Rel(repoPath, f)
			_, err := runGitCmd(repoPath, "git", "add", rel)
			if err != nil {
				return nil, err
			}
		}
	}
	_, err := runGitCmd(repoPath, "git", "add", "-u")
	if err != nil {
		return nil, err
	}

	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	hash, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorEmail,
			Email: authorEmail,
			When:  getNow(),
		},
	})
	if err != nil {
		return nil, err
	}

	c, err := r.CommitObject(hash)
	if err != nil {
		return nil, err
	}
	return tocommitEntry(c), nil
}

// RevertToCommit reverts all changes done since commitHash and creates a single new commit with these reversions.
func revertToCommit(repoPath string, commitHash string, newCommitMessage string, authorEmail string) (*commitEntry, error) {
	// Note go-git doesn't support revert, so we have to use the cli tool (and assume it's installed)
	_, err := runGitCmd(repoPath, "git", "revert", "--no-commit", fmt.Sprintf("%s..HEAD", commitHash))
	if err != nil {
		return nil, err
	}

	// The revert command didn't auto-commit but staged all the necessary changes
	// (would otherwise create one commit for each reverted commit)
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	hash, err := w.Commit(newCommitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorEmail,
			Email: authorEmail,
			When:  getNow(),
		},
	})
	if err != nil {
		return nil, err
	}

	c, err := r.CommitObject(hash)
	if err != nil {
		return nil, err
	}
	return tocommitEntry(c), nil
}

func push(repoPath string, remoteURL string, remoteUser string, remotePassword string) error {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return err
	}

	err = r.DeleteRemote("origin")
	if err != nil && !errors.Is(err, git.ErrRemoteNotFound) {
		return err
	}

	remote := config.RemoteConfig{
		Name: "origin",
		URLs: []string{remoteURL},
	}

	_, err = r.CreateRemote(&remote)
	if err != nil {
		return err
	}

	var auth *http.BasicAuth
	if remoteUser != "" {
		auth = &http.BasicAuth{
			Username: remoteUser,
			Password: remotePassword,
		}
	}

	return r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
}

func runGitCmd(repoPath string, cmdAndArgs ...string) (string, error) {
	cmd := exec.Command(cmdAndArgs[0], cmdAndArgs[1:]...)
	cmd.Dir = repoPath
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error running git command, %s, stderr: %s", err.Error(), stderr.String())
	}
	return stdout.String(), nil
}

var matchAny = func(c *commitEntry) bool {
	return true
}

// ListCommits returns all commits in the passed folder, ordered from newest to oldest.
func listCommits(repoPath string) ([]*commitEntry, error) {
	return doFindCommits(repoPath, matchAny, false)
}

// FindCommits returns all commits in the passed folder that match the predicate, ordered from newest to oldest.
func findCommits(repoPath string, predicate func(*commitEntry) bool) ([]*commitEntry, error) {
	return doFindCommits(repoPath, predicate, false)
}

// FindNewestCommit returns the newest commit in the passed folder that matches the predicate, or nil.
func findNewestCommit(repoPath string, predicate func(*commitEntry) bool) (*commitEntry, error) {
	found, err := doFindCommits(repoPath, predicate, true)
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

func doFindCommits(repoPath string, predicate func(*commitEntry) bool, stopAtFirst bool) ([]*commitEntry, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	cIter, err := r.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	var found []*commitEntry
	for {
		c, err := cIter.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		ce := tocommitEntry(c)
		if predicate(ce) {
			found = append(found, ce)
			if stopAtFirst {
				break
			}
		}
	}
	return found, nil
}

// commitEntry encapsulates a Git commit.
type commitEntry struct {
	Hash        string
	Timestamp   time.Time
	Message     string
	AuthorEmail string
	Files       []string
}

func tocommitEntry(c *object.Commit) *commitEntry {
	fIter, _ := c.Files()
	var files []string
	for {
		f, err := fIter.Next()
		if err == io.EOF {
			break
		}
		files = append(files, f.Name)
	}

	return &commitEntry{
		Hash:        c.Hash.String(),
		Timestamp:   c.Author.When,
		Message:     c.Message,
		AuthorEmail: c.Author.Email,
		Files:       files,
	}
}
