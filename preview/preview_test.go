package preview

import (
	"encoding/json"
	"testing"

	"github.com/sandro-h/sibylgo/instances"
	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
)

func TestCompileOverview(t *testing.T) {
	todos, _ := parse.File(tu.FullTestdataPath("overview.input"))

	overview := tu.ToJSON(compileTopLevelMomentsOverview(todos))

	tu.AssertGoldenOutput(t, "TestCompileOverview", "overview.output.json", overview)
}

func TestCompileReminders(t *testing.T) {
	todos, _ := parse.File(tu.FullTestdataPath("reminder_preview.input"))

	todays, _ := compileReminders(todos, tu.Dt("17.04.2021"))

	tu.AssertGoldenOutput(t, "TestCompileOverview", "reminder_preview.output.json", ToJSONWithNormalizedTime(todays))
}

func ToJSONWithNormalizedTime(insts []*instances.Instance) string {
	var normalized []*NormalizedTimeInst
	for _, t := range insts {
		normalized = append(normalized, (*NormalizedTimeInst)(t))
	}
	return tu.ToJSON(normalized)
}

type NormalizedTimeInst instances.Instance

func (i *NormalizedTimeInst) MarshalJSON() ([]byte, error) {
	type Alias NormalizedTimeInst
	return json.Marshal(&struct {
		*Alias
		Start string `json:"start"`
		End   string `json:"end"`
	}{
		Alias: (*Alias)(i),
		Start: i.Start.Format("2006-01-02 15:04:05"),
		End:   i.End.Format("2006-01-02 15:04:05"),
	})
}
