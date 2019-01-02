package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateInSmallerRange(t *testing.T) {
	mom, _ := parseMom("[] bla (18.06.2016-25.06.2016)")
	insts := GenerateInstances(mom, dt("20.06.2016"), dt("22.06.2016"))
	assertInstanceDates(t, insts, "20.06.2016", "22.06.2016")
	assert.False(t, insts[0].endsInRange)
}

func TestGenerateInLargerRange(t *testing.T) {
	mom, _ := parseMom("[] bla (18.06.2016-25.06.2016)")
	insts := GenerateInstances(mom, dt("01.06.2016"), dt("01.08.2016"))
	assertInstanceDates(t, insts, "18.06.2016", "25.06.2016")
	assert.True(t, insts[0].endsInRange)
}

func TestGenerateUnbounded(t *testing.T) {
	mom, _ := parseMom("[] bla")
	insts := GenerateInstances(mom, dt("20.06.2016"), dt("22.06.2016"))
	assertInstanceDates(t, insts, "20.06.2016", "22.06.2016")
	assert.False(t, insts[0].endsInRange)
}

func TestGenerateOutOfRange(t *testing.T) {
	mom, _ := parseMom("[] bla (18.06.2016-25.06.2016)")
	insts := GenerateInstances(mom, dt("01.07.2016"), dt("13.07.2016"))
	assert.Equal(t, 0, len(insts))
}

func TestGenerateChildren(t *testing.T) {
	todos, _ := ParseString(`
[] 1 (18.06.2016-25.06.2016)
	[] 1.1 (20.06.2016-23.06.2016)
	[] 1.2 (18.06.2016-19.06.2016)
`)
	insts := GenerateInstances(todos.moments[0], dt("01.06.2016"), dt("01.08.2016"))
	assertInstanceDates(t, insts[0].subInstances,
		"20.06.2016", "23.06.2016",
		"18.06.2016", "19.06.2016")
}

func TestGenerateChildrenCutOffByParent(t *testing.T) {
	todos, _ := ParseString(`
[] 1 (18.06.2016-25.06.2016)
	[] 1.1 (20.06.2016-30.06.2016)
	[] 1.2 (01.07.2016-05.07.2016)
`)
	// 1.2 should not show up at all
	insts := GenerateInstances(todos.moments[0], dt("01.06.2016"), dt("01.08.2016"))
	assertInstanceDates(t, insts[0].subInstances,
		"20.06.2016", "25.06.2016")
}

func TestGenerateChildrenCutOffByRange(t *testing.T) {
	todos, _ := ParseString(`
[] 1 (18.06.2016-25.06.2016)
	[] 1.1 (20.06.2016-23.06.2016)
`)
	insts := GenerateInstances(todos.moments[0], dt("01.06.2016"), dt("21.06.2016"))
	assertInstanceDates(t, insts[0].subInstances,
		"20.06.2016", "21.06.2016")
}

func assertInstanceDates(t *testing.T, insts []*MomentInstance, dates ...string) {
	assert.Equal(t, len(dates)/2, len(insts))
	for i := 0; i < len(dates); i += 2 {
		inst := insts[i/2]
		assert.Equal(t, dates[i], dts(inst.start))
		assert.Equal(t, dates[i+1], dts(inst.end))
	}
}

func dt(s string) time.Time {
	d, _ := time.Parse("02.01.2006", s)
	return d
}

func dts(t time.Time) string {
	return t.Format("02.01.2006")
}
