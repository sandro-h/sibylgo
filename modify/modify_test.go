package modify

import (
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/stretchr/testify/assert"
	"testing"
)

var insertTestInput = `[] hello

---------------
 cat 1
---------------
[] foo
[x] bar
	some commet
	[] bar1
	[] bar2
---------------
 cat 2
---------------

[] zonk
`

var upsertTestInput = `---------------
 cat 1
---------------
[] foo #foo
[x] bar #bar
	some commet
	[] bar1
	[] bar2
---------------
 cat 2
---------------

[] zonk #zonk
`

var toInsert []moment.Moment
var toUpsert []moment.Moment

func resetToInsert() {
	toInsert = make([]moment.Moment, 0)
	toInsert = append(toInsert,
		moment.NewSingleMoment("a new thing 1"),
		moment.NewSingleMoment("a new thing 2", moment.NewSingleMoment("a new sub thing 2.1")))
	toInsert[0].AddComment(&moment.CommentLine{Content: "my comment"})
	toInsert[0].AddComment(&moment.CommentLine{Content: "haha"})
	toInsert[1].GetSubMoments()[0].SetDone(true)
}

func resetToUpsert() {
	toUpsert = make([]moment.Moment, 0)
	toUpsert = append(toUpsert,
		moment.NewSingleMoment("a new thing 1"),
		moment.NewSingleMoment("a new thing 2", moment.NewSingleMoment("a new sub thing 2.1")))
	toUpsert[0].SetID(&moment.Identifier{Value: "foo"})
	toUpsert[1].SetID(&moment.Identifier{Value: "bar"})
	toUpsert[0].AddComment(&moment.CommentLine{Content: "my comment"})
	toUpsert[0].AddComment(&moment.CommentLine{Content: "haha"})
	toUpsert[1].GetSubMoments()[0].SetDone(true)
}

func TestAppendWithoutCategories(t *testing.T) {
	resetToInsert()

	modified, err := Append(insertTestInput, toInsert)

	assert.Nil(t, err)
	assert.Equal(t, `[] hello
[] a new thing 1
	my comment
	haha
[] a new thing 2
	[x] a new sub thing 2.1

---------------
 cat 1
---------------
[] foo
[x] bar
	some commet
	[] bar1
	[] bar2
---------------
 cat 2
---------------

[] zonk
`, modified)
}

func TestAppendSeparateCategories(t *testing.T) {
	resetToInsert()
	toInsert[0].SetCategory(&moment.Category{Name: "cat 1"})

	modified, err := Append(insertTestInput, toInsert)

	fmt.Printf("%s\n", modified)
	assert.Nil(t, err)
	assert.Equal(t, `[] hello
[] a new thing 2
	[x] a new sub thing 2.1

---------------
 cat 1
---------------
[] foo
[x] bar
	some commet
	[] bar1
	[] bar2
[] a new thing 1
	my comment
	haha
---------------
 cat 2
---------------

[] zonk
`, modified)
}

func TestAppendIntoEmptyCategory(t *testing.T) {
	insertTestInput := `---------------
 cat 1
---------------

---------------
 cat 2
---------------

[] zonk
`
	resetToInsert()
	toInsert[0].SetCategory(&moment.Category{Name: "cat 1"})

	modified, err := Append(insertTestInput, toInsert)

	assert.Nil(t, err)
	assert.Equal(t, `[] a new thing 2
	[x] a new sub thing 2.1
---------------
 cat 1
---------------
[] a new thing 1
	my comment
	haha

---------------
 cat 2
---------------

[] zonk
`, modified)
}

func TestPrependWithoutCategories(t *testing.T) {
	resetToInsert()

	modified, err := Prepend(insertTestInput, toInsert)

	assert.Nil(t, err)
	assert.Equal(t, `[] a new thing 1
	my comment
	haha
[] a new thing 2
	[x] a new sub thing 2.1
[] hello

---------------
 cat 1
---------------
[] foo
[x] bar
	some commet
	[] bar1
	[] bar2
---------------
 cat 2
---------------

[] zonk
`, modified)
}

func TestPrependSeparateCategories(t *testing.T) {
	resetToInsert()
	toInsert[0].SetCategory(&moment.Category{Name: "cat 1"})

	modified, err := Prepend(insertTestInput, toInsert)

	fmt.Printf("%s\n", modified)
	assert.Nil(t, err)
	assert.Equal(t, `[] a new thing 2
	[x] a new sub thing 2.1
[] hello

---------------
 cat 1
---------------
[] a new thing 1
	my comment
	haha
[] foo
[x] bar
	some commet
	[] bar1
	[] bar2
---------------
 cat 2
---------------

[] zonk
`, modified)
}

func TestPrependIntoEmptyCategory(t *testing.T) {
	insertTestInput := `---------------
 cat 1
---------------

---------------
 cat 2
---------------

[] zonk
`
	resetToInsert()
	toInsert[0].SetCategory(&moment.Category{Name: "cat 1"})

	modified, err := Prepend(insertTestInput, toInsert)

	assert.Nil(t, err)
	assert.Equal(t, `[] a new thing 2
	[x] a new sub thing 2.1
---------------
 cat 1
---------------
[] a new thing 1
	my comment
	haha

---------------
 cat 2
---------------

[] zonk
`, modified)
}

func TestMissingCategory(t *testing.T) {
	resetToInsert()
	toInsert[0].SetCategory(&moment.Category{Name: "nonexistent cat"})

	_, err := Append(insertTestInput, toInsert)

	assert.NotNil(t, err)
	assert.Equal(t, "Content is missing necessary categories to insert moments: [nonexistent cat]", err.Error())
}

func TestUpsert(t *testing.T) {
	resetToUpsert()

	modified, err := Upsert(upsertTestInput, toUpsert, false)

	assert.Nil(t, err)
	assert.Equal(t, `---------------
 cat 1
---------------
[] a new thing 1 #foo
	my comment
	haha
[] a new thing 2 #bar
	[x] a new sub thing 2.1
---------------
 cat 2
---------------

[] zonk #zonk
`, modified)
}

func TestUpsertAtEndOfContent(t *testing.T) {
	var customToUpsert []moment.Moment
	customToUpsert = append(customToUpsert, toUpsert[0])
	customToUpsert[0].SetID(&moment.Identifier{Value: "zonk"})

	modified, err := Upsert(upsertTestInput, customToUpsert, false)

	assert.Nil(t, err)
	assert.Equal(t, `---------------
 cat 1
---------------
[] foo #foo
[x] bar #bar
	some commet
	[] bar1
	[] bar2
---------------
 cat 2
---------------

[] a new thing 1 #zonk
	my comment
	haha
`, modified)
}

func TestUpsertWithNewMoment(t *testing.T) {
	resetToUpsert()
	toUpsert = append(toUpsert,
		moment.NewSingleMoment("a new thing 3"))
	toUpsert[len(toUpsert)-1].SetID(&moment.Identifier{Value: "newid123"})
	toUpsert[len(toUpsert)-1].SetCategory(&moment.Category{Name: "cat 1"})

	modified, err := Upsert(upsertTestInput, toUpsert, false)

	assert.Nil(t, err)
	assert.Equal(t, `---------------
 cat 1
---------------
[] a new thing 1 #foo
	my comment
	haha
[] a new thing 2 #bar
	[x] a new sub thing 2.1
[] a new thing 3 #newid123
---------------
 cat 2
---------------

[] zonk #zonk
`, modified)
}

func TestUpsertWithoutID(t *testing.T) {
	resetToUpsert()
	toUpsert[0].SetID(nil)

	_, err := Upsert(upsertTestInput, toUpsert, false)

	assert.NotNil(t, err)
	assert.Equal(t, "Moment 'a new thing 1' doesn't have an ID", err.Error())
}

func TestUpsertWithDuplicateID(t *testing.T) {
	resetToUpsert()
	toUpsert[0].SetID(&moment.Identifier{Value: "bar"})

	_, err := Upsert(upsertTestInput, toUpsert, false)

	assert.NotNil(t, err)
	assert.Equal(t, "Duplicate moment ID 'bar'", err.Error())
}
