package stringify

import (
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
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
	doneMarker := ""
	if m.IsDone() {
		doneMarker = "x"
	}

	idSuffix := ""
	if m.GetID() != nil {
		idSuffix = fmt.Sprintf(" #%s", m.GetID().Value)
	}

	dateSuffix := stringifyDate(m)

	prioritySuffix := ""
	if m.GetPriority() > 0 {
		// TODO
		panic("priority stringify not implemented")
	}

	res := fmt.Sprintf("%s[%s] %s%s%s%s\n", indent, doneMarker, m.GetName(), prioritySuffix, dateSuffix, idSuffix)
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
			panic("date stringify not implemented")
		}
		if v.End != nil && (v.Start == nil || v.End.DocCoords != v.Start.DocCoords) {
			// TODO
			panic("date stringify not implemented")
		}
	case *moment.RecurMoment:
		// TODO
		panic("date stringify not implemented")
	}

	if m.GetTimeOfDay() != nil {
		// TODO
		panic("date stringify not implemented")
	}

	return ""
}
