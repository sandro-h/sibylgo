package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sandro-h/sibylgo/format"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/reminder"
	"net/http"
	"os"
	"time"
)

var port = flag.Int("port", 8082, "REST port. Default: 8082")
var mailHost = flag.String("mailHost", "", "STMP host for sending mail reminders.")
var mailPort = flag.Int("mailPort", 0, "STMP port for sending mail reminders.")
var mailUser = flag.String("mailUser", "", "User name for STMP auth for sending mail reminders.")
var mailPassword = flag.String("mailPassword", "", "Password for STMP auth for sending mail reminders.")
var mailFrom = flag.String("mailFrom", "", "E-mail address to use as sender for sending mail reminders.")
var mailTo = flag.String("mailTo", "", "E-mail address to which to send mail reminders.")
var todoFile = flag.String("todoFile", "", "Todo file to monitor for reminders.")

func main() {
	flag.Parse()

	if *mailTo != "" {
		startMailReminders()
	}

	router := mux.NewRouter()
	router.HandleFunc("/format", formatMoments).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%d", "localhost", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func startMailReminders() {
	assertStringFlagSet("todoFile", todoFile)
	assertStringFlagSet("mailHost", mailHost)
	assertIntFlagSet("mailPort", mailPort)
	assertStringFlagSet("mailFrom", mailFrom)
	assertStringFlagSet("mailTo", mailTo)

	host := reminder.MailHostProperties{Host: *mailHost, Port: *mailPort, User: *mailUser, Password: *mailPassword}
	p := reminder.NewMailReminderProcessForSMTP(*todoFile, host, *mailFrom, *mailTo)
	go p.CheckInfinitely()
}

func assertStringFlagSet(name string, s *string) {
	if *s == "" {
		failUnsetFlag(name)
	}
}

func assertIntFlagSet(name string, s *int) {
	if *s == 0 {
		failUnsetFlag(name)
	}
}

func failUnsetFlag(name string) {
	flag.Usage()
	fmt.Fprintf(os.Stderr, "%s must be set\n", name)
	os.Exit(1)
}

func formatMoments(w http.ResponseWriter, r *http.Request) {
	tm := time.Now()
	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
	todos, err := parse.ParseReader(reader)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	res := format.FormatVSCode(todos)
	fmt.Fprintf(w, res)
	fmt.Printf("response time: %dms\n", int(time.Now().Sub(tm)/time.Millisecond))
}
