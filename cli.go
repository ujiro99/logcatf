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
		parser Parser
		format = kingpin.Arg("format", Message["helpFormat"]).String()

		// TODO not implemented yet.
		//trigger = kingpin.Flag("trigger", "Regex for triggering a command.").Short('t').String()
		//command = kingpin.Flag("command", "Command.").Short('c').String()
	)

	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Version(Version)
	kingpin.Parse() // parse args here

	if *format == "" {
		*format = DefaultFormat
	}

	for line := range lines.Lines(cli.inStream) {
		logFormat := FindFormat(line)
		if logFormat == "" {
			continue
		}
		parser = GetParser(logFormat)
		item := parser.Parse(line)

		// output
		fmt.Println(item.Format(*format))
	}
	log.Println("finished")

	return ExitCodeOK
}
