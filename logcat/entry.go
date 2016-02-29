package logcat

var allKeys = []string{"time", "pid", "tid", "priority", "tag", "message"}

// Entry represents a line of logcat log.
type Entry map[string]string

// Keys returns existing keys in this Entry
func (item *Entry) Keys() []string {
	keys := make([]string, len(*item))
	j := 0
	for _, k := range allKeys {
		if _, ok := (*item)[k]; ok {
			keys[j] = k
			j = j + 1
		}
	}
	return keys
}

// Values returns existing values int this Entry
func (item *Entry) Values() []string {
	values := make([]string, len(*item))
	j := 0
	for _, k := range allKeys {
		if v, ok := (*item)[k]; ok {
			values[j] = v
			j = j + 1
		}
	}
	return values
}
