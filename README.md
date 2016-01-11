# Logcatf  [![Build Status](https://travis-ci.org/ujiro99/logcatf.svg?branch=master)](https://travis-ci.org/ujiro99/logcatf)  [![Coverage Status](https://coveralls.io/repos/ujiro99/logcatf/badge.svg?branch=master&service=github)](https://coveralls.io/github/ujiro99/logcatf?branch=master)

A Command line tool for format Android Logcat.

![ScreenShot](./screenshot.png?raw=true "Optional Title")

## Examples

```bash
# show time, pid and message formatted.
$ adb logcat -v time | logcatf "%time %4i %message"

# output to csv format.
$ adb logcat -v threadtime | logcatf "%t %a %p %i %I %m" --to-csv > logcat.csv

# get screencap on Exception
$ adb logcat -v time | logcatf -o "MY_APP.*Error" -c "adb shell screencap -p /sdcard/a.png"
```

## Basic Usage

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

Other Flags:

|  Flag | description              |
|------:|-------------------------:|
|   %a  | left-align               |
|  %8a  | min-width 8, right-align |
| %-8a  | min-width 8, left-align  |
| %8.8a | width 8, right-align      |


Default Format:

    "%t %p %a: %m"

## Options

### execute commands

You can execute other commands when a keyword matched to Logcat.

    -o, --on=ON              regex to trigger a COMMAND.
    -c, --command=COMMAND    COMMAND will be executed on regex matched.
                             In COMMAND, you can use parsed logcat as shell variables.

    
* on Linux, You need to escape a dollar sign.

```
ex) -o "MY_APP.*Error" -c "echo \${message} > error.log"  # linux, mac
    -o "MY_APP.*Error" -c "echo %message% > error.log"    # Windows
```

* Command's stdout is redirected to stderr of logcatf.
* You can use multiple -o / -c pairs.
 

### Output CSV

output to CSV format.

        --to-csv             output to CSV format. double-quote will be escaped.
        --encode=ENCODE      output character encode. { utf-8 | shift-jis }

### Color

specify output Color.

```
    --color            enable ANSI color. ( defalult = false )
    --color-v=COLOR    - color for verbose.
    --color-d=COLOR    - color for debug.
    --color-i=COLOR    - color for information.
    --color-w=COLOR    - color for warning.
    --color-e=COLOR    - color for error.
    --color-f=COLOR    - color for fatal.
```

* This function uses [mitchellh/colorstring](https://github.com/mitchellh/colorstring).
* In format string, you can use color tags. 

```
ex) $ adb logcat | logcatf "%t [invert] %a [reset] [_white_] %m" --color --color-i "cyan"
```

Available Color Tags:

|Foreground|             |Background |                  |
|:--------|:-------------|:----------|:-----------------|
|black    | dark_gray    |\_black\_  | \_dark_gray\_    |
|red      | light_red    |\_red\_    | \_light_red\_    |
|green    | light_green  |\_green\_  | \_light_green\_  |
|yellow   | light_yellow |\_yellow\_ | \_light_yellow\_ |
|blue     | light_blue   |\_blue\_   | \_light_blue\_   |
|magenta  | light_magenta|\_magenta\_| \_light_magenta\_|
|cyan     | light_cyan   |\_cyan\_   | \_light_cyan\_   |
|white    | light_gray   |\_white\_  | \_light_gray\_   |
|default  |              |\_default\_|                  |

|Attributes|           |
:----------|:----------|
|bold      |blink_slow |
|dim       |blink_fast |
|underline |invert     |

| Reset      |
|:-----------|
| reset      |
| reset_bold |


## Install

You can get binary from github release page.

[-> Release Page](https://github.com/ujiro99/logcatf/releases)

or, use `go get`:

```bash
$ go get github.com/ujiro99/logcatf
```

## Contribution

1. Fork ([https://github.com/ujiro99/logcatf/fork](https://github.com/ujiro99/logcatf/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[ujiro99](https://github.com/ujiro99)

