package generate

import (
	"github.com/sandro-h/sibylgo/moment"
	"time"
)

func GenerateInstances(mom moment.Moment, from time.Time, to time.Time) []*moment.MomentInstance {
	return generateInstances(mom, from, to, true)
}

func GenerateInstancesWithoutSubs(mom moment.Moment, from time.Time, to time.Time) []*moment.MomentInstance {
	return generateInstances(mom, from, to, false)
}

func generateInstances(mom moment.Moment, from time.Time, to time.Time, inclSubs bool) []*moment.MomentInstance {
	insts := mom.CreateInstances(from, to)
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
