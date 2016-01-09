package main

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
	"time"
)

func Test_IfMatch(t *testing.T) {
	line := "aaa"

	e := executor{
		Stdout:  new(bytes.Buffer),
		trigger: regexp.MustCompile("a"),
	}
	res := e.IfMatch(&line)
	if &e != res {
		t.Errorf("trigger should mathes line.")
	}
}

func Test_IfMatch_not_match(t *testing.T) {
	line := "bbb"

	e := executor{
		Stdout:  new(bytes.Buffer),
		trigger: regexp.MustCompile("a"),
	}
	res := e.IfMatch(&line)
	if &e == res {
		t.Errorf("IfMatch should returns emptyExecutor.")
	}
}

func Test_Exec(t *testing.T) {
	expect := "test"
	command := "echo $message"
	item := LogcatItem{"message": expect}

	out := new(bytes.Buffer)
	e := executor{
		Stdout:  out,
		command: &command,
	}
	e.Exec(&item)
	<-time.After(time.Second / 100)

	if !strings.Contains(out.String(), expect) {
		t.Errorf("expect %s to eq \"%s\"", out.String(), expect)
	}
}

func Test_Exec_empty(t *testing.T) {
	expect := "test"
	item := LogcatItem{"message": expect}

	out := new(bytes.Buffer)
	e := emptyExecutor{}
	e.Exec(&item)
	<-time.After(time.Second / 100)

	if strings.Contains(out.String(), expect) {
		t.Errorf("emptyExecutor should not execute command")
	}
}
