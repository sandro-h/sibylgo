package main

import (
	t "time"
)

const Days = t.Hour * 24

func getLowerBound(t1 *t.Time, t2 *t.Time) t.Time {
	return getBound(t1, t2, false)
}

func getUpperBound(t1 *t.Time, t2 *t.Time) t.Time {
	return getBound(t1, t2, true)
}

func getBound(t1 *t.Time, t2 *t.Time, lowerOrUpper bool) t.Time {
	if t1 != nil && t2 != nil {
		if !lowerOrUpper {
			if t1.Before(*t2) {
				return *t1
			} else {
				return *t2
			}
		} else {
			if t1.After(*t2) {
				return *t1
			} else {
				return *t2
			}
		}
	} else if t1 != nil {
		return *t1
	} else {
		return *t2
	}
}

func setWeekday(dt t.Time, wd t.Weekday) t.Time {
	di := int(wd - dt.Weekday())
	return dt.AddDate(0, 0, di)
}

func setToStartOfDay(dt t.Time) t.Time {
	y, m, d := dt.Date()
	return t.Date(y, m, d, 0, 0, 0, 0, t.Local)
}

func setToEndOfDay(dt t.Time) t.Time {
	y, m, d := dt.Date()
	return t.Date(y, m, d, 23, 59, 59, 999999999, t.Local)
}
