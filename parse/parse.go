package parse

import (
	"io"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/sandro-h/sibylgo/moment"
)

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
	parserState.todos.MomentsByID = make(map[string]moment.Moment)
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
	if line.HasPrefix(ParseConfig.GetCategoryDelim()) {
		err = p.handleCategoryLine(line)
	} else if line.HasPrefix(ParseConfig.GetLBracket()) {
		err = p.handleMomentLine(line)
	}
	//fmt.Printf("%s\n", line.content)
	return err
}

func (p *parserState) handleCategoryLine(line *Line) error {
	ok, catLine := p.scanner.ScanAndLine()
	if !ok {
		return newParseError(catLine, "Expected a category name after category delimiter")
	}

	p.curCategory = parseCategory(catLine)
	p.todos.Categories = append(p.todos.Categories, p.curCategory)

	ok, nxt := p.scanner.ScanAndLine()
	if !ok {
		return newParseError(catLine,
			"Expected a delimiter after category %s, but reached end",
			p.curCategory.Name)
	}
	if !nxt.HasPrefix(ParseConfig.GetCategoryDelim()) {
		return newParseError(nxt,
			"Expected a delimiter after category %s, got %s",
			p.curCategory.Name, nxt.Content())
	}

	return nil
}

func parseCategory(line *Line) *moment.Category {
	lineVal := line.Content()

	col, lineVal := parseCategoryColor(lineVal)
	prio, lineVal := parsePriority(lineVal)

	return &moment.Category{
		Name:      lineVal,
		Priority:  prio,
		Color:     col,
		DocCoords: moment.DocCoords{LineNumber: line.LineNumber(), Offset: line.Offset(), Length: line.Length()}}
}

func parseCategoryColor(lineVal string) (string, string) {
	if !strings.HasSuffix(lineVal, "]") {
		return "", lineVal
	}
	p := LastRuneIndex(lineVal, "[")
	colStr := lineVal[p+1 : len(lineVal)-1]
	return colStr, strings.TrimSpace(lineVal[:p])
}

func (p *parserState) handleMomentLine(line *Line) error {
	mom, err := p.parseFullMoment(line, line.TrimmedContent(), "")
	if err != nil {
		return err
	}
	p.todos.Moments = append(p.todos.Moments, mom)
	if mom.GetID() != nil {
		p.todos.MomentsByID[mom.GetID().Value] = mom
	}
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
	nextIndent := indent + ParseConfig.GetIndent()
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
	if strings.HasPrefix(lineVal, ParseConfig.GetLBracket()) {
		subMom, err := p.parseFullMoment(line, lineVal, indent+ParseConfig.GetIndent())
		if err != nil {
			return err
		}
		mom.AddSubMoment(subMom)
	} else {
		// Assume it's a comment
		comment := &moment.CommentLine{
			Content: lineVal,
			DocCoords: moment.DocCoords{LineNumber: line.LineNumber(),
				Offset: line.Offset() + len(indent+ParseConfig.GetIndent()),
				Length: utf8.RuneCountInString(lineVal)}}
		mom.AddComment(comment)
	}
	return nil
}
