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
	router.HandleFunc("/format", Format).Methods("POST")

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

	// fmt.Printf("hello from sibylgo\n")
	// tm := time.Now()
	// todos, err := ParseFile("todo.txt")
	// dur := time.Now().Sub(tm)
	// fmt.Printf("%dms\n", int(dur/time.Millisecond))
	// if err != nil {
	// 	panic(err)
	// }

	// for _, c := range todos.categories {
	// 	fmt.Printf("%s\n", c)
	// }
	// for _, m := range todos.moments {
	// 	printMom(m, "")
	// }
}

func Format(w http.ResponseWriter, r *http.Request) {
	reader := base64.NewDecoder(base64.StdEncoding, r.Body)
	todos, err := ParseReader(reader)
	if err != nil {
		// TODO
	}
	for _, m := range todos.moments {
		printMom(m, "")
	}
	res := FormatVSCode(todos)
	fmt.Printf(res)
	fmt.Fprintf(w, res)
}

func printMom(m Moment, indent string) {
	fmt.Printf("%s%s\n", indent, m)
	for _, s := range m.GetSubMoments() {
		printMom(s, indent+"  ")
	}
}
