package reminder

import (
	"github.com/sandro-h/sibylgo/instances"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"sort"
	"time"
)

// CompileRemindersForTodayAndThisWeek returns a list of moments that are due today, and a list of moments
// that are due this week.
func CompileRemindersForTodayAndThisWeek(todos *moment.Todos, today time.Time) ([]*instances.Instance, []*instances.Instance) {
	todaysReminders := CompileMomentsEndingInRange(todos, util.SetToStartOfDay(today), util.SetToEndOfDay(today))
	weeksReminders := CompileMomentsEndingInRange(todos, util.SetToStartOfWeek(today), util.SetToEndOfWeek(today))
	sort.Sort(byStartDate(weeksReminders))
	return todaysReminders, weeksReminders
}

// CompileMomentsEndingInRange returns a list of moments that are due in the given time range.
func CompileMomentsEndingInRange(todos *moment.Todos, from time.Time, to time.Time) []*instances.Instance {
	insts := instances.GenerateFiltered(todos, from, to, func(mom *instances.Instance) bool { return !mom.Done })
	return FilterMomentsEndingInRange(insts)
}

// FilterMomentsEndingInRange keeps only the moments and sub moments that have the EndsInRange
// flag set.
func FilterMomentsEndingInRange(insts []*instances.Instance) []*instances.Instance {
	// Explicitly make it a 0-len array, otherwise it's 'nil' and will be converted
	// to null by the JSON encoder.
	res := make([]*instances.Instance, 0)
	for _, m := range insts {
		subs := FilterMomentsEndingInRange(m.SubInstances)
		if len(subs) > 0 || m.EndsInRange {
			c := m.CloneShallow()
			c.SubInstances = subs
			res = append(res, c)
		}
	}
	return res
}

type byStartDate []*instances.Instance

func (a byStartDate) Len() int           { return len(a) }
func (a byStartDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byStartDate) Less(i, j int) bool { return a[i].Start.Before(a[j].Start) }
