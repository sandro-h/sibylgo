package parse

import (
	"strings"

	"github.com/sandro-h/sibylgo/moment"
)

// parseMoment parses a moment from the line. It only parses the moment of this current line
// and none of the sub moments or comments appear on subsequent lines.
func parseMoment(line *Line, lineVal string) moment.Moment {
	id, lineVal := parseID(line, lineVal)
	mom, lineVal := parseBaseMoment(line, lineVal)
	mom.SetID(id)

	state, lineVal := parseStateMark(line, lineVal)
	if state == nil {
		return nil
	}
	mom.SetWorkState(*state)

	prio, lineVal := parsePriority(lineVal)
	mom.SetPriority(prio)

	mom.SetName(lineVal)

	return mom
}

func parseID(line *Line, lineVal string) (*moment.Identifier, string) {
	idPos := strings.LastIndex(lineVal, " #")
	if idPos < 0 {
		return nil, lineVal
	}

	untrimmedPos := LastRuneIndex(line.Content(), " #") + 1
	idStr := strings.TrimSpace(lineVal[idPos+2:])
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

func parseStateMark(line *Line, lineVal string) (*moment.WorkState, string) {
	rBracketPos := 0
	innerContent := ' '
	// [1:] to skip left bracket
	for i, c := range lineVal[1:] {
		if c == ParseConfig.GetRBracket() {
			// +1 because [1:]
			rBracketPos = i + 1
			break
		}
		if c != ' ' && c != '\t' {
			if innerContent == ' ' {
				innerContent = c
			} else {
				return nil, lineVal
			}

		}
	}

	if rBracketPos == 0 {
		return nil, lineVal
	}

	var state moment.WorkState
	switch innerContent {
	case ' ':
		state = moment.NewState
	case ParseConfig.GetDoneMark():
		state = moment.DoneState
	case ParseConfig.GetInProgressMark():
		state = moment.InProgressState
	case ParseConfig.GetWaitingMark():
		state = moment.WaitingState
	default:
		return nil, lineVal
	}

	return &state, strings.TrimSpace(lineVal[rBracketPos+1:])
}
