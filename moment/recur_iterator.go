package moment

import (
	"github.com/sandro-h/sibylgo/util"
	t "time"
)

type RecurIterator struct {
	recurrence Recurrence
	from       t.Time
	until      t.Time
	cur        t.Time
	next       t.Time
}

func NewRecurIterator(recurrence Recurrence, from t.Time, until t.Time) *RecurIterator {
	it := &RecurIterator{
		recurrence: recurrence,
		from:       from,
		until:      until,
		cur:        from.AddDate(0, 0, -1)}
	it.prepareNext()
	return it
}

func (it *RecurIterator) HasNext() bool {
	return !it.next.After(it.until)
}

func (it *RecurIterator) Next() t.Time {
	res := it.next
	it.prepareNext()
	return res
}

func (it *RecurIterator) prepareNext() {
	switch it.recurrence.Recurrence {
	case RE_DAILY:
		it.next = getNextDaily(it.cur, it.recurrence.RefDate.Time)
	case RE_WEEKLY:
		it.next = getNextWeekly(it.cur, it.recurrence.RefDate.Time)
	case RE_MONTHLY:
		it.next = getNextMonthly(it.cur, it.recurrence.RefDate.Time)
	case RE_YEARLY:
		it.next = getNextYearly(it.cur, it.recurrence.RefDate.Time)
	}
	it.cur = it.next
}

func getNextDaily(after t.Time, ref t.Time) t.Time {
	return after.AddDate(0, 0, 1)
}

func getNextWeekly(after t.Time, ref t.Time) t.Time {
	dt := util.SetWeekday(after, ref.Weekday())
	if !dt.After(after) {
		dt = dt.AddDate(0, 0, 7)
	}
	return dt
}

func getNextMonthly(after t.Time, ref t.Time) t.Time {
	y, m, _ := after.Date()
	dt := t.Date(y, m, ref.Day(), 0, 0, 0, 0, t.Local)
	if !dt.After(after) {
		dt = dt.AddDate(0, 1, 0)
	}
	return dt
}

func getNextYearly(after t.Time, ref t.Time) t.Time {
	_, m, d := ref.Date()
	dt := t.Date(after.Year(), m, d, 0, 0, 0, 0, t.Local)
	if !dt.After(after) {
		dt = dt.AddDate(1, 0, 0)
	}
	return dt
}
