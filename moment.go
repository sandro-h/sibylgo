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

type SingleMoment struct {
	BaseMoment
	start *Date
	end   *Date
}

func (m *SingleMoment) String() string {
	startStr := "nil"
	endStr := "nil"
	if m.start != nil {
		startStr = m.start.time.Format("02.01.2006 15:04") +
			fmt.Sprintf(" (%d:%d)", m.start.offset, m.start.length)
	}
	if m.end != nil {
		endStr = m.end.time.Format("02.01.2006 15:04") +
			fmt.Sprintf(" (%d:%d)", m.end.offset, m.end.length)
	}
	return fmt.Sprintf("SingleMom{name: %s, done: %t prio: %d, start: %s, end: %s, coords: %s}",
		m.name, m.done, m.priority, startStr, endStr, m.DocCoords.String())
}

type Date struct {
	time   t.Time
	offset int
	length int
}

type CommentLine struct {
	content string
}

type DocCoords struct {
	lineNumber int
	offset     int
	length     int
}

func (c *DocCoords) String() string {
	return fmt.Sprintf("%d:%d:%d", c.lineNumber, c.offset, c.length)
}
