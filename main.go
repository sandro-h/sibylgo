package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sandro-h/sibylgo/backup"
	"github.com/sandro-h/sibylgo/calendar"
	"github.com/sandro-h/sibylgo/cleanup"
	"github.com/sandro-h/sibylgo/extsources"
	"github.com/sandro-h/sibylgo/format"
	"github.com/sandro-h/sibylgo/instances"
	"github.com/sandro-h/sibylgo/modify"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/outlook"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/reminder"
	"github.com/sandro-h/sibylgo/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

var showVersion = flag.Bool("version", false, "Show version")
var configFile = flag.String("config", "", "Path to config yml file. By default uses sibylgo.yml in same directory as this executable, if it exists.")
var todoFile string
var extSourcesProcess *extsources.ExternalSourcesProcess

func main() {
	flag.Parse()

	log.SetFormatter(&SimpleFormatter{})

	if *showVersion {
		log.Infof("%s.%s\n", buildVersion, buildNumber)
		return
	}

	fmt.Printf("%s\n", ascii)
	log.Infof("Version %s.%s (%s)\n", buildVersion, buildNumber, buildRevision)

	cfg := loadConfig()
	log.SetLevel(getConfigLogLevel(cfg))
	todoFile = cfg.GetString("todoFile", "")
	if todoFile != "" {
		log.Infof("Using todo file %s\n", todoFile)
		startDailyBackupProcess(todoFile)
	}

	if cfg.HasKey("mailTo") {
		startMailReminders(cfg)
	}

	if cfg.HasKey("external_sources") {
		extSrcConfig := cfg.GetSubConfig("external_sources")
		if todoFile == "" {
			panic("Cannot run external sources without todoFile set")
		}
		startExternalSources(todoFile, extSrcConfig)
	}

	if cfg.HasKey("outlook_events") {
		outlookConfig := cfg.GetSubConfig("outlook_events")
		if todoFile == "" {
			panic("Cannot run outlook events without todoFile set")
		}
		startOutlookEvents(todoFile, outlookConfig)
	}

	startRestServer(cfg)

	handleUserCommands()
}

func loadConfig() *util.Config {
	absoluteCfgFile := *configFile
	if absoluteCfgFile == "" {
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		absoluteCfgFile = filepath.Join(dir, "sibylgo.yml")
	} else {
		// We don't care if the default config file doesn't exist,
		// but if a user set a config file explicitly, we should inform them
		// if the file doesn't actually exist.
		if _, err := os.Stat(absoluteCfgFile); os.IsNotExist(err) {
			panic(fmt.Sprintf("Config file %s set with -config does not exist.\n", absoluteCfgFile))
		}
	}

	cfg := &util.Config{}
	log.Infof("%s\n", absoluteCfgFile)
	if _, err := os.Stat(absoluteCfgFile); !os.IsNotExist(err) {
		cfg, err = util.LoadConfig(absoluteCfgFile)
		if err != nil {
			panic(err)
		}
	}

	return cfg
}

func getConfigLogLevel(cfg *util.Config) log.Level {
	switch strings.ToLower(cfg.GetString("log_level", "info")) {
	case "debug":
		return log.DebugLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}

func startMailReminders(cfg *util.Config) {
	cfg.GetStringOrFail("todoFile")
	mailHost := cfg.GetStringOrFail("mailHost")
	mailPort := cfg.GetIntOrFail("mailPort")
	mailFrom := cfg.GetStringOrFail("mailFrom")
	mailTo := cfg.GetStringOrFail("mailTo")
	mailUser := cfg.GetString("mailUser", "")
	mailPassword := cfg.GetString("mailPassword", "")

	host := reminder.MailHostProperties{Host: mailHost, Port: mailPort, User: mailUser, Password: mailPassword}
	p := reminder.NewMailReminderProcessForSMTP(todoFile, host, mailFrom, mailTo)
	go p.CheckInfinitely()
	log.Info("Started mail reminders\n")
}

func startExternalSources(todoFile string, extSrcConfig *util.Config) {
	extSourcesProcess = extsources.NewExternalSourcesProcess(todoFile, extSrcConfig)
	go extSourcesProcess.CheckInfinitely()
	log.Info("Started external sources\n")
}

func startOutlookEvents(todoFile string, outlookConfig *util.Config) {
	if outlookConfig.GetBool("enabled", false) {
		go outlook.CheckInfinitely(todoFile, 5*time.Second)
		log.Info("Started outlook syncing\n")
	}
}

func startDailyBackupProcess(todoFile string) {
	dailyBackupFunc := func() {
		for {
			backup.CheckAndMakeDailyBackup(todoFile)
			time.Sleep(5 * time.Minute)
		}
	}
	go dailyBackupFunc()
	log.Info("Started daily backup\n")
}

func startRestServer(cfg *util.Config) {
	host := cfg.GetString("host", "localhost")
	port := cfg.GetInt("port", 8082)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := mux.NewRouter()
	router.HandleFunc("/format", formatMoments).Methods("POST")
	router.HandleFunc("/folding", foldMoments).Methods("POST")
	router.HandleFunc("/clean", wrapForRequest(clean)).Methods("POST")
	router.HandleFunc("/trash", wrapForRequest(trash)).Methods("POST")
	router.HandleFunc("/moments", getCalendarEntries).Methods("GET")
	router.HandleFunc("/moments", insertMoment).Methods("POST")
	router.HandleFunc("/reminders/{date}/weekly", getWeeklyReminders).Methods("GET")

	srv := &http.Server{
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(router),
		Addr:         fmt.Sprintf("%s:%d", host, port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	go srv.ListenAndServe()
	log.Infof("Started REST server on %s:%d\n", host, port)
}

func wrapForRequest(fn func()) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fn()
	}
}

func handleUserCommands() {
	log.Info()
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
		case "extsrc":
			runExtSources()
		default:
			log.Infof("Unknown command\n")
		}
		fmt.Print("Command> ")
	}
}

func printCommands() {
	fmt.Print("help   - show this\n")
	fmt.Print("quit   - end app\n")
	fmt.Print("clean  - move done moments to end of todo file\n")
	fmt.Print("trash  - move done moments to trash file\n")
	fmt.Print("extsrc - manually run external source collection\n")
}

func clean() {
	if todoFile == "" {
		log.Errorf("Cannot clean without todoFile set\n")
		return
	}

	backup.Save(todoFile, "Backup before cleaning")
	err := cleanup.MoveDoneToEndOfFile(todoFile, true)
	if err != nil {
		log.Infof("Error cleaning up: %s\n", err)
	} else {
		log.Infof("Moved done to end of: %s\n", todoFile)
	}
}

func trash() {
	if todoFile == "" {
		log.Error("Cannot clean without todoFile set\n")
		return
	}

	trashFile := removeExt(todoFile) + "-trash.txt"

	backup.Save(todoFile, "Backup before trashing")
	err := cleanup.MoveDoneToTrashFile(todoFile, trashFile, true)
	if err != nil {
		log.Errorf("Error trashing: %s", err)
	} else {
		log.Infof("Trashed: %s\n", todoFile)
		log.Infof("Moved done moments to: %s\n", trashFile)
	}
}

func runExtSources() {
	if extSourcesProcess == nil {
		log.Info("External sources not initialized. Are they configured?")
		return
	}
	extSourcesProcess.CheckOnce()
}

func removeExt(s string) string {
	return s[:len(s)-len(filepath.Ext(s))]
}

func formatMoments(w http.ResponseWriter, r *http.Request) {
	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
	todos, err := parse.Reader(reader)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	res := format.ForVSCode(todos)
	fmt.Fprintf(w, res)
}

func foldMoments(w http.ResponseWriter, r *http.Request) {
	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
	todos, err := parse.Reader(reader)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	res := format.FoldForVSCode(todos)
	fmt.Fprintf(w, res)
}

func getCalendarEntries(w http.ResponseWriter, r *http.Request) {
	start, err := parseDate(r.FormValue("start"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	end, err := parseDate(r.FormValue("end"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	todos, err := parse.File(todoFile)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	entries := calendar.CompileCalendarEntries(todos, start, end)
	setJSONContentType(w)
	json.NewEncoder(w).Encode(entries)
}

func insertMoment(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "name parameter not set", 400)
		return
	}
	category := r.FormValue("category")
	mom := moment.NewSingleMoment(name)
	if category != "" {
		mom.SetCategory(&moment.Category{Name: category})
	}

	log.Infof("Inserting '%s' into category '%s'\n", name, category)
	backup.Save(todoFile, "Backup before programmatically inserting moment")
	err := modify.PrependInFile(todoFile, []moment.Moment{mom})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.WriteHeader(http.StatusCreated)
	setJSONContentType(w)
	w.Write([]byte("{\"message\": \"Inserted moment\"}"))
}

func getWeeklyReminders(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	date, err := parseDate(vars["date"])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	todos, err := parse.File(todoFile)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	todays, weeks := reminder.CompileRemindersForTodayAndThisWeek(todos, date)
	res := map[string][]*instances.Instance{
		"today": todays,
		"week":  weeks}
	setJSONContentType(w)
	json.NewEncoder(w).Encode(res)
}

func setJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func parseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, time.Local)
}
