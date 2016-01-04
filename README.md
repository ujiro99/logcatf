# Logcatf  [![Build Status](https://travis-ci.org/ujiro99/logcatf.svg?branch=master)](https://travis-ci.org/ujiro99/logcatf)  [![Coverage Status](https://coveralls.io/repos/ujiro99/logcatf/badge.svg?branch=master&service=github)](https://coveralls.io/github/ujiro99/logcatf?branch=master)

A Command line tool to format Android Logcat.


## Usage

```bash
$ adb logcat -v time | logcatf "%t, %m"
```

Available Format:

| format | long ver. |
|:------:|:---------:|
|   %t   | %time     |
|   %a   | %tag      |
|   %p   | %priority |
|   %i   | %pid      |
|   %I   | %tid      |
|   %m   | %message  |


Default Format:

    "%time %tag %priority %message"

## Other options

You can execute a command when a keyword  matched to Logcat.

    -o, --on=ON            regex to trigger a COMMAND.
    -c, --command=COMMAND  COMMAND will be executed on regex mathed. 
                           In COMMAND, you can use parsed logcat as shell variables. 
                         
    ex) -o "MY_APP.*Error" -c "echo \${message} > error.log" 
    
    * You need to escape a dollar sign.


## Examples

```bash
# show only time and message.
$ adb logcat -v time | logcatf "%time %message"

# output to csv format.
$ adb logcat -v threadtime | logcatf "%t, %a, %p, %i, %I, %m" > logcat.csv

# get screencap on Exception
$ adb logcat -v threadtime | "%t %m" -o "MY_APP.*Error" -c "adb shell screencap -p /sdcard/a.png"
```


## Install

To install, use `go get`:

```bash
$ go get github.com/Ujiro99/logcatf
```

## Contribution

1. Fork ([https://github.com/Ujiro99/logcatf/fork](https://github.com/Ujiro99/logcatf/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[Ujiro99](https://github.com/Ujiro99)

