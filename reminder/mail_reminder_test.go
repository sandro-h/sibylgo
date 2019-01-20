package reminder

import (
	"fmt"
	tu "github.com/sandro-h/sibylgo/testutil"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var testLastSentFile = filepath.Join(os.TempDir(), "mail_reminder_test_lastsent.txt")

func TestMain(m *testing.M) {
	startup()
	retCode := m.Run()
	os.Exit(retCode)
}

func startup() {
	getNow = func() time.Time { return tu.Dt("04.01.2019") }
	os.Remove(testLastSentFile)
}

func TestEmptyDailyReminder(t *testing.T) {
	defer os.Remove(testLastSentFile)
	todoFile := writeTodoFile("")
	var rcvTitle string
	var rcvContent string
	p := createTestReminderProcess(todoFile, &rcvTitle, &rcvContent)
	p.CheckOnce()

	assert.Equal(t, "TODOs for Friday, 4 Jan 2019", rcvTitle)
	assert.Equal(t, `<ul>
<li>None</li>
</ul>
`, rcvContent)
}

func TestDailyReminder(t *testing.T) {
	defer os.Remove(testLastSentFile)
	todoFile := writeTodoFile(`
[] foo (5.1.19)
[] bar (4.1.19)
[] zon
	[] ran (every friday)
[] other
`)
	var rcvTitle string
	var rcvContent string
	p := createTestReminderProcess(todoFile, &rcvTitle, &rcvContent)
	p.CheckOnce()

	assert.Equal(t, "TODOs for Friday, 4 Jan 2019", rcvTitle)
	assert.Equal(t, `<ul>
<li><b>bar</b></li>
<li>zon<ul>
<li><b>ran</b></li>
</ul>
</li>
</ul>
`, rcvContent)
}

func TestNoRepeatOnSameDay(t *testing.T) {
	defer os.Remove(testLastSentFile)
	todoFile := writeTodoFile("")
	var rcvTitle string
	var rcvContent string
	p := createTestReminderProcess(todoFile, &rcvTitle, &rcvContent)

	p.CheckOnce()
	rcvTitle = ""
	p.CheckOnce()

	assert.Equal(t, "", rcvTitle)
}

func TestRepeatOnNextDay(t *testing.T) {
	defer os.Remove(testLastSentFile)
	todoFile := writeTodoFile("")
	var rcvTitle string
	var rcvContent string
	p := createTestReminderProcess(todoFile, &rcvTitle, &rcvContent)

	p.CheckOnce()
	rcvTitle = ""
	rcvContent = ""
	getNow = func() time.Time { return tu.Dt("05.01.2019") }
	p.CheckOnce()

	assert.Equal(t, "TODOs for Saturday, 5 Jan 2019", rcvTitle)
}

func TestTimedReminder(t *testing.T) {
	defer os.Remove(testLastSentFile)
	getNow = func() time.Time { return tu.Dtt("05.01.2019 13:02") }
	setLastSentFileToToday()

	todoFile := writeTodoFile(`
[] foo (5.1.19 13:15)
`)
	var rcvTitle string
	var rcvContent string
	p := createTestReminderProcess(todoFile, &rcvTitle, &rcvContent)
	p.CheckOnce()

	assert.Equal(t, "Reminder for foo in 13min", rcvTitle)
	assert.Equal(t, "foo starts at 13:15", rcvContent)
}

func TestTimedReminderTooEarly(t *testing.T) {
	defer os.Remove(testLastSentFile)
	getNow = func() time.Time { return tu.Dtt("05.01.2019 12:59") }
	setLastSentFileToToday()

	todoFile := writeTodoFile(`
[] foo (5.1.19 13:15)
`)
	var rcvTitle string
	var rcvContent string
	p := createTestReminderProcess(todoFile, &rcvTitle, &rcvContent)
	p.CheckOnce()

	assert.Equal(t, "", rcvTitle)
	assert.Equal(t, "", rcvContent)
}

func TestTimedReminderTooLate(t *testing.T) {
	defer os.Remove(testLastSentFile)
	// 13:12 means with check interval of 5min, reminder would've already been sent
	// in a previous check
	getNow = func() time.Time { return tu.Dtt("05.01.2019 13:12") }
	setLastSentFileToToday()

	todoFile := writeTodoFile(`
[] foo (5.1.19 13:15)
`)
	var rcvTitle string
	var rcvContent string
	p := createTestReminderProcess(todoFile, &rcvTitle, &rcvContent)
	p.CheckOnce()

	assert.Equal(t, "", rcvTitle)
	assert.Equal(t, "", rcvContent)
}

func writeTodoFile(todos string) string {
	path := filepath.Join(os.TempDir(), "mail_reminder_test_todo.txt")
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	defer file.Close()
	fmt.Fprintf(file, todos)
	return path
}

func createTestReminderProcess(todoFile string, rcvTitle *string, rcvContent *string) *MailReminderProcess {
	p := NewMailReminderProcess(todoFile,
		func(title string, content string) error {
			*rcvTitle = title
			*rcvContent = content
			return nil
		})
	p.LastSentFile = testLastSentFile
	return p
}

func setLastSentFileToToday() {
	todoFile := writeTodoFile("")
	var rcvTitle string
	var rcvContent string
	p := createTestReminderProcess(todoFile, &rcvTitle, &rcvContent)
	p.CheckOnce()
}
