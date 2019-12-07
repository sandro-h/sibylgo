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
	if noCatEnd.bottom == -1 {
		for _, m := range noCat {
			res += stringify.Moment(m)
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
				res += stringify.Moment(m)
			}
		} else if ln == catEnds[k].bottom {
			list, ok := byCategory[catEnds[k].name]
			if ok {
				for _, m := range list {
					res += stringify.Moment(m)
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

// Upsert updates moment if they exist in the todo file content, otherwise inserts them.
// To find existing moments, the moment ID must be set.
func Upsert(content string, toUpsert []moment.Moment) (string, error) {
	toReplace, toInsert, err := groupByReplaceAndInsert(content, toUpsert)
	if err != nil {
		return "", err
	}

	res := ""
	if len(toReplace) > 0 {
		ln := 0
		k := 0
		scanner := bufio.NewScanner(strings.NewReader(content))
		for scanner.Scan() {
			line := scanner.Text()

			if k < len(toReplace) && toReplace[k].oldLineRange.contains(ln) {
				if ln == toReplace[k].oldLineRange.endLine {
					res += stringify.Moment(toReplace[k].new)
					k++
				}
			} else {
				res += line + "\n"
			}

			ln++
		}
	} else {
		res = content
	}

	if len(toInsert) > 0 {
		return Insert(res, toInsert)
	}
	return res, nil
}

func groupByReplaceAndInsert(content string, toUpsert []moment.Moment) ([]replacement, []moment.Moment, error) {
	toUpsertMap := make(map[string]moment.Moment)
	for _, m := range toUpsert {
		if m.GetID() == nil {
			return nil, nil, fmt.Errorf("Moment '%s' doesn't have an ID", m.GetName())
		}
		id := m.GetID().Value
		if _, exists := toUpsertMap[id]; exists {
			return nil, nil, fmt.Errorf("Duplicate moment ID '%s'", id)
		}
		toUpsertMap[id] = m
	}

	todos, err := parse.String(content)
	if err != nil {
		return nil, nil, err
	}

	var toReplace []replacement
	for _, m := range todos.Moments {
		if m.GetID() == nil {
			continue
		}

		new, found := toUpsertMap[m.GetID().Value]
		if found {
			toReplace = append(toReplace, replacement{m, new, getFullLineRange(m)})
			delete(toUpsertMap, m.GetID().Value)
		}
	}

	toInsert := make([]moment.Moment, len(toUpsertMap))
	i := 0
	for _, m := range toUpsertMap {
		toInsert[i] = m
		i++
	}

	return toReplace, toInsert, nil
}

type replacement struct {
	old          moment.Moment
	new          moment.Moment
	oldLineRange *lineRange
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
			if curRange.contains(ln) {
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

func (r *lineRange) contains(ln int) bool {
	return ln >= r.startLine && ln <= r.endLine
}
