package cleanup

import (
	"bufio"
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var getNow = func() time.Time {
	return time.Now()
}

func CleanupDoneFromFile(todoFilePath string, trashFilePath string, onlyTopLevel bool) error {
	b, err := ioutil.ReadFile(todoFilePath)
	if err != nil {
		return err
	}
	s := string(b)

	todoFile, err := os.OpenFile(todoFilePath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer todoFile.Close()

	trashFile, err := os.OpenFile(trashFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer trashFile.Close()

	cleanupDone(s, onlyTopLevel,
		func(line string) { fmt.Fprintf(todoFile, "%s\n", line) },
		func(line string) { fmt.Fprintf(trashFile, "%s\n", line) },
		func() {
			fmt.Fprint(trashFile, "\n------------------\n")
			fmt.Fprintf(trashFile, "  Trash from %s\n", getNow().Format("02.01.2006 15:04:05"))
			fmt.Fprint(trashFile, "------------------\n")
		})

	return nil
}

func CleanupDoneFromFileToEnd(todoFilePath string, onlyTopLevel bool) error {
	b, err := ioutil.ReadFile(todoFilePath)
	if err != nil {
		return err
	}
	s := string(b)

	todoFile, err := os.OpenFile(todoFilePath, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer todoFile.Close()

	var deleted []string
	cleanupDone(s, onlyTopLevel,
		func(line string) { fmt.Fprintf(todoFile, "%s\n", line) },
		func(line string) { deleted = append(deleted, line) },
		nil)

	for _, d := range deleted {
		fmt.Fprintf(todoFile, "%s\n", d)
	}

	return nil
}

func CleanupDoneFromString(content string, onlyTopLevel bool) (string, string, error) {
	kept := ""
	deleted := ""
	err := cleanupDone(content, onlyTopLevel,
		func(line string) { addLine(&kept, line) },
		func(line string) { addLine(&deleted, line) },
		nil)
	if err != nil {
		return "", "", err
	}
	return kept, deleted, nil
}

func addLine(s *string, l string) {
	if *s != "" {
		*s += "\n"
	}
	*s += l
}

func cleanupDone(content string, onlyTopLevel bool,
	keepFunc func(string), deleteFunc func(string), firstDeleteFunc func()) error {
	todos, err := parse.ParseString(content)
	if err != nil {
		return err
	}

	toDel := computeDoneLines(todos.Moments, onlyTopLevel)

	scanner := bufio.NewScanner(strings.NewReader(content))
	ln := 0
	k := 0
	var curRange *LineRange
	firstDelete := true
	if len(toDel) > 0 {
		curRange = &toDel[0]
	}
	for scanner.Scan() {
		line := scanner.Text()
		delete := false
		if curRange != nil {
			if ln >= curRange.startLine && ln <= curRange.endLine {
				delete = true
			}
			if ln == curRange.endLine {
				if k < len(toDel)-1 {
					k++
					curRange = &toDel[k]
				} else {
					curRange = nil
				}
			}
		}
		if delete {
			if firstDelete {
				firstDelete = false
				if firstDeleteFunc != nil {
					firstDeleteFunc()
				}
			}
			if deleteFunc != nil {
				deleteFunc(line)
			}
		} else {
			if keepFunc != nil {
				keepFunc(line)
			}
		}

		ln++
	}

	return nil
}

func computeDoneLines(moms []moment.Moment, onlyTopLevel bool) []LineRange {
	var toDel []LineRange
	for _, m := range moms {
		if m.IsDone() {
			toDel = append(toDel, getFullLineRange(m))
		} else if !onlyTopLevel {
			subDels := computeDoneLines(m.GetSubMoments(), onlyTopLevel)
			toDel = append(toDel, subDels...)
		}
	}
	return toDel
}

func getFullLineRange(mom moment.Moment) LineRange {
	return LineRange{mom.GetDocCoords().LineNumber, mom.GetBottomLineNumber()}
}

type LineRange struct {
	startLine int
	endLine   int
}
