package format

import (
	"fmt"
	"github.com/sandro-h/sibylgo/moment"
)

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
