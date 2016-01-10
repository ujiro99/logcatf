package main

import (
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	log "github.com/Sirupsen/logrus"
)

// Executor executes a command if neccesary.
type Executor interface {

	// IfMatch watches logcat line.
	// - if line matches trigger, returns executer.
	// - if line not matches trigger, returns emptyExecuter.
	IfMatch(line *string) Executor

	// Exec command.
	Exec(item *LogcatItem)
}

// implements Executer.
type executor struct {
	trigger *regexp.Regexp
	command *string
	Stdout  io.Writer
}

// check a line. implements Executer.
func (e *executor) IfMatch(line *string) Executor {
	if !e.trigger.MatchString(*line) {
		return &emptyExecutor{}
	}
	log.Debugf("--execute on: \"%s\" ", e.trigger.String())
	return e
}

// execute command. implements Executer.
func (e *executor) Exec(item *LogcatItem) {
	log.Debugf("--command start: \"%s\"", *e.command)

	for _, k := range item.Keys() {
		os.Setenv(k, (*item)[k])
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(os.Getenv("COMSPEC"), "/c", *e.command)
	} else {
		cmd = exec.Command(os.Getenv("SHELL"), "-c", *e.command)
	}
	cmd.Stdout = e.Stdout
	err := cmd.Start()
	if err != nil {
		log.Error(err)
	}
	// cmd.Wait()

	log.Debugf("--command finish: \"%s\"", *e.command)
}

// implements Executer. but not execute anything.
type emptyExecutor struct {
}

// don't execute command.
func (e *emptyExecutor) IfMatch(line *string) Executor {
	return &emptyExecutor{}
}

// don't execute command.
func (e *emptyExecutor) Exec(item *LogcatItem) {
}
