package modify

import (
	"bufio"
	"github.com/sandro-h/sibylgo/moment"
	"strings"
)

// Delete removes moments from the todo content. It returns
// the content without the removed moments lines and all the removed moment
// lines.
func Delete(content string, toDel []moment.Moment) (string, string) {
	kept := ""
	deleted := ""

	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	curDelRangeIndex := 0
	var curDelRange *lineRange
	prevLineWasDeleted := false
	if len(toDel) > 0 {
		curDelRange = getFullLineRange(toDel[0])
	}
	for scanner.Scan() {
		line := scanner.Text()
		delete := false
		delete, curDelRange, curDelRangeIndex = updateDeleteState(lineNum, curDelRange, curDelRangeIndex, toDel)
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

		lineNum++
	}

	return kept, deleted
}

func updateDeleteState(lineNum int, curDelRange *lineRange, curDelRangeIndex int, toDel []moment.Moment) (bool, *lineRange, int) {
	delete := false
	if curDelRange != nil {
		if curDelRange.contains(lineNum) {
			delete = true
		}
		if lineNum == curDelRange.endLine {
			// Switch to next delete range
			if curDelRangeIndex < len(toDel)-1 {
				curDelRangeIndex++
				curDelRange = getFullLineRange(toDel[curDelRangeIndex])
			} else {
				curDelRange = nil
			}
		}
	}

	return delete, curDelRange, curDelRangeIndex
}

func addLine(s *string, l string) {
	if *s != "" {
		*s += "\n"
	}
	*s += l
}
