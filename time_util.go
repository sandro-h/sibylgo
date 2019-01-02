package main

import (
	t "time"
)

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
