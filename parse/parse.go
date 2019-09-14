package parse

import (
	"github.com/sandro-h/sibylgo/moment"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

const categoryDelim = "------"
const doneLBracket = "["
const doneRBracket = ']'
const doneMark = 'x'
const doneMarkUpper = 'X'
const priorityMark = '!'
const indentChar = "\t"

type parserState struct {
	todos       *moment.Todos
	curCategory *moment.Category
	scanner     *LineScanner
}

// File parses a text file into a Todos object.
func File(path string) (*moment.Todos, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parse(NewFileLineScanner(file))
}

// String parses a string into a Todos object. The string
// is usually the content of a text file and therefore contains
// one or more lines.
func String(str string) (*moment.Todos, error) {
	return parse(NewLineStringScanner(str))
}

// Reader parses the contents returned by the given reader into
// a Todos object.
func Reader(reader io.Reader) (*moment.Todos, error) {
	return parse(NewLineScanner(reader))
}

func parse(scanner *LineScanner) (*moment.Todos, error) {
	parserState := parserState{todos: &moment.Todos{}, scanner: scanner}
	for parserState.scanner.Scan() {
		err := parserState.handleLine(parserState.scanner.Line())
		if err != nil {
			return nil, err
		}
	}

	if err := parserState.scanner.Err(); err != nil {
		return nil, err
	}

	return parserState.todos, nil
}

func (p *parserState) handleLine(line *Line) error {
	if line.IsEmpty() {
		return nil
	}
	var err error
	if line.HasPrefix(categoryDelim) {
		err = p.handleCategoryLine(line)
	} else if line.HasPrefix(doneLBracket) {
		err = p.handleMomentLine(line)
	}
	//fmt.Printf("%s\n", line.content)
	return err
}

func (p *parserState) handleCategoryLine(line *Line) error {
	ok, catLine := p.scanner.ScanAndLine()

	p.curCategory = parseCategory(catLine)
	p.todos.Categories = append(p.todos.Categories, p.curCategory)

	ok, nxt := p.scanner.ScanAndLine()
	if !ok {
		return newParseError(catLine,
			"Expected a delimiter after category %s, but reached end",
			p.curCategory.Name)
	}
	if !nxt.HasPrefix(categoryDelim) {
		return newParseError(nxt,
			"Expected a delimiter after category %s, got %s",
			p.curCategory.Name, nxt.Content())
	}

	return nil
}

func parseCategory(line *Line) *moment.Category {
	lineVal := line.Content()

	prio, lineVal := parsePriority(lineVal)

	return &moment.Category{
		Name:      lineVal,
		Priority:  prio,
		DocCoords: moment.DocCoords{LineNumber: line.LineNumber(), Offset: line.Offset(), Length: line.Length()}}
}

func (p *parserState) handleMomentLine(line *Line) error {
	mom, err := p.parseFullMoment(line, line.TrimmedContent(), "")
	if err != nil {
		return err
	}
	p.todos.Moments = append(p.todos.Moments, mom)
	return nil
}

func (p *parserState) parseFullMoment(line *Line, lineVal string, indent string) (moment.Moment, error) {
	mom, err := parseMoment(line, lineVal)
	if err != nil {
		return nil, err
	}
	mom.SetCategory(p.curCategory)

	err = p.parseCommentsAndSubMoments(mom, indent)
	if err != nil {
		return nil, err
	}

	return mom, nil
}

func (p *parserState) parseCommentsAndSubMoments(mom moment.Moment, indent string) error {
	nextIndent := indent + indentChar
	for p.scanner.Scan() {
		line := p.scanner.Line()
		if line.HasPrefix(nextIndent) {
			p.handleSubLine(mom, line, line.Content()[len(nextIndent):], indent)
		} else if line.IsEmpty() && len(mom.GetComments()) > 0 {
			// special case: treat empty line between comments as a comment
			comment := &moment.CommentLine{
				Content:   "",
				DocCoords: moment.DocCoords{LineNumber: line.LineNumber(), Offset: line.Offset(), Length: 0}}
			mom.AddComment(comment)
		} else {
			p.scanner.Unscan()
			break
		}
	}

	// Remove trailing empty comments
	lc := mom.GetLastComment()
	for lc != nil && len(lc.Content) == 0 {
		mom.RemoveLastComment()
		lc = mom.GetLastComment()
	}

	return nil
}

func (p *parserState) handleSubLine(mom moment.Moment, line *Line, lineVal string, indent string) error {
	if strings.HasPrefix(lineVal, doneLBracket) {
		subMom, err := p.parseFullMoment(line, lineVal, indent+indentChar)
		if err != nil {
			return err
		}
		mom.AddSubMoment(subMom)
	} else {
		// Assume it's a comment
		comment := &moment.CommentLine{
			Content: lineVal,
			DocCoords: moment.DocCoords{LineNumber: line.LineNumber(),
				Offset: line.Offset() + len(indent+indentChar),
				Length: utf8.RuneCountInString(lineVal)}}
		mom.AddComment(comment)
	}
	return nil
}