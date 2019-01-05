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
	todos, _ := parse.ParseString(`
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
	todos, _ := parse.ParseString(`
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
