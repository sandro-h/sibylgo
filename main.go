package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sandro-h/sibylgo/calendar"
	"github.com/sandro-h/sibylgo/cleanup"
	"github.com/sandro-h/sibylgo/format"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/reminder"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var ascii = `.▄▄ · ▪  ▄▄▄▄·  ▄· ▄▌▄▄▌  
▐█ ▀. ██ ▐█ ▀█▪▐█▪██▌██•  
▄▀▀▀█▄▐█·▐█▀▀█▄▐█▌▐█▪██▪  
▐█▄▪▐█▐█▌██▄▪▐█ ▐█▀·.▐█▌▐▌
 ▀▀▀▀ ▀▀▀·▀▀▀▀   ▀ • .▀▀▀ `

var buildVersion = "0.0.0"
var buildNumber = "0"
var buildRevision = "-"

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
	fmt.Printf("%s\n", ascii)
	fmt.Printf("Version %s.%s (%s)\n", buildVersion, buildNumber, buildRevision)

	if *mailTo != "" {
		startMailReminders()
	}

	startRestServer()

	handleUserCommands()
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
	fmt.Print("Started mail reminders\n")
}

func startRestServer() {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := mux.NewRouter()
	router.HandleFunc("/format", formatMoments).Methods("POST")
	router.HandleFunc("/moments", getCalendarEntries).Methods("GET")
	router.HandleFunc("/reminders/{date}/weekly", getWeeklyReminders).Methods("GET")

	srv := &http.Server{
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(router),
		Addr:         fmt.Sprintf("%s:%d", "localhost", *port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go srv.ListenAndServe()
	fmt.Print("Started REST server\n")
}

func handleUserCommands() {
	printCommands()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Command> ")
	for scanner.Scan() {
		cmd := scanner.Text()
		switch cmd {
		case "help":
			printCommands()
		case "quit":
			return
		case "clean":
			clean()
		case "trash":
			trash()
		default:
			fmt.Printf("Unknown command\n")
		}
		fmt.Print("Command> ")
	}
}

func printCommands() {
	fmt.Print("help  - show this\n")
	fmt.Print("quit  - end app\n")
	fmt.Print("clean - move done moments to end of todo file\n")
	fmt.Print("trash - move done moments to trash file\n")
}

func clean() {
	assertStringFlagSet("todoFile", todoFile)

	err := cleanup.CleanupDoneFromFileToEnd(*todoFile, true)
	if err != nil {
		fmt.Printf("Error cleaning up: %s", err)
	} else {
		fmt.Printf("Moved done to end of: %s\n", *todoFile)
	}
}

func trash() {
	if *todoFile == "" {
		fmt.Print("No -todoFile defined\n")
		return
	}

	trashFile := removeExt(*todoFile) + "-trash.txt"

	err := cleanup.CleanupDoneFromFile(*todoFile, trashFile, true)
	if err != nil {
		fmt.Printf("Error trashing: %s", err)
	} else {
		fmt.Printf("Trashed: %s\n", *todoFile)
		fmt.Printf("Moved done moments to: %s\n", trashFile)
	}
}

func removeExt(s string) string {
	return s[:len(s)-len(filepath.Ext(s))]
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

func getCalendarEntries(w http.ResponseWriter, r *http.Request) {
	start, err := time.ParseInLocation("2006-01-02", r.FormValue("start"), time.Local)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	end, err := time.ParseInLocation("2006-01-02", r.FormValue("end"), time.Local)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	todos, err := parse.ParseFile(*todoFile)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	entries := calendar.CompileCalendarEntries(todos, start, end)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func getWeeklyReminders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date, err := time.ParseInLocation("2006-01-02", vars["date"], time.Local)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	todos, err := parse.ParseFile(*todoFile)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	todays, weeks := reminder.CompileRemindersForTodayAndThisWeek(todos, date)
	res := map[string][]*moment.MomentInstance{
		"today": todays,
		"week":  weeks}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
