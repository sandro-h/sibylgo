package testutil

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Dt parses a date string (02.01.2006) into a Time at 00:00 Local time.
func Dt(s string) time.Time {
	d, _ := time.ParseInLocation("02.01.2006", s, time.Local)
	return d
}

// DtUtc parses a date string (02.01.2006) into a Time at 00:00 UTC.
func DtUtc(s string) time.Time {
	d, _ := time.Parse("02.01.2006", s)
	return d
}

// Dtt parses a date and time string (02.01.2006 15:04) into a Time.
func Dtt(s string) time.Time {
	d, _ := time.ParseInLocation("02.01.2006 15:04", s, time.Local)
	return d
}

// Dts formats a Time into a date string (02.01.2006)
func Dts(t time.Time) string {
	return t.Format("02.01.2006")
}

// Tts formats a Time into a time string (15:04:05)
func Tts(t time.Time) string {
	return t.Format("15:04:05")
}

// Dtts formats a Time into a date and time string (02.01.2006 15:04)
func Dtts(t time.Time) string {
	return t.Format("02.01.2006 15:04:05")
}

// MakeTempDir creates a temporary directory with name dirName.
// Don't forget to call defer DeleteTempDir to get rid of it at the end of the test.
func MakeTempDir(dirName string) string {
	dir := filepath.Join(os.TempDir(), dirName)
	err := os.Mkdir(dir, 0755)
	if err != nil {
		panic(err)
	}
	return dir
}

// DeleteTempDir deteles the temporary directory and all its content.
func DeleteTempDir(dir string) {
	os.RemoveAll(dir)
}

// ToJSON pretty prints the passed object as JSON.
func ToJSON(any interface{}) string {
	bytes, err := json.MarshalIndent(any, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(bytes)
}
