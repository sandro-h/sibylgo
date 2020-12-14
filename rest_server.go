package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sandro-h/sibylgo/backup"
	"github.com/sandro-h/sibylgo/calendar"
	"github.com/sandro-h/sibylgo/cleanup"
	"github.com/sandro-h/sibylgo/format"
	"github.com/sandro-h/sibylgo/instances"
	"github.com/sandro-h/sibylgo/modify"
	"github.com/sandro-h/sibylgo/moment"
	"github.com/sandro-h/sibylgo/parse"
	"github.com/sandro-h/sibylgo/reminder"
	"github.com/sandro-h/sibylgo/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func startRestServer(cfg *util.Config) {
	host := cfg.GetString("host", "localhost")
	port := cfg.GetInt("port", 8082)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router := mux.NewRouter()
	router.HandleFunc("/format", formatMoments).Methods("POST")
	router.HandleFunc("/folding", foldMoments).Methods("POST")
	router.HandleFunc("/clean", clean).Methods("POST")
	router.HandleFunc("/trash", trash).Methods("POST")
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

func formatMoments(w http.ResponseWriter, r *http.Request) {
	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
	todos, err := parse.Reader(reader)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	res := format.ForVSCode(todos)
	fmt.Fprint(w, res)
}

func foldMoments(w http.ResponseWriter, r *http.Request) {
	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
	todos, err := parse.Reader(reader)
	if err != nil {
		http.Error(w, err.Error(), 400)
	}
	res := format.FoldForVSCode(todos)
	fmt.Fprint(w, res)
}

func getCalendarEntries(w http.ResponseWriter, r *http.Request) {
	start, err := util.ParseISODate(r.FormValue("start"))
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	end, err := util.ParseISODate(r.FormValue("end"))
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
	date, err := util.ParseISODate(vars["date"])
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

func clean(w http.ResponseWriter, r *http.Request) {
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

func trash(w http.ResponseWriter, r *http.Request) {
	if todoFile == "" {
		log.Error("Cannot clean without todoFile set\n")
		return
	}

	trashFile := util.RemoveExtension(todoFile) + "-trash.txt"

	backup.Save(todoFile, "Backup before trashing")
	err := cleanup.MoveDoneToTrashFile(todoFile, trashFile, true)
	if err != nil {
		log.Errorf("Error trashing: %s", err)
	} else {
		log.Infof("Trashed: %s\n", todoFile)
		log.Infof("Moved done moments to: %s\n", trashFile)
	}
}

func setJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
