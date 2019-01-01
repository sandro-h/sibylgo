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
		printMom(m, "")
	}
}

func printMom(m Moment, indent string) {
	fmt.Printf("%s%s\n", indent, m)
	for _, s := range m.GetSubMoments() {
		printMom(s, indent+"  ")
	}
}
