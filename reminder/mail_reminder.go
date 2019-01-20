package reminder

import (
	"bufio"
	"fmt"
	"github.com/sandro-h/sibylgo/generate"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/util"
	"gopkg.in/gomail.v2"
	"os"
	"path/filepath"
	"time"
)

var getNow = func() time.Time {
	return time.Now()
}

const defaultLastSentFile = "sibylgo_lastsent.txt"

type SendMailFunction func(string, string) error

type MailReminderProcess struct {
	todoFilePath  string
	sendMailFunc  SendMailFunction
	LastSentFile  string
	checkInterval time.Duration
	reminderTime  time.Duration
}

func NewMailReminderProcess(todoFilePath string, sendMailFunc SendMailFunction) *MailReminderProcess {
	return &MailReminderProcess{todoFilePath, sendMailFunc,
		filepath.Join(os.TempDir(), defaultLastSentFile),
		5 * time.Minute,
		15 * time.Minute}
}

func NewMailReminderProcessForSMTP(todoFilePath string, host MailHostProperties, from string, to string) *MailReminderProcess {
	return NewMailReminderProcess(todoFilePath,
		func(subject string, body string) error {
			return sendMail(host, from, to, subject, body)
		})
}

func (p *MailReminderProcess) CheckInfinitely() {
	for {
		p.CheckOnce()
		time.Sleep(p.checkInterval)
	}
}

func (p *MailReminderProcess) CheckOnce() {
	now := getNow()
	today := util.SetToStartOfDay(now)

	insts, err := p.loadTodaysMoments(today)
	if err != nil {
		fmt.Printf("Could not load moments for reminders: %s\n", err.Error())
		return
	}

	p.checkDailyReminder(today, insts)
	p.checkTimedReminders(now, insts)
}

func (p *MailReminderProcess) loadLastDaySent() time.Time {
	errDt := getNow().AddDate(0, 0, -1)
	path := p.LastSentFile
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return errDt
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	if !sc.Scan() || sc.Err() != nil {
		return errDt
	}
	dt, err := time.ParseInLocation("2006-01-02", sc.Text(), time.Local)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return errDt
	}
	return dt
}

func (p *MailReminderProcess) saveLastDaySent(dt time.Time) {
	path := p.LastSentFile
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, "%s\n", dt.Format("2006-01-02"))
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
}

func (p *MailReminderProcess) loadTodaysMoments(today time.Time) ([]*moment.MomentInstance, error) {
	todos, err := parse.ParseFile(p.todoFilePath)
	if err != nil {
		return nil, err
	}
	insts := generate.GenerateInstancesFiltered(todos, today, util.SetToEndOfDay(today),
		func(mom *moment.MomentInstance) bool { return !mom.Done })
	return insts, nil
}

func (p *MailReminderProcess) checkDailyReminder(today time.Time, insts []*moment.MomentInstance) {
	lastDaySent := p.loadLastDaySent()
	if today.After(lastDaySent) {
		fmt.Printf("Sending daily reminder for %s\n", today)
		err := p.sendDailyReminder(today, insts)
		if err != nil {
			fmt.Printf("Could not send reminder: %s\n", err.Error())
			return
		}
		lastDaySent = today
		p.saveLastDaySent(today)
	}
}

func (p *MailReminderProcess) sendDailyReminder(today time.Time, insts []*moment.MomentInstance) error {
	subject := fmt.Sprintf("TODOs for %s", today.Format("Monday, 2 Jan 2006"))
	content := ""
	addMomentsEndingInRange(&content, insts)
	return p.sendMailFunc(subject, content)
}

func addMomentsEndingInRange(content *string, insts []*moment.MomentInstance) {
	found := false
	*content += "<ul>\n"
	for _, m := range insts {
		subsEnding := hasSubsEndingInRange(m)
		if m.EndsInRange || subsEnding {
			found = true
			*content += "<li>"
			if m.EndsInRange {
				*content += "<b>"
			}
			*content += m.Name
			if m.EndsInRange {
				*content += "</b>"
			}
			if subsEnding {
				addMomentsEndingInRange(content, m.SubInstances)
			}
			*content += "</li>\n"
		}
	}
	if !found {
		*content += "<li>None</li>\n"
	}
	*content += "</ul>\n"
}

func hasSubsEndingInRange(m *moment.MomentInstance) bool {
	for _, s := range m.SubInstances {
		if s.EndsInRange || hasSubsEndingInRange(s) {
			return true
		}
	}
	return false
}

func (p *MailReminderProcess) checkTimedReminders(now time.Time, insts []*moment.MomentInstance) {
	upcoming := p.findUpcomingTimedMoments(now, p.reminderTime, p.checkInterval, insts)
	for _, m := range upcoming {
		fmt.Print(m.Delta)
		subject := fmt.Sprintf("Reminder for %s in %.0fmin", m.Name, m.Delta.Minutes())
		content := fmt.Sprintf("%s starts at %s", m.Name, m.TimeOfDay.Format("15:04"))
		p.sendMailFunc(subject, content)
	}
}

type Upcoming struct {
	Name      string
	TimeOfDay time.Time
	Delta     time.Duration
}

func (p *MailReminderProcess) findUpcomingTimedMoments(now time.Time, dur time.Duration,
	checkInterval time.Duration, insts []*moment.MomentInstance) []Upcoming {
	var res []Upcoming
	for _, i := range insts {
		if i.TimeOfDay != nil {
			delta := i.TimeOfDay.Sub(now)
			if delta <= dur && delta+checkInterval > dur {
				res = append(res, Upcoming{i.Name, *i.TimeOfDay, delta})
			}
		}
		res = append(res, p.findUpcomingTimedMoments(now, dur, checkInterval, i.SubInstances)...)
	}
	return res
}

type MailHostProperties struct {
	Host     string
	Port     int
	User     string
	Password string
}

func sendMail(host MailHostProperties, from string, to string, subject string, body string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewPlainDialer(host.Host, host.Port, host.User, host.Password)

	return d.DialAndSend(m)
}
