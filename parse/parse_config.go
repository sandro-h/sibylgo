package parse

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/sandro-h/sibylgo/util"
)

const defaultCategoryDelim = "------"

const defaultIndent = "\t"

const defaultLBracket = "["

const defaultRBracket = "]"

const defaultPriorityMark = "!"

const defaultInProgressMark = "p"

const defaultWaitingMark = "w"

const defaultDoneMark = "x"

var defaultDateFormats []string = []string{"02.01.06", "02.01.2006", "2.1.06", "2.1.2006"}

const defaultTimeFormat = "15:04"

var defaultWeekDays []string = []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

const defaultDailyPattern = "(?i)(every day|today)"

const defaultWeeklyPattern = "(?i)every (monday|tuesday|wednesday|thursday|friday|saturday|sunday)"

const defaultNthWeeklyPattern = "(?i)every (2nd|3rd|4th) (monday|tuesday|wednesday|thursday|friday|saturday|sunday)"

var defaultNths []string = []string{"2nd", "3rd", "4th"}

const defaultMonthlyPattern = `(?i)every (\d{1,2})\.?$`

const defaultYearlyPattern = `(?i)every (\d{1,2})\.(\d{1,2})\.?$`

type parseConfig struct {
	categoryDelim    string
	indent           string
	lBracket         string
	rBracket         *rune
	priorityMark     *rune
	inProgressMark   *rune
	waitingMark      *rune
	doneMark         *rune
	dateFormats      []string
	timeFormat       string
	weekDays         map[string]time.Weekday
	dailyPattern     *regexp.Regexp
	weeklyPattern    *regexp.Regexp
	nthWeeklyPattern *regexp.Regexp
	nths             map[string]int
	monthlyPattern   *regexp.Regexp
	yearlyPattern    *regexp.Regexp
	BackingCfg       *util.Config
}

// ParseConfig defines how moments are parsed.
var ParseConfig *parseConfig = &parseConfig{
	BackingCfg: &util.Config{},
}

// ResetConfig resets the configuration, clearing all cached values.
func ResetConfig() {
	ParseConfig = &parseConfig{
		BackingCfg: ParseConfig.BackingCfg,
	}
}

func (c *parseConfig) GetCategoryDelim() string {
	if c.categoryDelim == "" {
		c.categoryDelim = c.BackingCfg.GetString("category_delim", defaultCategoryDelim)
	}

	return c.categoryDelim
}

func (c *parseConfig) SetCategoryDelim(delim string) {
	c.categoryDelim = delim
}

func (c *parseConfig) GetIndent() string {
	if c.indent == "" {
		c.indent = c.BackingCfg.GetString("indent", defaultIndent)
	}

	return c.indent
}

func (c *parseConfig) SetIndent(indent string) {
	c.indent = indent
}

func (c *parseConfig) GetLBracket() string {
	if c.lBracket == "" {
		c.lBracket = c.BackingCfg.GetString("lbracket", defaultLBracket)
	}

	return c.lBracket
}

func (c *parseConfig) SetLBracket(bracket string) {
	c.lBracket = bracket
}

func (c *parseConfig) GetRBracket() rune {
	if c.rBracket == nil {
		r, _ := utf8.DecodeRuneInString(c.BackingCfg.GetString("rbracket", defaultRBracket))
		c.rBracket = &r
	}

	return *c.rBracket
}

func (c *parseConfig) SetRBracket(bracket rune) {
	c.rBracket = &bracket
}

func (c *parseConfig) GetPriorityMark() byte {
	if c.priorityMark == nil {
		r, _ := utf8.DecodeRuneInString(c.BackingCfg.GetString("priority_mark", defaultPriorityMark))
		c.priorityMark = &r
	}

	return byte(*c.priorityMark)
}

func (c *parseConfig) SetPriorityMark(mark rune) {
	c.priorityMark = &mark
}

func (c *parseConfig) GetInProgressMark() rune {
	if c.inProgressMark == nil {
		r, _ := utf8.DecodeRuneInString(c.BackingCfg.GetString("inprogress_mark", defaultInProgressMark))
		c.inProgressMark = &r
	}

	return *c.inProgressMark
}

func (c *parseConfig) SetInProgressMark(mark rune) {
	c.inProgressMark = &mark
}

func (c *parseConfig) GetWaitingMark() rune {
	if c.waitingMark == nil {
		r, _ := utf8.DecodeRuneInString(c.BackingCfg.GetString("waiting_mark", defaultWaitingMark))
		c.waitingMark = &r
	}

	return *c.waitingMark
}

func (c *parseConfig) SetWaitingMark(mark rune) {
	c.waitingMark = &mark
}

func (c *parseConfig) GetDoneMark() rune {
	if c.doneMark == nil {
		r, _ := utf8.DecodeRuneInString(c.BackingCfg.GetString("done_mark", defaultDoneMark))
		c.doneMark = &r
	}

	return *c.doneMark
}

func (c *parseConfig) SetDoneMark(mark rune) {
	c.doneMark = &mark
}

func (c *parseConfig) GetDateFormats() []string {
	if c.dateFormats == nil {
		c.dateFormats = c.BackingCfg.GetStringList("date_formats", defaultDateFormats)
	}

	return c.dateFormats
}

func (c *parseConfig) SetDateFormats(dateFormats []string) {
	c.dateFormats = dateFormats
}

func (c *parseConfig) GetTimeFormat() string {
	if c.timeFormat == "" {
		c.timeFormat = c.BackingCfg.GetString("time_format", defaultTimeFormat)
	}

	return c.timeFormat
}

func (c *parseConfig) SetTimeFormat(timeFormat string) {
	c.timeFormat = timeFormat
}

func (c *parseConfig) GetWeekDays() map[string]time.Weekday {
	if c.weekDays == nil {
		weekDayList := c.BackingCfg.GetStringList("week_days", defaultWeekDays)
		c.SetWeekDaysFromList(weekDayList)
	}

	return c.weekDays
}

// SetWeekDaysFromList sets the week days. Must start with Sunday!
func (c *parseConfig) SetWeekDaysFromList(weekDayList []string) {
	c.weekDays = make(map[string]time.Weekday)
	for i, d := range weekDayList {
		c.weekDays[strings.ToLower(d)] = time.Weekday(i)
	}
}

func (c *parseConfig) GetDailyPattern() *regexp.Regexp {
	if c.dailyPattern == nil {
		patternStr := c.BackingCfg.GetString("daily_pattern", defaultDailyPattern)
		c.SetDailyPattern(patternStr)
	}

	return c.dailyPattern
}

func (c *parseConfig) SetDailyPattern(patternStr string) {
	c.dailyPattern = parsePattern(patternStr)
}

func (c *parseConfig) GetWeeklyPattern() *regexp.Regexp {
	if c.weeklyPattern == nil {
		patternStr := c.BackingCfg.GetString("weekly_pattern", defaultWeeklyPattern)
		c.SetWeeklyPattern(patternStr)
	}

	return c.weeklyPattern
}

func (c *parseConfig) SetWeeklyPattern(patternStr string) {
	c.weeklyPattern = parsePattern(patternStr)
}

func (c *parseConfig) GetNthWeeklyPattern() *regexp.Regexp {
	if c.nthWeeklyPattern == nil {
		patternStr := c.BackingCfg.GetString("nth_weekly_pattern", defaultNthWeeklyPattern)
		c.SetNthWeeklyPattern(patternStr)
	}

	return c.nthWeeklyPattern
}

func (c *parseConfig) SetNthWeeklyPattern(patternStr string) {
	c.nthWeeklyPattern = parsePattern(patternStr)
}

func (c *parseConfig) GetNths() map[string]int {
	if c.nths == nil {
		nths := c.BackingCfg.GetStringList("nths", defaultNths)
		c.SetNthsFromList(nths)
	}

	return c.nths
}

// SetWeekDaysFromList sets the week days. Must start with Sunday!
func (c *parseConfig) SetNthsFromList(nths []string) {
	c.nths = make(map[string]int)
	for i, nth := range nths {
		c.nths[strings.ToLower(nth)] = 2 + i
	}
}

func (c *parseConfig) GetMonthlyPattern() *regexp.Regexp {
	if c.monthlyPattern == nil {
		patternStr := c.BackingCfg.GetString("monthly_pattern", defaultMonthlyPattern)
		c.SetMonthlyPattern(patternStr)
	}

	return c.monthlyPattern
}

func (c *parseConfig) SetMonthlyPattern(patternStr string) {
	c.monthlyPattern = parsePattern(patternStr)
}

func (c *parseConfig) GetYearlyPattern() *regexp.Regexp {
	if c.yearlyPattern == nil {
		patternStr := c.BackingCfg.GetString("yearly_pattern", defaultYearlyPattern)
		c.SetYearlyPattern(patternStr)
	}

	return c.yearlyPattern
}

func (c *parseConfig) SetYearlyPattern(patternStr string) {
	c.yearlyPattern = parsePattern(patternStr)
}

func parsePattern(patternStr string) *regexp.Regexp {
	if !strings.HasPrefix(patternStr, "(?i)") {
		patternStr = "(?i)" + patternStr
	}
	pattern, err := regexp.Compile(patternStr)
	if err != nil {
		panic(err)
	}
	return pattern
}
