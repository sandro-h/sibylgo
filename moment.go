package main

import (
	t "time"
)

type Moment interface{}

type Todos struct {
	categories []*Category
	moments    []*Moment
}

type Category struct {
	name string
}

type SingleMoment struct {
	name       string
	start      t.Time
	end        t.Time
	done       bool
	category   *Category
	priority   int
	comments   []*CommentLine
	subMoments []*Moment
}

type CommentLine struct {
	content string
}
