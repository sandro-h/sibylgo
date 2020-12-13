package reminder

import (
	"bufio"
	"fmt"
	"github.com/sandro-h/sibylgo/instances"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/util"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"os"
	"path/filepath"
	"time"
)

var getNow = func() time.Time {
	return time.Now()
}

const defaultLastSentFile = "sibylgo_lastsent.txt"

// SendMailFunction takes an e-mail subject and content and sends
// it to a predefined recipient.
type SendMailFunction func(string, string) error

// MailReminderProcess checks the moments in the given todo file
// and sends a reminder mail for those moments that are due on the current day
// or due within a couple of minutes (if they have a TimeOfDay set).
type MailReminderProcess struct {
	todoFilePath  string
	sendMailFunc  SendMailFunction
	LastSentFile  string
	checkInterval time.Duration
	reminderTime  time.Duration
}

// NewMailReminderProcess creates a MailReminderProcess that uses the given sendMailFunc to send the
// reminder mails.
func NewMailReminderProcess(todoFilePath string, sendMailFunc SendMailFunction) *MailReminderProcess {
	return &MailReminderProcess{todoFilePath, sendMailFunc,
		filepath.Join(os.TempDir(), defaultLastSentFile),
		5 * time.Minute,
		15 * time.Minute}
}

// NewMailReminderProcessForSMTP creates a MailReminderProjcess that uses SMTP to send reminder mails to the given
// recipient.
func NewMailReminderProcessForSMTP(todoFilePath string, host MailHostProperties, from string, to string) *MailReminderProcess {
	return NewMailReminderProcess(todoFilePath,
		func(subject string, body string) error {
			return sendMail(host, from, to, subject, body)
		})
}

// CheckInfinitely repeatedly checks for reminders in the check interval.
// This method blocks indefinitely and should be run as a go routine.
func (p *MailReminderProcess) CheckInfinitely() {
	for {
		p.CheckOnce()
		time.Sleep(p.checkInterval)
	}
}

// CheckOnce does a single check for reminders and sends them if found.
func (p *MailReminderProcess) CheckOnce() {
	now := getNow()
	today := util.SetToStartOfDay(now)

	insts, err := p.loadTodaysMoments(today)
	if err != nil {
		log.Errorf("Could not load moments for reminders: %s\n", err.Error())
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
		log.Errorf("%s\n", err.Error())
		return errDt
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	if !sc.Scan() || sc.Err() != nil {
		return errDt
	}
	dt, err := time.ParseInLocation("2006-01-02", sc.Text(), time.Local)
	if err != nil {
		log.Errorf("%s\n", err.Error())
		return errDt
	}
	return dt
}

func (p *MailReminderProcess) saveLastDaySent(dt time.Time) {
	path := p.LastSentFile
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Errorf("%s\n", err.Error())
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, "%s\n", dt.Format("2006-01-02"))
	if err != nil {
		log.Errorf("%s\n", err.Error())
	}
}

func (p *MailReminderProcess) loadTodaysMoments(today time.Time) ([]*instances.Instance, error) {
	todos, err := parse.File(p.todoFilePath)
	if err != nil {
		return nil, err
	}
	insts := instances.GenerateFiltered(todos, today, util.SetToEndOfDay(today),
		func(mom *instances.Instance) bool { return !mom.Done })
	return insts, nil
}

func (p *MailReminderProcess) checkDailyReminder(today time.Time, insts []*instances.Instance) {
	lastDaySent := p.loadLastDaySent()
	if today.After(lastDaySent) {
		log.Infof("Sending daily reminder for %s\n", today)
		err := p.sendDailyReminder(today, insts)
		if err != nil {
			log.Errorf("Could not send reminder: %s\n", err.Error())
			return
		}
		p.saveLastDaySent(today)
	}
}

func (p *MailReminderProcess) sendDailyReminder(today time.Time, insts []*instances.Instance) error {
	subject := fmt.Sprintf("TODOs for %s", today.Format("Monday, 2 Jan 2006"))
	content := ""
	ending := FilterMomentsEndingInRange(insts)
	addMomentHTML(&content, ending)
	return p.sendMailFunc(subject, content)
}

func addMomentHTML(content *string, insts []*instances.Instance) {
	*content += "<ul>\n"
	for _, m := range insts {
		*content += "<li>"
		if m.EndsInRange {
			*content += "<b>"
		}
		*content += m.Name
		if m.EndsInRange {
			*content += "</b>"
		}
		if len(m.SubInstances) > 0 {
			addMomentHTML(content, m.SubInstances)
		}
		*content += "</li>\n"
	}
	if len(insts) == 0 {
		*content += "<li>None</li>\n"
	}
	*content += "</ul>\n"
}

func hasSubsEndingInRange(m *instances.Instance) bool {
	for _, s := range m.SubInstances {
		if s.EndsInRange || hasSubsEndingInRange(s) {
			return true
		}
	}
	return false
}

func (p *MailReminderProcess) checkTimedReminders(now time.Time, insts []*instances.Instance) {
	upcoming := p.findUpcomingTimedMoments(now, p.reminderTime, p.checkInterval, insts)
	for _, m := range upcoming {
		subject := fmt.Sprintf("Reminder for %s in %.0fmin", m.Name, m.Delta.Minutes())
		content := fmt.Sprintf("%s starts at %s", m.Name, m.TimeOfDay.Format("15:04"))
		p.sendMailFunc(subject, content)
	}
}

type upcoming struct {
	Name      string
	TimeOfDay time.Time
	Delta     time.Duration
}

func (p *MailReminderProcess) findUpcomingTimedMoments(now time.Time, dur time.Duration,
	checkInterval time.Duration, insts []*instances.Instance) []upcoming {
	var res []upcoming
	for _, i := range insts {
		if i.TimeOfDay != nil {
			delta := i.TimeOfDay.Sub(now)
			if delta <= dur && delta+checkInterval > dur {
				res = append(res, upcoming{i.Name, *i.TimeOfDay, delta})
			}
		}
		res = append(res, p.findUpcomingTimedMoments(now, dur, checkInterval, i.SubInstances)...)
	}
	return res
}

// MailHostProperties defines the mail host used for SMTP.
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
