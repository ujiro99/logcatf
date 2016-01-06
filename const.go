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
   adb logcat -v time | logcatf "%t %4i %m"
   adb logcat -v threadtime | logcatf "%t, %a, %m" > logcat.csv
   adb logcat -v threadtime | logcatf "%t" -o "Exception" -c "adb shell screencap -p /data/local/tmp/screen.png"
 - Default Format:
   "%time %tag %priority %message"`,
	"helpTrigger":           "regex to trigger a COMMAND.",
	"helpCommand":           "COMMAND will be executed on regex mathed. In COMMAND, you can use parsed logcat as shell variables. ex) `\\${message}` You need to escape a dollar mark.",
	"helpEncode":            "output character encode.",
	"helpToCsv":             "output to CSV format. double-quote will be escaped.",
	"msgUnavailableKeyword": "error: %s is not available. Please check `Availavle Keyword:` on help.",
	"msgDuplicatedKeyword":  "error: %s or %s is duplicated.",
}

var (
	// UTF8 represents encode `utf-8`
	UTF8 = "utf-8"
)
