package instances

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"time"
)

// Instance is an actual moment at a point or range in time based on the moment definition.
// For example, a weekly recurring moment definition can yield multiple instances, one for every week
// in the given time range.
type Instance struct {
	Name         string           `json:"name"`
	Start        time.Time        `json:"start"`
	End          time.Time        `json:"end"`
	TimeOfDay    *time.Time       `json:"timeOfDay"`
	Priority     int              `json:"priority"`
	Category     *moment.Category `json:"-"`
	Done         bool             `json:"done"`
	EndsInRange  bool             `json:"endsInRange"`
	SubInstances []*Instance      `json:"subInstances"`
}

// CloneShallow creates a clone of the moment instances without its sub instances.
func (m *Instance) CloneShallow() *Instance {
	c := Instance{
		Name:        m.Name,
		Start:       m.Start,
		Priority:    m.Priority,
		Done:        m.Done,
		EndsInRange: m.EndsInRange}
	if m.TimeOfDay != nil {
		cp := *m.TimeOfDay
		c.TimeOfDay = &cp
	}
	return &c
}

// MomentFilterFunc takes a moment instance and returns true if it should be used,
// false if not. This means it filters on a generated instance and not on the moment
// definition (for example, it could use the effective instance timestamp).
type MomentFilterFunc func(*Instance) bool

// Generate generates moment instances for the given moment in the given time range.
// See Instance for more information.
func Generate(mom moment.Moment, from time.Time, to time.Time) []*Instance {
	return generateInstances(mom, from, to, true, nil)
}

// GenerateWithoutSubs generates moment instances for the given moment in the given time range
// but without any of the given moment's sub moments. For more information on instance generation, see
// Instances.
func GenerateWithoutSubs(mom moment.Moment, from time.Time, to time.Time) []*Instance {
	return generateInstances(mom, from, to, false, nil)
}

// GenerateFiltered generates moment instances for the given moment in the given time range and only
// keeping the moments that match the given filter function. Note that sub-moments of a skipped moment
// are not evaluated, so even if a filter would let a particular sub moment pass, it won't be present if its
// parent was skipped.
func GenerateFiltered(todos *moment.Todos, from time.Time, to time.Time, filter MomentFilterFunc) []*Instance {
	return generateInstancesFiltered(todos, from, to, filter, true)
}

// GenerateFilteredWithoutSubs generates moment instances for the given moment in the given time range and only
// keeping the moments that match the given filter function and without any of the given moment's sub moments.
func GenerateFilteredWithoutSubs(todos *moment.Todos, from time.Time, to time.Time, filter MomentFilterFunc) []*Instance {
	return generateInstancesFiltered(todos, from, to, filter, false)
}

func generateInstancesFiltered(todos *moment.Todos, from time.Time, to time.Time,
	filter MomentFilterFunc, inclSubs bool) []*Instance {
	var insts []*Instance
	for _, mom := range todos.Moments {
		insts = append(insts, generateInstances(mom, from, to, true, filter)...)
	}
	return insts
}

func generateInstances(mom moment.Moment, from time.Time, to time.Time,
	inclSubs bool, filter MomentFilterFunc) []*Instance {
	insts := createInstances(mom, from, to)
	if filter != nil {
		insts = filterInstances(insts, filter)
	}
	// Sub moments:
	if inclSubs {
		for _, inst := range insts {
			var subInsts []*Instance
			for _, sub := range mom.GetSubMoments() {
				subInsts = append(subInsts,
					generateInstances(sub, inst.Start, inst.End, inclSubs, filter)...)
			}
			inst.SubInstances = subInsts
		}
	}
	return insts
}

func filterInstances(insts []*Instance, filter MomentFilterFunc) []*Instance {
	var result []*Instance
	for _, i := range insts {
		if filter(i) {
			result = append(result, i)
		}
	}
	return result
}

func createInstances(mom moment.Moment, from time.Time, to time.Time) []*Instance {
	switch v := mom.(type) {
	case *moment.SingleMoment:
		return createSingleInstances(v, from, to)
	case *moment.RecurMoment:
		return createRecurInstances(v, from, to)
	}
	return nil
}

func createSingleInstances(mom *moment.SingleMoment, from time.Time, to time.Time) []*Instance {
	start := util.GetUpperBound(&from, dateTm(mom.Start))
	end := util.GetLowerBound(&to, dateTm(mom.End))
	if end.Before(start) {
		// Not actually in range
		return nil
	}

	inst := Instance{Name: mom.GetName(), Start: start, End: end}
	inst.Priority = mom.GetPriority()
	inst.Category = mom.GetCategory()
	inst.Done = mom.IsDone()
	inst.EndsInRange = mom.End != nil && !mom.End.Time.After(end)
	if mom.TimeOfDay != nil {
		tm := util.SetTime(start, mom.TimeOfDay.Time)
		inst.TimeOfDay = &tm
	}
	return []*Instance{&inst}
}

func createRecurInstances(mom *moment.RecurMoment, from time.Time, to time.Time) []*Instance {
	var insts []*Instance
	for it := NewRecurIterator(mom.Recurrence, from, to); it.HasNext(); {
		start := it.Next()
		inst := Instance{Name: mom.GetName(), Start: start, End: util.SetToEndOfDay(start)}
		inst.Priority = mom.GetPriority()
		inst.Category = mom.GetCategory()
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
