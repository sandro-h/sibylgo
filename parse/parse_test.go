package parse

import (
	"strings"
	"testing"

	"github.com/sandro-h/sibylgo/moment"
	"github.com/stretchr/testify/assert"
)

func TestSubMoments(t *testing.T) {
	todos, _ := String(`
[] 1
	[] 1.1
	[] 1.2
		[] 1.2.1
			[] 1.2.1.1
	[] 1.3
[] 2
	`)

	assert.Equal(t, 2, len(todos.Moments))
	assertMomentExists(t, todos, "1")
	assertMomentExists(t, todos, "1/1.1")
	assertMomentExists(t, todos, "1/1.2")
	assertMomentExists(t, todos, "1/1.2/1.2.1")
	assertMomentExists(t, todos, "1/1.2/1.2.1/1.2.1.1")
	assertMomentExists(t, todos, "1/1.3")
	assertMomentExists(t, todos, "2")
}

func TestComments(t *testing.T) {
	todos, _ := String(`
[] 1
	some comment
	more more more
[] 2
	second comment
	`)

	assertComments(t, todos, "1", "some comment", "more more more")
	assertComments(t, todos, "2", "second comment")
	assertDocCoords(t, 2, 7, 12, momentByPath(todos, "1").GetComment(0).DocCoords)
	assertDocCoords(t, 3, 21, 14, momentByPath(todos, "1").GetComment(1).DocCoords)
	assertDocCoords(t, 5, 42, 14, momentByPath(todos, "2").GetComment(0).DocCoords)
}

func TestCommentsWithMoreIndentation(t *testing.T) {
	todos, _ := String(`
[] 1
	some comment
		more tabs should still work
	back to one indent
	`)

	assertComments(t, todos, "1",
		"some comment",
		"	more tabs should still work",
		"back to one indent")
}

func TestSubMomentComments(t *testing.T) {
	todos, _ := String(`
[] 1
	some comment
	[] 1.1
		sub comment
		other sub comment
	`)

	assertComments(t, todos, "1", "some comment")
	assertComments(t, todos, "1/1.1", "sub comment", "other sub comment")
	assertDocCoords(t, 4, 30, 11, momentByPath(todos, "1/1.1").GetComment(0).DocCoords)
	assertDocCoords(t, 5, 44, 17, momentByPath(todos, "1/1.1").GetComment(1).DocCoords)
}

func TestMoreCommentsAfterSubMoments(t *testing.T) {
	todos, _ := String(`
[] 1
	some comment
	[] 1.1
	[] 1.2
	back to comments
	`)

	assertComments(t, todos, "1", "some comment", "back to comments")
}

func TestCommentsWithEmptyLines(t *testing.T) {
	todos, _ := String(`
[] 1
	some comment

	more more more
	`)

	assertComments(t, todos, "1", "some comment", "", "more more more")
}

func TestEmptyTrailingComments(t *testing.T) {
	todos, _ := String(`
[] 1
	some comment
	more more more


[] 2
	`)

	assertComments(t, todos, "1", "some comment", "more more more")
}

func TestCategory(t *testing.T) {
	todos, _ := String(`
[] 1
------------------
 a cat
------------------
[] 2
	[] 2.1
		[] 2.1.1
	`)

	assert.Nil(t, momentByPath(todos, "1").GetCategory())
	assert.Equal(t, 1, len(todos.Categories))
	assert.Equal(t, "a cat", todos.Categories[0].Name)
	assert.Equal(t, "a cat", momentByPath(todos, "2").GetCategory().Name)
	assert.Equal(t, "a cat", momentByPath(todos, "2/2.1").GetCategory().Name)
	assert.Equal(t, "a cat", momentByPath(todos, "2/2.1/2.1.1").GetCategory().Name)
}

func TestPriorityCategory(t *testing.T) {
	todos, _ := String(`
------------------
 a cat!!
------------------
[] 1
	`)

	assert.Equal(t, 1, len(todos.Categories))
	assert.Equal(t, "a cat", todos.Categories[0].Name)
	assert.Equal(t, 2, todos.Categories[0].Priority)
}

func TestColorCategory(t *testing.T) {
	todos, _ := String(`
------------------
 a cat! [green]
------------------
[] 1
	`)

	assert.Equal(t, 1, len(todos.Categories))
	assert.Equal(t, "a cat", todos.Categories[0].Name)
	assert.Equal(t, 1, todos.Categories[0].Priority)
	assert.Equal(t, "green", todos.Categories[0].Color)
}

func TestBadCategory(t *testing.T) {
	_, err := String(`
------------------
 a cat
 invalid more stuff
[] 1
	`)

	assert.Contains(t, err.Error(), "Expected a delimiter after category a cat")
}

func TestCategoryWithDifferentConfig(t *testing.T) {
	defer ResetConfig()
	ParseConfig.SetCategoryDelim("=====")

	todos, _ := String(`
[] 1
==================
 a cat
==================
[] 2
	`)

	assert.Nil(t, momentByPath(todos, "1").GetCategory())
	assert.Equal(t, 1, len(todos.Categories))
	assert.Equal(t, "a cat", todos.Categories[0].Name)
}

func TestUnicodeMoments(t *testing.T) {
	// Non-unicode version for range references
	// 	todos, _ := String(`
	// [] ao
	// 	hehe aa
	// 	[] aba
	// 		heyyy aa
	// 		gobobob
	// 	`)
	todos, _ := String(`
[] äö
	hehe ää
	[] äbä
		héyyy ää
		gobobob
		`)

	assertComments(t, todos, "äö", "hehe ää")
	assertComments(t, todos, "äö/äbä", "héyyy ää", "gobobob")
	assertDocCoords(t, 1, 1, 5, momentByPath(todos, "äö").GetDocCoords())
	assertDocCoords(t, 3, 16, 7, momentByPath(todos, "äö/äbä").GetDocCoords())
	assertDocCoords(t, 4, 26, 8, momentByPath(todos, "äö/äbä").GetComment(0).DocCoords)
	assertDocCoords(t, 5, 37, 7, momentByPath(todos, "äö/äbä").GetComment(1).DocCoords)
}

func TestMomentsByID(t *testing.T) {
	todos, _ := String(`
[] 1 #id1
	[] 1.1
[] 2 #id2
[] 3
	`)

	assert.Equal(t, 2, len(todos.MomentsByID))
	assert.Equal(t, "1", todos.MomentsByID["id1"].GetName())
	assert.Equal(t, "2", todos.MomentsByID["id2"].GetName())
	assert.Nil(t, todos.MomentsByID["not-existing-id"])
}

func TestDifferentIndent(t *testing.T) {
	defer ResetConfig()
	ParseConfig.SetIndent("  ")

	todos, _ := String(`
[] 1
  some comment
  [] 1.1
    sub comment
    other sub comment
	`)

	assertComments(t, todos, "1", "some comment")
	assertComments(t, todos, "1/1.1", "sub comment", "other sub comment")
	assertDocCoords(t, 4, 34, 11, momentByPath(todos, "1/1.1").GetComment(0).DocCoords)
	assertDocCoords(t, 5, 50, 17, momentByPath(todos, "1/1.1").GetComment(1).DocCoords)
}

func assertMomentExists(t *testing.T, todos *moment.Todos, path string) moment.Moment {
	mom := momentByPath(todos, path)
	if mom == nil {
		assert.Fail(t, "Expected moment with path "+path+" to exist.")
	}
	return mom
}

func assertComments(t *testing.T, todos *moment.Todos, path string, expected ...string) {
	mom := assertMomentExists(t, todos, path)
	comms := mom.GetComments()
	assert.Equal(t, len(expected), len(comms), "Number of comments")
	for i, e := range expected {
		assert.Equal(t, e, comms[i].Content, "Comment is same")
	}
}

func assertDocCoords(t *testing.T, lineNum int, offset int, len int, coords moment.DocCoords) {
	assert.Equal(t, lineNum, coords.LineNumber)
	assert.Equal(t, offset, coords.Offset)
	assert.Equal(t, len, coords.Length)
}

func momentByPath(todos *moment.Todos, path string) moment.Moment {
	parts := strings.Split(path, "/")
	return navigateToMoment(todos.Moments, &parts, 0)
}

func navigateToMoment(moms []moment.Moment, parts *[]string, index int) moment.Moment {
	for _, m := range moms {
		if m.GetName() == (*parts)[index] {
			index++
			if index == len(*parts) {
				return m
			}
			return navigateToMoment(m.GetSubMoments(), parts, index)
		}
	}
	return nil
}
