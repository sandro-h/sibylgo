package outlook

import (
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateEventCommand(t *testing.T) {
	cases := [][]string{
		[]string{"[] bla (04.12.2020)"}, []string{"add", "-l", "sibyl", "-s", "bla", "-d", "2020-12-04"},
		[]string{"[] foo (04.12.2020 8:00)"}, []string{"add", "-l", "sibyl", "-s", "foo", "-d", "2020-12-04", "-t", "08:00", "-e", "09:00"},
	}

	for i := 0; i < len(cases); i += 2 {
		momString := cases[i][0]
		expectedArgs := cases[i+1]

		todos, _ := parse.String(momString)
		cmd := getCreateEventCommand(todos.Moments[0].(*moment.SingleMoment))
		args := cmd[1:]
		assert.Equal(t, expectedArgs, args, momString)
	}
}

func TestRemoveEventCommand(t *testing.T) {
	cases := [][]string{
		[]string{"[] bla (04.12.2020)"}, []string{"remove", "-s", "bla"},
		[]string{"[] foo (04.12.2020 8:00)"}, []string{"remove", "-s", "foo"},
	}

	for i := 0; i < len(cases); i += 2 {
		momString := cases[i][0]
		expectedArgs := cases[i+1]

		todos, _ := parse.String(momString)
		cmd := getRemoveEventCommand(todos.Moments[0].(*moment.SingleMoment))
		args := cmd[1:]
		assert.Equal(t, expectedArgs, args, momString)
	}
}

func TestParseListOutput(t *testing.T) {
	output := `
04.12.2020 00:00:00;04.13.2020 00:00:00;True;bla
04.12.2020 08:00:00;04.12.2020 09:00:00;False;foo
04.12.2020 00:00:00;04.13.2020 00:00:00;True;some ;subject ;with ;semicolons
`

	moms, err := parseListOutput(output)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(moms))
	assertMom(t, "bla", "04.12.2020", "", moms[0])
	assertMom(t, "foo", "04.12.2020", "08:00", moms[1])
	assertMom(t, "some ;subject ;with ;semicolons", "04.12.2020", "", moms[2])
}

func TestParseListOutputWithErrors(t *testing.T) {
	_, err := parseListOutput("invaliddate;invaliddate;True;bla")
	assert.NotNil(t, err)

	_, err = parseListOutput("04.12.2020 00:00:00;04.12.2020 00:00:00;invalidbool;bla")
	assert.NotNil(t, err)
}

func assertMom(t *testing.T, name string, start string, timeOfDay string, actual *moment.SingleMoment) {
	assert.Equal(t, name, actual.GetName())
	assert.Equal(t, tu.Dt(start), actual.Start.Time)

	tod := ""
	if actual.TimeOfDay != nil {
		tod = actual.TimeOfDay.Time.Format("15:04")

	}
	assert.Equal(t, timeOfDay, tod)
}
