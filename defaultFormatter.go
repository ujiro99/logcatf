package main

import (
	"fmt"
	"regexp"
	"strings"
)

const flagAll = -1

var (
	formatRegex  = regexp.MustCompile(`%.*?(time|pid|tid|priority|tag|message)`)
	sformatRegex = regexp.MustCompile(`%.*?(t|i|I|p|a|m)`)
	flagRegex    = regexp.MustCompile(`t|i|I|p|a|m`)
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
	// Format create a formatted string.
	Format(item *LogcatEntry) string
	// Verify checks a format
	Verify() error
	// Normarize convert long keys to short keys.
	Normarize()
}

// implements Formatter
type defaultFormatter struct {
	format *string
}

// Format implements Formatter
// replace %* keywords to real value. use short format.
//   ex) "%t %a" => "12-28 19:01:14.073 GLSUser"
func (f *defaultFormatter) Format(item *LogcatEntry) string {
	matches := sformatRegex.FindAllStringSubmatch(*f.format, len(formatMap))
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
	format := sformatRegex.ReplaceAllStringFunc(*f.format, func(str string) string {
		return flagRegex.ReplaceAllString(str, "s")
	})
	return fmt.Sprintf(format, formatArgs...)
}

// Verify implements Formatter
func (f *defaultFormatter) Verify() error {
	err := FormatError{}
	var msg []string
	msg = f.verifyUnabailavleKeyword(*f.format)
	if msg != nil {
		err.errors = append(err.errors, msg...)
	}
	msg = f.verifyDuplicatedKeyword(*f.format)
	if msg != nil {
		err.errors = append(err.errors, msg...)
	}
	if len(err.errors) == 0 {
		return nil
	}
	return &err
}

func (f *defaultFormatter) verifyUnabailavleKeyword(format string) []string {
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

func (f *defaultFormatter) verifyDuplicatedKeyword(format string) []string {
	// find unavailable keyword.
	duplicated := [][]string{}
	for k, v := range formatMap {
		count := strings.Count(format, "%"+v)
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

// FormatError has error message of parameter.
type FormatError struct {
	errors []string
}

// Error returns all error message.
func (e *FormatError) Error() string {
	if e.errors != nil {
		return strings.Join(e.errors, "\n")
	}
	return ""
}

// Normarize implements Formatter
func (f *defaultFormatter) Normarize() {
	newFormat := *f.format
	for long, short := range formatMap {
		newFormat = strings.Replace(newFormat, long, short, flagAll)
	}
	res := f.replaceEscape(newFormat)
	f.format = &res
}

func (f *defaultFormatter) replaceEscape(format string) string {
	format = strings.Replace(format, "\\t", "\t", flagAll)
	format = strings.Replace(format, "\\n", "\n", flagAll)
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
