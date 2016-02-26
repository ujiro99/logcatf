package main

import (
	"errors"
	"regexp"
)

// Parser parse logcat to struct.
type Parser interface {
	Parse(line string) (LogcatEntry, error)
}

// logcatParser implements Parser
type logcatParser struct{}

var (
	// logcat formats.
	// Parser find out a format according to this order.
	formats = []string{
		"brief",
		"process",
		"tag",
		"time",
		"threadtime",
		"ddms_save",
		"raw",
	}

	// key list for analyze logcat.
	paramOrderMap = map[string][]string{
		"brief":      []string{"priority", "tag", "pid", "message"},
		"process":    []string{"priority", "pid", "message", "tag"},
		"tag":        []string{"priority", "tag", "message"},
		"time":       []string{"time", "priority", "tag", "pid", "message"},
		"threadtime": []string{"time", "pid", "tid", "priority", "tag", "message"},
		"ddms_save":  []string{"time", "priority", "tag", "pid", "message"},
		"raw":        []string{"message"},
	}

	// regex pattrens for analyze logcat.
	patternMap = map[string]*regexp.Regexp{
		"brief":      regexp.MustCompile(`^([VDIWEAF])\/(.*?)\s*\(\s*(\d+)\):\s(.*?)\s*$`),
		"process":    regexp.MustCompile(`^([VDIWEAF])\(\s*(\d+)\)\s(.*?)\s*\((.*)\)$`),
		"tag":        regexp.MustCompile(`^([VDIWEAF])\/(.*?)\s*:\s(.*?)\s*$`),
		"time":       regexp.MustCompile(`^(\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\.\d+):*\s([VDIWEAF])\/(.*?)\s*\(\s*(\d+)\):\s(.*?)\s*$`),
		"threadtime": regexp.MustCompile(`^(\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\.\d+)\s*(\d+)\s*(\d+)\s([VDIWEAF])\s(.*?)\s*:\s(.*?)\s*$`),
		"ddms_save":  regexp.MustCompile(`^(\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\.\d+):*\s(VERBOSE|DEBUG|ERROR|WARN|INFO|ASSERT)\/(.*?)\((\s*\d+)\):\s(.*?)\s*$`),
		"raw":        regexp.MustCompile(`^(.*?)\s*$`),
	}
)

// Parse a line of logcat.
func (p *logcatParser) Parse(line string) (LogcatEntry, error) {
	logFormat := p.findFormat(line)
	if logFormat == "" {
		return nil, errors.New("Parse failed: " + line)
	}
	item := p.search(line, logFormat)
	return item, nil
}

// search keyword using regex pattern.
func (p *logcatParser) search(line string, format string) LogcatEntry {
	pattern := patternMap[format]
	match := pattern.FindStringSubmatch(line)
	item := LogcatEntry{}
	for i, val := range match[1:] {
		key := paramOrderMap[format][i]
		item[key] = val
	}
	return item
}

// findFormat analyze a line and find out logcat format.
func (p *logcatParser) findFormat(line string) string {
	for _, format := range formats {
		re, ok := patternMap[format]
		if ok && re.MatchString(line) {
			return format
		}
	}
	return ""
}
