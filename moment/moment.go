package moment

import (
	"github.com/sandro-h/sibylgo/util"
	"time"
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

	CreateInstances(from time.Time, to time.Time) []*MomentInstance
}

type Todos struct {
	Categories []*Category
	Moments    []Moment
}

type Category struct {
	Name     string
	Priority int
	DocCoords
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
	Start *Date
	End   *Date
}

func (m *SingleMoment) CreateInstances(from time.Time, to time.Time) []*MomentInstance {
	start := util.GetUpperBound(&from, dateTm(m.Start))
	end := util.GetLowerBound(&to, dateTm(m.End))
	if end.Before(start) {
		// Not actually in range
		return nil
	}

	inst := MomentInstance{Start: start, End: end}
	inst.EndsInRange = m.End != nil && !m.End.Time.After(end)
	return []*MomentInstance{&inst}
}

type RecurMoment struct {
	BaseMoment
	Recurrence Recurrence
}

func (m *RecurMoment) CreateInstances(from time.Time, to time.Time) []*MomentInstance {
	var insts []*MomentInstance
	for it := NewRecurIterator(m.Recurrence, from, to); it.HasNext(); {
		start := it.Next()
		inst := MomentInstance{Start: start, End: util.SetToEndOfDay(start)}
		inst.EndsInRange = true
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
	Recurrence int
	RefDate    *Date
}

type Date struct {
	Time time.Time
	DocCoords
}

func dateTm(dt *Date) *time.Time {
	if dt == nil {
		return nil
	}
	return &dt.Time
}

type CommentLine struct {
	Content string
	DocCoords
}

type DocCoords struct {
	LineNumber int
	Offset     int
	Length     int
}

type MomentInstance struct {
	Start        time.Time
	End          time.Time
	EndsInRange  bool
	SubInstances []*MomentInstance
}
