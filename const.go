package main

const (
	// Name of this command.
	Name string = "logcatf"
	// Version of this command.
	Version string = "0.1.0"
	// DefaultFormat will be used if format wasn't specified.
	DefaultFormat string = "%Time %Tag %Priority %Message"
)

// Message has message strings.
var Message = map[string]string{
	"helpFormat": `Format of output
	- Available keyword:
	  %Time, %Tag, %Priority, %Pid, %Tid, %Message

	- Example:
	  "%Time %Message"
	  "%Time, %Tag, %Priority, %Pid, %Tid, %Message"

	- Default:
	  "%Time %Tag %Priority %Message"`,
}
