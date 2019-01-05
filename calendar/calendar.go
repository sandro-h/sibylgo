package calendar

import (
	"github.com/sandro-h/sibylgo/generate"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"sort"
	"time"
)

const dateFormat = "2006-01-02"

type Entry struct {
	Title string `json:"title"`
	Start string `json:"start"`
	End   string `json:"end"`
}

func NewEntry(inst *moment.MomentInstance) Entry {
	return Entry{
		Title: inst.Name,
		Start: inst.Start.Format(dateFormat),
		End:   inst.End.AddDate(0, 0, 1).Format(dateFormat)} // +1 because fullcalendar is non-inclusive
}

func CompileCalendarEntries(todos *moment.Todos, from time.Time, until time.Time) []Entry {
	from = util.SetToStartOfDay(from)
	until = util.SetToStartOfDay(until)

	insts := generate.GenerateInstancesFilteredWithoutSubs(todos, from, until,
		func(mom *moment.MomentInstance) bool { return !mom.Done && mom.EndsInRange })
	sort.Sort(ByPriority(insts))

	var entries []Entry
	for _, inst := range insts {
		entries = append(entries, NewEntry(inst))
	}
	return entries
}

type ByPriority []*moment.MomentInstance

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool { return a[i].Priority > a[j].Priority }
