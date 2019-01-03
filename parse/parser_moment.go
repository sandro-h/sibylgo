package parse

import (
	"github.com/sandro-h/sibylgo/moment"
	"strings"
)

func ParseMoment(line *Line, lineVal string) (moment.Moment, error) {
	mom, lineVal := parseBaseMoment(line, lineVal)

	done, lineVal, err := parseDoneMark(line, lineVal)
	if err != nil {
		return nil, err
	}
	mom.SetDone(done)

	prio, lineVal := parsePriority(lineVal)
	mom.SetPriority(prio)

	mom.SetName(lineVal)

	return mom, nil
}

func parseBaseMoment(line *Line, lineVal string) (moment.Moment, string) {
	re, newLineVal := parseRecurMoment(line, lineVal)
	if re != nil {
		return re, newLineVal
	}
	return parseSingleMoment(line, lineVal)
}

func parseRecurMoment(line *Line, lineVal string) (*moment.RecurMoment, string) {
	if !strings.HasSuffix(lineVal, ")") {
		return nil, lineVal
	}
	re, newLineVal := parseRecurrence(line, lineVal)
	if re != nil {
		mom := &moment.RecurMoment{Recurrence: *re}
		mom.DocCoords = moment.DocCoords{line.LineNumber(), line.Offset(), line.Length()}
		return mom, newLineVal
	}
	return nil, lineVal
}

func parseSingleMoment(line *Line, lineVal string) (*moment.SingleMoment, string) {
	var start *moment.Date
	var end *moment.Date
	if strings.HasSuffix(lineVal, ")") {
		start, end, lineVal = parseTimeSuffix(line, lineVal)
	}
	mom := &moment.SingleMoment{Start: start, End: end}
	mom.DocCoords = moment.DocCoords{line.LineNumber(), line.Offset(), line.Length()}
	return mom, lineVal
}

func parseDoneMark(line *Line, lineVal string) (bool, string, error) {
	rBracketPos := 0
	done := false
	for i, c := range lineVal {
		if c == doneRBracket {
			rBracketPos = i
			break
		}
		if c == doneMark || c == doneMarkUpper {
			done = true
		}
	}
	if rBracketPos == 0 {
		return false, "", newParseError(line, "Expected closing %c for moment line %s", doneMark, line.Content())
	}
	return done, strings.TrimSpace(lineVal[rBracketPos+1:]), nil
}
