package main

import (
	"fmt"
	"io"
	"regexp"
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

// CLI is the command line object
type CLI struct {
	inStream             io.Reader
	outStream, errStream io.Writer
}

// CLIParameter represents parameters to execute command.
type CLIParameter struct {
	format,
	trigger,
	command *string
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {

	cliParam := initParameter(args)
	err := verifyParameter(cliParam)
	if err != nil {
		fmt.Fprintln(cli.errStream, err.Error())
		return ExitCodeError
	}

	// exec parse and format
	for line := range lines.Lines(cli.inStream) {
		item := Parse(line)
		if item == nil {
			continue
		}
		output := item.Format(*cliParam.format)
		fmt.Fprintln(cli.outStream, output)
	}

	log.Debugf("finished")
	return ExitCodeOK
}

func initParameter(args []string) *CLIParameter {
	// setup kingpin & parse args
	var (
		app     = kingpin.New(Name, Message["commandDescription"])
		format  = app.Arg("format", Message["helpFormat"]).Default(DefaultFormat).String()
		trigger = app.Flag("trigger", Message["helpTrigger"]).Short('t').String()
		command = app.Flag("command", Message["helpCommand"]).Short('c').String()
	)
	app.HelpFlag.Short('h')
	app.Version(Version)
	kingpin.MustParse(app.Parse(args[1:]))

	log.WithFields(log.Fields{"format": *format, "trigger": *trigger, "command": *command}).Debug("Cli Parameter initialized.")
	return &CLIParameter{
		format:  format,
		trigger: trigger,
		command: command,
	}
}

func verifyParameter(param *CLIParameter) error {
	// find unavailable keyword.
	noLimit := -1
	removed := formatRegexps.ReplaceAllString(*param.format, "")
	keyRegexp := regexp.MustCompile(`%\w+`)
	matches := keyRegexp.FindAllString(removed, noLimit)

	if len(matches) == 0 {
		return nil // no probrem!
	}

	// return error message.
	err := ParameterError{}
	for _, match := range matches {
		errMsg := fmt.Sprintf(Message["msgUnavailableKeyword"], match)
		err.errors = append(err.errors, errMsg)
	}
	return &err
}

// ParameterError has error message of parameter.
type ParameterError struct {
	errors []string
}

// Error returns all error message.
func (e *ParameterError) Error() string {
	if e.errors != nil {
		return strings.Join(e.errors, "\n")
	}
	return ""
}
