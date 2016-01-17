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
	executors            []Executor
}

var (
	formatter Formatter
	parser    Parser
	writer    io.Writer
	fmtc      Colorizer
)

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {

	// initialize
	err := cli.initialize(args)
	if err != nil {
		fmt.Fprintln(cli.errStream, err.Error())
		log.Debug(err.Error())
		return ExitCodeError
	}

	// let's start
	for line := range lines.Lines(cli.inStream) {
		item := cli.parseLine(line)
		cli.execute(line, item)
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
	fmtc.Fprintln(writer, output, item)
	return item
}

func (cli *CLI) initialize(args []string) error {
	// setup kingpin & parse args
	var (
		app      = kingpin.New(Name, Message["commandDescription"])
		format   = app.Arg("format", Message["helpFormat"]).Default(DefaultFormat).String()
		triggers = app.Flag("on", Message["helpTrigger"]).Short('o').RegexpList()
		commands = app.Flag("command", Message["helpCommand"]).Short('c').Strings()
		encode   = app.Flag("encode", Message["helpEncode"]).Default(UTF8).String()
		toCsv    = app.Flag("to-csv", Message["helpToCsv"]).Bool()
		color    = app.Flag("color", Message["helpToColor"]).Bool()

		colorV = app.Flag("color-v", Message["helpToColorV"]).PlaceHolder("COLOR").String()
		colorD = app.Flag("color-d", Message["helpToColorD"]).PlaceHolder("COLOR").String()
		colorI = app.Flag("color-i", Message["helpToColorI"]).PlaceHolder("COLOR").String()
		colorW = app.Flag("color-w", Message["helpToColorW"]).PlaceHolder("COLOR").String()
		colorE = app.Flag("color-e", Message["helpToColorE"]).PlaceHolder("COLOR").String()
		colorF = app.Flag("color-f", Message["helpToColorF"]).PlaceHolder("COLOR").String()
	)
	app.HelpFlag.Short('h')
	app.Version(Version)
	kingpin.MustParse(app.Parse(args[1:]))

	// initialize executor
	if *triggers == nil {
		// if trigger not exists, not execute anything.
		cli.executors = []Executor{&emptyExecutor{}}
	} else {
		if len(*triggers) != len(*commands) {
			return &ParameterError{Message["msgCommandNumMismatch"]}
		}
		es := []Executor{}
		for i, t := range *triggers {
			es = append(es, &executor{
				trigger: t,
				command: &(*commands)[i],
				Stdout:  cli.errStream,
			})
		}
		cli.executors = es
	}

	// initialize colorizer
	config := ColorConfig{
		"V": *colorV,
		"D": *colorD,
		"I": *colorI,
		"W": *colorW,
		"E": *colorE,
		"F": *colorF,
	}
	fmtc = Colorizer{}
	fmtc.SetUp(*color, config)
	newFormat := fmtc.ReplaceColorCode(*format)

	// initialize formatter
	if *toCsv {
		formatter = &csvFormatter{
			&defaultFormatter{format: &newFormat},
		}
	} else {
		formatter = &defaultFormatter{format: &newFormat}
	}
	// convert format (long => short)
	formatter.Normarize()

	// initialize writer
	switch *encode {
	case ShiftJIS:
		writer = transform.NewWriter(cli.outStream, japanese.ShiftJIS.NewEncoder())
	case EUCJP:
		writer = transform.NewWriter(cli.outStream, japanese.EUCJP.NewEncoder())
	case ISO2022JP:
		writer = transform.NewWriter(cli.outStream, japanese.ISO2022JP.NewEncoder())
	default:
		writer = cli.outStream
	}

	// initialize parser
	parser = &logcatParser{}

	log.WithFields(log.Fields{
		"format":  *format,
		"trigger": *triggers,
		"command": *commands}).Debug("Parameter initialized.")

	return formatter.Verify()
}

// execute calls multiple executors.
func (cli *CLI) execute(line string, item LogcatItem) {
	for _, e := range cli.executors {
		e.IfMatch(line).Exec(item)
	}
}

// ParameterError has error message of parameter.
type ParameterError struct {
	msg string
}

// Error returns all error message.
func (e *ParameterError) Error() string {
	return e.msg
}
