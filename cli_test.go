package main

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"os"
	"strings"
	"testing"
	"time"
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
	//args := []string{"./logcatf", "-o", "SourceFile:135", "-c", "sleep 0.3; say hello"}
	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
}

func TestRun_formatEscapedCharactor(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   930 I my_app  : message")
	out := new(bytes.Buffer)
	cli.outStream = out

	args := strings.Split(`./logcatf %time\t%tag\t%priority\n%message`, " ")
	status := cli.Run(args)

	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expect := "12-28 18:54:07.180\tmy_app\tI\nmessage\n"
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

	expect := fmt.Sprintf(Message["msgUnavailableKeyword"], "%level")
	str := err.String()
	if !strings.Contains(str, expect) {
		t.Errorf("expected \"%s\" to eq \"%s\"", str, expect)
	}
}

func TestRun_formatError_DuplicatedKeyword(t *testing.T) {
	err := new(bytes.Buffer)
	cli := newCli()
	cli.errStream = err

	args := strings.Split(`./logcatf "%t\t%p\t%t\t%priority"`, " ")
	status := cli.Run(args)

	if status != ExitCodeError {
		t.Errorf("expected %d to eq %d", status, ExitCodeError)
	}

	expect1 := fmt.Sprintf(Message["msgDuplicatedKeyword"], "%time", "%t")
	expect2 := fmt.Sprintf(Message["msgDuplicatedKeyword"], "%priority", "%p")
	str := err.String()
	if !strings.Contains(str, expect1) {
		t.Errorf("expected \"%s\" contains \"%s\"", str, expect1)
	}
	if !strings.Contains(str, expect2) {
		t.Errorf("expected \"%s\" contains \"%s\"", str, expect2)
	}
}

func TestRun_parameterError_exec_mismatch(t *testing.T) {
	err := new(bytes.Buffer)
	cli := newCli()
	cli.errStream = err

	args := []string{"./logcatf",
		"-o", "my_app.*first", "-c", "echo 1st $message",
		"-o", "my_app.*third",
	}
	status := cli.Run(args)

	if status != ExitCodeError {
		t.Errorf("expected %d to eq %d", status, ExitCodeError)
	}

	expect := fmt.Sprintf(Message["msgmsgCommandNumMismatch"])
	str := err.String()
	if !strings.Contains(str, expect) {
		t.Errorf("expect: \"%s\"\nresult: \"%s\"", expect, str)
	}
}

func TestRun_formatShort(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   931 I my_app  : message")
	out := new(bytes.Buffer)
	cli.outStream = out

	args := []string{"./logcatf", "%t %i %I %p %a: %m"}
	status := cli.Run(args)

	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expect := "12-28 18:54:07.180 930 931 I my_app: message\n"
	str := out.String()
	if str != expect {
		t.Errorf("expected \"%s\" to eq \"%s\"", str, expect)
	}
}

func TestRun_toCsv(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   931 I my_app  : hello, world.")
	out := new(bytes.Buffer)
	cli.outStream = out

	args := []string{"./logcatf", "%t %m", "--to-csv"}
	status := cli.Run(args)

	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expect := "12-28 18:54:07.180,\"hello, world.\""
	str := out.String()
	if !strings.Contains(str, expect) {
		t.Errorf("\nexpect: %s\nresult: %s", expect, str)
	}
}

func TestRun_execCommand(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   931 I my_app  : message")
	out := new(bytes.Buffer)
	cli.errStream = out

	args := []string{"./logcatf", "%t %i %I %p %a: %m", "-o", "my_app.*message", "-c", "echo test"}
	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	<-time.After(time.Second / 100)
	expect := "test"
	str := out.String()
	if !strings.Contains(str, expect) {
		t.Errorf("expected \"%s\" to eq \"%s\"", str, expect)
	}
}

func TestRun_execCommand_async(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180 0 1 I tag: first line\n12-28 18:54:07.180 0 1 I tag: second line")
	out := new(bytes.Buffer)
	cli.outStream = out

	args := []string{"./logcatf", "%t %i %I %p %a: %m", "-o", "first line", "-c", "sleep 0.2; echo finish"}
	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	<-time.After(time.Second / 10)
	expect := "second line"
	if !strings.Contains(out.String(), expect) {
		t.Error("command did not works asynchronously")
	}
}

func TestRun_execCommand_outputParsed(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   931 I my_app  : message")
	err := new(bytes.Buffer)
	cli.errStream = err

	args := []string{"./logcatf", "-o", "my_app.*message", "-c", "echo parsed: $time"}
	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	<-time.After(time.Second / 10)
	expect := "parsed: 12-28 18:54:07.180"
	str := err.String()
	if !strings.Contains(str, expect) {
		t.Errorf("$time was not expanded.")
	}
}

func TestRun_execCommand_multiple(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("" +
		"12-28 18:54:07.180   930   931 I my_app  : first\n" +
		"12-28 18:54:07.180   930   931 I my_app  : second\n" +
		"12-28 18:54:07.180   930   931 I my_app  : third\n")
	//err := new(bytes.Buffer)
	// TODO couldn't test this.
	// checked that exec works correctly, using os.Stdout.
	//cli.errStream = err
	//cli.errStream = os.Stdout
	args := []string{"./logcatf",
		"-o", "my_app.*first", "-c", "echo 1st $message",
		"-o", "my_app.*third", "-c", "echo 3rd $message",
	}
	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}
	// <-time.After(time.Second / 10)
	// expect := "1st first_line\n2nd second_line"
	// str := err.String()
	// if !strings.Contains(str, expect) {
	// 	t.Errorf("\nresult: %s\nexpect: %s", err.String(), expect)
	// }
}

func TestRun_encode_shiftjis(t *testing.T) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180     0     1 I logcat_messages:テスト")
	out := new(bytes.Buffer)
	cli.outStream = out

	args := strings.Split(`./logcatf "%m" --encode shift-jis`, " ")
	status := cli.Run(args)

	if status == ExitCodeError {
		t.Errorf("expected %d to eq %d", status, ExitCodeError)
	}

	str := out.String()
	expect := toShiftJis("logcat_messages:テスト")
	if !strings.Contains(str, expect) {
		t.Errorf("expected \"%s\" to eq \"%s\"", str, expect)
	}
}

func toShiftJis(str string) string {
	buf := new(bytes.Buffer)
	w := transform.NewWriter(buf, japanese.ShiftJIS.NewEncoder())
	w.Write([]byte(str))
	return buf.String()
}

func BenchmarkDefault(b *testing.B) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   931 I my_app  : message")
	args := []string{"./logcatf", "%t %i %I %p %a: %m"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli.Run(args)
	}
}

func BenchmarkEncodeShiftJis(b *testing.B) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   931 I my_app  : message")
	args := []string{"./logcatf", "%t %i %I %p %a: %m", "--encode", "shift-jis"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli.Run(args)
	}
}

func BenchmarkToCsv(b *testing.B) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   931 I my_app  : message")
	args := []string{"./logcatf", "%t %i %I %p %a: %m", "--to-csv"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli.Run(args)
	}
}

func BenchmarkExecCommand(b *testing.B) {
	cli := newCli()
	cli.inStream = strings.NewReader("12-28 18:54:07.180   930   931 I my_app  : message")
	args := []string{"./logcatf", "-o", "my_app.*message", "-c", "echo test"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cli.Run(args)
	}
}
