package modify

import (
	"bufio"
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/stringify"
	"strings"
)

// Upsert updates moment if they exist in the todo content, otherwise
// appends or prepends them (depending on prepend flag).
// To find existing moments, the moment ID must be set.
func Upsert(content string, toUpsert []moment.Moment, prepend bool) (string, error) {
	toReplace, toInsert, err := partitionByReplaceAndInsert(content, toUpsert)
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

func partitionByReplaceAndInsert(content string, toUpsert []moment.Moment) ([]replacement, []moment.Moment, error) {
	toUpsertMap := make(map[string]moment.Moment)
	for _, m := range toUpsert {
		if m.GetID() == nil {
			return nil, nil, fmt.Errorf("moment '%s' doesn't have an ID", m.GetName())
		}
		id := m.GetID().Value
		if _, exists := toUpsertMap[id]; exists {
			return nil, nil, fmt.Errorf("duplicate moment ID '%s'", id)
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
