package cleanup

import (
	"fmt"
	"github.com/sandro-h/sibylgo/modify"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/util"
	"time"
)

var getNow = func() time.Time {
	return time.Now()
}

// MoveDoneToTrashFile moves all done moments in the todo file to a fixed trash file
func MoveDoneToTrashFile(todoFilePath string, trashFilePath string, onlyTopLevel bool) error {
	rawTodoContent, err := util.ReadFile(todoFilePath)
	if err != nil {
		return err
	}

	done, err := computeDoneLinesFromContent(rawTodoContent, onlyTopLevel)
	if err != nil {
		return err
	}
	if done == nil {
		return nil
	}

	kept, deleted := modify.Delete(rawTodoContent, done)
	header := fmt.Sprintf(`
------------------
  Trash from %s
------------------
`, getNow().Format("02.01.2006 15:04:05"))

	util.WriteFile(todoFilePath, kept)
	util.WriteFile(trashFilePath, header+deleted)

	return nil
}

// MoveDoneToEndOfFile moves all done moments in the todo file to the end of that file.
func MoveDoneToEndOfFile(todoFilePath string, onlyTopLevel bool) error {
	rawTodoContent, err := util.ReadFile(todoFilePath)
	if err != nil {
		return err
	}

	done, err := computeDoneLinesFromContent(rawTodoContent, onlyTopLevel)
	if err != nil {
		return err
	}
	if done == nil {
		return nil
	}

	kept, deleted := modify.Delete(rawTodoContent, done)
	util.WriteFile(todoFilePath, kept+"\n"+deleted)

	return nil
}

// SeparateDoneFromString separates the moments from the raw content string into
// done and not done moments and returns them.
func SeparateDoneFromString(content string, onlyTopLevel bool) (string, string, error) {
	done, err := computeDoneLinesFromContent(content, onlyTopLevel)
	if err != nil {
		return "", "", err
	}
	if done == nil {
		return content, "", nil
	}

	kept, deleted := modify.Delete(content, done)
	return kept, deleted, nil
}

func addLine(s *string, l string) {
	if *s != "" {
		*s += "\n"
	}
	*s += l
}

func computeDoneLinesFromContent(content string, onlyTopLevel bool) ([]moment.Moment, error) {
	todos, err := parse.String(content)
	if err != nil {
		return nil, err
	}
	return computeDoneLines(todos.Moments, onlyTopLevel), nil
}

func computeDoneLines(moms []moment.Moment, onlyTopLevel bool) []moment.Moment {
	var toDel []moment.Moment
	for _, m := range moms {
		if m.IsDone() {
			toDel = append(toDel, m)
		} else if !onlyTopLevel {
			subDels := computeDoneLines(m.GetSubMoments(), onlyTopLevel)
			toDel = append(toDel, subDels...)
		}
	}
	return toDel
}
