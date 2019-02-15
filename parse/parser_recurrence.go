package parse

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const dailyPattern = "every day"

var weeklyPattern, _ = regexp.Compile("(?i)every (monday|tuesday|wednesday|thursday|friday|saturday|sunday)")
var monthlyPattern, _ = regexp.Compile("(?i)every (\\d{1,2})\\.?$")
var yearlyPattern, _ = regexp.Compile("(?i)every (\\d{1,2})\\.(\\d{1,2})\\.?$")

// expected lineVal: .*(\s+<recur>\s+)
func parseRecurrence(line *Line, lineVal string) (*moment.Recurrence, *moment.Date, string) {
	p := strings.LastIndex(lineVal, "(")
	untrimmedPos := LastRuneIndex(line.Content(), "(") + 1
	reStr := lineVal[p+1 : len(lineVal)-1]
	timeOfDay, reStr := parseTimeSuffix(line, reStr)
	if timeOfDay != nil {
		timeOfDay.Offset += line.Offset() + untrimmedPos
	}

	var re *moment.Recurrence
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
		return nil, nil, lineVal
	}
	return setDocCoords(re, line.LineNumber(), line.Offset()+untrimmedPos, len(reStr)),
		timeOfDay,
		strings.TrimSpace(lineVal[:p])
}

func tryParseDaily(reStr string) *moment.Recurrence {
	if strings.EqualFold(reStr, dailyPattern) {
		return &moment.Recurrence{
			Recurrence: moment.RE_DAILY,
			RefDate:    &moment.Date{Time: time.Now()}}
	}
	return nil
}

func tryParseWeekly(reStr string) *moment.Recurrence {
	matches := weeklyPattern.FindStringSubmatch(reStr)
	if matches != nil {
		wd := parseWeekday(matches[1])
		dt := util.SetWeekday(time.Now(), wd)
		return &moment.Recurrence{
			Recurrence: moment.RE_WEEKLY,
			RefDate:    &moment.Date{Time: dt}}
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

func tryParseMonthly(reStr string) *moment.Recurrence {
	matches := monthlyPattern.FindStringSubmatch(reStr)
	if matches != nil {
		day, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil
		}
		y, m, _ := time.Now().Date()
		dt := time.Date(y, m, day, 0, 0, 0, 0, time.Local)
		return &moment.Recurrence{
			Recurrence: moment.RE_MONTHLY,
			RefDate:    &moment.Date{Time: dt}}
	}
	return nil
}

func tryParseYearly(reStr string) *moment.Recurrence {
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
		return &moment.Recurrence{
			Recurrence: moment.RE_YEARLY,
			RefDate:    &moment.Date{Time: dt}}
	}
	return nil
}

func setDocCoords(re *moment.Recurrence, lineNumber int, offset int, length int) *moment.Recurrence {
	re.RefDate.LineNumber = lineNumber
	re.RefDate.Offset = offset
	re.RefDate.Length = length
	return re
}
