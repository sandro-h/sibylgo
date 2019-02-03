package util

import (
	"time"
)

const Days = time.Hour * 24

func GetLowerBound(t1 *time.Time, t2 *time.Time) time.Time {
	return getBound(t1, t2, false)
}

func GetUpperBound(t1 *time.Time, t2 *time.Time) time.Time {
	return getBound(t1, t2, true)
}

func getBound(t1 *time.Time, t2 *time.Time, lowerOrUpper bool) time.Time {
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

func SetWeekday(dt time.Time, wd time.Weekday) time.Time {
	di := int(wd - dt.Weekday())
	return dt.AddDate(0, 0, di)
}

func SetToStartOfDay(dt time.Time) time.Time {
	y, m, d := dt.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.Local)
}

func SetToEndOfDay(dt time.Time) time.Time {
	y, m, d := dt.Date()
	return time.Date(y, m, d, 23, 59, 59, 999999999, time.Local)
}

func SetToStartOfWeek(dt time.Time) time.Time {
	// shift so monday=0, sunday=6
	wd := (dt.Weekday() + 6) % 7
	return SetToStartOfDay(dt.AddDate(0, 0, int(-wd)))
}

func SetToEndOfWeek(dt time.Time) time.Time {
	// shift so monday=0, sunday=6
	wd := (dt.Weekday() + 6) % 7
	return SetToEndOfDay(dt.AddDate(0, 0, int(6-wd)))
}

func SetTime(dt time.Time, tm time.Time) time.Time {
	y, m, d := dt.Date()
	return time.Date(y, m, d, tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond(), time.Local)
}