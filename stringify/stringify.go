package stringify

import (
	"fmt"

	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
)

// Todos converts the moments into the same string content used in a todo file.
func Todos(todos *moment.Todos) string {

	res := ""
	var lastCat string
	for _, m := range todos.Moments {
		cat := m.GetCategory()
		if cat != nil && cat.Name != lastCat {
			res += stringifyCategory(cat)
			lastCat = cat.Name
		}
		res += Moment(m)
	}
	return res
}

// Moment converts the moment to the same string content used in a todo file.
func Moment(m moment.Moment) string {
	return stringifyMoment(m, false, "")
}

func stringifyCategory(c *moment.Category) string {
	// TODO
	return ""
}

func stringifyMoment(m moment.Moment, parentDone bool, indent string) string {
	stateMarker := ""
	if m.IsDone() {
		stateMarker = string(parse.ParseConfig.GetDoneMark())
	} else {
		switch m.GetWorkState() {
		case moment.InProgressState:
			stateMarker = string(parse.ParseConfig.GetInProgressMark())
		case moment.WaitingState:
			stateMarker = string(parse.ParseConfig.GetWaitingMark())
		}
	}

	idSuffix := ""
	if m.GetID() != nil {
		idSuffix = fmt.Sprintf(" #%s", m.GetID().Value)
	}

	dateSuffix := stringifyDate(m)

	prioritySuffix := ""
	if m.GetPriority() > 0 {
		// TODO
		panicNotImplemented()
	}

	res := fmt.Sprintf("%s%c%s%c %s%s%s%s\n",
		indent,
		parse.ParseConfig.GetLBracket(),
		stateMarker,
		parse.ParseConfig.GetRBracket(),
		m.GetName(),
		prioritySuffix,
		dateSuffix,
		idSuffix)
	for _, c := range m.GetComments() {
		res += fmt.Sprintf("%s%s\n", indent+"\t", c.Content)
	}

	for _, s := range m.GetSubMoments() {
		res += stringifyMoment(s, parentDone || m.IsDone(), indent+"\t")
	}
	return res
}

func stringifyDate(m moment.Moment) string {
	switch v := m.(type) {
	case *moment.SingleMoment:
		if v.Start != nil {
			// TODO
			panicNotImplemented()
		}
		if v.End != nil && (v.Start == nil || v.End.DocCoords != v.Start.DocCoords) {
			// TODO
			panicNotImplemented()
		}
	case *moment.RecurMoment:
		// TODO
		panicNotImplemented()
	}

	if m.GetTimeOfDay() != nil {
		// TODO
		panicNotImplemented()
	}

	return ""
}

func panicNotImplemented() {
	panic("not implemented")
}
