package main

import (
	"fmt"
	t "time"
)

type Moment interface {
	SetName(name string)
	SetCategory(cat *Category)
	SetDone(done bool)
	SetPriority(prio int)
	AddSubMoment(sub Moment)
	AddComment(com *CommentLine)
	RemoveLastComment()

	GetName() string
	GetPriority() int
	IsDone() bool
	GetComments() []*CommentLine
	GetComment(index int) *CommentLine
	GetSubMoments() []Moment
	GetLastComment() *CommentLine
}

type Todos struct {
	categories []*Category
	moments    []Moment
}

type Category struct {
	name     string
	priority int
	DocCoords
}

func (c *Category) String() string {
	return fmt.Sprintf("Category{name: %s, prio: %d, coords: %s}", c.name, c.priority, c.DocCoords.String())
}

type BaseMoment struct {
	name       string
	done       bool
	priority   int
	category   *Category
	comments   []*CommentLine
	subMoments []Moment
	DocCoords
}

func (m *BaseMoment) SetCategory(cat *Category) {
	m.category = cat
}

func (m *BaseMoment) SetName(name string) {
	m.name = name
}

func (m *BaseMoment) SetDone(done bool) {
	m.done = done
}

func (m *BaseMoment) SetPriority(prio int) {
	m.priority = prio
}

func (m *BaseMoment) AddSubMoment(sub Moment) {
	m.subMoments = append(m.subMoments, sub)
}

func (m *BaseMoment) AddComment(com *CommentLine) {
	m.comments = append(m.comments, com)
}

func (m *BaseMoment) RemoveLastComment() {
	m.comments = m.comments[:len(m.comments)-1]
}

func (m *BaseMoment) GetName() string {
	return m.name
}

func (m *BaseMoment) GetPriority() int {
	return m.priority
}

func (m *BaseMoment) IsDone() bool {
	return m.done
}

func (m *BaseMoment) GetComments() []*CommentLine {
	return m.comments
}

func (m *BaseMoment) GetComment(index int) *CommentLine {
	return m.comments[index]
}

func (m *BaseMoment) GetSubMoments() []Moment {
	return m.subMoments
}

func (m *BaseMoment) GetLastComment() *CommentLine {
	if len(m.comments) == 0 {
		return nil
	}
	return m.comments[len(m.comments)-1]
}

type SingleMoment struct {
	BaseMoment
	start *Date
	end   *Date
}

func (m *SingleMoment) String() string {
	startStr := "nil"
	endStr := "nil"
	if m.start != nil {
		startStr = m.start.time.Format("02.01.06 15:04") +
			fmt.Sprintf(" (%d:%d)", m.start.offset, m.start.length)
	}
	if m.end != nil {
		endStr = m.end.time.Format("02.01.06 15:04") +
			fmt.Sprintf(" (%d:%d)", m.end.offset, m.end.length)
	}
	return fmt.Sprintf("SingleMom{name: %s, done: %t prio: %d, start: %s, end: %s, comms: %d, coords: %s}",
		m.name, m.done, m.priority, startStr, endStr, len(m.comments), m.DocCoords.String())
}

type Date struct {
	time t.Time
	DocCoords
}

type CommentLine struct {
	content string
	DocCoords
}

type DocCoords struct {
	lineNumber int
	offset     int
	length     int
}

func (c *DocCoords) String() string {
	return fmt.Sprintf("%d:%d:%d", c.lineNumber, c.offset, c.length)
}
