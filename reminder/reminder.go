package reminder

import (
	"github.com/sandro-h/sibylgo/generate"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"sort"
	"time"
)

func CompileRemindersForTodayAndThisWeek(todos *moment.Todos, today time.Time) ([]*moment.Instance, []*moment.Instance) {
	todaysReminders := CompileMomentsEndingInRange(todos, util.SetToStartOfDay(today), util.SetToEndOfDay(today))
	weeksReminders := CompileMomentsEndingInRange(todos, util.SetToStartOfWeek(today), util.SetToEndOfWeek(today))
	sort.Sort(ByStartDate(weeksReminders))
	return todaysReminders, weeksReminders
}

func CompileMomentsEndingInRange(todos *moment.Todos, from time.Time, to time.Time) []*moment.Instance {
	insts := generate.InstancesFiltered(todos, from, to, func(mom *moment.Instance) bool { return !mom.Done })
	return FilterMomentsEndingInRange(insts)
}

func FilterMomentsEndingInRange(insts []*moment.Instance) []*moment.Instance {
	// Explicitly make it a 0-len array, otherwise it's 'nil' and will be converted
	// to null by the JSON encoder.
	res := make([]*moment.Instance, 0)
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

type ByStartDate []*moment.Instance

func (a ByStartDate) Len() int           { return len(a) }
func (a ByStartDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStartDate) Less(i, j int) bool { return a[i].Start.Before(a[j].Start) }
