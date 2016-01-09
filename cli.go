package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/Maki-Daisuke/go-lines"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"gopkg.in/alecthomas/kingpin.v2"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// EventFunc watches logcat line.
// if line contains param.trigger, exec param.command.
type EventFunc func(param *CLIParameter, line *string, item *LogcatItem)

// CLI is the command line object
type CLI struct {
	inStream             io.Reader
	outStream, errStream io.Writer
	eventFunc            EventFunc
	eventTrigger         *regexp.Regexp
}

// CLIParameter represents parameters to execute command.
type CLIParameter struct {
	format,
	trigger,
	command *string
}

var (
	formatter Formatter
	parser    Parser
	writer    io.Writer
)

func init() {
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {

	param := cli.initParameter(args)
	err := cli.verifyParameter(param)
	if err != nil {
		fmt.Fprintln(cli.errStream, err.Error())
		log.Debug(err.Error())
		return ExitCodeError
	}

	// let's start
	for line := range lines.Lines(cli.inStream) {
		item := cli.parseLine(param, line)
		cli.eventFunc(param, &line, &item)
	}

	log.Debugf("run finished")
	return ExitCodeOK
}

// exec parse and format
func (cli *CLI) parseLine(param *CLIParameter, line string) LogcatItem {
	item := parser.Parse(line)
	if item == nil {
		return nil
	}
	output := formatter.Format(*param.format, &item)
	fmt.Fprintln(writer, output)
	return item
}

// dont execute command.
func (cli *CLI) execCommandNot(param *CLIParameter, line *string, item *LogcatItem) {
}

// execute command if
func (cli *CLI) execCommand(param *CLIParameter, line *string, item *LogcatItem) {
	if !cli.eventTrigger.MatchString(*line) {
		return
	}
	log.Debugf("--command start: \"%s\" on \"%s\"", *param.command, *line)

	for k := range formatMap {
		os.Setenv(k, (*item)[k])
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(os.Getenv("COMSPEC"), "/c", *param.command)
	} else {
		cmd = exec.Command(os.Getenv("SHELL"), "-c", *param.command)
	}
	cmd.Stdout = cli.errStream
	err := cmd.Start()
	if err != nil {
		log.Error(err)
	}
	// cmd.Wait()

	log.Debugf("--command finish: \"%s\"", *param.command)
}

func (cli *CLI) initParameter(args []string) *CLIParameter {
	// setup kingpin & parse args
	var (
		app     = kingpin.New(Name, Message["commandDescription"])
		format  = app.Arg("format", Message["helpFormat"]).Default(DefaultFormat).String()
		trigger = app.Flag("on", Message["helpTrigger"]).Short('o').String()
		command = app.Flag("command", Message["helpCommand"]).Short('c').String()
		encode  = app.Flag("encode", Message["helpEncode"]).String()
		toCsv   = app.Flag("toCsv", Message["helpToCsv"]).Bool()
	)
	app.HelpFlag.Short('h')
	app.Version(Version)
	kingpin.MustParse(app.Parse(args[1:]))

	// if trigger not exists, not execute anything.
	if *trigger == "" {
		cli.eventFunc = cli.execCommandNot
	} else {
		cli.eventFunc = cli.execCommand
		cli.eventTrigger = regexp.MustCompile(*trigger)
	}

	if *toCsv {
		formatter = &csvFormatter{}
	} else {
		formatter = &defaultFormatter{}
	}

	if *encode == "shift-jis" {
		writer = transform.NewWriter(cli.outStream, japanese.ShiftJIS.NewEncoder())
	} else {
		writer = cli.outStream
	}

	// convert format (long => short)
	normarized := formatter.Normarize(*format)
	format = &normarized
	parser = &logcatParser{}

	log.WithFields(log.Fields{"format": *format, "trigger": *trigger, "command": *command}).Debug("Parameter initialized.")
	return &CLIParameter{
		format:  format,
		trigger: trigger,
		command: command,
	}
}

func (cli *CLI) verifyParameter(param *CLIParameter) error {
	return formatter.Verify(*param.format)
}
