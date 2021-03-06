package cleanup

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var testInput = `
[] foo
[x] bar
	some commet
	[] bar1
	[] bar2
[] gib
	[x] ja
[x] haba
	comments1
	comments2
	comments3
[] yo`

func TestCleanupDoneTopLevel(t *testing.T) {
	kept, deleted, _ := SeparateDoneFromString(testInput, true)
	assert.Equal(t, `[] foo
[] gib
	[x] ja
[] yo`, kept)
	assert.Equal(t, `[x] bar
	some commet
	[] bar1
	[] bar2
[x] haba
	comments1
	comments2
	comments3`, deleted)
}

func TestCleanupDoneAll(t *testing.T) {
	kept, deleted, _ := SeparateDoneFromString(testInput, false)
	assert.Equal(t, `[] foo
[] gib
[] yo`, kept)
	assert.Equal(t, `[x] bar
	some commet
	[] bar1
	[] bar2
	[x] ja
[x] haba
	comments1
	comments2
	comments3`, deleted)
}

func TestMoveDoneToTrashFile(t *testing.T) {
	var testfiles []string
	writeTestFile(&testfiles, "todo.txt", testInput)
	writeTestFile(&testfiles, "trash.txt", `
[] existing trash`)
	defer deleteTestFiles(&testfiles)

	testTime := "13.01.2019 12:02:42"
	getNow = func() time.Time {
		t, _ := time.ParseInLocation("02.01.2006 15:04:05", testTime, time.Local)
		return t
	}

	MoveDoneToTrashFile(testfiles[0], testfiles[1], true)

	cleanedTodo := readFile(testfiles[0])
	trash := readFile(testfiles[1])

	assert.Equal(t, `[] foo
[] gib
	[x] ja
[] yo
`, cleanedTodo)
	assert.Equal(t, `
[] existing trash
------------------
  Trash from 13.01.2019 12:02:42
------------------
[x] bar
	some commet
	[] bar1
	[] bar2
[x] haba
	comments1
	comments2
	comments3
`, trash)
}

func TestMoveDoneToEndOfFile(t *testing.T) {
	var testfiles []string
	writeTestFile(&testfiles, "todo.txt", testInput)
	defer deleteTestFiles(&testfiles)

	MoveDoneToEndOfFile(testfiles[0], true)

	cleanedTodo := readFile(testfiles[0])

	assert.Equal(t, `[] foo
[] gib
	[x] ja
[] yo
[x] bar
	some commet
	[] bar1
	[] bar2
[x] haba
	comments1
	comments2
	comments3
`, cleanedTodo)
}

func TestTrimEmptyLinesAfterCleanup(t *testing.T) {
	input := `
[] foo

[x] bar
	some commet
	[] bar1
	[] bar2

[] gib
	[x] ja

[x] haba
	comments1
	comments2
	comments3


	
[] yo`
	kept, _, _ := SeparateDoneFromString(input, true)
	assert.Equal(t, `[] foo

[] gib
	[x] ja

[] yo`, kept)
}

func writeTestFile(testfiles *[]string, filename string, content string) {
	p := getTestFilePath(filename)
	os.Mkdir(filepath.Dir(p), 0755)
	os.WriteFile(p, []byte(content), 0644)
	*testfiles = append(*testfiles, p)
}

func deleteTestFiles(testfiles *[]string) {
	for _, f := range *testfiles {
		os.Remove(f)
	}
}

func readFile(path string) string {
	b, _ := os.ReadFile(path)
	return string(b)
}

func getTestFilePath(filename string) string {
	return filepath.Join(os.TempDir(), "sibylgo_cleanup_test", filename)
}
