package generate

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateInSmallerRange(t *testing.T) {
	todos, _ := parse.ParseString("[] bla (18.06.2016-25.06.2016)")
	insts := GenerateInstances(todos.Moments[0], tu.Dt("20.06.2016"), tu.Dt("22.06.2016"))
	assertInstanceDates(t, insts, "20.06.2016", "22.06.2016")
	assert.False(t, insts[0].EndsInRange)
}

func TestGenerateInLargerRange(t *testing.T) {
	todos, _ := parse.ParseString("[] bla (18.06.2016-25.06.2016)")
	insts := GenerateInstances(todos.Moments[0], tu.Dt("01.06.2016"), tu.Dt("01.08.2016"))
	assertInstanceDates(t, insts, "18.06.2016", "25.06.2016")
	assert.True(t, insts[0].EndsInRange)
}

func TestGenerateUnbounded(t *testing.T) {
	todos, _ := parse.ParseString("[] bla")
	insts := GenerateInstances(todos.Moments[0], tu.Dt("20.06.2016"), tu.Dt("22.06.2016"))
	assertInstanceDates(t, insts, "20.06.2016", "22.06.2016")
	assert.False(t, insts[0].EndsInRange)
}

func TestGenerateOutOfRange(t *testing.T) {
	todos, _ := parse.ParseString("[] bla (18.06.2016-25.06.2016)")
	insts := GenerateInstances(todos.Moments[0], tu.Dt("01.07.2016"), tu.Dt("13.07.2016"))
	assert.Equal(t, 0, len(insts))
}

func TestGenerateChildren(t *testing.T) {
	todos, _ := parse.ParseString(`
[] 1 (18.06.2016-25.06.2016)
	[] 1.1 (20.06.2016-23.06.2016)
	[] 1.2 (18.06.2016-19.06.2016)
`)
	insts := GenerateInstances(todos.Moments[0], tu.Dt("01.06.2016"), tu.Dt("01.08.2016"))
	assertInstanceDates(t, insts[0].SubInstances,
		"20.06.2016", "23.06.2016",
		"18.06.2016", "19.06.2016")
}

func TestGenerateChildrenCutOffByParent(t *testing.T) {
	todos, _ := parse.ParseString(`
[] 1 (18.06.2016-25.06.2016)
	[] 1.1 (20.06.2016-30.06.2016)
	[] 1.2 (01.07.2016-05.07.2016)
`)
	// 1.2 should not show up at all
	insts := GenerateInstances(todos.Moments[0], tu.Dt("01.06.2016"), tu.Dt("01.08.2016"))
	assertInstanceDates(t, insts[0].SubInstances,
		"20.06.2016", "25.06.2016")
}

func TestGenerateChildrenCutOffByRange(t *testing.T) {
	todos, _ := parse.ParseString(`
[] 1 (18.06.2016-25.06.2016)
	[] 1.1 (20.06.2016-23.06.2016)
`)
	insts := GenerateInstances(todos.Moments[0], tu.Dt("01.06.2016"), tu.Dt("21.06.2016"))
	assertInstanceDates(t, insts[0].SubInstances,
		"20.06.2016", "21.06.2016")
}

func TestGenerateRecurring(t *testing.T) {
	todos, _ := parse.ParseString("[] bla (every day)")
	insts := GenerateInstances(todos.Moments[0], tu.Dt("20.06.2016"), tu.Dt("22.06.2016"))
	assertInstanceDates(t, insts,
		"20.06.2016", "20.06.2016",
		"21.06.2016", "21.06.2016",
		"22.06.2016", "22.06.2016")
}

func TestGenerateRecurringNotInRange(t *testing.T) {
	todos, _ := parse.ParseString("[] bla (every 23.)")
	insts := GenerateInstances(todos.Moments[0], tu.Dt("20.06.2016"), tu.Dt("22.06.2016"))
	assert.Equal(t, 0, len(insts))
}

func TestGenerateRecurringAsChildren(t *testing.T) {
	todos, _ := parse.ParseString(`
[] 1 (18.06.2016-25.06.2016)
	[] 1.1 (every 20.)
	[] 1.2 (every day)
`)
	insts := GenerateInstances(todos.Moments[0], tu.Dt("01.06.2016"), tu.Dt("20.06.2016"))
	assertInstanceDates(t, insts[0].SubInstances,
		"20.06.2016", "20.06.2016",
		"18.06.2016", "18.06.2016",
		"19.06.2016", "19.06.2016",
		"20.06.2016", "20.06.2016")
}

func TestGenerateRecurringWithChildren(t *testing.T) {
	todos, _ := parse.ParseString(`
[] 1 (every 20.)
	[] 1.1 (every 20.6)
	[] 1.2 (20.7.2016)
	[] 1.3 (21.6.2016)
`)
	insts := GenerateInstances(todos.Moments[0], tu.Dt("01.06.2016"), tu.Dt("30.07.2016"))
	assertInstanceDates(t, insts[0].SubInstances,
		"20.06.2016", "20.06.2016") // 1.1
	assertInstanceDates(t, insts[1].SubInstances,
		"20.07.2016", "20.07.2016") // 1.2
}

func assertInstanceDates(t *testing.T, insts []*moment.MomentInstance, dates ...string) {
	assert.Equal(t, len(dates)/2, len(insts))
	for i := 0; i < len(dates); i += 2 {
		inst := insts[i/2]
		assert.Equal(t, dates[i], tu.Dts(inst.Start))
		assert.Equal(t, dates[i+1], tu.Dts(inst.End))
	}
}
