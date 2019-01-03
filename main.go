package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sandro-h/sibylgo/format"
	"github.com/sandro-h/sibylgo/parse"
	"net/http"
	"time"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/format", formatMoments).Methods("POST")

	srv := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%s:%d", "localhost", 8082),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func formatMoments(w http.ResponseWriter, r *http.Request) {
	tm := time.Now()
	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
	todos, err := parse.ParseReader(reader)
	if err != nil {
		// TODO
	}
	res := format.FormatVSCode(todos)
	fmt.Fprintf(w, res)
	fmt.Printf("response time: %dms\n", int(time.Now().Sub(tm)/time.Millisecond))
}
