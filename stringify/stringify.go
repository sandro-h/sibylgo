package stringify

import (
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
)

// ForTodoFile converts the moments into the same string content used in a todo file.
func ForTodoFile(todos *moment.Todos) string {

	res := ""
	var lastCat string
	for _, m := range todos.Moments {
		cat := m.GetCategory()
		if cat != nil && cat.Name != lastCat {
			res += formatCategory(cat)
			lastCat = cat.Name
		}
		res += FormatMoment(m)
	}
	return res
}

// FormatMoment converts the moment to the same string content used in a todo file.
func FormatMoment(m moment.Moment) string {
	return formatMoment(m, false, "")
}

func formatCategory(c *moment.Category) string {
	// TODO
	return ""
}

func formatMoment(m moment.Moment, parentDone bool, indent string) string {
	doneMarker := ""
	if m.IsDone() {
		doneMarker = "x"
	}
	res := fmt.Sprintf("%s[%s] %s\n", indent, doneMarker, m.GetName())
	for _, c := range m.GetComments() {
		res += fmt.Sprintf("%s%s\n", indent+"\t", c.Content)
	}

	// TODO

	for _, s := range m.GetSubMoments() {
		res += formatMoment(s, parentDone || m.IsDone(), indent+"\t")
	}
	return res
}
