package main

import (
	"time"
)

// LogcatEntry represents a line of logcat log.
type LogcatEntry map[string]string

// Keys returns valid keys for LogcatEntry
func (item *LogcatEntry) Keys() []string {
	return []string{
		"time", "pid", "tid", "priority", "tag", "message",
	}
}

// Time return parsed time.
func (item *LogcatEntry) Time() time.Time {
	t, _ := time.Parse("12-28 18:54:08.043", (*item)["time"])
	return t
}

// Pid.
func (item *LogcatEntry) Pid() string {
	return (*item)["pid"]
}

// Tid.
func (item *LogcatEntry) Tid() string {
	return (*item)["tid"]
}

// Priority.
func (item *LogcatEntry) Priority() string {
	return (*item)["priority"]
}

// Tag.
func (item *LogcatEntry) Tag() string {
	return (*item)["tag"]
}

// Message.
func (item *LogcatEntry) Message() string {
	return (*item)["message"]
}
