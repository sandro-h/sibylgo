package modify

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
)

func modifyInFile(todoFile string, modifyFunc func(string) (string, error)) error {
	content, err := util.ReadFile(todoFile)
	if err != nil {
		return err
	}

	updatedContent, err := modifyFunc(content)
	if err != nil {
		return err
	}

	err = util.WriteFile(todoFile, updatedContent)
	if err != nil {
		return err
	}

	return nil
}

type lineRange struct {
	startLine int
	endLine   int
}

func (r *lineRange) contains(ln int) bool {
	return ln >= r.startLine && ln <= r.endLine
}

func getFullLineRange(mom moment.Moment) *lineRange {
	return &lineRange{mom.GetDocCoords().LineNumber, mom.GetBottomLineNumber()}
}
