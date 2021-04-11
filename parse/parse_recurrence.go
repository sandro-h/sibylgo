package parse

import (
	"strconv"
	"strings"
	"time"

	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
)

var getNow = func() time.Time {
	return time.Now()
}

// expected lineVal: .*(\s+<recur>\s+)
func parseRecurrence(line *Line, lineVal string) (*moment.Recurrence, *moment.Date, string) {
	p := strings.LastIndex(lineVal, "(")
	if p < 0 {
		return nil, nil, lineVal
	}
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
		re = tryParseNWeekly(reStr)
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
	if ParseConfig.GetDailyPattern().MatchString(reStr) {
		return &moment.Recurrence{
			Recurrence: moment.RecurDaily,
			RefDate:    &moment.Date{Time: getNow()}}
	}
	return nil
}

func tryParseWeekly(reStr string) *moment.Recurrence {
	matches := ParseConfig.GetWeeklyPattern().FindStringSubmatch(reStr)
	if matches != nil {
		wd := parseWeekday(matches[1])
		dt := util.SetWeekday(getNow(), wd)
		return &moment.Recurrence{
			Recurrence: moment.RecurWeekly,
			RefDate:    &moment.Date{Time: dt}}
	}
	return nil
}

func tryParseNWeekly(reStr string) *moment.Recurrence {
	matches := ParseConfig.GetNthWeeklyPattern().FindStringSubmatch(reStr)
	if matches != nil {
		n, re := parseNth(matches[1])
		if n < 0 {
			return nil
		}

		wd := parseWeekday(matches[2])
		dt := util.SetWeekday(getNow(), wd)
		weekOffset := util.EpochWeek(dt) % n
		dt = dt.AddDate(0, 0, -7*weekOffset)
		return &moment.Recurrence{
			Recurrence: re,
			RefDate:    &moment.Date{Time: dt}}
	}
	return nil
}

func parseWeekday(str string) time.Weekday {
	day, ok := ParseConfig.GetWeekDays()[strings.ToLower(str)]
	if !ok {
		return -1
	}
	return day
}

func parseNth(str string) (int, int) {
	nth, ok := ParseConfig.GetNths()[strings.ToLower(str)]
	if !ok || nth > 4 {
		return -1, -1
	}

	return nth, moment.RecurBiWeekly + (nth - 2)
}

func tryParseMonthly(reStr string) *moment.Recurrence {
	matches := ParseConfig.GetMonthlyPattern().FindStringSubmatch(reStr)
	if matches != nil {
		day, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil
		}
		y, m, _ := getNow().Date()
		dt := time.Date(y, m, day, 0, 0, 0, 0, time.Local)
		return &moment.Recurrence{
			Recurrence: moment.RecurMonthly,
			RefDate:    &moment.Date{Time: dt}}
	}
	return nil
}

func tryParseYearly(reStr string) *moment.Recurrence {
	matches := ParseConfig.GetYearlyPattern().FindStringSubmatch(reStr)
	if matches != nil {
		day, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil
		}
		month, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil
		}
		y := getNow().Year()
		dt := time.Date(y, time.Month(month), day, 0, 0, 0, 0, time.Local)
		return &moment.Recurrence{
			Recurrence: moment.RecurYearly,
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
