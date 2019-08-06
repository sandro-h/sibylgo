package generate

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"time"
)

// RecurIterator generates timestamps within a time range based on a recurrence type.
type RecurIterator struct {
	recurrence moment.Recurrence
	from       time.Time
	until      time.Time
	cur        time.Time
	next       time.Time
}

// NewRecurIterator creates a new recurrence iterator for the given recurrence type in the
// given time range.
func NewRecurIterator(recurrence moment.Recurrence, from time.Time, until time.Time) *RecurIterator {
	it := &RecurIterator{
		recurrence: recurrence,
		from:       from,
		until:      until,
		cur:        from.AddDate(0, 0, -1)}
	it.prepareNext()
	return it
}

// HasNext returns true if there are more timestamps the iterator can generate in the time range.
func (it *RecurIterator) HasNext() bool {
	return !it.next.After(it.until)
}

// Next returns the next timestamp in the time range.
func (it *RecurIterator) Next() time.Time {
	res := it.next
	it.prepareNext()
	return res
}

func (it *RecurIterator) prepareNext() {
	switch it.recurrence.Recurrence {
	case moment.RecurDaily:
		it.next = getNextDaily(it.cur)
	case moment.RecurWeekly:
		it.next = getNextWeekly(it.cur, it.recurrence.RefDate.Time)
	case moment.RecurMonthly:
		it.next = getNextMonthly(it.cur, it.recurrence.RefDate.Time)
	case moment.RecurYearly:
		it.next = getNextYearly(it.cur, it.recurrence.RefDate.Time)
	}
	it.cur = it.next
}

func getNextDaily(after time.Time) time.Time {
	return after.AddDate(0, 0, 1)
}

func getNextWeekly(after time.Time, ref time.Time) time.Time {
	dt := util.SetWeekday(after, ref.Weekday())
	if !dt.After(after) {
		dt = dt.AddDate(0, 0, 7)
	}
	return dt
}

func getNextMonthly(after time.Time, ref time.Time) time.Time {
	y, m, _ := after.Date()
	dt := time.Date(y, m, ref.Day(), 0, 0, 0, 0, time.Local)
	if !dt.After(after) {
		dt = dt.AddDate(0, 1, 0)
	}
	return dt
}

func getNextYearly(after time.Time, ref time.Time) time.Time {
	_, m, d := ref.Date()
	dt := time.Date(after.Year(), m, d, 0, 0, 0, 0, time.Local)
	if !dt.After(after) {
		dt = dt.AddDate(1, 0, 0)
	}
	return dt
}
