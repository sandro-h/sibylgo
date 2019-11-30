package modify

import (
	"bufio"
	"github.com/sandro-h/sibylgo/moment"
	"strings"
)

// Delete removes moments from the todo file content. It returns
// the content without the removed moments lines and all the removed moment
// lines.
func Delete(content string, toDel []moment.Moment) (string, string) {
	kept := ""
	deleted := ""

	scanner := bufio.NewScanner(strings.NewReader(content))
	ln := 0
	k := 0
	var curRange *lineRange
	prevLineWasDeleted := false
	if len(toDel) > 0 {
		curRange = getFullLineRange(toDel[0])
	}
	for scanner.Scan() {
		line := scanner.Text()
		delete := false
		// Check if line is part of current to-delete range.
		if curRange != nil {
			if ln >= curRange.startLine && ln <= curRange.endLine {
				delete = true
			}
			if ln == curRange.endLine {
				if k < len(toDel)-1 {
					k++
					curRange = getFullLineRange(toDel[k])
				} else {
					curRange = nil
				}
			}
		}
		// Check if line is empty right after a deleted line -> trim superfluous empty lines.
		if prevLineWasDeleted && strings.TrimSpace(line) == "" {
			delete = true
		}
		// Delete or keep the line
		if delete {
			prevLineWasDeleted = true
			addLine(&deleted, line)
		} else {
			prevLineWasDeleted = false
			addLine(&kept, line)
		}

		ln++
	}

	return kept, deleted
}

func getFullLineRange(mom moment.Moment) *lineRange {
	return &lineRange{mom.GetDocCoords().LineNumber, mom.GetBottomLineNumber()}
}

func addLine(s *string, l string) {
	if *s != "" {
		*s += "\n"
	}
	*s += l
}

type lineRange struct {
	startLine int
	endLine   int
}
