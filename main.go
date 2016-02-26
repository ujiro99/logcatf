package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/mattn/go-colorable"
)

func init() {
	log.SetLevel(log.ErrorLevel)
	log.SetOutput(colorable.NewColorableStderr())
}

func main() {
	cli := &CLI{
		inStream:  os.Stdin,
		outStream: colorable.NewColorableStdout(),
		errStream: colorable.NewColorableStderr(),
	}
	os.Exit(cli.Run(os.Args))
}
