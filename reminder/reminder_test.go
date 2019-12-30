package reminder

import (
	"github.com/sandro-h/sibylgo/instances"
	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterMomentsEndingInRange(t *testing.T) {
	todos, _ := parse.String(`
[] foo
[] bar (2.2.19)
	[] hello
[] world
	[] yes (2.2.19)
	[] mop
[] more (4.2.19)`)
	insts := instances.GenerateFiltered(todos, tu.Dt("01.02.2019"), tu.Dt("03.02.2019"), nil)
	res := FilterMomentsEndingInRange(insts)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "bar", res[0].Name)
	assert.Equal(t, 0, len(res[0].SubInstances))
	assert.Equal(t, "world", res[1].Name)
	assert.Equal(t, 1, len(res[1].SubInstances))
	assert.Equal(t, "yes", res[1].SubInstances[0].Name)
}

func TestCompileMomentsEndingInRange(t *testing.T) {
	todos, _ := parse.String(`
[] foo
[] bar (2.2.19)
	[] hello
[] world
	[] yes (2.2.19)
	[] mop
[x] bar2 (2.2.19)
[] world2
	[x] yes (2.2.19)
	[] mop
[] more (4.2.19)`)
	res := CompileMomentsEndingInRange(todos, tu.Dt("01.02.2019"), tu.Dt("03.02.2019"))
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "bar", res[0].Name)
	assert.Equal(t, 0, len(res[0].SubInstances))
	assert.Equal(t, "world", res[1].Name)
	assert.Equal(t, 1, len(res[1].SubInstances))
	assert.Equal(t, "yes", res[1].SubInstances[0].Name)
}

func TestCompileRemindersForTodayAndThisWeek(t *testing.T) {
	todos, _ := parse.String(`
[] foo
[] start of week (28.01.2019)
[] today (30.01.2019)
[] this week (02.02.2019)
[] end of week (03.02.2019)
[] next week (04.02.2019)`)
	todays, weeks := CompileRemindersForTodayAndThisWeek(todos, tu.Dt("30.01.2019"))
	assert.Equal(t, 1, len(todays))
	assert.Equal(t, "today", todays[0].Name)
	assert.Equal(t, 4, len(weeks))
	assert.Equal(t, "start of week", weeks[0].Name)
	assert.Equal(t, "today", weeks[1].Name)
	assert.Equal(t, "this week", weeks[2].Name)
	assert.Equal(t, "end of week", weeks[3].Name)
}
