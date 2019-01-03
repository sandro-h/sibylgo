package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/format", format).Methods("POST")

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

func format(w http.ResponseWriter, r *http.Request) {
	tm := time.Now()
	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
	todos, err := ParseReader(reader)
	if err != nil {
		// TODO
	}
	// for _, m := range todos.moments {
	// 	printMom(m, "")
	// }
	res := FormatVSCode(todos)
	//fmt.Printf(res)
	fmt.Fprintf(w, res)
	fmt.Printf("response time: %dms\n", int(time.Now().Sub(tm)/time.Millisecond))
}

func printMom(m Moment, indent string) {
	fmt.Printf("%s%s\n", indent, m)
	for _, s := range m.GetSubMoments() {
		printMom(s, indent+"  ")
	}
}
