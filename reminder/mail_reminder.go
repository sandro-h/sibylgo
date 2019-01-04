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
	todoFilePath string
	sendMailFunc SendMailFunction
	LastSentFile string
}

func NewMailReminderProcess(todoFilePath string, sendMailFunc SendMailFunction) *MailReminderProcess {
	return &MailReminderProcess{todoFilePath, sendMailFunc,
		filepath.Join(os.TempDir(), defaultLastSentFile)}
}

func NewMailReminderProcessForSMTP(todoFilePath string, host MailHostProperties, from string, to string) *MailReminderProcess {
	return &MailReminderProcess{todoFilePath,
		func(subject string, body string) error {
			return sendMail(host, from, to, subject, body)
		},
		filepath.Join(os.TempDir(), defaultLastSentFile)}
}

func (p *MailReminderProcess) CheckInfinitely() {
	for {
		p.CheckOnce()
		time.Sleep(10 * time.Minute)
	}
}

func (p *MailReminderProcess) CheckOnce() {
	lastDaySent := p.loadLastDaySent()
	fmt.Printf("Last sent %s\n", lastDaySent)
	today := util.SetToStartOfDay(getNow())
	if today.After(lastDaySent) {
		fmt.Printf("Sending reminder for %s\n", today)
		err := p.sendReminderForToday(today)
		if err != nil {
			fmt.Printf("Could not send reminder: %s\n", err.Error())
			return
		}
		lastDaySent = today
		p.saveLastDaySent(today)
	}
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
	dt, err := time.Parse("2006-01-02", sc.Text())
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

func (p *MailReminderProcess) sendReminderForToday(today time.Time) error {
	todos, err := parse.ParseFile(p.todoFilePath)
	if err != nil {
		return err
	}
	subject := fmt.Sprintf("TODOs for %s", today.Format("Monday, 2 Jan 2006"))
	content := compileMomentsForToday(today, todos)
	return p.sendMailFunc(subject, content)
}

func compileMomentsForToday(today time.Time, todos *moment.Todos) string {
	insts := generate.GenerateInstancesFiltered(todos, today, util.SetToEndOfDay(today),
		func(mom moment.Moment) bool { return !mom.IsDone() })
	content := ""
	addMomentsEndingInRange(&content, insts)
	return content
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
