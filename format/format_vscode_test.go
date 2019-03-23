package format

import (
	"fmt"
	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestFormatCat(t *testing.T) {
	todos, _ := parse.ParseString(`
------------------
 cat1
------------------

[] bla

------------------
 cat2
------------------

[] foo
	`)

	format := FormatVSCode(todos)
	assert.Equal(t, `20,25,cat
73,78,cat
46,52,mom
99,105,mom
`, format)
}

func TestFormatDueSoon(t *testing.T) {
	yesterday := tu.Dts(time.Now().AddDate(0, 0, -1))
	in2Days := tu.Dts(time.Now().AddDate(0, 0, 2))
	in10Days := tu.Dts(time.Now().AddDate(0, 0, 10))
	in11Days := tu.Dts(time.Now().AddDate(0, 0, 11))
	input := `
[] bla1 ($in2Days)
[] bla2 ($in2Days-$in10Days)
[] bla3 ($in2Days-$in11Days)
[] bla4 (every day)
[] bla5 ($yesterday)
	`
	input = strings.Replace(input, "$yesterday", yesterday, -1)
	input = strings.Replace(input, "$in2Days", in2Days, -1)
	input = strings.Replace(input, "$in10Days", in10Days, -1)
	input = strings.Replace(input, "$in11Days", in11Days, -1)
	todos, _ := parse.ParseString(input)

	format := FormatVSCode(todos)

	fmt.Printf("%s\n", format)
	assertUntils(t, format,
		"until2",
		"until10",
		"until0")
}

func TestDueSoonDaylightSavings(t *testing.T) {
	// Scenario:
	// It's the 23.3
	// Daylight savings time happens on 31.3
	// The moment date 3.4. will be seen as 263 hours away, instead of 264 (=11 days)
	getNow = func() time.Time { return tu.Dt("23.03.2019") }

	input := `
[] bla1 (3.4.19)
[] bla2 (2.4.19)`
	todos, _ := parse.ParseString(input)

	format := FormatVSCode(todos)

	fmt.Printf("%s\n", format)
	assert.Equal(t, `1,17,mom
10,16,date
18,34,mom.until10
27,33,date
`, format)
}

func assertUntils(t *testing.T, format string, expected ...string) {
	lines := strings.Split(format, "\n")
	k := 0
	for _, l := range lines {
		if strings.Contains(l, "until") {
			if k >= len(expected) {
				assert.Failf(t, "", "Got more than expected %d untils", len(expected))
			} else if strings.Contains(l, expected[k]) {
				k++
			} else {
				assert.Failf(t, "", "Expected next until to be %s, but was %s", expected[k], l)
			}

		}
	}
	if k < len(expected) {
		assert.Failf(t, "", "Expected %d untils, got %d", len(expected), k)
	}
}
