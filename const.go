package main

const (
	// Name of this command.
	Name string = "logcatf"
	// Version of this command.
	Version string = "0.4.0"
	// DefaultFormat will be used if format wasn't specified.
	DefaultFormat string = "%time [invert] %priority [reset] %tag: %message"
	// UTF8 represents encode `utf-8`
	UTF8 = "utf-8"
	// ShiftJIS represents encode `shift-jis`
	ShiftJIS = "shift-jis"
)

// Message has message strings.
var Message = map[string]string{
	"commandDescription": "A Command line tool for format Android Logcat",
	"helpFormat": `Format of output
 - Available keyword:
     %t (%time), %a (%tag), %p (%priority)
     %i (%pid), %I (%tid), %m (%message)
 - Example:
   adb logcat -v time | logcatf "%t %4i %m"
   adb logcat -v time | logcatf "%t %a %m" --to-csv > logcat.csv
   adb logcat -v | logcatf -o "Exception" -c "adb shell screencap -p /data/local/tmp/s.png"
 - Default Format:
   "%t [invert] %p [reset] %a: %m"`,
	"helpTrigger":           "regex to trigger a COMMAND.",
	"helpCommand":           "COMMAND will be executed on regex mathed. In COMMAND, you can use parsed logcat as environment variables. ex) `\\${message}` ** You need to escape a dollar mark. **",
	"helpEncode":            "output character encode.",
	"helpToCsv":             "output to CSV format. double-quote will be escaped.",
	"helpToColor":           "enable ANSI color. In format, you can use color tags. ex) [blue], [_red_], [bold], [reset], and more.",
	"helpToColorV":          "- color for verbose.",
	"helpToColorD":          "- color for debug.",
	"helpToColorI":          "- color for information.",
	"helpToColorW":          "- color for warning.",
	"helpToColorE":          "- color for error.",
	"helpToColorF":          "- color for fatal.",
	"msgUnavailableKeyword": "error: %s is not available. Please check `Availavle Keyword:` on help.",
	"msgDuplicatedKeyword":  "error: %s or %s is duplicated.",
	"msgCommandNumMismatch": "error: number of on and command flags are mismatch.",
}
