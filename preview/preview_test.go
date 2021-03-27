package preview

import (
	"testing"

	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
)

func TestCompileOverview(t *testing.T) {
	todos, _ := parse.File(tu.FullTestdataPath("overview.input"))

	overview := tu.ToJSON(compileTopLevelMomentsOverview(todos))

	tu.AssertGoldenOutput(t, "TestCompileOverview", "overview.output.json", overview)
}
