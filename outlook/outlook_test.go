package outlook

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
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
[] ranged moment (-4.12.20)
[] recurring moment (every day)`,
			outlook: "",
			added:   []string{"bla", "foo"},
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

func names(moms []*moment.SingleMoment) []string {
	var res []string
	for _, m := range moms {
		res = append(res, m.GetName())
	}
	return res
}
