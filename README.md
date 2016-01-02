# Logcatf
A command to format Android Logcat.

## Usage

```bash
$ adb logcat -v time | logcatf "%time, %message"
```

Available Keyword:

    %time, %tag, %priority, %pid, %tid, %message

Example:

```bash
# show only time and message.
$ adb logcat -v time | logcatf "%time %message"

# output to csv format.
$ adb logcat -v threadtime | logcatf "%time, %tag, %priority, %pid, %tid, %message" > logcat.csv
```

Default Format:

    "%time %tag %priority %message"

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

