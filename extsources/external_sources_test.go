package extsources

import (
	"fmt"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/sandro-h/sibylgo/util"
	"github.com/stretchr/testify/assert"
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

[] bink`

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
  bb_token: aba1234
  category: Today`

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
	assert.Equal(t, originalTodos, updatedTodo)
}
