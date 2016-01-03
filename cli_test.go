package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

const (
	logPrefix = "./test/logcat."
	logExt    = ".txt"
)

func newCli() *CLI {
	inStream, outStream, errStream := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	return &CLI{
		inStream:  inStream,
		outStream: outStream,
		errStream: errStream,
	}
}

func TestRun_formatDefault(t *testing.T) {
	fp, err := os.Open(logPrefix + "threadtime" + logExt)
	if err != nil {
		t.Errorf("os.Open: %v", err)
	}
	defer fp.Close()

	cli := newCli()
	cli.inStream = bufio.NewReader(fp)
	args := strings.Split("./logcatf", " ")
	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
}

func TestRun_formatEscapedCharactor(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   930 I auditd  : type=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295")
	out := new(bytes.Buffer)
	cli.outStream = out

	args := strings.Split(`./logcatf %time\t%tag\t%priority\n%message`, " ")
	status := cli.Run(args)

	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expect := "12-28 18:54:07.180\tauditd\tI\ntype=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295\n"
	str := out.String()
	if str != expect {
		t.Errorf("expected \"%s\" to eq \"%s\"", str, expect)
	}
}

func TestRun_formatError_UnavailableKeyword(t *testing.T) {
	err := new(bytes.Buffer)
	cli := newCli()
	cli.errStream = err

	args := strings.Split(`./logcatf "%time\t%level"`, " ")
	status := cli.Run(args)

	if status != ExitCodeError {
		t.Errorf("expected %d to eq %d", status, ExitCodeError)
	}

	expect := fmt.Sprintf(Message["msgUnavailableKeyword"], "%level") + "\n"
	str := err.String()
	if str != expect {
		t.Errorf("expected \"%s\" to eq \"%s\"", str, expect)
	}
}

func TestRun_formatShort(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   931 I auditd  : type=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295")
	out := new(bytes.Buffer)
	cli.outStream = out

	args := []string{"./logcatf", "%t %i %I %p %a: %m"}
	status := cli.Run(args)

	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expect := "12-28 18:54:07.180 930 931 I auditd: type=1403 audit(0.0:2): policy loaded auid=4294967295 ses=4294967295\n"
	str := out.String()
	if str != expect {
		t.Errorf("expected \"%s\" to eq \"%s\"", str, expect)
	}
}
