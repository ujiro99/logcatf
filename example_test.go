package main

import (
	"fmt"
)

func ExampleLogcatItem_Format_time_message() {
	item := &LogcatItem{
		"time":    "12-28 19:01:14.073",
		"message": "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	fmt.Println(formatter.Format("%t, %m", item))
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
	fmt.Println(formatter.Format("%t\t%a\t%p\n%m", item))
	// Output:
	// 12-28 19:01:14.073	GLSUser	W
	// at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}

func ExampleLogcatItem_Format_priority_missing() {
	item := &LogcatItem{
		"time":    "12-28 19:01:14.073",
		"message": "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	fmt.Println(formatter.Format("%p, %t, %m", item))
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
	fmt.Println(formatter.Format("%t, %a, %i, %I, %p, %m", item))
	// Output:
	// 12-28 19:01:14.073, GLSUser, 1836, 2720, W, at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}
