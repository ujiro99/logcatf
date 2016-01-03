package main

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
)

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

	// regex pattrens for analyze logcat.
	patterns = map[string]*regexp.Regexp{
		"brief":      regexp.MustCompile(`^(?P<priority>[VDIWEAF])\/(?P<tag>.*?)\s*\(\s*(?P<pid>\d+)\):\s(?P<message>.*?)\s*$`),
		"process":    regexp.MustCompile(`^(?P<priority>[VDIWEAF])\(\s*(?P<pid>\d+)\)\s(?P<message>.*?)\s*\((?P<tag>.*)\)$`),
		"tag":        regexp.MustCompile(`^(?P<priority>[VDIWEAF])\/(?P<tag>.*?)\s*:\s(?P<message>.*?)\s*$`),
		"time":       regexp.MustCompile(`^(?P<time>\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\.\d+):*\s(?P<priority>[VDIWEAF])\/(?P<tag>.*?)\s*\(\s*(?P<pid>\d+)\):\s(?P<message>.*?)\s*$`),
		"threadtime": regexp.MustCompile(`^(?P<time>\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\.\d+)\s*(?P<pid>\d+)\s*(?P<tid>\d+)\s(?P<priority>[VDIWEAF])\s(?P<tag>.*?)\s*:\s(?P<message>.*?)\s*$`),
		"ddms_save":  regexp.MustCompile(`^(?P<time>\d{2}-\d{2}\s\d{2}:\d{2}:\d{2}\.\d+):*\s(?P<priority>VERBOSE|DEBUG|ERROR|WARN|INFO|ASSERT)\/(?P<tag>.*?)\((?P<pid>\s*\d+)\):\s(?P<message>.*?)\s*$`),
		"raw":        regexp.MustCompile(`^(?P<message>.*?)\s*$`),
	}
)

// Parse a line of logcat.
func Parse(line string) LogcatItem {
	logFormat := findFormat(line)
	if logFormat == "" {
		log.Warnf("Parse failed: %s", line)
		return nil
	}
	match := search(line, patterns[logFormat])
	item := LogcatItem{}
	for key, value := range match {
		item[key] = value
	}

	return item
}

// search keyword using regex pattern.
func search(line string, pattern *regexp.Regexp) map[string]string {
	match := pattern.FindStringSubmatch(line)
	result := make(map[string]string, len(keys))
	for i, name := range pattern.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	return result
}

// findFormat analyze a line and find out logcat format.
func findFormat(line string) string {
	for _, format := range formats {
		re, ok := patterns[format]
		if ok && re.MatchString(line) {
			return format
		}
	}

	return ""
}
