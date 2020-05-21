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

// Append inserts moments into the todo content. The moments are inserted at the end
// of whichever category is set for them. The categories set for the moments must all exist in the content already,
// new categories are not created. If no category is set, the moment will be appended into the "none"
// category at the start of the content.
func Append(content string, toInsert []moment.Moment) (string, error) {
	return insert(content, toInsert, false)
}

// Prepend inserts moments into the todo content. The moments are inserted at the start
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

func validateMissingInsertCategories(momsByCategory *map[string][]moment.Moment, catBoundaries *[]categoryBoundary) error {
	missingCats := make(map[string]bool)
	for c := range *momsByCategory {
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
