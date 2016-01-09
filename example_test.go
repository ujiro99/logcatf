package main

import (
	"fmt"
	"os"
)

func ExampleLogcatItem_Format_time_message() {
	item := &LogcatItem{
		"time":    "12-28 19:01:14.073",
		"message": "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	format := "%t, %m"
	formatter = &defaultFormatter{format: &format}
	fmt.Println(formatter.Format(item))
	// Output:
	// 12-28 19:01:14.073, at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}

func ExampleLogcatItem_Format_escapedCharactor() {
	item := &LogcatItem{
		"time":     "12-28 19:01:14.073",
		"tag":      "GLSUser",
		"priority": "W",
		"message":  "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	format := "%t\t%a\t%p\n%m"
	formatter = &defaultFormatter{format: &format}
	fmt.Println(formatter.Format(item))
	// Output:
	// 12-28 19:01:14.073	GLSUser	W
	// at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}

func ExampleLogcatItem_Format_priority_missing() {
	item := &LogcatItem{
		"time":    "12-28 19:01:14.073",
		"message": "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	format := "%p, %t, %m"
	formatter = &defaultFormatter{format: &format}
	fmt.Println(formatter.Format(item))
	// Output:
	// , 12-28 19:01:14.073, at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}

func ExampleLogcatItem_Format_all() {
	item := &LogcatItem{
		"time":     "12-28 19:01:14.073",
		"pid":      "1836",
		"tid":      "2720",
		"priority": "W",
		"tag":      "GLSUser",
		"message":  "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	format := "%t, %a, %i, %I, %p, %m"
	formatter = &defaultFormatter{format: &format}
	fmt.Println(formatter.Format(item))
	// Output:
	// 12-28 19:01:14.073, GLSUser, 1836, 2720, W, at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}

func ExampleLogcatItem_Format_remainFrags() {
	item := &LogcatItem{
		"time":     "12-28 19:01:14.073",
		"pid":      "1",
		"tid":      "2",
		"priority": "W",
		"tag":      "Tag",
	}
	format := "%t, %8a, %4i, %-4I,%2p"
	formatter = &defaultFormatter{format: &format}
	fmt.Println(formatter.Format(item))
	// Output:
	// 12-28 19:01:14.073,      Tag,    1, 2   , W
}

func ExampleLogcatItem_Format_toCsv() {
	item := &LogcatItem{
		"pid":     "1",
		"message": "aaa\"bbb\"ccc",
	}

	cli := newCli()
	cli.outStream = os.Stdout
	format := "%m %i"
	f := csvFormatter{format: &format}
	fmt.Println(f.Format(item))
	// Output:
	// "aaa""bbb""ccc",1
}
