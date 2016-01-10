package main

import (
	"fmt"
	"io"

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

var (
	formatter Formatter
	parser    Parser
	writer    io.Writer
)

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {

	cli.initialize(args)
	err := cli.verifyParameter()
	if err != nil {
		fmt.Fprintln(cli.errStream, err.Error())
		log.Debug(err.Error())
		return ExitCodeError
	}

	// let's start
	for line := range lines.Lines(cli.inStream) {
		item := cli.parseLine(line)
		cli.executor.IfMatch(&line).Exec(&item)
	}

	log.Debugf("run finished")
	return ExitCodeOK
}

// exec parse and format
func (cli *CLI) parseLine(line string) LogcatItem {
	item := parser.Parse(line)
	if item == nil {
		return nil
	}
	output := formatter.Format(&item)
	fmt.Fprintln(writer, output)
	return item
}

func (cli *CLI) initialize(args []string) {
	// setup kingpin & parse args
	var (
		app     = kingpin.New(Name, Message["commandDescription"])
		format  = app.Arg("format", Message["helpFormat"]).Default(DefaultFormat).String()
		trigger = app.Flag("on", Message["helpTrigger"]).Short('o').Regexp()
		command = app.Flag("command", Message["helpCommand"]).Short('c').String()
		encode  = app.Flag("encode", Message["helpEncode"]).String()
		toCsv   = app.Flag("to-csv", Message["helpToCsv"]).Bool()
	)
	app.HelpFlag.Short('h')
	app.Version(Version)
	kingpin.MustParse(app.Parse(args[1:]))

	// if trigger not exists, not execute anything.
	if *trigger == nil {
		cli.executor = &emptyExecutor{}
	} else {
		cli.executor = &executor{
			trigger: *trigger,
			command: command,
			Stdout:  cli.errStream,
		}
	}

	// initialize formatter
	if *toCsv {
		formatter = &csvFormatter{
			&defaultFormatter{format: format},
		}
	} else {
		formatter = &defaultFormatter{format: format}
	}
	// convert format (long => short)
	formatter.Normarize()

	// initialize writer
	if *encode == "shift-jis" {
		writer = transform.NewWriter(cli.outStream, japanese.ShiftJIS.NewEncoder())
	} else {
		writer = cli.outStream
	}

	// initialize parser
	parser = &logcatParser{}

	log.WithFields(log.Fields{"format": *format, "trigger": *trigger, "command": *command}).Debug("Parameter initialized.")
}

func (cli *CLI) verifyParameter() error {
	return formatter.Verify()
}
