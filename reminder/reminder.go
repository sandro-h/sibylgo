package reminder

import (
	"github.com/sandro-h/sibylgo/generate"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"sort"
	"time"
)

func CompileRemindersForTodayAndThisWeek(todos *moment.Todos, today time.Time) ([]*moment.MomentInstance, []*moment.MomentInstance) {
	todaysReminders := CompileMomentsEndingInRange(todos, util.SetToStartOfDay(today), util.SetToEndOfDay(today))
	weeksReminders := CompileMomentsEndingInRange(todos, util.SetToStartOfWeek(today), util.SetToEndOfWeek(today))
	sort.Sort(ByStartDate(weeksReminders))
	return todaysReminders, weeksReminders
}

func CompileMomentsEndingInRange(todos *moment.Todos, from time.Time, to time.Time) []*moment.MomentInstance {
	insts := generate.GenerateInstancesFiltered(todos, from, to, func(mom *moment.MomentInstance) bool { return !mom.Done })
	return FilterMomentsEndingInRange(insts)
}

func FilterMomentsEndingInRange(insts []*moment.MomentInstance) []*moment.MomentInstance {
	var res []*moment.MomentInstance
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

type ByStartDate []*moment.MomentInstance

func (a ByStartDate) Len() int           { return len(a) }
func (a ByStartDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStartDate) Less(i, j int) bool { return a[i].Start.Before(a[j].Start) }
