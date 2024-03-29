package backup

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"os"
	"path/filepath"
	"testing"
	"time"

	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/sandro-h/sibylgo/util"
	"github.com/stretchr/testify/assert"
)

func TestIsRepoInitiated_ProperRepo(t *testing.T) {
	pwd, _ := os.Getwd()
	dir := filepath.Dir(pwd)

	initiated := isRepoInitiated(dir)

	assert.True(t, initiated)
}

func TestIsRepoInitiated_NotARepo(t *testing.T) {
	dir, _ := os.Getwd()

	initiated := isRepoInitiated(dir)

	assert.False(t, initiated)
}

func TestInitRepo(t *testing.T) {
	repoPath := tu.MakeTempDir("sibyl_git_backup_test")
	defer tu.DeleteTempDir(repoPath)
	assert.False(t, isRepoInitiated(repoPath))

	err := initRepo(repoPath)

	assert.NoError(t, err)
	assert.True(t, isRepoInitiated(repoPath))
}

func TestEnableGitEncryption(t *testing.T) {
	repoPath := tu.MakeTempDir("sibyl_git_backup_test")
	defer tu.DeleteTempDir(repoPath)
	assert.False(t, isRepoInitiated(repoPath))
	initRepo(repoPath)

	err := EnableGitEncryption(repoPath, "/bin/bla")

	assert.NoError(t, err)
	assert.True(t, util.Exists(filepath.Join(repoPath, ".gitattributes")))
	assert.True(t, util.Exists(filepath.Join(repoPath, ".git", "config")))
	gitConfig, _ := util.ReadFile(filepath.Join(repoPath, ".git", "config"))
	expectedGitConfig := `
[filter "sibylgo_filter"]
	clean = /bin/bla --encrypt
	smudge = /bin/bla --decrypt

[diff "sibylgo_filter"]
	textconv = /bin/bla --decrypt
`
	assert.Equal(t, expectedGitConfig, gitConfig)
}

func TestGitEncryption(t *testing.T) {
	repoPath := tu.MakeTempDir("sibyl_git_backup_test")
	defer tu.DeleteTempDir(repoPath)
	assert.False(t, isRepoInitiated(repoPath))

	var sibylgoExecutableName string
	if runtime.GOOS == "windows" {
		sibylgoExecutableName = "sibylgo.exe"
	} else {
		sibylgoExecutableName = "sibylgo"
	}
	sibylgoExecutable, _ := filepath.Abs(filepath.Join("..", sibylgoExecutableName))
	assert.True(t, util.Exists(sibylgoExecutable), "Sibylgo executable at %s does not exist. Build it first", sibylgoExecutable)

	configFile := repoPath + "/sibylgo.yml"
	util.WriteFile(configFile, `
backup:
  encrypt_password: password123
`)
	sibylgoExecutable += " --config " + configFile

	initRepo(repoPath)
	EnableGitEncryption(repoPath, sibylgoExecutable)

	// Debugging in case of failure:
	gitattributes, _ := util.ReadFile(repoPath + "/.gitattributes")
	gitconfig, _ := util.ReadFile(repoPath + "/.git/config")
	fmt.Println(".gitattributes:\n" + gitattributes)
	fmt.Println(".git/config:\n" + gitconfig)

	fileContent := "hello world!"
	file1 := filepath.Join(repoPath, "file1.txt")
	os.WriteFile(file1, []byte(fileContent), 0644)

	// When
	_, err := commit(repoPath, "A test commit", "test@example.com", file1)

	// Then
	assert.NoError(t, err)
	committedFileContent, _ := runGitCmd(repoPath, "git", "show", "HEAD:file1.txt")
	assert.NotEmpty(t, committedFileContent)
	assert.NotEqual(t, fileContent, committedFileContent)

	var b bytes.Buffer
	cryptor := AnsibleCryptor{Password: "password123"}
	cryptor.DecryptContent(strings.NewReader(committedFileContent), &b)
	assert.Equal(t, fileContent, b.String())
}

func TestCommit(t *testing.T) {
	// Given
	startTime := secondPrecision(time.Now())
	repoPath := tu.MakeTempDir("sibyl_git_backup_test")
	defer tu.DeleteTempDir(repoPath)
	err := initRepo(repoPath)
	assert.NoError(t, err)

	file1 := filepath.Join(repoPath, "file1.txt")
	os.WriteFile(file1, []byte("hello world!"), 0644)
	file2 := filepath.Join(repoPath, "file2.txt")
	os.WriteFile(file2, []byte("zomk!"), 0644)

	// When
	commit, err := commit(repoPath, "A test commit", "test@example.com", file1, file2)

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, commit.Hash)
	assert.Equal(t, "A test commit", commit.Message)
	assert.Equal(t, "file1.txt", commit.Files[0])
	assert.Equal(t, "file2.txt", commit.Files[1])

	commits, err := listCommits(repoPath)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(commits))
	assert.NotNil(t, commits[0].Hash)
	assert.Equal(t, "A test commit", commits[0].Message)
	assert.Equal(t, "file1.txt", commits[0].Files[0])
	assert.Equal(t, "file2.txt", commits[0].Files[1])
	assert.True(t, !commits[0].Timestamp.Before(startTime),
		"Commit timestamp is greater or equal to %s, but was %s", startTime, commits[0].Timestamp)
}

func TestCommit_NoChanges(t *testing.T) {
	// Given
	repoPath := tu.MakeTempDir("sibyl_git_backup_test")
	defer tu.DeleteTempDir(repoPath)
	err := initRepo(repoPath)
	assert.NoError(t, err)

	file1 := filepath.Join(repoPath, "file1.txt")
	os.WriteFile(file1, []byte("hello world!"), 0644)
	file2 := filepath.Join(repoPath, "file2.txt")
	os.WriteFile(file2, []byte("zomk!"), 0644)

	_, err = commit(repoPath, "A test commit", "test@example.com", file1, file2)
	assert.NoError(t, err)

	// When
	_, err = commit(repoPath, "A test commit 2", "test@example.com", file1, file2)

	// Then
	assert.NoError(t, err)
	commits, err := listCommits(repoPath)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(commits))
	assert.NotEqual(t, commits[0].Hash, commits[1].Hash)
	assert.Equal(t, "A test commit 2", commits[0].Message)
	assert.Equal(t, "file1.txt", commits[0].Files[0])
	assert.Equal(t, "file2.txt", commits[0].Files[1])
	assert.Equal(t, "A test commit", commits[1].Message)
	assert.Equal(t, "file1.txt", commits[1].Files[0])
	assert.Equal(t, "file2.txt", commits[1].Files[1])
}

func TestRevert(t *testing.T) {
	// Given
	repoPath := tu.MakeTempDir("sibyl_git_backup_test")
	defer tu.DeleteTempDir(repoPath)
	err := initRepo(repoPath)
	assert.NoError(t, err)

	file1 := filepath.Join(repoPath, "file1.txt")
	for i := 0; i < 5; i++ {
		os.WriteFile(file1, []byte(fmt.Sprintf("Content %d", i)), 0644)
		commit(repoPath, fmt.Sprintf("Commit %d", i), "test@example.com", file1)
	}
	commits, _ := listCommits(repoPath)
	assert.Equal(t, 5, len(commits))
	todoContent, _ := util.ReadFile(file1)
	assert.Equal(t, "Content 4", todoContent)

	// When
	assert.Equal(t, "Commit 2", commits[2].Message)
	revertCommit, err := revertToCommit(repoPath, commits[2].Hash, "Revert commit", "test@example.com")

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "Revert commit", revertCommit.Message)
	commitsAfterRevert, _ := listCommits(repoPath)
	assert.Equal(t, 6, len(commitsAfterRevert))
	todoContentAfterRevert, _ := util.ReadFile(file1)
	assert.Equal(t, "Content 2", todoContentAfterRevert)

}

func secondPrecision(dt time.Time) time.Time {
	return time.Date(dt.Year(), dt.Month(), dt.Day(),
		dt.Hour(), dt.Minute(), dt.Second(), 0,
		dt.Location())
}
