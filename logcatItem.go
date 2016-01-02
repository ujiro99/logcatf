package main

import (
	"fmt"
	"regexp"
	"time"
)

var (
	keys = []string{
		"Time", "Pid", "Tid", "Priority", "Tag", "Message",
	}
	formatRegexps = regexp.MustCompile(`%(Time)|%(Pid)|%(Tid)|%(Priority)|%(Tag)|%(Message)`)
)

// LogcatItem represents a line of logcat log.
type LogcatItem map[string]string

// Time return parsed time.
func (item *LogcatItem) Time() time.Time {
	t, _ := time.Parse("12-28 18:54:08.043", (*item)["Time"])
	return t
}

// Format returns formatted string.
func (item *LogcatItem) Format(format string) string {
	matches := formatRegexps.FindAllStringSubmatch(format, len(keys))
	formatArgs := make([]interface{}, 0, len(matches))
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
