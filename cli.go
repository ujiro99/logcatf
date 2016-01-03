package main

import (
	"fmt"
	"io"

	"github.com/Maki-Daisuke/go-lines"
	log "github.com/Sirupsen/logrus"
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
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {

	var (
		parser  Parser
		app     = kingpin.New(Name, Message["commandDescription"])
		format  = app.Arg("format", Message["helpFormat"]).Default(DefaultFormat).String()
		trigger = app.Flag("trigger", Message["helpTrigger"]).Short('t').String()
		command = app.Flag("command", Message["helpCommand"]).Short('c').String()
	)

	app.HelpFlag.Short('h')
	app.Version(Version)
	kingpin.MustParse(app.Parse(args[1:]))

	log.WithFields(log.Fields{
		"format":  *format,
		"trigger": *trigger,
		"command": *command,
	}).Debug("Parsed")

	for line := range lines.Lines(cli.inStream) {
		logFormat := FindFormat(line)
		if logFormat == "" {
			continue
		}
		parser = GetParser(logFormat)
		item := parser.Parse(line)

		// output
		fmt.Fprintln(cli.outStream, item.Format(*format))
	}
	log.Debugf("finished")

	return ExitCodeOK
}
