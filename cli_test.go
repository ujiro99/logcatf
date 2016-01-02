package main

import (
	"bufio"
	"bytes"
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
	cli := newCli()
	args := strings.Split("./logcatf", " ")

	filename := logPrefix + "threadtime" + logExt
	fp, err := os.Open(filename)
	if err != nil {
		t.Errorf("os.Open: %v", err)
	}
	defer fp.Close()

	cli.inStream = bufio.NewReader(fp)
	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
}
