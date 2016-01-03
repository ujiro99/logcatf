package main

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

const flagAll = -1

var (
	formatRegex  = regexp.MustCompile(`%(time)|%(pid)|%(tid)|%(priority)|%(tag)|%(message)`)
	sformatRegex = regexp.MustCompile(`%(t)|%(i)|%(I)|%(p)|%(a)|%(m)`)
	formatMap    = map[string]string{
		"time":     "t",
		"pid":      "i",
		"tid":      "I",
		"priority": "p",
		"tag":      "a",
		"message":  "m",
	}
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

// replace keyword to real value. use short format.
// ex) "%t %p" => "12-28 19:01:14.073 GLSUser"
func (item *LogcatItem) replaceKeyword(format string) string {
	matches := sformatRegex.FindAllStringSubmatch(format, len(formatMap))
	formatArgs := make([]interface{}, 0, len(matches))
	// find matched keyword and store value on item
	for _, match := range matches {
		for _, m := range match {
			key, ok := findKey(formatMap, m)
			if ok {
				formatArgs = append(formatArgs, (*item)[key])
				break
			}
		}
	}
	format = sformatRegex.ReplaceAllString(format, "%s")
	return fmt.Sprintf(format, formatArgs...)
}

func (item *LogcatItem) replaceEscape(format string) string {
	result := format
	result = strings.Replace(result, "\\t", "\t", flagAll)
	result = strings.Replace(result, "\\n", "\n", flagAll)
	return result
}

func verifyFormat(format string) error {
	// find unavailable keyword.
	removed := formatRegex.ReplaceAllString(format, "")
	removed = sformatRegex.ReplaceAllString(removed, "")
	keyRegexp := regexp.MustCompile(`%\w+`)
	matches := keyRegexp.FindAllString(removed, flagAll)

	if len(matches) == 0 {
		return nil // no probrem!
	}

	// return error message.
	err := ParameterError{}
	for _, match := range matches {
		errMsg := fmt.Sprintf(Message["msgUnavailableKeyword"], match)
		err.errors = append(err.errors, errMsg)
	}
	return &err
}

// ParameterError has error message of parameter.
type ParameterError struct {
	errors []string
}

// Error returns all error message.
func (e *ParameterError) Error() string {
	if e.errors != nil {
		return strings.Join(e.errors, "\n")
	}
	return ""
}

// convert long keys in format to short keys.
func normalizeFormat(format string) string {
	for long, short := range formatMap {
		tmp := strings.Replace(format, long, short, flagAll)
		format = tmp
	}
	return format
}

func findKey(m map[string]string, value string) (string, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	return "", false
}
