package main

import (
	"bytes"
	"encoding/csv"
	"runtime"
	"strings"
)

// implements Formatter
type csvFormatter struct {
	*defaultFormatter
	buf *bytes.Buffer
	w   *csv.Writer
}

func NewCsvFormatter(format string) Formatter {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	res := &csvFormatter{&defaultFormatter{format: &format}, buf, writer}
	if runtime.GOOS == Windows {
		res.w.UseCRLF = true
	}
	return res
}

// Format implements Formatter
// replace %* keywords to real value. use short format.
//   ex) "%t %a" => "12-28 19:01:14.073", GLSUser
func (f *csvFormatter) Format(item *LogcatEntry) string {
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

	f.w.Write(args)
	f.w.Flush()
	res := strings.TrimRight(f.buf.String(), "\r\n")
	f.buf.Reset()
	return res
}
