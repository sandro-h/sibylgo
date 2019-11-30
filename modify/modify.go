package modify

import (
	"bufio"
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/stringify"
	"github.com/sandro-h/sibylgo/util"
	"strings"
)

// Insert inserts moments into the todo file content. The moments are
// inserted at the very end of the file, unless they are assigned to an
// existing category, then they're inserted at the and of the respective category.
func Insert(content string, toInsert []moment.Moment) (string, error) {
	byCategory, noCat := groupByCategory(toInsert)

	res := ""
	todos, err := parse.String(content)
	if err != nil {
		return "", err
	}
	catEnds, noCatEnd := findCategoryEnds(todos)

	err = validateMissingInsertCategories(&byCategory, &catEnds)
	if err != nil {
		return "", err
	}

	res = ""
	fmt.Printf("noCatEnd.bottom=%d\n", noCatEnd.bottom)
	if noCatEnd.bottom == -1 {
		for _, m := range noCat {
			res += stringify.FormatMoment(m)
		}
	}
	ln := 0
	k := 0
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		res += line + "\n"
		if noCatEnd.bottom != -1 && ln == noCatEnd.bottom {
			for _, m := range noCat {
				res += stringify.FormatMoment(m)
			}
		} else if ln == catEnds[k].bottom {
			list, ok := byCategory[catEnds[k].name]
			if ok {
				for _, m := range list {
					res += stringify.FormatMoment(m)
				}
			}
		}

		ln++
	}

	return res, nil
}

func groupByCategory(moms []moment.Moment) (map[string][]moment.Moment, []moment.Moment) {
	var noCat []moment.Moment
	byCategory := make(map[string][]moment.Moment)
	for _, m := range moms {
		cat := m.GetCategory()
		if cat == nil {
			noCat = append(noCat, m)
		} else {
			list, _ := byCategory[cat.Name]
			list = append(list, m)
			byCategory[cat.Name] = list
		}
	}
	return byCategory, noCat
}

func findCategoryEnds(todos *moment.Todos) ([]categoryEnd, categoryEnd) {
	noCatEnd := categoryEnd{"noCat", -1}
	var catEnds []categoryEnd
	for _, c := range todos.Categories {
		catEnd := categoryEnd{
			c.Name,
			c.LineNumber + 1}
		catEnds = append(catEnds, catEnd)
	}

	var lastCat string
	previousMomBottom := -1
	k := 0
	for _, m := range todos.Moments {
		cat := m.GetCategory()
		if cat != nil && cat.Name != lastCat {
			// Cat change
			if lastCat != "" {
				catEnds[k].bottom = previousMomBottom
				k++
			} else {
				noCatEnd.bottom = previousMomBottom
			}
			lastCat = cat.Name
		}
		previousMomBottom = m.GetBottomLineNumber()
	}

	return catEnds, noCatEnd
}

func validateMissingInsertCategories(byCategory *map[string][]moment.Moment, catEnds *[]categoryEnd) error {
	missingCats := make(map[string]bool)
	for c := range *byCategory {
		missingCats[c] = true
	}
	for _, c := range *catEnds {
		delete(missingCats, c.name)
	}
	if len(missingCats) > 0 {
		return fmt.Errorf("Content is missing necessary categories to insert moments: %s", util.Keys(missingCats))
	}
	return nil
}

type categoryEnd struct {
	name   string
	bottom int
}

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
