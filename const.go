package main

const (
	// Name of this command.
	Name string = "logcatf"
	// Version of this command.
	Version string = "0.1.0"
	// DefaultFormat will be used if format wasn't specified.
	DefaultFormat string = "%time %tag %priority %message"
)

// Message has message strings.
var Message = map[string]string{
	"commandDescription": "A command to format Android Logcat",
	"helpFormat": `Format of output
	- Available keyword:
	  %time, %tag, %priority, %pid, %tid, %message

	- Example:
	  "%time %message"
	  "%time, %tag, %priority, %pid, %tid, %message"

	- Default:
	  "%time %tag %priority %message"`,
}
