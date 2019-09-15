package calendar

import (
	"bytes"
	"encoding/json"
	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalendar(t *testing.T) {
	todos, _ := parse.String(`
[] foo (5.1.19)
[] bar (every wednesday)
[x] done (4.1.19)
`)
	entries := CompileCalendarEntries(todos, tu.Dt("31.12.2018"), tu.Dt("06.01.2019"))
	buf := bytes.NewBufferString("")
	json.NewEncoder(buf).Encode(entries)
	assert.Equal(t, `[{"title":"foo","start":"2019-01-05","end":"2019-01-06"},{"title":"bar","start":"2019-01-02","end":"2019-01-03"}]
`, buf.String())
}

func TestCalendarPriority(t *testing.T) {
	todos, _ := parse.String(`
[] foo (5.1.19)
[] bar! (every wednesday)
[x] done (4.1.19)
`)
	entries := CompileCalendarEntries(todos, tu.Dt("31.12.2018"), tu.Dt("06.01.2019"))
	buf := bytes.NewBufferString("")
	json.NewEncoder(buf).Encode(entries)
	assert.Equal(t, `[{"title":"bar","start":"2019-01-02","end":"2019-01-03"},{"title":"foo","start":"2019-01-05","end":"2019-01-06"}]
`, buf.String())
}

func TestCalendarColors(t *testing.T) {
	todos, _ := parse.String(`
[] foo (10.9.19)
------------------
a cat [green]
------------------
[] bar (13.9.19)
`)
	entries := CompileCalendarEntries(todos, tu.Dt("09.09.2019"), tu.Dt("15.09.2019"))
	buf := bytes.NewBufferString("")
	json.NewEncoder(buf).Encode(entries)
	assert.Equal(t, `[{"title":"foo","start":"2019-09-10","end":"2019-09-11"},{"title":"bar","start":"2019-09-13","end":"2019-09-14","color":"green"}]
`, buf.String())
}
