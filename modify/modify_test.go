package modify

import (
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var testInput = `[] hello

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

var toInsert []moment.Moment

func TestMain(m *testing.M) {
	toInsert = append(toInsert,
		moment.NewSingleMoment("a new thing 1"),
		moment.NewSingleMoment("a new thing 2", moment.NewSingleMoment("a new sub thing 2.1")))
	toInsert[0].AddComment(&moment.CommentLine{Content: "my comment"})
	toInsert[0].AddComment(&moment.CommentLine{Content: "haha"})
	toInsert[1].GetSubMoments()[0].SetDone(true)

	retCode := m.Run()
	os.Exit(retCode)
}

func TestInsert(t *testing.T) {
	modified, err := Insert(testInput, toInsert)
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

func TestInsertSeparateCategories(t *testing.T) {
	toInsert[0].SetCategory(&moment.Category{Name: "cat 1"})
	modified, err := Insert(testInput, toInsert)
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

func TestInsertIntoEmptyCategory(t *testing.T) {
	testInput := `---------------
 cat 1
---------------

---------------
 cat 2
---------------

[] zonk
`
	toInsert[0].SetCategory(&moment.Category{Name: "cat 1"})
	modified, err := Insert(testInput, toInsert)
	fmt.Printf("%s\n", modified)
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
	toInsert[0].SetCategory(&moment.Category{Name: "nonexistent cat"})
	_, err := Insert(testInput, toInsert)
	assert.NotNil(t, err)
	assert.Equal(t, "Content is missing necessary categories to insert moments: [nonexistent cat]", err.Error())
}
