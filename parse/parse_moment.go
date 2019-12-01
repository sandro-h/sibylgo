package parse

import (
	"github.com/sandro-h/sibylgo/moment"
	"strings"
)

// parseMoment parses a moment from the line. It only parses the moment of this current line
// and none of the sub moments or comments appear on subsequent lines.
func parseMoment(line *Line, lineVal string) (moment.Moment, error) {
	id, lineVal := parseID(line, lineVal)
	mom, lineVal := parseBaseMoment(line, lineVal)
	mom.SetID(id)

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

func parseID(line *Line, lineVal string) (*moment.Identifier, string) {
	idPos := strings.LastIndex(lineVal, " #")
	if idPos < 0 {
		return nil, lineVal
	}

	untrimmedPos := LastRuneIndex(line.Content(), " #") + 1
	idStr := strings.TrimSpace(lineVal[idPos+2 : len(lineVal)])
	id := moment.Identifier{Value: idStr,
		DocCoords: moment.DocCoords{LineNumber: line.LineNumber(), Offset: line.Offset() + untrimmedPos, Length: len(idStr) + 1}}
	return &id, strings.TrimSpace(lineVal[:idPos])
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
	re, timeOfDay, newLineVal := parseRecurrence(line, lineVal)
	if re != nil {
		mom := &moment.RecurMoment{Recurrence: *re}
		mom.TimeOfDay = timeOfDay
		mom.DocCoords = moment.DocCoords{LineNumber: line.LineNumber(), Offset: line.Offset(), Length: line.Length()}
		return mom, newLineVal
	}
	return nil, lineVal
}

func parseSingleMoment(line *Line, lineVal string) (*moment.SingleMoment, string) {
	var start *moment.Date
	var end *moment.Date
	var timeOfDay *moment.Date
	if strings.HasSuffix(lineVal, ")") {
		start, end, timeOfDay, lineVal = parseDateSuffix(line, lineVal)
	}
	mom := &moment.SingleMoment{Start: start, End: end}
	mom.TimeOfDay = timeOfDay
	mom.DocCoords = moment.DocCoords{LineNumber: line.LineNumber(), Offset: line.Offset(), Length: line.Length()}
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
