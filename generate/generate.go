package generate

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"time"
)

type MomentFilterFunc func(moment.Moment) bool

func GenerateInstances(mom moment.Moment, from time.Time, to time.Time) []*moment.MomentInstance {
	return generateInstances(mom, from, to, true, nil)
}

func GenerateInstancesFiltered(todos *moment.Todos, from time.Time, to time.Time,
	filter MomentFilterFunc) []*moment.MomentInstance {
	var insts []*moment.MomentInstance
	for _, mom := range todos.Moments {
		insts = append(insts, generateInstances(mom, from, to, true, filter)...)
	}
	return insts
}

func GenerateInstancesWithoutSubs(mom moment.Moment, from time.Time, to time.Time) []*moment.MomentInstance {
	return generateInstances(mom, from, to, false, nil)
}

func generateInstances(mom moment.Moment, from time.Time, to time.Time,
	inclSubs bool, filter MomentFilterFunc) []*moment.MomentInstance {
	if filter != nil && !filter(mom) {
		return nil
	}
	insts := createInstances(mom, from, to)
	// Sub moments:
	if inclSubs {
		for _, inst := range insts {
			var subInsts []*moment.MomentInstance
			for _, sub := range mom.GetSubMoments() {
				subInsts = append(subInsts,
					generateInstances(sub, inst.Start, inst.End, inclSubs, filter)...)
			}
			inst.SubInstances = subInsts
		}
	}
	return insts
}

func createInstances(mom moment.Moment, from time.Time, to time.Time) []*moment.MomentInstance {
	switch v := mom.(type) {
	case *moment.SingleMoment:
		return createSingleInstances(v, from, to)
	case *moment.RecurMoment:
		return createRecurInstances(v, from, to)
	}
	return nil
}

func createSingleInstances(mom *moment.SingleMoment, from time.Time, to time.Time) []*moment.MomentInstance {
	start := util.GetUpperBound(&from, dateTm(mom.Start))
	end := util.GetLowerBound(&to, dateTm(mom.End))
	if end.Before(start) {
		// Not actually in range
		return nil
	}

	inst := moment.MomentInstance{Name: mom.GetName(), Start: start, End: end}
	inst.EndsInRange = mom.End != nil && !mom.End.Time.After(end)
	return []*moment.MomentInstance{&inst}
}

func createRecurInstances(mom *moment.RecurMoment, from time.Time, to time.Time) []*moment.MomentInstance {
	var insts []*moment.MomentInstance
	for it := NewRecurIterator(mom.Recurrence, from, to); it.HasNext(); {
		start := it.Next()
		inst := moment.MomentInstance{Name: mom.GetName(), Start: start, End: util.SetToEndOfDay(start)}
		inst.EndsInRange = true
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
