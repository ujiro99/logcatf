package main

import (
	"fmt"
)

func ExampleLogcatItem_Format_time_message() {
	item := &LogcatItem{
		"Time":    "12-28 19:01:14.073",
		"Message": "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	fmt.Println(item.Format("%Time, %Message"))
	// Output:
	// 12-28 19:01:14.073, at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}

func ExampleLogcatItem_Format_priority_missing() {
	item := &LogcatItem{
		"Time":    "12-28 19:01:14.073",
		"Message": "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	fmt.Println(item.Format("%Priority, %Time, %Message"))
	// Output:
	// , 12-28 19:01:14.073, at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}

func ExampleLogcatItem_Format_all() {
	item := &LogcatItem{
		"Time":     "12-28 19:01:14.073",
		"Pid":      "1836",
		"Tid":      "2720",
		"Priority": "W",
		"Tag":      "GLSUser",
		"Message":  "at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)",
	}
	fmt.Println(item.Format("%Time, %Tag, %Pid, %Tid, %Priority, %Message"))
	// Output:
	// 12-28 19:01:14.073, GLSUser, 1836, 2720, W, at com.google.android.gms.auth.be.appcert.b.a(SourceFile:43)
}
