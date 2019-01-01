package main

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSubMoments(t *testing.T) {
	todos, _ := ParseString(`
[] 1
	[] 1.1
	[] 1.2
		[] 1.2.1
			[] 1.2.1.1
	[] 1.3
[] 2
	`)

	assert.Equal(t, 2, len(todos.moments))
	assertMomentExists(t, todos, "1")
	assertMomentExists(t, todos, "1/1.1")
	assertMomentExists(t, todos, "1/1.2")
	assertMomentExists(t, todos, "1/1.2/1.2.1")
	assertMomentExists(t, todos, "1/1.2/1.2.1/1.2.1.1")
	assertMomentExists(t, todos, "1/1.3")
	assertMomentExists(t, todos, "2")
}

func TestComments(t *testing.T) {
	todos, _ := ParseString(`
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
	todos, _ := ParseString(`
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
	todos, _ := ParseString(`
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
	todos, _ := ParseString(`
[] 1
	some comment
	[] 1.1
	[] 1.2
	back to comments
	`)

	assertComments(t, todos, "1", "some comment", "back to comments")
}

func TestCommentsWithEmptyLines(t *testing.T) {
	todos, _ := ParseString(`
[] 1
	some comment

	more more more
	`)

	assertComments(t, todos, "1", "some comment", "", "more more more")
}

func TestEmptyTrailingComments(t *testing.T) {
	todos, _ := ParseString(`
[] 1
	some comment
	more more more


[] 2
	`)

	assertComments(t, todos, "1", "some comment", "more more more")
}

func TestCategory(t *testing.T) {
	todos, _ := ParseString(`
[] 1
------------------
 a cat
------------------
[] 2
	[] 2.1
		[] 2.1.1
	`)

	assert.Nil(t, momentByPath(todos, "1").GetCategory())
	assert.Equal(t, 1, len(todos.categories))
	assert.Equal(t, "a cat", todos.categories[0].name)
	assert.Equal(t, "a cat", momentByPath(todos, "2").GetCategory().name)
	assert.Equal(t, "a cat", momentByPath(todos, "2/2.1").GetCategory().name)
	assert.Equal(t, "a cat", momentByPath(todos, "2/2.1/2.1.1").GetCategory().name)
}

func TestPriorityCategory(t *testing.T) {
	todos, _ := ParseString(`
------------------
 a cat!!
------------------
[] 1
	`)

	assert.Equal(t, 1, len(todos.categories))
	assert.Equal(t, "a cat", todos.categories[0].name)
	assert.Equal(t, 2, todos.categories[0].priority)
}

func TestBadCategory(t *testing.T) {
	_, err := ParseString(`
------------------
 a cat
 invalid more stuff
[] 1
	`)

	assert.Contains(t, err.Error(), "Expected a delimiter after category a cat")
}

func assertMomentExists(t *testing.T, todos *Todos, path string) Moment {
	mom := momentByPath(todos, path)
	if mom == nil {
		assert.Fail(t, "Expected moment with path "+path+" to exist.")
	}
	return mom
}

func assertComments(t *testing.T, todos *Todos, path string, expected ...string) {
	mom := assertMomentExists(t, todos, path)
	comms := mom.GetComments()
	assert.Equal(t, len(expected), len(comms), "Number of comments")
	for i, e := range expected {
		assert.Equal(t, e, comms[i].content, "Comment is same")
	}
}

func assertDocCoords(t *testing.T, lineNum int, offset int, len int, coords DocCoords) {
	assert.Equal(t, lineNum, coords.lineNumber)
	assert.Equal(t, offset, coords.offset)
	assert.Equal(t, len, coords.length)
}

func momentByPath(todos *Todos, path string) Moment {
	parts := strings.Split(path, "/")
	return navigateToMoment(todos.moments, &parts, 0)
}

func navigateToMoment(moms []Moment, parts *[]string, index int) Moment {
	for _, m := range moms {
		if m.GetName() == (*parts)[index] {
			index++
			if index == len(*parts) {
				return m
			} else {
				return navigateToMoment(m.GetSubMoments(), parts, index)
			}
		}
	}
	return nil
}
