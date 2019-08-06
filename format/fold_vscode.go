package format

import (
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
)

// FoldForVSCode returns a string of lines with the line ranges of all top-level moments,
// so they can be folded in Visual Studio Code.
func FoldForVSCode(todos *moment.Todos) string {
	res := ""
	for _, m := range todos.Moments {
		a := m.GetDocCoords().LineNumber
		b := m.GetBottomLineNumber()
		if b > a {
			res += fmt.Sprintf("%d-%d\n", a, b)
		}
	}
	return res
}
