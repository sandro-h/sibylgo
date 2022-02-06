package backup

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"time"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/format/index"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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

// Commit stages and commits the passed files in the passed folder.
// Also commits if none of the passed files changed.
func commit(repoPath string, message string, authorEmail string, files ...string) (*commitEntry, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	w, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		rel, _ := filepath.Rel(repoPath, f)
		_, err := w.Add(rel)
		if err != nil && !errors.Is(err, index.ErrEntryNotFound) {
			return nil, err
		}
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
	cmd := exec.Command("git", "revert", "--no-commit", fmt.Sprintf("%s..HEAD", commitHash))
	cmd.Dir = repoPath
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error running git command, %s, stderr: %s", err.Error(), stderr.String())
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
