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

const noCatIdentifier = "__noCat__"

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

	catBoundaries, toInsertByCategory, err := evaluateCategoriesForInsert(content, toInsert)
	if err != nil {
		return "", err
	}

	res := ""
	// Special case if no-cat is currently empty: put all no-cat inserts at the beginning of the content.
	if catBoundaries[0].getBound(prepend) == -1 {
		for _, m := range toInsertByCategory[noCatIdentifier] {
			res += stringify.Moment(m)
		}
		catBoundaries = catBoundaries[1:]
	}
	ln := 0
	k := 0
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		res += line + "\n"

		if k < len(catBoundaries) && ln == catBoundaries[k].getBound(prepend) {
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

func evaluateCategoriesForInsert(content string, toInsert []moment.Moment) ([]categoryBoundary, map[string][]moment.Moment, error) {
	toInsertByCategory := groupByCategory(toInsert)

	todos, err := parse.String(content)
	if err != nil {
		return nil, nil, err
	}
	catBoundaries := findCategoryBoundaries(todos)

	err = validateMissingInsertCategories(&toInsertByCategory, &catBoundaries)
	if err != nil {
		return nil, nil, err
	}

	return catBoundaries, toInsertByCategory, nil
}

func groupByCategory(moms []moment.Moment) map[string][]moment.Moment {
	byCategory := make(map[string][]moment.Moment)
	for _, m := range moms {
		cat := m.GetCategory()
		catName := noCatIdentifier
		if cat != nil {
			catName = cat.Name
		}

		list, _ := byCategory[catName]
		list = append(list, m)
		byCategory[catName] = list
	}
	return byCategory
}

// findCategoryBoundaries returns the start and end line number of each category.
// The end line number is the end of the last moment in the category, or
// the end of the category definition, if the category is empty.
// The first boundary entry is for the no-category.
func findCategoryBoundaries(todos *moment.Todos) []categoryBoundary {
	var catBoundaries []categoryBoundary
	catBoundaries = append(catBoundaries, categoryBoundary{noCatIdentifier, -1, -1})
	for _, c := range todos.Categories {
		catBoundary := categoryBoundary{
			c.Name,
			c.LineNumber + 1,
			c.LineNumber + 1}
		catBoundaries = append(catBoundaries, catBoundary)
	}

	lastCat := noCatIdentifier
	previousMomBottom := -1
	k := 0
	for _, m := range todos.Moments {
		cat := m.GetCategory()
		if cat != nil && cat.Name != lastCat {
			// Cat change
			catBoundaries[k].bottom = previousMomBottom
			k++
			lastCat = cat.Name
		}
		previousMomBottom = m.GetBottomLineNumber()
	}

	catBoundaries[len(catBoundaries)-1].bottom = previousMomBottom

	return catBoundaries
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
		return fmt.Errorf("content is missing necessary categories to insert moments: %s", util.Keys(missingCats))
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
