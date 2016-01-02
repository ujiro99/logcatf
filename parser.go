package main

import (
	"regexp"

	log "github.com/Sirupsen/logrus"
)

// Parser can parse logcat.
type Parser interface {
	// Parse parses a line of logcat.
	Parse(line string) LogcatItem
}

type logcatParser struct {
	logFormat string
	pattern   *regexp.Regexp
}

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

	// cache of parser.
	parserCache = make(map[string]Parser, len(formats))
)

// Parse implements Parser interface.
func (p *logcatParser) Parse(line string) LogcatItem {
	item := LogcatItem{}
	match := p.search(line)

	for key, value := range match {
		item[key] = value
	}

	return item
}

func (p *logcatParser) search(line string) map[string]string {
	match := p.pattern.FindStringSubmatch(line)
	result := make(map[string]string, len(keys))
	for i, name := range p.pattern.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	return result
}

// GetParser create and return parser
func GetParser(format string) Parser {
	log.WithFields(log.Fields{"format": format}).Info("set parser type")

	parser, ok := parserCache[format]

	if ok {
		return parser
	}

	parser = &logcatParser{
		logFormat: format,
		pattern:   patterns[format],
	}
	parserCache[format] = parser
	return parser
}

// FindFormat analyze a line and find out logcat format.
func FindFormat(line string) string {
	log.WithFields(log.Fields{"line": line}).Debug("FindFormat start")

	for _, format := range formats {
		re, ok := patterns[format]
		if ok && re.MatchString(line) {
			log.WithFields(log.Fields{"format": format}).Debug("FindFormat finish")
			return format
		}
	}

	return ""
}
