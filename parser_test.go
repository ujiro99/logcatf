package main

import (
	"bufio"
	"github.com/Maki-Daisuke/go-lines"
	"os"
	"testing"
)

var (
	parsedHas = map[string]map[string]bool{
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
		"brief":      "I/auditd  (  930): type=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295",
		"process":    "I(  930) type=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295  (auditd)",
		"tag":        "I/auditd  : type=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295",
		"time":       "12-28 18:54:07.180 I/auditd  (  930): type=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295",
		"threadtime": "12-28 18:54:07.180   930   930 I auditd  : type=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295",
		"raw":        "type=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295",
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
	for expect, log := range logPatterns {
		finded := FindFormat(log)
		if expect != finded {
			t.Errorf("expected %s to eq %s", expect, finded)
		}
	}
}

func TestGetParser(t *testing.T) {
	for format := range logPatterns {
		parser := GetParser(format)
		if parser == nil {
			t.Errorf("Parser did not created on %s.", format)
		}
	}
}

func TestParse_allFormat(t *testing.T) {
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
			logFormat := FindFormat(line)
			parser := GetParser(logFormat)
			item := parser.Parse(line)
			expect := parsedHas[logFormat]
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

func TestParse_removeTailSpace(t *testing.T) {
	expects := map[string]string{
		"time":     "12-28 18:54:07.180",
		"pid":      "930",
		"tid":      "931",
		"priority": "I",
		"tag":      "auditd",
		"message":  "  test Message",
	}

	for format, log := range logPatternsHasSpace {
		parser := GetParser(format)
		item := parser.Parse(log)
		for key, expect := range expects {
			if parsedHas[format][key] && expect != item[key] {
				t.Errorf("on %s, \"%s\" must eq \"%s\"", format, item[key], expect)
			}
		}
	}
}

func TestParse_hasTabInMessage(t *testing.T) {
	expects := map[string]string{
		"message": "	test Message",
	}

	for format, log := range logPatternsHasTab {
		parser := GetParser(format)
		item := parser.Parse(log)
		for key, expect := range expects {
			if parsedHas[format][key] && expect != item[key] {
				t.Errorf("on %s, \"%s\" must eq \"%s\"", format, item[key], expect)
			}
		}
	}
}
