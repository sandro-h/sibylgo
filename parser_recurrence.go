package main

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

const dailyPattern = "every day"

var weeklyPattern, _ = regexp.Compile("(?i)every (monday|tuesday|wednesday|thursday|friday|saturday|sunday)")
var monthlyPattern, _ = regexp.Compile("(?i)every (\\d{1,2})\\.?$")
var yearlyPattern, _ = regexp.Compile("(?i)every (\\d{1,2})\\.(\\d{1,2})\\.?$")

func parseRecurrence(line *Line, lineVal string) (*Recurrence, string) {
	p := strings.LastIndex(lineVal, "(")
	reStr := lineVal[p+1 : len(lineVal)-1]
	untrimmedPos := strings.LastIndex(line.Content(), "(") + 1

	var re *Recurrence
	re = tryParseDaily(reStr)
	if re == nil {
		re = tryParseWeekly(reStr)
	}
	if re == nil {
		re = tryParseMonthly(reStr)
	}
	if re == nil {
		re = tryParseYearly(reStr)
	}

	if re == nil {
		return nil, lineVal
	}
	return setDocCoords(re, line.LineNumber(), line.Offset()+untrimmedPos, len(reStr)),
		strings.TrimSpace(lineVal[:p])
}

func tryParseDaily(reStr string) *Recurrence {
	if strings.EqualFold(reStr, dailyPattern) {
		return &Recurrence{
			recurrence: RE_DAILY,
			refDate:    &Date{time: time.Now()}}
	}
	return nil
}

func tryParseWeekly(reStr string) *Recurrence {
	matches := weeklyPattern.FindStringSubmatch(reStr)
	if matches != nil {
		wd := parseWeekday(matches[1])
		dt := setWeekday(time.Now(), wd)
		return &Recurrence{
			recurrence: RE_WEEKLY,
			refDate:    &Date{time: dt}}
	}
	return nil
}

func parseWeekday(str string) time.Weekday {
	switch strings.ToLower(str) {
	case "sunday":
		return time.Sunday
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	}
	return -1
}

func tryParseMonthly(reStr string) *Recurrence {
	matches := monthlyPattern.FindStringSubmatch(reStr)
	if matches != nil {
		day, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil
		}
		y, m, _ := time.Now().Date()
		dt := time.Date(y, m, day, 0, 0, 0, 0, time.Local)
		return &Recurrence{
			recurrence: RE_MONTHLY,
			refDate:    &Date{time: dt}}
	}
	return nil
}

func tryParseYearly(reStr string) *Recurrence {
	matches := yearlyPattern.FindStringSubmatch(reStr)
	if matches != nil {
		day, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil
		}
		month, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil
		}
		y := time.Now().Year()
		dt := time.Date(y, time.Month(month), day, 0, 0, 0, 0, time.Local)
		return &Recurrence{
			recurrence: RE_YEARLY,
			refDate:    &Date{time: dt}}
	}
	return nil
}

func setDocCoords(re *Recurrence, lineNumber int, offset int, length int) *Recurrence {
	re.refDate.lineNumber = lineNumber
	re.refDate.offset = offset
	re.refDate.length = length
	return re
}
