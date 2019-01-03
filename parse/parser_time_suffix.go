package parse

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"strings"
	t "time"
)

var dateFormats = [...]string{
	"02.01.06",
	"02.01.2006",
	"2.1.06",
	"2.1.2006"}

func parseTimeSuffix(line *Line, lineVal string) (*moment.Date, *moment.Date, string) {
	p := strings.LastIndex(lineVal, "(")
	untrimmedPos := strings.LastIndex(line.Content(), "(") + 1
	dtStr := lineVal[p+1 : len(lineVal)-1]
	dsTrimLen := countStartWhitespaces(dtStr)
	dtStr = strings.TrimSpace(dtStr)

	dashPos := strings.IndexRune(dtStr, '-')
	var start *moment.Date
	var end *moment.Date
	if dashPos >= 0 {
		start, end = parseTimeSuffixRanged(dtStr, dashPos)
	} else {
		start = parseTimeSuffixSingle(dtStr)
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
		return start, end, strings.TrimSpace(lineVal[:p])
	}

	return nil, nil, lineVal
}

func finalizeDocCoords(dt *moment.Date, lineNumber int, offsetDelta int) {
	if dt != nil {
		dt.LineNumber = lineNumber
		dt.Offset += offsetDelta
	}
}

func parseTimeSuffixSingle(lineVal string) *moment.Date {
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

func parseTimeSuffixRanged(lineVal string, dashPos int) (*moment.Date, *moment.Date) {
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
				Offset: len(startStr) + 1 + countStartWhitespaces(endStr),
				Length: lengthWithoutStartEndWhitespaces(endStr)}}
	}

	return start, end
}

func parseDate(str string) (bool, t.Time) {
	str = strings.TrimSpace(str)
	for _, fmt := range dateFormats {
		tm, err := t.Parse(fmt, str)
		if err == nil {
			return true, tm
		}
	}
	return false, t.Unix(0, 0)
}
