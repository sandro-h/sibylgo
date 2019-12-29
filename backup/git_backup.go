package backup

import (
	"bytes"
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"io"
	"os/exec"
	"path/filepath"
	"time"
)

// IsRepoInitiated returns true if the passed folder is a git repository.
func IsRepoInitiated(repoPath string) bool {
	_, err := git.PlainOpen(repoPath)
	if err != nil {
		return false
	}
	return true
}

// InitRepo initiates a new non-bare Git repo in the passed folder. Fails if there already is a git repository.
func InitRepo(repoPath string) error {
	_, err := git.PlainInit(repoPath, false)
	return err
}

// Commit stages and commits the passed files in the passed folder.
// Also commits if none of the passed files changed.
func Commit(repoPath string, message string, authorEmail string, files ...string) (*CommitEntry, error) {
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
		if err != nil {
			return nil, err
		}
	}

	hash, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorEmail,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}

	c, err := r.CommitObject(hash)
	if err != nil {
		return nil, err
	}
	return toCommitEntry(c), nil
}

// RevertToCommit reverts all changes done since commitHash and creates a single new commit with these reversions.
func RevertToCommit(repoPath string, commitHash string, newCommitMessage string, authorEmail string) (*CommitEntry, error) {
	// Note go-git doesn't support revert, so we have to use the cli tool (and assume it's installed)
	cmd := exec.Command("git", "revert", "--no-commit", fmt.Sprintf("%s..HEAD", commitHash))
	cmd.Dir = repoPath
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("Error running git command, %s, stderr: %s", err.Error(), stderr.String())
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
			When:  time.Now(),
		},
	})
	if err != nil {
		return nil, err
	}

	c, err := r.CommitObject(hash)
	if err != nil {
		return nil, err
	}
	return toCommitEntry(c), nil
}

var matchAny = func(c *CommitEntry) bool {
	return true
}

// ListCommits returns all commits in the passed folder, ordered from newest to oldest.
func ListCommits(repoPath string) ([]*CommitEntry, error) {
	return findCommits(repoPath, matchAny, false)
}

// FindCommits returns all commits in the passed folder that match the predicate, ordered from newest to oldest.
func FindCommits(repoPath string, predicate func(*CommitEntry) bool) ([]*CommitEntry, error) {
	return findCommits(repoPath, predicate, false)
}

// ListNewestCommit returns the newest commit in the passed folder. If there are no commits, nil is returned.
func ListNewestCommit(repoPath string) (*CommitEntry, error) {
	found, err := findCommits(repoPath, matchAny, true)
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

// FindNewestCommit returns the newest commit in the passed folder that matches the predicate, or nil.
func FindNewestCommit(repoPath string, predicate func(*CommitEntry) bool) (*CommitEntry, error) {
	found, err := findCommits(repoPath, predicate, true)
	if err != nil {
		return nil, err
	}
	if len(found) == 0 {
		return nil, nil
	}
	return found[0], nil
}

func findCommits(repoPath string, predicate func(*CommitEntry) bool, stopAtFirst bool) ([]*CommitEntry, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	cIter, err := r.Log(&git.LogOptions{})
	if err != nil {
		return nil, err
	}

	var found []*CommitEntry
	for {
		c, err := cIter.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		ce := toCommitEntry(c)
		if predicate(ce) {
			found = append(found, ce)
			if stopAtFirst {
				break
			}
		}
	}
	return found, nil
}

// CommitEntry encapsulates a Git commit.
type CommitEntry struct {
	Hash        string
	Timestamp   time.Time
	Message     string
	AuthorEmail string
	Files       []string
}

func toCommitEntry(c *object.Commit) *CommitEntry {
	fIter, _ := c.Files()
	var files []string
	for {
		f, err := fIter.Next()
		if err == io.EOF {
			break
		}
		files = append(files, f.Name)
	}

	return &CommitEntry{
		Hash:        c.Hash.String(),
		Timestamp:   c.Author.When,
		Message:     c.Message,
		AuthorEmail: c.Author.Email,
		Files:       files,
	}
}
