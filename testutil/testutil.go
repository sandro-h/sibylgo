package testutil

import (
	"time"
)

func Dt(s string) time.Time {
	d, _ := time.ParseInLocation("02.01.2006", s, time.Local)
	return d
}

func Dtt(s string) time.Time {
	d, _ := time.ParseInLocation("02.01.2006 15:04", s, time.Local)
	return d
}

func Dts(t time.Time) string {
	return t.Format("02.01.2006")
}

func Tts(t time.Time) string {
	return t.Format("15:04:05")
}
