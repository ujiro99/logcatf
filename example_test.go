package main

import (
	"fmt"
)

func ExampleLogcatItem_Format_time_message() {
	item := &LogcatItem{
		"time":    "12-28 19:01:14.073",
		"message": "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	fmt.Println(item.Format("%time, %message"))
	// Output:
	// 12-28 19:01:14.073, at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}

func ExampleLogcatItem_Format_priority_missing() {
	item := &LogcatItem{
		"time":    "12-28 19:01:14.073",
		"message": "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	fmt.Println(item.Format("%priority, %time, %message"))
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
	fmt.Println(item.Format("%time, %tag, %pid, %tid, %priority, %message"))
	// Output:
	// 12-28 19:01:14.073, GLSUser, 1836, 2720, W, at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}
