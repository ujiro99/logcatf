# logcatf
A command to format Android Logcat.

## Usage

```bash
$ adb logcat | logcatf "%Time, %Message"
```

Available keyword:

    %Time, %Tag, %Priority, %Pid, %Tid, %Message

Example:

    "%Time %Message"
    "%Time, %Tag, %Priority, %Pid, %Tid, %Message"

Default:

    "%Time %Tag %Priority %Message"

## Install

To install, use `go get`:

```bash
$ go get -d github.com/Ujiro99/logcatf
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

