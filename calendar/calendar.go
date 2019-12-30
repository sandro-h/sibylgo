package calendar

import (
	"github.com/sandro-h/sibylgo/instances"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"sort"
	"time"
)

const dateFormat = "2006-01-02"

// Entry contains data for a single calendar entry.
type Entry struct {
	Title string `json:"title"`
	Start string `json:"start"`
	End   string `json:"end"`
	Color string `json:"color,omitempty"`
}

// NewEntry maps a moment instance to a new calendar entry.
func NewEntry(inst *instances.Instance) Entry {
	entry := Entry{
		Title: inst.Name,
		Start: inst.Start.Format(dateFormat),
		End:   inst.End.AddDate(0, 0, 1).Format(dateFormat)} // +1 because fullcalendar is non-inclusive
	if inst.Category != nil {
		entry.Color = inst.Category.Color
	}
	return entry
}

// CompileCalendarEntries maps all moment instances that are not done and end in the given
// time range to calendar entries.
func CompileCalendarEntries(todos *moment.Todos, from time.Time, until time.Time) []Entry {
	from = util.SetToStartOfDay(from)
	until = util.SetToStartOfDay(until)

	insts := instances.GenerateFilteredWithoutSubs(todos, from, until,
		func(mom *instances.Instance) bool { return !mom.Done && mom.EndsInRange })
	sort.Sort(byPriority(insts))

	var entries []Entry
	for _, inst := range insts {
		entries = append(entries, NewEntry(inst))
	}
	return entries
}

type byPriority []*instances.Instance

func (a byPriority) Len() int           { return len(a) }
func (a byPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool { return a[i].Priority > a[j].Priority }
