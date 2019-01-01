package main

import (
	"strings"
	t "time"
)

var dateFormats = [...]string{
	"02.01.06",
	"02.01.2006",
	"2.1.06",
	"2.1.2006"}

func parseTimeSuffix(line *Line, lineVal string) (*Date, *Date, string) {
	p := strings.LastIndex(lineVal, "(")
	untrimmedPos := strings.LastIndex(line.Content(), "(") + 1
	dtStr := lineVal[p+1 : len(lineVal)-1]
	dsTrimLen := countStartWhitespaces(dtStr)
	dtStr = strings.TrimSpace(dtStr)

	dashPos := strings.IndexRune(dtStr, '-')
	var start *Date
	var end *Date
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
		end.time = end.time.Add(24*t.Hour - t.Nanosecond)
	}

	if start != nil || end != nil {
		// Success
		finalizeDocCoords(start, line.LineNumber(), line.Offset()+untrimmedPos+dsTrimLen)
		finalizeDocCoords(end, line.LineNumber(), line.Offset()+untrimmedPos+dsTrimLen)
		return start, end, strings.TrimSpace(lineVal[:p])
	}

	return nil, nil, lineVal
}

func finalizeDocCoords(dt *Date, lineNumber int, offsetDelta int) {
	if dt != nil {
		dt.lineNumber = lineNumber
		dt.offset += offsetDelta
	}
}

func parseTimeSuffixSingle(lineVal string) *Date {
	ok, tm := parseDate(lineVal)
	if !ok {
		return nil
	}
	return &Date{
		time: tm,
		DocCoords: DocCoords{
			offset: countStartWhitespaces(lineVal),
			length: lengthWithoutStartEndWhitespaces(lineVal)}}
}

func parseTimeSuffixRanged(lineVal string, dashPos int) (*Date, *Date) {
	var start *Date
	var end *Date
	startStr := lineVal[:dashPos]
	endStr := lineVal[dashPos+1:]

	if startStr != "" {
		ok, tm := parseDate(startStr)
		if !ok {
			return nil, nil
		}
		start = &Date{
			time: tm,
			DocCoords: DocCoords{
				offset: countStartWhitespaces(startStr),
				length: lengthWithoutStartEndWhitespaces(startStr)}}
	}

	if endStr != "" {
		ok, tm := parseDate(endStr)
		if !ok {
			return nil, nil
		}
		end = &Date{
			time: tm,
			DocCoords: DocCoords{
				offset: len(startStr) + 1 + countStartWhitespaces(endStr),
				length: lengthWithoutStartEndWhitespaces(endStr)}}
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
