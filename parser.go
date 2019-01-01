package main

import (
	"os"
	"strings"
)

const categoryDelim = "------"
const doneLBracket = "["
const doneRBracket = ']'
const doneMark = 'x'
const doneMarkUpper = 'X'
const priorityMark = '!'
const indentChar = "\t"

type Parser struct {
	todos       *Todos
	curCategory *Category
	scanner     *LineScanner
}

func ParseFile(path string) (*Todos, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parse(NewFileLineScanner(file))
}

func ParseString(str string) (*Todos, error) {
	return parse(NewLineStringScanner(str))
}

func parse(scanner *LineScanner) (*Todos, error) {
	parser := Parser{todos: &Todos{}, scanner: scanner}
	for parser.scanner.Scan() {
		err := parser.handleLine(parser.scanner.Line())
		if err != nil {
			return nil, err
		}
	}

	if err := parser.scanner.Err(); err != nil {
		return nil, err
	}

	return parser.todos, nil
}

func (p *Parser) handleLine(line *Line) error {
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

func (p *Parser) handleCategoryLine(line *Line) error {
	ok, catLine := p.scanner.ScanAndLine()

	p.curCategory = parseCategory(catLine)
	p.todos.categories = append(p.todos.categories, p.curCategory)

	ok, nxt := p.scanner.ScanAndLine()
	if !ok {
		return newParseError(catLine,
			"Expected a delimiter after category %s, but reached end",
			p.curCategory.name)
	}
	if !nxt.HasPrefix(categoryDelim) {
		return newParseError(nxt,
			"Expected a delimiter after category %s, got %s",
			p.curCategory.name, nxt.Content())
	}

	return nil
}

func parseCategory(line *Line) *Category {
	lineVal := line.Content()

	prio, lineVal := parsePriority(lineVal)

	return &Category{
		name:      lineVal,
		priority:  prio,
		DocCoords: DocCoords{line.LineNumber(), line.Offset(), line.Length()}}
}

func (p *Parser) handleMomentLine(line *Line) error {
	mom, err := p.parseFullMoment(line, line.TrimmedContent(), "")
	if err != nil {
		return err
	}
	p.todos.moments = append(p.todos.moments, mom)
	return nil
}

func (p *Parser) parseFullMoment(line *Line, lineVal string, indent string) (Moment, error) {
	mom, err := parseMoment(line, lineVal, "")
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

func (p *Parser) parseCommentsAndSubMoments(mom Moment, indent string) error {
	nextIndent := indent + indentChar
	for p.scanner.Scan() {
		line := p.scanner.Line()
		if line.HasPrefix(nextIndent) {
			p.handleSubLine(mom, line, line.Content()[len(nextIndent):], indent)
		} else if line.IsEmpty() && len(mom.GetComments()) > 0 {
			// special case: treat empty line between comments as a comment
			comment := &CommentLine{
				content:   "",
				DocCoords: DocCoords{line.LineNumber(), line.Offset(), 0}}
			mom.AddComment(comment)
		} else {
			p.scanner.Unscan()
			break
		}
	}

	// Remove trailing empty comments
	lc := mom.GetLastComment()
	for lc != nil && len(lc.content) == 0 {
		mom.RemoveLastComment()
		lc = mom.GetLastComment()
	}

	return nil
}

func (p *Parser) handleSubLine(mom Moment, line *Line, lineVal string, indent string) error {
	if strings.HasPrefix(lineVal, doneLBracket) {
		subMom, err := p.parseFullMoment(line, lineVal, indent+indentChar)
		if err != nil {
			return err
		}
		mom.AddSubMoment(subMom)
	} else {
		// Assume it's a comment
		comment := &CommentLine{
			content:   lineVal,
			DocCoords: DocCoords{line.LineNumber(), line.Offset() + len(indent+indentChar), len(lineVal)}}
		mom.AddComment(comment)
	}
	return nil
}
