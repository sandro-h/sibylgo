package generate

import (
	"github.com/sandro-h/sibylgo/moment"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIterateDaily(t *testing.T) {
	it := NewRecurIterator(re(moment.RecurDaily, "02.01.2019"),
		tu.Dt("10.01.2019"), tu.Dt("13.01.2019"))
	assertIterations(t, it,
		"10.01.2019",
		"11.01.2019",
		"12.01.2019",
		"13.01.2019")
}

func TestIterateWeekly(t *testing.T) {
	it := NewRecurIterator(re(moment.RecurWeekly, "02.01.2019"), // wednesday
		tu.Dt("10.01.2019"), tu.Dt("31.01.2019"))
	assertIterations(t, it,
		"16.01.2019",
		"23.01.2019",
		"30.01.2019")
}

func TestIterateMonthly(t *testing.T) {
	it := NewRecurIterator(re(moment.RecurMonthly, "02.01.2019"),
		tu.Dt("10.01.2019"), tu.Dt("30.04.2019"))
	assertIterations(t, it,
		"02.02.2019",
		"02.03.2019",
		"02.04.2019")
}

func TestIterateYearly(t *testing.T) {
	it := NewRecurIterator(re(moment.RecurYearly, "02.01.2019"),
		tu.Dt("10.01.2019"), tu.Dt("30.04.2022"))
	assertIterations(t, it,
		"02.01.2020",
		"02.01.2021",
		"02.01.2022")
}

func re(re int, d string) moment.Recurrence {
	return moment.Recurrence{Recurrence: re, RefDate: &moment.Date{Time: tu.Dt(d)}}
}

func assertIterations(t *testing.T, it *RecurIterator, expected ...string) {
	var vals []time.Time
	for it.HasNext() {
		vals = append(vals, it.Next())
	}
	assert.Equal(t, len(expected), len(vals))
	for i, v := range vals {
		assert.Equal(t, expected[i], tu.Dts(v))
	}
}
