package preview

import (
	"time"

	"github.com/sandro-h/sibylgo/calendar"
	"github.com/sandro-h/sibylgo/instances"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/reminder"
	"github.com/sandro-h/sibylgo/util"
)

var getNow = func() time.Time {
	return time.Now()
}

// Create creates a preview of the passed todos. A preview is an overview with things like:
// * Moments due today / due this week
// * All top-level moments by category
// * Week's calendar
// The preview is displayed as HTML in the VSCode extension.
func Create(todos *moment.Todos) Preview {
	now := getNow()

	overview := compileTopLevelMomentsOverview(todos)
	todays, weeks := compileReminders(todos, now)
	entries := calendar.CompileCalendarEntries(todos, util.SetToStartOfWeek(now), util.SetToEndOfWeek(now).AddDate(0, 0, 1))

	return Preview{
		Today:    todays,
		Week:     weeks,
		Overview: overview,
		Calendar: entries}
}

func compileTopLevelMomentsOverview(todos *moment.Todos) jsonTodos {
	var overview jsonTodos
	var curCat *jsonCategory
	for _, m := range todos.Moments {
		if !m.IsDone() {
			catName := "_none"
			if m.GetCategory() != nil {
				catName = m.GetCategory().Name
			}
			if curCat == nil || catName != curCat.Name {
				curCat = &jsonCategory{Name: catName}
				overview.Categories = append(overview.Categories, curCat)
			}
			curCat.Moments = append(curCat.Moments, toJSONMoment(m))
		}
	}
	return overview
}

func compileReminders(todos *moment.Todos, now time.Time) ([]*instances.Instance, []*instances.Instance) {
	todays, weeks := reminder.CompileRemindersForTodayAndThisWeek(todos, now)
	return flattenReminders("", todays), flattenReminders("", weeks)
}

func flattenReminders(parentPath string, insts []*instances.Instance) []*instances.Instance {
	// Explicitly make it a 0-len array, otherwise it's 'nil' and will be converted
	// to null by the JSON encoder.
	res := make([]*instances.Instance, 0)
	for _, i := range insts {
		if i.EndsInRange {
			flatInst := i.CloneShallow()
			flatInst.Name = parentPath + flatInst.Name
			res = append(res, flatInst)
		}
		subFlattened := flattenReminders(parentPath+i.Name+"/", i.SubInstances)
		res = append(res, subFlattened...)
	}
	return res
}

// Preview holds the contents of the preview.
type Preview struct {
	Today    []*instances.Instance `json:"today"`
	Week     []*instances.Instance `json:"week"`
	Overview jsonTodos             `json:"overview"`
	Calendar []calendar.Entry      `json:"calendar"`
}

type jsonTodos struct {
	Categories []*jsonCategory `json:"categories"`
}

type jsonCategory struct {
	Name    string       `json:"name"`
	Moments []jsonMoment `json:"moments"`
}

type jsonMoment struct {
	Name      string           `json:"name"`
	WorkState moment.WorkState `json:"workState"`
	DocCoords moment.DocCoords `json:"docCoords"`
}

func toJSONMoment(mom moment.Moment) jsonMoment {
	return jsonMoment{
		Name:      mom.GetName(),
		WorkState: mom.GetWorkState(),
		DocCoords: mom.GetDocCoords(),
	}
}
