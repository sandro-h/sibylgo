package main

import (
	"os"
	"strings"
	"unicode"
)

const categoryDelim = "------"
const doneLBracket = "["
const doneRBracket = ']'
const doneMark = 'x'
const doneMarkUpper = 'X'
const priorityMark = '!'

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

	parser := Parser{todos: &Todos{}, scanner: NewFileLineScanner(file)}
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
	moment, err := parseMoment(line)
	if err != nil {
		return err
	}
	moment.SetCategory(p.curCategory)
	p.todos.moments = append(p.todos.moments, moment)
	return nil
}

func parseMoment(line *Line) (Moment, error) {
	lineVal := line.TrimmedContent()

	mom, lineVal := parseBaseMoment(line, lineVal)

	done, lineVal, err := parseDoneMark(line, lineVal)
	if err != nil {
		return nil, err
	}
	mom.SetDone(done)

	prio, lineVal := parsePriority(lineVal)
	mom.SetPriority(prio)

	mom.SetName(lineVal)

	return mom, nil
}

func parseBaseMoment(line *Line, lineVal string) (Moment, string) {
	// TODO: check recurring
	return parseSingleMoment(line, lineVal)
}

func parseDoneMark(line *Line, lineVal string) (bool, string, error) {
	rBracketPos := 0
	done := false
	for i, c := range lineVal {
		if c == doneRBracket {
			rBracketPos = i
			break
		}
		if !unicode.IsSpace(c) && (c == doneMark || c == doneMarkUpper) {
			done = true
		}
	}
	if rBracketPos == 0 {
		return false, "", newParseError(line, "Expected closing %c for moment line %s", doneMark, line.Content())
	}
	return done, strings.TrimSpace(lineVal[rBracketPos+1:]), nil
}

func parsePriority(str string) (int, string) {
	prio := 0
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] != priorityMark {
			break
		}
		prio++
	}
	return prio, strings.TrimSpace(str[0 : len(str)-prio])
}
