package main

import (
	log "github.com/Sirupsen/logrus"
	"regexp"
)

// Parser parse logcat to struct.
type Parser interface {
	Parse(line string) LogcatItem
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
func (p *logcatParser) Parse(line string) LogcatItem {
	logFormat := p.findFormat(line)
	if logFormat == "" {
		log.Warnf("Parse failed: %s", line)
		return nil
	}
	item := p.search(line, patterns[logFormat])
	return item
}

// search keyword using regex pattern.
func (p *logcatParser) search(line string, pattern *regexp.Regexp) LogcatItem {
	match := pattern.FindStringSubmatch(line)
	item := LogcatItem{}
	for i, key := range pattern.SubexpNames() {
		if i != 0 {
			item[key] = match[i]
		}
	}
	return item
}

// findFormat analyze a line and find out logcat format.
func (p *logcatParser) findFormat(line string) string {
	for _, format := range formats {
		re, ok := patterns[format]
		if ok && re.MatchString(line) {
			return format
		}
	}

	return ""
}
