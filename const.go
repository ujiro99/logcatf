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
     %t (%time), %a (%tag), %p (%priority)
     %i (%pid), %I (%tid), %m (%message)
 - Example:
   adb logcat -v time | "%time %message"
   adb logcat -v threadtime | "%t, %a, %p, %i, %I, %m" > logcat.csv
 - Default Format:
   "%time %tag %priority %message"`,
	"helpTrigger":           "Regex for triggering a command.",
	"helpCommand":           "Command to be executed.",
	"msgUnavailableKeyword": "error: %s is not available. Please check `Availavle Keyword:` on help.",
	"msgDuplicatedKeyword":  "error: %s or %s is duplicated.",
}
