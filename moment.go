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
	GetCategory() *Category
	IsDone() bool
	GetComments() []*CommentLine
	GetComment(index int) *CommentLine
	GetSubMoments() []Moment
	GetLastComment() *CommentLine
	GetDocCoords() DocCoords

	CreateInstances(from t.Time, to t.Time) []*MomentInstance
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

func (m *BaseMoment) GetCategory() *Category {
	return m.category
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

func (m *BaseMoment) GetDocCoords() DocCoords {
	return m.DocCoords
}

type SingleMoment struct {
	BaseMoment
	start *Date
	end   *Date
}

func (m *SingleMoment) CreateInstances(from t.Time, to t.Time) []*MomentInstance {
	start := getUpperBound(&from, dateTm(m.start))
	end := getLowerBound(&to, dateTm(m.end))
	if end.Before(start) {
		// Not actually in range
		return nil
	}

	inst := MomentInstance{start: start, end: end}
	inst.endsInRange = m.end != nil && !m.end.time.After(end)
	return []*MomentInstance{&inst}
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

type RecurMoment struct {
	BaseMoment
	recurrence Recurrence
}

func (m *RecurMoment) CreateInstances(from t.Time, to t.Time) []*MomentInstance {
	var insts []*MomentInstance
	for it := NewRecurIterator(m.recurrence, from, to); it.HasNext(); {
		start := it.Next()
		inst := MomentInstance{start: start, end: setToEndOfDay(start)}
		inst.endsInRange = true
		insts = append(insts, &inst)
	}
	return insts
}

const (
	RE_DAILY = iota
	RE_WEEKLY
	RE_MONTHLY
	RE_YEARLY
)

type Recurrence struct {
	recurrence int
	refDate    *Date
}

type Date struct {
	time t.Time
	DocCoords
}

func dateTm(dt *Date) *t.Time {
	if dt == nil {
		return nil
	}
	return &dt.time
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

type MomentInstance struct {
	start        t.Time
	end          t.Time
	endsInRange  bool
	subInstances []*MomentInstance
}

func GenerateInstances(mom Moment, from t.Time, to t.Time) []*MomentInstance {
	return generateInstances(mom, from, to, true)
}

func GenerateInstancesWithoutSubs(mom Moment, from t.Time, to t.Time) []*MomentInstance {
	return generateInstances(mom, from, to, false)
}

func generateInstances(mom Moment, from t.Time, to t.Time, inclSubs bool) []*MomentInstance {
	insts := mom.CreateInstances(from, to)
	// Sub moments:
	if inclSubs {
		for _, inst := range insts {
			var subInsts []*MomentInstance
			for _, sub := range mom.GetSubMoments() {
				subInsts = append(subInsts, generateInstances(sub, inst.start, inst.end, inclSubs)...)
			}
			inst.subInstances = subInsts
		}
	}
	return insts
}
