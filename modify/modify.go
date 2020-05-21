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

// Append inserts moments into the todo file content. The moments are inserted at the end
// of whichever category is set for them. The categories set for the moments must all exist in the content already,
// new categories are not created. If no category is set, the moment will be appended into the "none"
// category at the start of the content.
func Append(content string, toInsert []moment.Moment) (string, error) {
	return insert(content, toInsert, false)
}

// Prepend inserts moments into the todo file content. The moments are inserted at the start
// of whichever category is set for them. The categories set for the moments must all exist in the content already,
// new categories are not created. If no category is set, the moment will be prepended into the "none"
// category at the start of the content.
func Prepend(content string, toInsert []moment.Moment) (string, error) {
	return insert(content, toInsert, true)
}

// PrependInFile inserts moments into the todo file. The moments are inserted at the start
// of whichever category is set for them. The categories set for the moments must all exist in the todo file already,
// new categories are not created. If no category is set, the moment will be prepended into the "none"
// category at the start of the todo file.
func PrependInFile(todoFile string, toInsert []moment.Moment) error {
	return modifyInFile(todoFile, func(content string) (string, error) {
		return Prepend(content, toInsert)
	})
}

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

func insert(content string, toInsert []moment.Moment, prepend bool) (string, error) {
	toInsertByCategory, noCat := groupByCategory(toInsert)

	res := ""
	todos, err := parse.String(content)
	if err != nil {
		return "", err
	}
	catBoundaries, noCatBoundary := findCategoryBoundaries(todos)

	err = validateMissingInsertCategories(&toInsertByCategory, &catBoundaries)
	if err != nil {
		return "", err
	}

	res = ""
	if noCatBoundary.getBound(prepend) == -1 {
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
		if noCatBoundary.getBound(prepend) != -1 && ln == noCatBoundary.getBound(prepend) {
			for _, m := range noCat {
				res += stringify.Moment(m)
			}
		} else if k < len(catBoundaries) && ln == catBoundaries[k].getBound(prepend) {

			list, ok := toInsertByCategory[catBoundaries[k].name]
			if ok {
				for _, m := range list {
					res += stringify.Moment(m)
				}
			}
			k++
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

func findCategoryBoundaries(todos *moment.Todos) ([]categoryBoundary, categoryBoundary) {
	noCatBoundary := categoryBoundary{"noCat", -1, -1}
	var catBoundaries []categoryBoundary
	for _, c := range todos.Categories {
		catBoundary := categoryBoundary{
			c.Name,
			c.LineNumber + 1,
			c.LineNumber + 1}
		catBoundaries = append(catBoundaries, catBoundary)
	}

	var lastCat string
	previousMomBottom := -1
	k := 0
	for _, m := range todos.Moments {
		cat := m.GetCategory()
		if cat != nil && cat.Name != lastCat {
			// Cat change
			if lastCat != "" {
				catBoundaries[k].bottom = previousMomBottom
				k++
			} else {
				noCatBoundary.bottom = previousMomBottom
			}
			lastCat = cat.Name
		}
		previousMomBottom = m.GetBottomLineNumber()
	}
	if len(catBoundaries) > 0 {
		catBoundaries[len(catBoundaries)-1].bottom = previousMomBottom
	} else {
		noCatBoundary.bottom = previousMomBottom
	}

	return catBoundaries, noCatBoundary
}

func validateMissingInsertCategories(byCategory *map[string][]moment.Moment, catBoundaries *[]categoryBoundary) error {
	missingCats := make(map[string]bool)
	for c := range *byCategory {
		missingCats[c] = true
	}
	for _, c := range *catBoundaries {
		delete(missingCats, c.name)
	}
	if len(missingCats) > 0 {
		return fmt.Errorf("Content is missing necessary categories to insert moments: %s", util.Keys(missingCats))
	}
	return nil
}

type categoryBoundary struct {
	name   string
	top    int
	bottom int
}

func (b categoryBoundary) getBound(prepend bool) int {
	if prepend {
		return b.top
	}
	return b.bottom
}

// Upsert updates moment if they exist in the todo file content, otherwise
// appends or prepends them (depending on prepend flag).
// To find existing moments, the moment ID must be set.
func Upsert(content string, toUpsert []moment.Moment, prepend bool) (string, error) {
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
		return insert(res, toInsert, prepend)
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
	// Ensure the insert order is maintained.
	i := 0
	for _, m := range toUpsert {
		id := m.GetID().Value
		if _, exists := toUpsertMap[id]; exists {
			toInsert[i] = m
			i++
		}
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
