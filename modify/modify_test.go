package modify

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ModifyFunc func(orig string, modifyData []moment.Moment) (string, error)

type testCase struct {
	name       string
	origFile   string
	updateFile string
	finalFile  string
	modifyFunc ModifyFunc
}

var wrappedUpsert = func(orig string, modifyData []moment.Moment) (string, error) { return Upsert(orig, modifyData, false) }

var testCases = [...]testCase{
	testCase{"AppendWithoutCategories", "insert.orig", "append_without_categories.input", "append_without_categories.output", Append},
	testCase{"AppendSeparateCategories", "insert.orig", "append_separate_categories.input", "append_separate_categories.output", Append},
	testCase{"AppendSecondCategory", "insert.orig", "append_second_category.input", "append_second_category.output", Append},
	testCase{"AppendIntoEmptyCategory", "insert_empty_cats.orig", "append_empty_category.input", "append_empty_category.output", Append},

	testCase{"PrependWithoutCategories", "insert.orig", "append_without_categories.input", "prepend_without_categories.output", Prepend},
	testCase{"PrependSeparateCategories", "insert.orig", "append_separate_categories.input", "prepend_separate_categories.output", Prepend},
	testCase{"PrependIntoEmptyCategory", "insert_empty_cats.orig", "append_empty_category.input", "prepend_empty_category.output", Prepend},

	testCase{"Upsert", "upsert.orig", "upsert.input", "upsert.output", wrappedUpsert},
	testCase{"UpsertAtEndOfContent", "upsert.orig", "upsert_end_of_content.input", "upsert_end_of_content.output", wrappedUpsert},
	testCase{"UpsertWithNewMoment", "upsert.orig", "upsert_with_new.input", "upsert_with_new.output", wrappedUpsert},
}

func TestAll(t *testing.T) {
	for _, tc := range testCases {
		orig := tu.ReadTestdata(t, tc.name, tc.origFile)
		toUpdate := parseTestdata(t, tc.name, tc.updateFile)
		modified, err := tc.modifyFunc(orig, toUpdate)
		assert.Nil(t, err, "Testcase: %s", tc.name)
		tu.AssertGoldenOutput(t, tc.name, tc.finalFile, modified)
	}
}

func TestMissingCategory(t *testing.T) {
	const testName = "TestMissingCategory"
	input := tu.ReadTestdata(t, testName, "insert.orig")
	toInsert := parseTestdata(t, testName, "append_without_categories.input")
	toInsert[0].SetCategory(&moment.Category{Name: "nonexistent cat"})

	_, err := Append(input, toInsert)

	assert.NotNil(t, err)
	assert.Equal(t, "Content is missing necessary categories to insert moments: [nonexistent cat]", err.Error())
}

func TestUpsertWithoutID(t *testing.T) {
	const testName = "TestUpsertWithoutID"
	input := tu.ReadTestdata(t, testName, "upsert.orig")
	toUpsert := parseTestdata(t, testName, "upsert.input")
	toUpsert[0].SetID(nil)

	_, err := Upsert(input, toUpsert, false)

	assert.NotNil(t, err)
	assert.Equal(t, "Moment 'a new thing 1' doesn't have an ID", err.Error())
}

func TestUpsertWithDuplicateID(t *testing.T) {
	const testName = "TestUpsertWithDuplicateID"
	input := tu.ReadTestdata(t, testName, "upsert.orig")
	toUpsert := parseTestdata(t, testName, "upsert.input")
	toUpsert[0].SetID(&moment.Identifier{Value: "bar"})

	_, err := Upsert(input, toUpsert, false)

	assert.NotNil(t, err)
	assert.Equal(t, "Duplicate moment ID 'bar'", err.Error())
}

func parseTestdata(t *testing.T, testName string, path string) []moment.Moment {
	update := tu.ReadTestdata(t, testName, path)
	toUpdate, err := parse.String(update)
	if err != nil {
		assert.Fail(t, "Failed to parse update content", "Testcase: %s, Error: %s, Content: %s", testName, err, update)
	}
	return toUpdate.Moments
}
