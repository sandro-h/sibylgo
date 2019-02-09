package moment

import (
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
	GetLastSubMoment() Moment
	GetLastComment() *CommentLine
	GetDocCoords() DocCoords
	GetTimeOfDay() *Date
	GetBottomLineNumber() int
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
	TimeOfDay  *Date
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

func (m *BaseMoment) GetLastSubMoment() Moment {
	if len(m.subMoments) == 0 {
		return nil
	}
	return m.subMoments[len(m.subMoments)-1]
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

func (m *BaseMoment) GetTimeOfDay() *Date {
	return m.TimeOfDay
}

func (m *BaseMoment) GetBottomLineNumber() int {
	max := m.GetDocCoords().LineNumber
	if m.GetLastComment() != nil && m.GetLastComment().LineNumber > max {
		max = m.GetLastComment().LineNumber
	}
	if m.GetLastSubMoment() != nil {
		subMax := m.GetLastSubMoment().GetBottomLineNumber()
		if subMax > max {
			max = subMax
		}
	}
	return max
}

type SingleMoment struct {
	BaseMoment
	Start *Date
	End   *Date
}

type RecurMoment struct {
	BaseMoment
	Recurrence Recurrence
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
	Name         string            `json:"name"`
	Start        time.Time         `json:"start"`
	End          time.Time         `json:"end"`
	TimeOfDay    *time.Time        `json:"timeOfDay"`
	Priority     int               `json:"priority"`
	Done         bool              `json:"done"`
	EndsInRange  bool              `json:"endsInRange"`
	SubInstances []*MomentInstance `json:"subInstances"`
}

// CloneShallow creates a clone of the moment instances without its sub instances.
func (m *MomentInstance) CloneShallow() *MomentInstance {
	c := MomentInstance{
		Name:        m.Name,
		Start:       m.Start,
		Priority:    m.Priority,
		Done:        m.Done,
		EndsInRange: m.EndsInRange}
	if m.TimeOfDay != nil {
		cp := *m.TimeOfDay
		c.TimeOfDay = &cp
	}
	return &c
}
