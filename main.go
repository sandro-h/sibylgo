package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Printf("hello from sibylgo\n")
	tm := time.Now()
	todos, err := ParseFile("todo.txt")
	dur := time.Now().Sub(tm)
	fmt.Printf("%dms\n", int(dur/time.Millisecond))
	if err != nil {
		panic(err)
	}

	for _, c := range todos.categories {
		fmt.Printf("%s\n", c)
	}
	for _, m := range todos.moments {
		fmt.Printf("%s\n", m)
	}
}
