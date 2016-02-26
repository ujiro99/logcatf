package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/Maki-Daisuke/go-lines"
)

var (
	parsedContains = map[string]map[string]bool{
		"brief":      map[string]bool{"time": false, "pid": true, "tid": false, "priority": true, "tag": true, "message": true},
		"process":    map[string]bool{"time": false, "pid": true, "tid": false, "priority": true, "tag": true, "message": true},
		"tag":        map[string]bool{"time": false, "pid": false, "tid": false, "priority": true, "tag": true, "message": true},
		"time":       map[string]bool{"time": true, "pid": true, "tid": false, "priority": true, "tag": true, "message": true},
		"threadtime": map[string]bool{"time": true, "pid": true, "tid": true, "priority": true, "tag": true, "message": true},
		"raw":        map[string]bool{"time": false, "pid": false, "tid": false, "priority": false, "tag": false, "message": true},
	}

	// TODO not supported yet.
	//"long", "ddms_save"
	logPatterns = map[string]string{
		"brief":      "I/tag_value  (  930): message_value",
		"process":    "I(  930) message_value  (tag_value)",
		"tag":        "I/tag_value  : message_value",
		"time":       "01-01 00:00:00.000 I/tag_value  (  930): message_value",
		"threadtime": "01-01 00:00:00.000   930   931 I tag_value  : message_value",
		"raw":        "message_value",
	}

	parsedValue = map[string]map[string]string{
		"brief":      map[string]string{"pid": "930", "priority": "I", "tag": "tag_value", "message": "message_value"},
		"process":    map[string]string{"pid": "930", "priority": "I", "tag": "tag_value", "message": "message_value"},
		"tag":        map[string]string{"priority": "I", "tag": "tag_value", "message": "message_value"},
		"time":       map[string]string{"time": "01-01 00:00:00.000", "pid": "930", "priority": "I", "tag": "tag_value", "message": "message_value"},
		"threadtime": map[string]string{"time": "01-01 00:00:00.000", "pid": "930", "tid": "931", "priority": "I", "tag": "tag_value", "message": "message_value"},
		"raw":        map[string]string{"message": "message_value"},
	}

	logPatternsHasSpace = map[string]string{
		"brief":      "I/auditd  (  930):   test Message  ",
		"process":    "I(  930)   test Message    (auditd)",
		"tag":        "I/auditd  :   test Message  ",
		"time":       "12-28 18:54:07.180 I/auditd  (  930):   test Message  ",
		"threadtime": "12-28 18:54:07.180   930   931 I auditd  :   test Message  ",
		"raw":        "  test Message  ",
	}

	logPatternsHasTab = map[string]string{
		"brief": "I/auditd  (  930): 	test Message",
		"process": "I(  930) 	test Message  (auditd)",
		"tag": "I/auditd  : 	test Message",
		"time": "12-28 18:54:07.180 I/auditd  (  930): 	test Message",
		"threadtime": "12-28 18:54:07.180   930   931 I auditd  : 	test Message",
		"raw": "	test Message",
	}
)

func TestFindFormat(t *testing.T) {
	parser := logcatParser{}
	for expect, log := range logPatterns {
		finded := parser.findFormat(log)
		if expect != finded {
			t.Errorf("expected %s to eq %s", expect, finded)
		}
	}
}

func TestParse_in_test_dir(t *testing.T) {
	parser := logcatParser{}
	for format := range logPatterns {

		filename := logPrefix + format + logExt
		fp, err := os.Open(filename)
		if err != nil {
			t.Errorf("os.Open: %v", err)
		}
		defer fp.Close()
		for line := range lines.Lines(bufio.NewReader(fp)) {
			// logfile has 2 formats, raw and some else.
			// "--------- beginning of *" must be used raw format.
			item, _ := parser.Parse(line)
			expect := parsedContains[parser.findFormat(line)]
			for key, hasValue := range expect {
				if hasValue {
					_, ok := item[key]
					if !ok {
						t.Errorf("%s must has value", key)
					}
				}
			}
		}
	}
}

func TestParse_value(t *testing.T) {
	parser := logcatParser{}
	for format, line := range logPatterns {
		item, _ := parser.Parse(line)
		expectItem := parsedValue[format]
		for key, val := range expectItem {
			if val != item[key] {
				t.Errorf("\nexpect: %s\nparsed: %s", val, item[key])
			}
		}
	}
}

func TestParse_removeTailSpace(t *testing.T) {
	expects := map[string]string{
		"time":     "12-28 18:54:07.180",
		"pid":      "930",
		"tid":      "931",
		"priority": "I",
		"tag":      "auditd",
		"message":  "  test Message",
	}

	parser := logcatParser{}
	for format, log := range logPatternsHasSpace {
		item, _ := parser.Parse(log)
		for key, expect := range expects {
			if parsedContains[format][key] && expect != item[key] {
				t.Errorf("on %s, \"%s\" must eq \"%s\"", format, item[key], expect)
			}
		}
	}
}

func TestParse_hasTabInMessage(t *testing.T) {
	expects := map[string]string{
		"message": "	test Message",
	}

	parser := logcatParser{}
	for format, log := range logPatternsHasTab {
		item, _ := parser.Parse(log)
		for key, expect := range expects {
			if parsedContains[format][key] && expect != item[key] {
				t.Errorf("on %s, \"%s\" must eq \"%s\"", format, item[key], expect)
			}
		}
	}
}

func BenchmarkParse(b *testing.B) {
	parser := logcatParser{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.Parse(logPatterns["threadtime"])
	}
}

func BenchmarkParse_FindFormat(b *testing.B) {
	parser := logcatParser{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.findFormat(logPatterns["threadtime"])
	}
}

func BenchmarkParse_Search(b *testing.B) {
	parser := logcatParser{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.search(logPatterns["threadtime"], "threadtime")
	}
}
