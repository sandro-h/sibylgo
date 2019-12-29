package extsources

import (
	"fmt"
	"github.com/sandro-h/sibylgo/backup"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/sandro-h/sibylgo/util"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
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

	updatedTodo, err := FetchAndApplyExternalSourceMoments(originalTodos, cfg)
	assert.Equal(t, todosWithPR, updatedTodo)

	ts.Close()
	ts = tu.MockSimpleJSONResponse(`{"count": 10}`)
	defer ts.Close()
	cfg, _ = util.LoadConfigString(fmt.Sprintf(testConfig, ts.URL))
	updatedTodo, err = FetchAndApplyExternalSourceMoments(updatedTodo, cfg)

	assert.Nil(t, err)
	assert.Equal(t, todosWithUpdatedPR, updatedTodo)
}

func TestExternalSources_DroppedMoment(t *testing.T) {
	ts := tu.MockSimpleJSONResponse(`{"count": 3}`)

	cfg, _ := util.LoadConfigString(fmt.Sprintf(testConfig, ts.URL))

	updatedTodo, err := FetchAndApplyExternalSourceMoments(originalTodos, cfg)
	assert.Equal(t, todosWithPR, updatedTodo)

	ts.Close()
	ts = tu.MockSimpleJSONResponse(`{"count": 0}`)
	defer ts.Close()
	cfg, _ = util.LoadConfigString(fmt.Sprintf(testConfig, ts.URL))
	updatedTodo, err = FetchAndApplyExternalSourceMoments(updatedTodo, cfg)

	assert.Nil(t, err)
	assert.Equal(t, originalTodos, updatedTodo+"\n")
}

func TestExternalSources_FileBackup(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_external_sources_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	util.WriteFile(todoFile, originalTodos)
	cfg, _ := util.LoadConfigString(testConfigWithDummies)

	// When
	p := NewExternalSourcesProcess(todoFile, cfg)
	p.CheckOnce()

	// Then
	updatedContent, _ := util.ReadFile(todoFile)
	assert.Equal(t, todosWithDummies, updatedContent)

	backups, _ := backup.ListBackups(todoFile)
	assert.Equal(t, 1, len(backups))
	assert.Equal(t, "Backup before applying external source changes", backups[0].Message)
}

func TestExternalSources_NoFileBackupIfNoChanges(t *testing.T) {
	// Given
	todoDir := tu.MakeTempDir("sibyl_external_sources_test")
	defer tu.DeleteTempDir(todoDir)
	todoFile := filepath.Join(todoDir, "todo.txt")
	// Note we're already writing the todofile with dummies here
	util.WriteFile(todoFile, todosWithDummies)
	cfg, _ := util.LoadConfigString(testConfigWithDummies)

	// When
	p := NewExternalSourcesProcess(todoFile, cfg)
	p.CheckOnce()

	// Then
	updatedContent, _ := util.ReadFile(todoFile)
	assert.Equal(t, todosWithDummies, updatedContent)

	backups, _ := backup.ListBackups(todoFile)
	assert.Equal(t, 0, len(backups))
}
