package generate

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"time"
)

func GenerateInstances(mom moment.Moment, from time.Time, to time.Time) []*moment.MomentInstance {
	return generateInstances(mom, from, to, true)
}

func GenerateInstancesWithoutSubs(mom moment.Moment, from time.Time, to time.Time) []*moment.MomentInstance {
	return generateInstances(mom, from, to, false)
}

func generateInstances(mom moment.Moment, from time.Time, to time.Time, inclSubs bool) []*moment.MomentInstance {
	insts := createInstances(mom, from, to)
	// Sub moments:
	if inclSubs {
		for _, inst := range insts {
			var subInsts []*moment.MomentInstance
			for _, sub := range mom.GetSubMoments() {
				subInsts = append(subInsts, generateInstances(sub, inst.Start, inst.End, inclSubs)...)
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

	inst := moment.MomentInstance{Start: start, End: end}
	inst.EndsInRange = mom.End != nil && !mom.End.Time.After(end)
	return []*moment.MomentInstance{&inst}
}

func createRecurInstances(mom *moment.RecurMoment, from time.Time, to time.Time) []*moment.MomentInstance {
	var insts []*moment.MomentInstance
	for it := NewRecurIterator(mom.Recurrence, from, to); it.HasNext(); {
		start := it.Next()
		inst := moment.MomentInstance{Start: start, End: util.SetToEndOfDay(start)}
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
