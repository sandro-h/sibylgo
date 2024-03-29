package extsources

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/sandro-h/sibylgo/backup"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/sandro-h/sibylgo/util"
	"github.com/stretchr/testify/assert"
)

const originalTodos = `------------------
 Today
------------------

[] bla bla
[] zonk
-------------------
 This week
-------------------

[] bink
`

const todosWithPR = `------------------
 Today
------------------

[] bla bla
[] zonk
[] Reviews - 3 PRs #ext_bbprs
-------------------
 This week
-------------------

[] bink
`

const todosWithUpdatedPR = `------------------
 Today
------------------

[] bla bla
[] zonk
[] Reviews - 10 PRs #ext_bbprs
-------------------
 This week
-------------------

[] bink
`

const testConfig = `
bitbucket_prs:
  bb_url: %s
  bb_user: myuser
  bb_token: aba1234
  category: Today`

const todosWithDummies = `------------------
 Today
------------------

[] bla bla
[] zonk
[] dummy1 #ext_id1
[] dummy2 #ext_id2
-------------------
 This week
-------------------

[] bink
`

// some deliberate extra newlines to verify the logic doesn't continually
// try to reformat the file even though the external entries already exist:
const todosOddlySpacedDummies = `------------------
 Today
------------------

[] bla bla
[] zonk
[] dummy1 #ext_id1

[] dummy2 #ext_id2


-------------------
 This week
-------------------

[] bink
`

const testConfigWithDummies = `
dummies:
  category: Today
  dummy_moments:
    - id1:dummy1
    - id2:dummy2`

func TestExternalSources_NewMoment(t *testing.T) {
	ts := tu.MockSimpleJSONResponse(`{"count": 3}`)
	defer ts.Close()

	cfg, _ := util.LoadConfigString(fmt.Sprintf(testConfig, ts.URL))

	updatedTodo, err := FetchAndApplyExternalSourceMoments(originalTodos, cfg)

	assert.Nil(t, err)
	assert.Equal(t, todosWithPR, updatedTodo)
}

func TestExternalSources_UpdatedMoment(t *testing.T) {
	ts := tu.MockSimpleJSONResponse(`{"count": 3}`)

	cfg, _ := util.LoadConfigString(fmt.Sprintf(testConfig, ts.URL))

	updatedTodo, _ := FetchAndApplyExternalSourceMoments(originalTodos, cfg)
	assert.Equal(t, todosWithPR, updatedTodo)

	ts.Close()
	ts = tu.MockSimpleJSONResponse(`{"count": 10}`)
	defer ts.Close()
	cfg, _ = util.LoadConfigString(fmt.Sprintf(testConfig, ts.URL))
	updatedTodo, err := FetchAndApplyExternalSourceMoments(updatedTodo, cfg)

	assert.Nil(t, err)
	assert.Equal(t, todosWithUpdatedPR, updatedTodo)
}

func TestExternalSources_DroppedMoment(t *testing.T) {
	ts := tu.MockSimpleJSONResponse(`{"count": 3}`)

	cfg, _ := util.LoadConfigString(fmt.Sprintf(testConfig, ts.URL))

	updatedTodo, _ := FetchAndApplyExternalSourceMoments(originalTodos, cfg)
	assert.Equal(t, todosWithPR, updatedTodo)

	ts.Close()
	ts = tu.MockSimpleJSONResponse(`{"count": 0}`)
	defer ts.Close()
	cfg, _ = util.LoadConfigString(fmt.Sprintf(testConfig, ts.URL))
	updatedTodo, err := FetchAndApplyExternalSourceMoments(updatedTodo, cfg)

	assert.Nil(t, err)
	assert.Equal(t, originalTodos, updatedTodo+"\n")
}

func TestExternalSources_FileBackup(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_external_sources_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, originalTodos)
	files := util.NewFileConfigFromTodoFile(todoFile)
	cfg, _ := util.LoadConfigString(testConfigWithDummies)

	// When
	p := NewExternalSourcesProcess(files, cfg)
	p.CheckOnce()

	// Then
	updatedContent, _ := util.ReadFile(todoFile)
	assert.Equal(t, todosWithDummies, updatedContent)

	backups, _ := backup.ListBackups(files.TodoDir)
	assert.Equal(t, 1, len(backups))
	assert.Equal(t, "Backup before applying external source changes", backups[0].Message)
}

func TestExternalSources_NoFileBackupIfNoChanges(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_external_sources_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	files := util.NewFileConfigFromTodoFile(todoFile)
	// Note we're already writing the todofile with dummies here
	util.WriteFile(todoFile, todosOddlySpacedDummies)
	cfg, _ := util.LoadConfigString(testConfigWithDummies)

	// When
	p := NewExternalSourcesProcess(files, cfg)
	p.CheckOnce()

	// Then
	updatedContent, _ := util.ReadFile(files.TodoFile)
	assert.Equal(t, todosOddlySpacedDummies, updatedContent)

	backups, _ := backup.ListBackups(files.TodoDir)
	assert.Equal(t, 0, len(backups))
}

func TestExternalSources_NoFileBackupIfNoChanges_NoTrailingNewline(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_external_sources_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	// Dummy todos, but no trailing newline at the end:
	todosWithoutNewline := regexp.MustCompile("\n$").ReplaceAllString(todosWithDummies, "")
	// Don't use util.WriteFile since it always adds a trailing newline
	os.WriteFile(todoFile, []byte(todosWithoutNewline), 0644)
	files := util.NewFileConfigFromTodoFile(todoFile)

	cfg, _ := util.LoadConfigString(testConfigWithDummies)

	// When
	p := NewExternalSourcesProcess(files, cfg)
	p.CheckOnce()

	// Then
	updatedContent, _ := util.ReadFile(files.TodoFile)
	assert.Equal(t, todosWithDummies, updatedContent)

	backups, _ := backup.ListBackups(files.TodoDir)
	assert.Equal(t, 0, len(backups))
}
