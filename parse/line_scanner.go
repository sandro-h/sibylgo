package parse

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

// LineScanner reads content line-by-line and allows undoing a single line read
// in case the current read line is not what we expected while parsing.
type LineScanner struct {
	reader      *bufio.Reader
	scanner     *bufio.Scanner
	init        bool
	undoing     bool
	hasNextLine bool
	prevLine    string
	nextLine    string
	line        *Line
	lineNumber  int
	offset      int
	lastOffset  int
	nlLen       int
}

// Line contains the string content of a line and it's coordinates in the document.
type Line struct {
	content    string
	lineNumber int
	offset     int
}

// NewLineScanner creates a new LineScanner for a generic Reader.
func NewLineScanner(reader io.Reader) *LineScanner {
	br := bufio.NewReader(reader)
	return &LineScanner{reader: br, scanner: bufio.NewScanner(br), nlLen: 1}
}

// NewFileLineScanner creates a new LineScanner for a text file.
func NewFileLineScanner(file *os.File) *LineScanner {
	return NewLineScanner(file)
}

// NewLineStringScanner creates a new LineScanner for a string.
func NewLineStringScanner(str string) *LineScanner {
	return NewLineScanner(strings.NewReader(str))
}

// Scan reads the next line and returns true if one was found.
// It returns false if the end of the content was reached.
func (s *LineScanner) Scan() bool {
	if !s.init {
		s.detectNewlineLength()
		s.updateNextLine()
		s.init = true
	}

	if !s.undoing && !s.hasNextLine {
		return false
	}

	var ret string
	if s.undoing {
		ret = s.prevLine
		s.undoing = false
	} else {
		ret = s.nextLine
		s.updateNextLine()
	}
	s.prevLine = ret
	s.line = &Line{
		content:    ret,
		lineNumber: s.lineNumber - 2,
		offset:     s.lastOffset - utf8.RuneCountInString(ret) - s.nlLen}
	return true
}

// Line returns the line read in the preceding Scan call.
func (s *LineScanner) Line() *Line {
	return s.line
}

// ScanAndLine is a combination of calling Scan() and then Line().
func (s *LineScanner) ScanAndLine() (bool, *Line) {
	if !s.Scan() {
		return false, nil
	}
	return true, s.Line()
}

// Unscan undoes the previous Scan() call, such that calling Scan() again
// will have the same effect as the previous Scan() call.
// Calling Unscan repeatedly has no effect, only the single immediately preceding
// Scan call can be undone.
func (s *LineScanner) Unscan() {
	s.undoing = true
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s *LineScanner) Err() error {
	return s.scanner.Err()
}

func (s *LineScanner) updateNextLine() {
	s.lineNumber++
	s.lastOffset = s.offset

	if s.scanner.Scan() {
		s.nextLine = s.scanner.Text()
		s.hasNextLine = true
		s.offset += utf8.RuneCountInString(s.nextLine) + s.nlLen
	} else {
		s.hasNextLine = false
	}
}

func (s *LineScanner) detectNewlineLength() {
	buf, _ := s.reader.Peek(512)
	for i := range buf {
		if buf[i] == '\n' {
			if i > 0 && buf[i-1] == '\r' {
				s.nlLen = 2
			} else {
				s.nlLen = 1
			}
			break
		}
	}
}

// LineNumber returns the line number of the line.
// The first line of a document has line number 0.
func (l *Line) LineNumber() int {
	return l.lineNumber
}

// Offset returns the total rune offset of the start of the line
// relative to the start of the document.
func (l *Line) Offset() int {
	return l.offset
}

// Length returns the number of runes of the line.
func (l *Line) Length() int {
	return utf8.RuneCountInString(l.content)
}

// Content returns the content string of the line.
func (l *Line) Content() string {
	return l.content
}

// TrimmedContent returns the content trimmed of spaces.
func (l *Line) TrimmedContent() string {
	return strings.TrimSpace(l.content)
}

// IsEmpty returns true if the trimmed content is empty. This means a line
// containing nothing but spaces is considered empty.
func (l *Line) IsEmpty() bool {
	return len(l.TrimmedContent()) == 0
}

// HasPrefix returns true if the line starts with the prefix string.
func (l *Line) HasPrefix(prefix string) bool {
	return strings.HasPrefix(l.content, prefix)
}

// HasRunePrefix returns true if the line starts with the prefix rune.
func (l *Line) HasRunePrefix(prefix rune) bool {
	return HasRunePrefix(l.content, prefix)
}
