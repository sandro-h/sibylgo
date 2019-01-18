package testutil

import (
	"time"
)

func Dt(s string) time.Time {
	d, _ := time.Parse("02.01.2006", s)
	return d
}

func Dts(t time.Time) string {
	return t.Format("02.01.2006")
}

func Tts(t time.Time) string {
	return t.Format("15:04:05")
}
