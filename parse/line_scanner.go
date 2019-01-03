package parse

import (
	"bufio"
	"io"
	"os"
	"strings"
)

type LineScanner struct {
	reader      *bufio.Reader
	scanner     *bufio.Scanner
	init        bool
	undoing     bool
	hasNextLine bool
	prevLine    string
	nextLine    string
	line        Line
	lineNumber  int
	offset      int
	lastOffset  int
	nlLen       int
}

type Line struct {
	content    string
	lineNumber int
	offset     int
}

func NewLineScanner(reader io.Reader) *LineScanner {
	br := bufio.NewReader(reader)
	return &LineScanner{reader: br, scanner: bufio.NewScanner(br), nlLen: 1}
}

func NewFileLineScanner(file *os.File) *LineScanner {
	return NewLineScanner(file)
}

func NewLineStringScanner(str string) *LineScanner {
	return NewLineScanner(strings.NewReader(str))
}

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
	s.line = Line{
		content:    ret,
		lineNumber: s.lineNumber - 2,
		offset:     s.lastOffset - len(ret) - s.nlLen}
	return true
}

func (s *LineScanner) Line() *Line {
	return &s.line
}

func (s *LineScanner) ScanAndLine() (bool, *Line) {
	if !s.Scan() {
		return false, nil
	} else {
		return true, s.Line()
	}
}

func (s *LineScanner) Unscan() {
	s.undoing = true
}

func (s *LineScanner) Err() error {
	return s.scanner.Err()
}

func (s *LineScanner) updateNextLine() {
	s.lineNumber++
	s.lastOffset = s.offset

	if s.scanner.Scan() {
		s.nextLine = s.scanner.Text()
		s.hasNextLine = true
		s.offset += len(s.nextLine) + s.nlLen
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

func (l *Line) LineNumber() int {
	return l.lineNumber
}

func (l *Line) Offset() int {
	return l.offset
}

func (l *Line) Length() int {
	return len(l.content)
}

func (l *Line) Content() string {
	return l.content
}

func (l *Line) TrimmedContent() string {
	return strings.TrimSpace(l.content)
}

func (l *Line) IsEmpty() bool {
	return len(l.TrimmedContent()) == 0
}

func (l *Line) HasPrefix(prefix string) bool {
	return strings.HasPrefix(l.content, prefix)
}
