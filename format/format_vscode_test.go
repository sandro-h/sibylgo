package format

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
)

func TestFormatCat(t *testing.T) {
	todos, _ := parse.String(`
------------------
 cat1
------------------

[] bla

------------------
 cat2
------------------

[] foo
	`)

	format := ForVSCode(todos)
	assert.Equal(t, `20,25,cat
73,78,cat
46,52,mom
99,105,mom
`, format)
}

func TestFormatMoments(t *testing.T) {
	input := `[] bla1
	[x] sub
[] bla2
	comments
	comments
[x] bla3
	comments
	comments
	`
	todos, _ := parse.String(input)

	format := ForVSCode(todos)
	assert.Equal(t, `0,7,mom
8,16,mom.done
17,24,mom
45,53,mom.done
55,63,com.done
65,73,com.done
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
	todos, _ := parse.String(input)

	format := ForVSCode(todos)

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
	todos, _ := parse.String(input)

	format := ForVSCode(todos)

	fmt.Printf("%s\n", format)
	assert.Equal(t, `1,17,mom
10,16,date
18,34,mom.until10
27,33,date
`, format)
}

func TestUnoptimizedFormat(t *testing.T) {
	// Not using parse.File because of CRLF differences impacting formatting ranges
	todos, _ := parse.String(tu.ReadTestdata(t, "TestUnoptimizedFormat", "optimized.input"))

	format := ForVSCode(todos)

	tu.AssertGoldenOutput(t, "TestUnoptimizedFormat", "unoptimized.output", format)
}

func TestOptimizedFormat(t *testing.T) {
	// Not using parse.File because of CRLF differences impacting formatting ranges
	raw := tu.ReadTestdata(t, "TestOptimizedFormat", "optimized.input")
	todos, _ := parse.String(raw)

	format := ForVSCodeOptimized(todos, raw)

	tu.AssertGoldenOutput(t, "TestOptimizedFormat", "optimized.output", format)
}

func TestOptimizedFormatWithCommentsAfterSubComment(t *testing.T) {
	input := `
[] bla1
	[] bla2
	comments
	comments

[] bla3
	`
	todos, _ := parse.String(input)

	format := ForVSCodeOptimized(todos, input)
	assert.Equal(t, `1,17,mom
39,46,mom
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
