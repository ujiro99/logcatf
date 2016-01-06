package main

import (
	"encoding/csv"
)

// implements Formatter
type csvFormatter struct {
	defaultFormatter
	Cli *CLI
}

// Format implements Formatter
// replace %* keywords to real value. use short format.
//   ex) "%t %a" => "12-28 19:01:14.073 GLSUser"
func (f *csvFormatter) Format(format string, item *LogcatItem) string {
	matches := sformatRegex.FindAllStringSubmatch(format, len(formatMap))
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

	writer := csv.NewWriter(f.Cli.outStream)
	writer.Write(args)
	writer.Flush()

	// output to StdOut has finished.
	return ""
}
