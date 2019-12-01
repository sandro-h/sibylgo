package moment

import (
	"time"
)

// Moment defines an interface for a generic moment in time with some significance (a todo, a event, etc).
type Moment interface {
	SetName(name string)
	SetID(id *Identifier)
	SetCategory(cat *Category)
	SetDone(done bool)
	SetPriority(prio int)
	AddSubMoment(sub Moment)
	AddComment(com *CommentLine)
	RemoveLastComment()

	GetName() string
	GetID() *Identifier
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

// Identifier is a string that uniquely identifies an element in a given todo file.
type Identifier struct {
	Value string
	DocCoords
}

// Category can be assigned to moments to categorize them.
type Category struct {
	Name     string
	Priority int
	Color    string
	DocCoords
}

// Todos defines a list of moments and moment categories
type Todos struct {
	Categories []*Category
	Moments    []Moment
}

// BaseMoment is the parent class of all moments and implements
// the Moment interface.
type BaseMoment struct {
	name       string
	id         *Identifier
	done       bool
	priority   int
	category   *Category
	comments   []*CommentLine
	subMoments []Moment
	TimeOfDay  *Date
	DocCoords
}

// SetCategory sets the category of the moment.
func (m *BaseMoment) SetCategory(cat *Category) {
	m.category = cat
}

// SetName sets the name of the moment.
func (m *BaseMoment) SetName(name string) {
	m.name = name
}

// SetID sets the optional identifier of the moment.
func (m *BaseMoment) SetID(id *Identifier) {
	m.id = id
}

// SetDone sets the done state of the moment.
func (m *BaseMoment) SetDone(done bool) {
	m.done = done
}

// SetPriority sets the priority of the moment. The higher the value, the higher the priority.
func (m *BaseMoment) SetPriority(prio int) {
	m.priority = prio
}

// AddSubMoment adds a sub moment to the moment.
func (m *BaseMoment) AddSubMoment(sub Moment) {
	m.subMoments = append(m.subMoments, sub)
}

// AddComment adds a comment to the moment.
func (m *BaseMoment) AddComment(com *CommentLine) {
	m.comments = append(m.comments, com)
}

// RemoveLastComment removes the last comment of the moment.
func (m *BaseMoment) RemoveLastComment() {
	m.comments = m.comments[:len(m.comments)-1]
}

// GetName returns the name of the moment.
func (m *BaseMoment) GetName() string {
	return m.name
}

// GetID returns the optional identifier of the moment.
func (m *BaseMoment) GetID() *Identifier {
	return m.id
}

// GetPriority returns the priority of the moment. The higher the value, the higher the priority.
func (m *BaseMoment) GetPriority() int {
	return m.priority
}

// GetCategory returns the category assigned to the moment.
func (m *BaseMoment) GetCategory() *Category {
	return m.category
}

// IsDone returns true of the moment is done.
func (m *BaseMoment) IsDone() bool {
	return m.done
}

// GetComments returns all comments of the moment.
func (m *BaseMoment) GetComments() []*CommentLine {
	return m.comments
}

// GetComment returns the comment at the given index of the moment.
func (m *BaseMoment) GetComment(index int) *CommentLine {
	return m.comments[index]
}

// GetSubMoments returns all sub moments of the moment.
func (m *BaseMoment) GetSubMoments() []Moment {
	return m.subMoments
}

// GetLastSubMoment returns the last sub moment of the moment or nil if the
// moment has no sub moments.
func (m *BaseMoment) GetLastSubMoment() Moment {
	if len(m.subMoments) == 0 {
		return nil
	}
	return m.subMoments[len(m.subMoments)-1]
}

// GetLastComment returns the last comment of the moment or nil if the
// moment has no comments.
func (m *BaseMoment) GetLastComment() *CommentLine {
	if len(m.comments) == 0 {
		return nil
	}
	return m.comments[len(m.comments)-1]
}

// GetDocCoords returns the document coordinates of the moment.
func (m *BaseMoment) GetDocCoords() DocCoords {
	return m.DocCoords
}

// GetTimeOfDay returns the time of day (0:00:00-23:59:59) of the moment,
// if defined.
func (m *BaseMoment) GetTimeOfDay() *Date {
	return m.TimeOfDay
}

// GetBottomLineNumber returns the highest line number in the text file associated
// with the moment. This could be the line number of the last comment or last sub moment.
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

// SingleMoment is a moment that only occurs once in time.
// It can happen at a single point in time, in a time range, or always (if no start or end is defined).
type SingleMoment struct {
	BaseMoment
	Start *Date
	End   *Date
}

// NewSingleMoment creates a new SingleMoment.
func NewSingleMoment(name string, subMoments ...Moment) *SingleMoment {
	m := SingleMoment{}
	m.SetName(name)
	for _, s := range subMoments {
		m.AddSubMoment(s)
	}
	return &m
}

// RecurMoment is a moment that re-occurs once or more.
// It can currently only be a single point in time, not a time range.
type RecurMoment struct {
	BaseMoment
	Recurrence Recurrence
}

const (
	// RecurDaily defines a daily recurrence.
	RecurDaily = iota
	// RecurWeekly defines a weekly recurrence.
	RecurWeekly
	// RecurMonthly defines a monthly recurrence.
	RecurMonthly
	// RecurYearly defines a yearly recurrence.
	RecurYearly
	// RecurBiWeekly defines a once every two weeks recurrence.
	RecurBiWeekly
	// RecurTriWeekly defines a once every three weeks recurrence.
	RecurTriWeekly
	// RecurQuadriWeekly defines a once every four weeks recurrence.
	RecurQuadriWeekly
)

// Recurrence defines a particular point in time that recurs.
// It therefore defines a reference or start date to pinpoint the recurring time.
// E.g. the 5th of a month, the Tuesday of a week, etc.
type Recurrence struct {
	Recurrence int
	RefDate    *Date
}

// Date defines a timestamp and the coordinates where it was defined in the text file.
type Date struct {
	Time time.Time
	DocCoords
}

// CommentLine defines a comment of a moment.
type CommentLine struct {
	Content string
	DocCoords
}

// DocCoords defines the exact location of some object (moment, date, category, etc) in the text file.
type DocCoords struct {
	// LineNumber is the line number on which this object starts. The first line of a document has line number 0.
	LineNumber int
	// Offset is the total rune offset relative to the start of the document at which this object starts.
	Offset int
	// Length is the length in runes of the object.
	Length int
}

// Instance is an actual moment at a point or range in time based on the moment definition.
// For example, a weekly recurring moment definition can yield multiple instances, one for every week
// in the given time range.
type Instance struct {
	Name         string      `json:"name"`
	Start        time.Time   `json:"start"`
	End          time.Time   `json:"end"`
	TimeOfDay    *time.Time  `json:"timeOfDay"`
	Priority     int         `json:"priority"`
	Category     *Category   `json:"-"`
	Done         bool        `json:"done"`
	EndsInRange  bool        `json:"endsInRange"`
	SubInstances []*Instance `json:"subInstances"`
}

// CloneShallow creates a clone of the moment instances without its sub instances.
func (m *Instance) CloneShallow() *Instance {
	c := Instance{
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
