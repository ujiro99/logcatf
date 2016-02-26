package main

import (
	"io"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/ujiro99/logcatf/logcat"
)

// Executor executes a command if neccesary.
type Executor interface {

	// IfMatch watches logcat line.
	// - if line matches trigger, returns executer.
	// - if line not matches trigger, returns emptyExecuter.
	IfMatch(line string) Executor

	// Exec command.
	Exec(item logcat.Entry)
}

// implements Executer.
type executor struct {
	trigger *regexp.Regexp
	command *string
	Stdout  io.Writer
}

var (
	empty = emptyExecutor{}

	flagMapLinux = map[string]string{
		"%t": "${time}",
		"%i": "${pid}",
		"%I": "${tid}",
		"%p": "${priority}",
		"%a": "${tag}",
		"%m": "${message}",
	}

	flagMapWindows = map[string]string{
		"%t": "%time%",
		"%i": "%pid%",
		"%I": "%tid%",
		"%p": "%priority%",
		"%a": "%tag%",
		"%m": "%message%",
	}
)

// check a line. implements Executer.
func (e *executor) IfMatch(line string) Executor {
	if line == "" || !e.trigger.MatchString(line) {
		return &empty
	}
	log.Debugf("--execute on: \"%s\" ", e.trigger.String())
	return e
}

// execute command. implements Executer.
func (e *executor) Exec(item logcat.Entry) {
	log.Debugf("--command start: \"%s\"", *e.command)

	if item != nil {
		for _, k := range item.Keys() {
			os.Setenv(k, item[k])
		}
	}

	var cmd *exec.Cmd
	if runtime.GOOS == Windows {
		command := e.replaceFlags(*e.command, flagMapWindows)
		e.command = &command
		cmd = exec.Command(os.Getenv("COMSPEC"), "/c", *e.command)
	} else {
		command := e.replaceFlags(*e.command, flagMapLinux)
		e.command = &command
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

// replace variable flags in command.
func (e *executor) replaceFlags(cmd string, flags map[string]string) string {
	for k, v := range flags {
		cmd = strings.Replace(cmd, k, v, flagAll)
	}
	return cmd
}

// implements Executer. but not execute anything.
type emptyExecutor struct {
}

// don't execute command.
func (e *emptyExecutor) IfMatch(line string) Executor {
	return e
}

// don't execute command.
func (e *emptyExecutor) Exec(item logcat.Entry) {
}
