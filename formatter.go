package main

import (
	"fmt"
	"regexp"
	"strings"
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

// Formatter of logcatItem
type Formatter interface {
	Format(format string, item *LogcatItem) string
	Verify(format string) error
	// convert long keys in format to short keys.
	Normarize(format string) string
}

// implements Formatter
type logcatFormatter struct{}

// Format implements Formatter
func (f *logcatFormatter) Format(format string, item *LogcatItem) string {
	result := format
	result = f.replaceKeyword(result, item)
	result = f.replaceEscape(result)
	return result
}

// replace keyword to real value. use short format.
// ex) "%t %p" => "12-28 19:01:14.073 GLSUser"
func (f *logcatFormatter) replaceKeyword(format string, item *LogcatItem) string {
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

func (f *logcatFormatter) replaceEscape(format string) string {
	result := format
	result = strings.Replace(result, "\\t", "\t", flagAll)
	result = strings.Replace(result, "\\n", "\n", flagAll)
	return result
}

// Verify implements Formatter
func (f *logcatFormatter) Verify(format string) error {
	err := ParameterError{}
	var msg []string
	msg = f.verifyUnabailavleKeyword(format)
	if msg != nil {
		err.errors = append(err.errors, msg...)
	}
	msg = f.verifyDuplicatedKeyword(format)
	if msg != nil {
		err.errors = append(err.errors, msg...)
	}
	if len(err.errors) == 0 {
		return nil
	}
	return &err
}

func (f *logcatFormatter) verifyUnabailavleKeyword(format string) []string {
	// find unavailable keyword.
	removed := formatRegex.ReplaceAllString(format, "")
	removed = sformatRegex.ReplaceAllString(removed, "")
	keyRegexp := regexp.MustCompile(`%\w+`)
	matches := keyRegexp.FindAllString(removed, flagAll)

	if len(matches) == 0 {
		return nil // no probrem!
	}

	// return error message.
	msgs := make([]string, 0, len(matches))
	for _, match := range matches {
		errMsg := fmt.Sprintf(Message["msgUnavailableKeyword"], match)
		msgs = append(msgs, errMsg)
	}
	return msgs
}

func (f *logcatFormatter) verifyDuplicatedKeyword(format string) []string {
	// find unavailable keyword.
	str := f.Normarize(format)
	duplicated := [][]string{}
	for k, v := range formatMap {
		count := strings.Count(str, "%"+v)
		if count > 1 {
			duplicated = append(duplicated, []string{k, v})
		}
	}

	// is duplicated... ?
	if len(duplicated) == 0 {
		return nil // no probrem!
	}

	// return error message.
	msgs := make([]string, 0, len(duplicated))
	for _, d := range duplicated {
		errMsg := fmt.Sprintf(Message["msgDuplicatedKeyword"], "%"+d[0], "%"+d[1])
		msgs = append(msgs, errMsg)
	}
	return msgs
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

// Normarize implements Formatter
func (f *logcatFormatter) Normarize(format string) string {
	for long, short := range formatMap {
		format = strings.Replace(format, long, short, flagAll)
	}
	return format
}

// TODO find out generic api for find key.
func findKey(m map[string]string, value string) (string, bool) {
	for k, v := range m {
		if v == value {
			return k, true
		}
	}
	return "", false
}
