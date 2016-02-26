package logcat

import (
	"time"
)

// Entry represents a line of logcat log.
type Entry map[string]string

// Keys returns valid keys for Entry
func (item *Entry) Keys() []string {
	return []string{
		"time", "pid", "tid", "priority", "tag", "message",
	}
}

// Time return parsed time.
func (item *Entry) Time() time.Time {
	t, _ := time.Parse("12-28 18:54:08.043", (*item)["time"])
	return t
}
