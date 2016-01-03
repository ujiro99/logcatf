package main

import (
	"fmt"
	"regexp"
	"time"
)

var (
	keys = []string{
		"time", "pid", "tid", "priority", "tag", "message",
	}
	formatRegexps  = regexp.MustCompile(`%(time)|%(pid)|%(tid)|%(priority)|%(tag)|%(message)`)
	tabRegexps     = regexp.MustCompile(`\\t`)
	newLineRegexps = regexp.MustCompile(`\\n`)
)

// LogcatItem represents a line of logcat log.
type LogcatItem map[string]string

// time return parsed time.
func (item *LogcatItem) time() time.Time {
	t, _ := time.Parse("12-28 18:54:08.043", (*item)["time"])
	return t
}

// Format returns formatted string.
func (item *LogcatItem) Format(format string) string {
	result := format
	result = item.replaceKeyword(result)
	result = item.replaceEscape(result)
	return result
}

func (item *LogcatItem) replaceKeyword(format string) string {
	matches := formatRegexps.FindAllStringSubmatch(format, len(keys))
	formatArgs := make([]interface{}, 0, len(matches))
	// find matched keyword
	for _, match := range matches {
		key := ""
		for j := 1; key == ""; j++ {
			key = match[j]
		}
		formatArgs = append(formatArgs, (*item)[key])
	}
	format = formatRegexps.ReplaceAllString(format, "%s")
	return fmt.Sprintf(format, formatArgs...)
}

func (item *LogcatItem) replaceEscape(format string) string {
	result := format
	result = tabRegexps.ReplaceAllString(result, "\t")
	result = newLineRegexps.ReplaceAllString(result, "\n")
	return result
}
