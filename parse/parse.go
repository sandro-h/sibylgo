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
		parserState.handleLine(parserState.scanner.Line())
	}

	if err := parserState.scanner.Err(); err != nil {
		return nil, err
	}

	return parserState.todos, nil
}

func (p *parserState) handleLine(line *Line) {
	if line.IsEmpty() {
		return
	}
	if line.HasPrefix(ParseConfig.GetCategoryDelim()) {
		p.handleCategoryLine(line)
	} else if line.HasRunePrefix(ParseConfig.GetLBracket()) {
		p.handleMomentLine(line)
	}
}

func (p *parserState) handleCategoryLine(line *Line) {
	ok, catLine := p.scanner.ScanAndLine()
	if !ok {
		// Expected a category name after category delimiter
		return
	}

	// Consume closing delimiter after category line
	ok, nextLine := p.scanner.ScanAndLine()
	if !ok || !nextLine.HasPrefix(ParseConfig.GetCategoryDelim()) {
		return
	}

	p.curCategory = parseCategory(catLine)
	p.todos.Categories = append(p.todos.Categories, p.curCategory)
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

func (p *parserState) handleMomentLine(line *Line) {
	mom := p.parseFullMoment(line, line.TrimmedContent(), 0)
	if mom == nil {
		return
	}
	p.todos.Moments = append(p.todos.Moments, mom)
	if mom.GetID() != nil {
		p.todos.MomentsByID[mom.GetID().Value] = mom
	}
}

func (p *parserState) parseFullMoment(line *Line, lineVal string, indent int) moment.Moment {
	mom := parseMoment(line, lineVal)
	if mom == nil {
		return nil
	}
	mom.SetCategory(p.curCategory)

	p.parseCommentsAndSubMoments(mom, indent)

	return mom
}

func (p *parserState) parseCommentsAndSubMoments(mom moment.Moment, indent int) {
	nextIndent := indent + ParseConfig.GetTabSize()
	for p.scanner.Scan() {
		line := p.scanner.Line()
		lineIndent, indentCharCnt := countIndent(line.content, ParseConfig.GetTabSize(), nextIndent)
		if lineIndent >= nextIndent {
			p.handleSubLine(mom, line, line.Content()[indentCharCnt:], indent)
		} else if line.IsEmpty() {
			if len(mom.GetComments()) > 0 {
				// special case: treat empty line between comments as a comment
				comment := &moment.CommentLine{
					Content:   "",
					DocCoords: moment.DocCoords{LineNumber: line.LineNumber(), Offset: line.Offset(), Length: 0}}
				mom.AddComment(comment)
			}
			// Otherwise just ignore the empty line
		} else {
			// "Unconsume" the line since it is probably meant for a parent moment up the recursion stack
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
}

func (p *parserState) handleSubLine(mom moment.Moment, line *Line, lineVal string, indent int) {
	if HasRunePrefix(lineVal, ParseConfig.GetLBracket()) {
		subMom := p.parseFullMoment(line, lineVal, indent+ParseConfig.GetTabSize())
		if subMom != nil {
			mom.AddSubMoment(subMom)
			return
		}
	}

	// Assume it's a comment
	_, indentCharCnt := countIndent(line.content, ParseConfig.GetTabSize(), indent+ParseConfig.GetTabSize())
	comment := &moment.CommentLine{
		Content: lineVal,
		DocCoords: moment.DocCoords{LineNumber: line.LineNumber(),
			Offset: line.Offset() + indentCharCnt,
			Length: utf8.RuneCountInString(lineVal)}}
	mom.AddComment(comment)
}
