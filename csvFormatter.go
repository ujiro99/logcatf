package main

import (
	"bytes"
	"encoding/csv"
)

// implements Formatter
type csvFormatter struct {
	defaultFormatter
	format *string
}

// Format implements Formatter
// replace %* keywords to real value. use short format.
//   ex) "%t %a" => "12-28 19:01:14.073", GLSUser
func (f *csvFormatter) Format(item *LogcatItem) string {
	matches := sformatRegex.FindAllStringSubmatch(*f.format, len(formatMap))
	args := make([]string, 0, len(matches))
	// find matched keyword and store value on item
	for _, match := range matches {
		for _, m := range match {
			key, ok := findKey(formatMap, m)
			if ok {
				args = append(args, (*item)[key])
				break
			}
		}
	}

	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Write(args)
	writer.Flush()
	return buf.String()
}
