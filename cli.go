package main

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/Maki-Daisuke/go-lines"
	log "github.com/Sirupsen/logrus"
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
}

// CLIParameter represents parameters to execute command.
type CLIParameter struct {
	format,
	trigger,
	command *string
}

var formatter Formatter

func init() {
	formatter = &logcatFormatter{}
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

	// convert format (long => short)
	normalized := formatter.Normarize(*param.format)
	param.format = &normalized

	for line := range lines.Lines(cli.inStream) {
		item := cli.parseLine(param, line)
		cli.eventFunc(param, &line, &item)
	}

	log.Debugf("finished")
	return ExitCodeOK
}

// exec parse and format
func (cli *CLI) parseLine(param *CLIParameter, line string) LogcatItem {
	item := Parse(line)
	if item == nil {
		return nil
	}
	output := formatter.Format(*param.format, &item)
	fmt.Fprintln(cli.outStream, output)
	return item
}

// dont execute command.
func (cli *CLI) execCommandNot(param *CLIParameter, line *string, item *LogcatItem) {
}

// execute command if
func (cli *CLI) execCommand(param *CLIParameter, line *string, item *LogcatItem) {
	if !strings.Contains(*line, *param.trigger) {
		return
	}
	log.Debugf("--command start: \"%s\" on \"%s\"", *param.command, *line)
	out, err := exec.Command("sh", "-c", *param.command).Output()
	if err != nil {
		log.Error(err)
	}
	fmt.Fprintf(cli.errStream, "%s", out)
	log.Debugf("--command finish: \"%s\"", *param.command)
}

func (cli *CLI) initParameter(args []string) *CLIParameter {
	// setup kingpin & parse args
	var (
		app     = kingpin.New(Name, Message["commandDescription"])
		format  = app.Arg("format", Message["helpFormat"]).Default(DefaultFormat).String()
		trigger = app.Flag("on", Message["helpTrigger"]).Short('o').String()
		command = app.Flag("command", Message["helpCommand"]).Short('c').String()
	)
	app.HelpFlag.Short('h')
	app.Version(Version)
	kingpin.MustParse(app.Parse(args[1:]))

	// if trigger not exists, not execute anything.
	if *trigger == "" {
		cli.eventFunc = cli.execCommandNot
	} else {
		cli.eventFunc = cli.execCommand
	}

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
