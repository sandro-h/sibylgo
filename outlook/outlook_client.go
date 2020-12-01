package outlook

import (
	"errors"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/util"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const outlookCliExe = "outlook_cli/outlook_cli.exe"

func createEvent(mom *moment.SingleMoment) error {
	if mom.Start == nil || mom.End == nil || util.SetToStartOfDay(mom.Start.Time) != util.SetToStartOfDay(mom.End.Time) {
		return errors.New("Only single, non-range moments are supported at the moment")
	}

	cmdAndArgs := getCreateEventCommand(mom)
	cmd := exec.Command(cmdAndArgs[0], cmdAndArgs[1:]...)
	err := cmd.Run()

	return err
}

func removeEvent(mom *moment.SingleMoment) error {
	cmdAndArgs := getRemoveEventCommand(mom)
	cmd := exec.Command(cmdAndArgs[0], cmdAndArgs[1:]...)
	err := cmd.Run()

	return err
}

func listEvents() ([]*moment.SingleMoment, error) {
	cmdAndArgs := getListEventsCommand()
	cmd := exec.Command(cmdAndArgs[0], cmdAndArgs[1:]...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	parsed, err := parseListOutput(string(output))
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func getCreateEventCommand(mom *moment.SingleMoment) []string {
	cmd := []string{
		outlookCliExe,
		"add",
		"-l", "sibyl",
		"-s", mom.GetName(),
		"-d", mom.Start.Time.Format("2006-01-02"),
	}

	if mom.TimeOfDay != nil {
		cmd = append(cmd,
			"-t", mom.TimeOfDay.Time.Format("15:04"),
			"-e", mom.TimeOfDay.Time.Add(1*time.Hour).Format("15:04"),
		)
	}
	return cmd
}

func getRemoveEventCommand(mom *moment.SingleMoment) []string {
	cmd := []string{
		outlookCliExe,
		"remove",
		"-s", mom.GetName(),
	}
	return cmd
}

func getListEventsCommand() []string {
	cmd := []string{
		outlookCliExe,
		"list",
		"-f", "[Location] = 'sibyl'",
	}
	return cmd
}

func parseListOutput(output string) ([]*moment.SingleMoment, error) {
	var res []*moment.SingleMoment
	lines := strings.Split(strings.ReplaceAll(output, "\r", ""), "\n")
	for _, l := range lines {
		parts := strings.SplitN(l, ";", 4)
		if len(parts) < 4 {
			continue
		}

		start, err := time.Parse("02.01.2006 15:04:05", parts[0])
		if err != nil {
			return nil, err
		}

		allDay, err := strconv.ParseBool(parts[2])
		if err != nil {
			return nil, err
		}

		name := parts[3]

		mom := moment.SingleMoment{}
		mom.SetName(name)
		mom.Start = &moment.Date{Time: util.SetToStartOfDay(start)}
		if !allDay {
			mom.TimeOfDay = &moment.Date{Time: start}
		}
		res = append(res, &mom)
	}

	return res, nil
}
