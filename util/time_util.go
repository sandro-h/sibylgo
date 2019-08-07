package util

import (
	"time"
)

// Days is the duration of a day.
const Days = time.Hour * 24

// GetLowerBound returns the earlier of two times. A nil time is considered infinitely late so won't be
// used if the other time is not nil.
func GetLowerBound(t1 *time.Time, t2 *time.Time) time.Time {
	return getBound(t1, t2, false)
}

// GetUpperBound returns the later of two times. A nil time is considered infinitely early so won't be
// used if the other time is not nil.
func GetUpperBound(t1 *time.Time, t2 *time.Time) time.Time {
	return getBound(t1, t2, true)
}

func getBound(t1 *time.Time, t2 *time.Time, lowerOrUpper bool) time.Time {
	if t1 != nil && t2 != nil {
		if !lowerOrUpper {
			if t1.Before(*t2) {
				return *t1
			}
			return *t2
		}
		if t1.After(*t2) {
			return *t1
		}
		return *t2
	} else if t1 != nil {
		return *t1
	} else {
		return *t2
	}
}

// SetWeekday changes the given datetime to occur on the given week day.
func SetWeekday(dt time.Time, wd time.Weekday) time.Time {
	di := int(wd - dt.Weekday())
	return dt.AddDate(0, 0, di)
}

// SetToStartOfDay changes the given datetime to the start of the day (0:00).
func SetToStartOfDay(dt time.Time) time.Time {
	y, m, d := dt.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

// SetToEndOfDay changes the given datetime to the end of the day (23:59).
func SetToEndOfDay(dt time.Time) time.Time {
	y, m, d := dt.Date()
	return time.Date(y, m, d, 23, 59, 59, 999999999, time.Local)
}

// SetToStartOfWeek changes the given datetime to the start of the week (Monday).
func SetToStartOfWeek(dt time.Time) time.Time {
	// shift so monday=0, sunday=6
	wd := (dt.Weekday() + 6) % 7
	return SetToStartOfDay(dt.AddDate(0, 0, int(-wd)))
}

// SetToEndOfWeek changes the given datetime to the end of the week (Sunday).
func SetToEndOfWeek(dt time.Time) time.Time {
	// shift so monday=0, sunday=6
	wd := (dt.Weekday() + 6) % 7
	return SetToEndOfDay(dt.AddDate(0, 0, int(6-wd)))
}

// SetTime changes the given datetime to the time of the day.
func SetTime(dt time.Time, tm time.Time) time.Time {
	y, m, d := dt.Date()
	return time.Date(y, m, d, tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond(), time.Local)
}
