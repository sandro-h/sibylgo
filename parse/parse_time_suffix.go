package parse

import (
	"github.com/sandro-h/sibylgo/moment"
	"strings"
	"time"
	"unicode/utf8"
)

// expected lineVal: .*<timeval>\s*
// e.g. 12.5.2019 13:15
func parseTimeSuffix(line *Line, lineVal string) (*moment.Date, string) {
	trimmed := strings.TrimSpace(lineVal)
	p := strings.LastIndex(trimmed, " ")
	if p < 0 || p == len(trimmed)-1 {
		return nil, lineVal
	}
	unip := utf8.RuneCountInString(trimmed[:p])
	tmStr := trimmed[p+1:]
	ok, tm := parseTime(tmStr)
	if !ok {
		return nil, lineVal
	}
	return &moment.Date{Time: tm,
			DocCoords: moment.DocCoords{
				Offset: unip + 1,
				Length: utf8.RuneCountInString(tmStr)}},
		lineVal[0:p]
}

func parseTime(str string) (bool, time.Time) {
	str = strings.TrimSpace(str)
	tm, err := time.ParseInLocation("15:04", str, time.Local)
	if err != nil {
		return false, time.Unix(0, 0)

	}
	return true, tm
}
