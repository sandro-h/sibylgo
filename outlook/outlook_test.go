package outlook

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testcase struct {
	current string
	outlook string
	added   []string
	updated []string
	removed []string
}

func TestComputeDiff(t *testing.T) {
	cases := []testcase{
		testcase{
			current: `
[] undated moment
[] bla (4.12.20)
[] foo (4.12.20 8:00)
[x] done moment (4.12.20)
[] due moment (-4.12.20)
[] ranged moment (2.12.20-4.12.20)
[] recurring moment (every day)`,
			outlook: "",
			added:   []string{"bla", "foo", "due moment"},
			updated: []string(nil),
			removed: []string(nil),
		},
		testcase{
			current: `
[] undated moment
[] bla (4.12.20)
[] foo (4.12.20 8:00)
[] recurring moment (every day)`,
			outlook: `
[] bla (4.12.20)
[] zonk (5.12.20)`,
			added:   []string{"foo"},
			updated: []string(nil),
			removed: []string{"zonk"},
		},
		testcase{
			current: `
[] undated moment
[] bla (5.12.20)
[] foo (4.12.20 8:30)
[] recurring moment (every day)`,
			outlook: `
[] bla (4.12.20)
[] foo (4.12.20 8:00)`,
			added:   []string(nil),
			updated: []string{"bla", "foo"},
			removed: []string(nil),
		},
	}

	for _, tc := range cases {
		todos, _ := parse.String(tc.current)
		outlookTodos, _ := parse.String(tc.outlook)

		currentMoms := filterEligibleForOutlook(todos.Moments)
		outlookMoms := filterEligibleForOutlook(outlookTodos.Moments)

		added, updated, removed := computeDiff(currentMoms, outlookMoms)
		assert.Equal(t, tc.added, names(added))
		assert.Equal(t, tc.updated, names(updated))
		assert.Equal(t, tc.removed, names(removed))
	}

}

func TestDueMoment_ConvertsToSingleDay(t *testing.T) {
	todos, _ := parse.String("[] due moment (-4.12.20 8:00)")

	currentMoms := filterEligibleForOutlook(todos.Moments)

	assert.True(t, moment.IsSingleDayMoment(currentMoms[0]))
	assert.Equal(t, "04.12.2020", tu.Dts(currentMoms[0].Start.Time))
	assert.Equal(t, "04.12.2020", tu.Dts(currentMoms[0].End.Time))
	assert.Equal(t, "08:00:00", tu.Tts(currentMoms[0].TimeOfDay.Time))
}

func TestDueMoment_EqualsSingleDayFromOutlook(t *testing.T) {
	todos, _ := parse.String("[] due moment (-4.12.20 8:00)")
	outlookTodos, _ := parse.String("[] due moment (4.12.20 8:00)")

	currentMoms := filterEligibleForOutlook(todos.Moments)
	outlookMoms := filterEligibleForOutlook(outlookTodos.Moments)

	added, updated, removed := computeDiff(currentMoms, outlookMoms)
	assert.Equal(t, []string(nil), names(added))
	assert.Equal(t, []string(nil), names(updated))
	assert.Equal(t, []string(nil), names(removed))
}

func names(moms []*moment.SingleMoment) []string {
	var res []string
	for _, m := range moms {
		res = append(res, m.GetName())
	}
	return res
}
