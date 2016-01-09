package main

import (
	"fmt"
	"io"
	"regexp"

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

// CLI is the command line object
type CLI struct {
	inStream             io.Reader
	outStream, errStream io.Writer
	executor             Executor
}

// CLIParameter represents parameters to execute command.
type CLIParameter struct {
	format *string
}

var (
	formatter Formatter
	parser    Parser
	writer    io.Writer
)

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
		cli.executor.IfMatch(&line).Exec(&item)
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

func (cli *CLI) initParameter(args []string) *CLIParameter {
	// setup kingpin & parse args
	var (
		app     = kingpin.New(Name, Message["commandDescription"])
		format  = app.Arg("format", Message["helpFormat"]).Default(DefaultFormat).String()
		trigger = app.Flag("on", Message["helpTrigger"]).Short('o').String()
		command = app.Flag("command", Message["helpCommand"]).Short('c').String()
		encode  = app.Flag("encode", Message["helpEncode"]).String()
		toCsv   = app.Flag("to-csv", Message["helpToCsv"]).Bool()
	)
	app.HelpFlag.Short('h')
	app.Version(Version)
	kingpin.MustParse(app.Parse(args[1:]))

	// if trigger not exists, not execute anything.
	if *trigger == "" {
		cli.executor = &emptyExecutor{}
	} else {
		cli.executor = &executor{
			trigger: regexp.MustCompile(*trigger),
			command: command,
			Stdout:  cli.errStream,
		}
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
		format: format,
	}
}

func (cli *CLI) verifyParameter(param *CLIParameter) error {
	return formatter.Verify(*param.format)
}
