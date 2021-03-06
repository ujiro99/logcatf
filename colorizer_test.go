package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/ujiro99/logcatf/logcat"
)

func TestRun_ReplaceColorCode(t *testing.T) {

	before := "%t [blue] %m"
	expect := "%t [LOGCATF_blue] %m"

	f := Colorizer{}
	after := f.ReplaceColorCode(before)

	if after != expect {
		t.Errorf("\nexpect: %s\nresult: %s", expect, after)
	}
}

func TestRun_Fprintln_enable(t *testing.T) {

	format := "%t %p [_blue_]%m"
	expect := "\033[31m%t %p \033[44m%m"

	item := logcat.Entry{
		"time":     "12-28 19:01:14.073",
		"priority": "F",
		"message":  "logcat_message",
	}
	w := new(bytes.Buffer)

	f := Colorizer{}
	cc := ColorConfig{}
	f.Init(true, cc)
	format = f.ReplaceColorCode(format)
	f.Fprintln(w, format, item)

	res := w.String()

	if !strings.Contains(res, expect) {
		t.Errorf("\nexpect: %s\nresult: %s", expect, res)
	}
}

func TestRun_Fprintln_disable(t *testing.T) {

	format := "%t %p [blue]%m"
	expect := "%t %p %m"

	item := logcat.Entry{
		"time":     "12-28 19:01:14.073",
		"priority": "F",
		"message":  "logcat_message",
	}
	w := new(bytes.Buffer)

	f := Colorizer{}
	f.Init(false, nil)
	format = f.ReplaceColorCode(format)
	f.Fprintln(w, format, item)

	res := w.String()

	if !strings.Contains(res, expect) {
		t.Errorf("\nexpect: %s\nresult: %s", expect, res)
	}
}
