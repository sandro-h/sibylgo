package generate

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"time"
)

// MomentFilterFunc takes a moment instance and returns true if it should be used,
// false if not. This means it filters on a generated instance and not on the moment
// definition (for example, it could use the effective instance timestamp).
type MomentFilterFunc func(*moment.Instance) bool

// Instances generates moment instances for the given moment in the given time range.
// See Instance for more information.
func Instances(mom moment.Moment, from time.Time, to time.Time) []*moment.Instance {
	return generateInstances(mom, from, to, true, nil)
}

// InstancesWithoutSubs generates moment instances for the given moment in the given time range
// but without any of the given moment's sub moments. For more information on instance generation, see
// Instances.
func InstancesWithoutSubs(mom moment.Moment, from time.Time, to time.Time) []*moment.Instance {
	return generateInstances(mom, from, to, false, nil)
}

// InstancesFiltered generates moment instances for the given moment in the given time range and only
// keeping the moments that match the given filter function. Note that sub-moments of a skipped moment
// are not evaluated, so even if a filter would let a particular sub moment pass, it won't be present if its
// parent was skipped.
func InstancesFiltered(todos *moment.Todos, from time.Time, to time.Time, filter MomentFilterFunc) []*moment.Instance {
	return generateInstancesFiltered(todos, from, to, filter, true)
}

// InstancesFilteredWithoutSubs generates moment instances for the given moment in the given time range and only
// keeping the moments that match the given filter function and without any of the given moment's sub moments.
func InstancesFilteredWithoutSubs(todos *moment.Todos, from time.Time, to time.Time, filter MomentFilterFunc) []*moment.Instance {
	return generateInstancesFiltered(todos, from, to, filter, false)
}

func generateInstancesFiltered(todos *moment.Todos, from time.Time, to time.Time,
	filter MomentFilterFunc, inclSubs bool) []*moment.Instance {
	var insts []*moment.Instance
	for _, mom := range todos.Moments {
		insts = append(insts, generateInstances(mom, from, to, true, filter)...)
	}
	return insts
}

func generateInstances(mom moment.Moment, from time.Time, to time.Time,
	inclSubs bool, filter MomentFilterFunc) []*moment.Instance {
	insts := createInstances(mom, from, to)
	if filter != nil {
		insts = filterInstances(insts, filter)
	}
	// Sub moments:
	if inclSubs {
		for _, inst := range insts {
			var subInsts []*moment.Instance
			for _, sub := range mom.GetSubMoments() {
				subInsts = append(subInsts,
					generateInstances(sub, inst.Start, inst.End, inclSubs, filter)...)
			}
			inst.SubInstances = subInsts
		}
	}
	return insts
}

func filterInstances(insts []*moment.Instance, filter MomentFilterFunc) []*moment.Instance {
	var result []*moment.Instance
	for _, i := range insts {
		if filter(i) {
			result = append(result, i)
		}
	}
	return result
}

func createInstances(mom moment.Moment, from time.Time, to time.Time) []*moment.Instance {
	switch v := mom.(type) {
	case *moment.SingleMoment:
		return createSingleInstances(v, from, to)
	case *moment.RecurMoment:
		return createRecurInstances(v, from, to)
	}
	return nil
}

func createSingleInstances(mom *moment.SingleMoment, from time.Time, to time.Time) []*moment.Instance {
	start := util.GetUpperBound(&from, dateTm(mom.Start))
	end := util.GetLowerBound(&to, dateTm(mom.End))
	if end.Before(start) {
		// Not actually in range
		return nil
	}

	inst := moment.Instance{Name: mom.GetName(), Start: start, End: end}
	inst.Priority = mom.GetPriority()
	inst.Done = mom.IsDone()
	inst.EndsInRange = mom.End != nil && !mom.End.Time.After(end)
	if mom.TimeOfDay != nil {
		tm := util.SetTime(start, mom.TimeOfDay.Time)
		inst.TimeOfDay = &tm
	}
	return []*moment.Instance{&inst}
}

func createRecurInstances(mom *moment.RecurMoment, from time.Time, to time.Time) []*moment.Instance {
	var insts []*moment.Instance
	for it := NewRecurIterator(mom.Recurrence, from, to); it.HasNext(); {
		start := it.Next()
		inst := moment.Instance{Name: mom.GetName(), Start: start, End: util.SetToEndOfDay(start)}
		inst.Priority = mom.GetPriority()
		inst.Done = mom.IsDone()
		inst.EndsInRange = true
		if mom.TimeOfDay != nil {
			tm := util.SetTime(start, mom.TimeOfDay.Time)
			inst.TimeOfDay = &tm
		}
		insts = append(insts, &inst)
	}
	return insts
}

func dateTm(dt *moment.Date) *time.Time {
	if dt == nil {
		return nil
	}
	return &dt.Time
}
