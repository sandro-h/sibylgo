package parse

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"strings"
	"time"
	"unicode/utf8"
)

var dateFormats = [...]string{
	"02.01.06",
	"02.01.2006",
	"2.1.06",
	"2.1.2006"}

// expected lineVal: .*(\s+<date>\s+)
func parseDateSuffix(line *Line, lineVal string) (*moment.Date, *moment.Date, *moment.Date, string) {
	p := strings.LastIndex(lineVal, "(")
	untrimmedPos := LastRuneIndex(line.Content(), "(") + 1
	dtStr := lineVal[p+1 : len(lineVal)-1]
	timeOfDay, dtStr := parseTimeSuffix(line, dtStr)
	finalizeDocCoords(timeOfDay, line.LineNumber(), line.Offset()+untrimmedPos)
	dsTrimLen := countStartWhitespaces(dtStr)
	dtStr = strings.TrimSpace(dtStr)

	dashPos := strings.Index(dtStr, "-")
	var start *moment.Date
	var end *moment.Date
	if dashPos >= 0 {
		start, end = parseDateSuffixRanged(dtStr, dashPos)
	} else {
		start = parseDateSuffixSingle(dtStr)
		if start != nil {
			endCopy := *start
			end = &endCopy
		}
	}

	if end != nil {
		// Set to very end of day
		end.Time = util.SetToEndOfDay(end.Time)
	}

	if start != nil || end != nil {
		// Success
		finalizeDocCoords(start, line.LineNumber(), line.Offset()+untrimmedPos+dsTrimLen)
		finalizeDocCoords(end, line.LineNumber(), line.Offset()+untrimmedPos+dsTrimLen)
		return start, end, timeOfDay, strings.TrimSpace(lineVal[:p])
	}

	return nil, nil, nil, lineVal
}

func finalizeDocCoords(dt *moment.Date, lineNumber int, offsetDelta int) {
	if dt != nil {
		dt.LineNumber = lineNumber
		dt.Offset += offsetDelta
	}
}

func parseDateSuffixSingle(lineVal string) *moment.Date {
	ok, tm := parseDate(lineVal)
	if !ok {
		return nil
	}
	return &moment.Date{
		Time: tm,
		DocCoords: moment.DocCoords{
			Offset: countStartWhitespaces(lineVal),
			Length: lengthWithoutStartEndWhitespaces(lineVal)}}
}

func parseDateSuffixRanged(lineVal string, dashPos int) (*moment.Date, *moment.Date) {
	var start *moment.Date
	var end *moment.Date
	startStr := lineVal[:dashPos]
	endStr := lineVal[dashPos+1:]

	if startStr != "" {
		ok, tm := parseDate(startStr)
		if !ok {
			return nil, nil
		}
		start = &moment.Date{
			Time: tm,
			DocCoords: moment.DocCoords{
				Offset: countStartWhitespaces(startStr),
				Length: lengthWithoutStartEndWhitespaces(startStr)}}
	}

	if endStr != "" {
		ok, tm := parseDate(endStr)
		if !ok {
			return nil, nil
		}
		end = &moment.Date{
			Time: tm,
			DocCoords: moment.DocCoords{
				Offset: utf8.RuneCountInString(startStr) + 1 + countStartWhitespaces(endStr),
				Length: lengthWithoutStartEndWhitespaces(endStr)}}
	}

	return start, end
}

func parseDate(str string) (bool, time.Time) {
	str = strings.TrimSpace(str)
	for _, fmt := range dateFormats {
		tm, err := time.ParseInLocation(fmt, str, time.Local)
		if err == nil {
			return true, tm
		}
	}
	return false, time.Unix(0, 0)
}
