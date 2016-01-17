package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/mitchellh/colorstring"
)

const (
	defaultColorsKey   = 0
	defaultColorsValue = 1
	prefix             = "LOGCATF_"
)

var defaultColors = [][]string{
	// Default foreground/background colors
	{"default", "39"},
	{"_default_", "49"},

	// Foreground colors
	{"black", "30"},
	{"red", "31"},
	{"green", "32"},
	{"yellow", "33"},
	{"blue", "34"},
	{"magenta", "35"},
	{"cyan", "36"},
	{"light_gray", "37"},
	{"dark_gray", "90"},
	{"light_red", "91"},
	{"light_green", "92"},
	{"light_yellow", "93"},
	{"light_blue", "94"},
	{"light_magenta", "95"},
	{"light_cyan", "96"},
	{"white", "97"},

	// Background colors
	{"_black_", "40"},
	{"_red_", "41"},
	{"_green_", "42"},
	{"_yellow_", "43"},
	{"_blue_", "44"},
	{"_magenta_", "45"},
	{"_cyan_", "46"},
	{"_light_gray_", "47"},
	{"_dark_gray_", "100"},
	{"_light_red_", "101"},
	{"_light_green_", "102"},
	{"_light_yellow_", "103"},
	{"_light_blue_", "104"},
	{"_light_magenta_", "105"},
	{"_light_cyan_", "106"},
	{"_white_", "107"},

	// Attributes
	{"bold", "1"},
	{"dim", "2"},
	{"underline", "4"},
	{"blink_slow", "5"},
	{"blink_fast", "6"},
	{"invert", "7"},
	{"hidden", "8"},

	// Reset to reset everything to their defaults
	{"reset", "0"},
	{"reset_bold", "21"},
}

var colorMap = map[string]string{
	// Default foreground/background colors
	prefix + "default":   "39",
	prefix + "_default_": "49",

	// Foreground colors
	prefix + "black":         "30",
	prefix + "red":           "31",
	prefix + "green":         "32",
	prefix + "yellow":        "33",
	prefix + "blue":          "34",
	prefix + "magenta":       "35",
	prefix + "cyan":          "36",
	prefix + "light_gray":    "37",
	prefix + "dark_gray":     "90",
	prefix + "light_red":     "91",
	prefix + "light_green":   "92",
	prefix + "light_yellow":  "93",
	prefix + "light_blue":    "94",
	prefix + "light_magenta": "95",
	prefix + "light_cyan":    "96",
	prefix + "white":         "97",

	// Background colors
	prefix + "_black_":         "40",
	prefix + "_red_":           "41",
	prefix + "_green_":         "42",
	prefix + "_yellow_":        "43",
	prefix + "_blue_":          "44",
	prefix + "_magenta_":       "45",
	prefix + "_cyan_":          "46",
	prefix + "_light_gray_":    "47",
	prefix + "_dark_gray_":     "100",
	prefix + "_light_red_":     "101",
	prefix + "_light_green_":   "102",
	prefix + "_light_yellow_":  "103",
	prefix + "_light_blue_":    "104",
	prefix + "_light_magenta_": "105",
	prefix + "_light_cyan_":    "106",
	prefix + "_white_":         "107",

	// Attributes
	prefix + "bold":       "1",
	prefix + "dim":        "2",
	prefix + "underline":  "4",
	prefix + "blink_slow": "5",
	prefix + "blink_fast": "6",
	prefix + "invert":     "7",
	prefix + "hidden":     "8",

	// Reset to reset everything to their defaults
	prefix + "reset":      "0",
	prefix + "reset_bold": "21",
}

// priorityColors represents default color of priority.
var priorityColors = map[string]string{
	"V": "default",
	"D": "default",
	"I": "default",
	"W": "yellow",
	"E": "red",
	"F": "red",
}

// Colorizer enables to print strings colored.
type Colorizer struct {
	context colorstring.Colorize
}

// ColorConfig has configuration for colorize.
type ColorConfig map[string]string

// Fprintln prints string according to color settings.
func (f *Colorizer) Fprintln(writer io.Writer, str string, item LogcatItem) {
	color, ok := priorityColors[item["priority"]]
	if ok {
		str = fmt.Sprintf("[%s%s]%s", prefix, color, str)
	}
	fmt.Fprintln(writer, f.context.Color(str))
}

// SetUp color settings.
func (f *Colorizer) Init(colorEnable bool, config ColorConfig) {

	for k, v := range config {
		if v == "" {
			continue
		}
		priorityColors[k] = v
		colorEnable = true
	}

	c := colorstring.Colorize{
		Colors:  colorMap,
		Disable: !colorEnable,
		Reset:   true,
	}
	f.context = c
}

// ReplaceColorCode replaces colorCode in format, to be safety.
func (f *Colorizer) ReplaceColorCode(format string) string {
	for _, pair := range defaultColors {
		oldKey := "[" + pair[defaultColorsKey] + "]"
		newKey := "[" + prefix + pair[defaultColorsKey] + "]"
		format = strings.Replace(format, oldKey, newKey, flagAll)
	}
	return format
}
