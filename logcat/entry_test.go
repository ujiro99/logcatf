package logcat

import "testing"

var (
	itemEmpty = &Entry{}

	item = &Entry{
		"time":    "12-28 19:01:14.073",
		"message": "this is test message",
	}

	itemFull = &Entry{
		"time":     "12-28 19:01:14.073",
		"pid":      "1",
		"tid":      "2",
		"priority": "W",
		"message":  "This is test message",
		"tag":      "Tag",
	}

	itemRaw = &Entry{
		"message": "this is test message",
		"format":  "raw",
	}
)

func TestEntry_Keys_Empty(t *testing.T) {
	keys := itemEmpty.Keys()
	if len(keys) != 0 {
		t.Errorf("itemEmpty must has length 0")
	}
}

func TestEntry_Keys(t *testing.T) {
	keys := item.Keys()
	if len(keys) != 2 {
		t.Error("item must has length 2")
	}
	if keys[0] != "time" {
		t.Error("keys[0] must be time")
	}
	if keys[1] != "message" {
		t.Error("keys[1] must be message")

	}
}

func TestEntry_Keys_Full(t *testing.T) {
	keys := itemFull.Keys()
	if len(keys) != 6 {
		t.Error("itemEmpty must has length 6")
	}
	for i, k := range keys {
		if k != allKeys[i] {
			t.Errorf("invalid key order: %s", k)
		}
	}
}

func TestEntry_Values_Empty(t *testing.T) {
	values := itemEmpty.Values()
	if len(values) != 0 {
		t.Errorf("itemEmpty must has length 0")
	}
}

func TestEntry_Values(t *testing.T) {
	values := item.Values()
	if len(values) != 2 {
		t.Errorf("item must has length 2")
	}
}

func TestEntry_Values_Full(t *testing.T) {
	values := itemFull.Values()
	if len(values) != 6 {
		t.Errorf("item must has length 6")
	}
}

func TestEntry_Values_Raw(t *testing.T) {
	values := itemRaw.Values()
	if len(values) != 1 {
		t.Errorf("item must has length 1")
	}
}

func TestEntry_Format(t *testing.T) {
	if itemRaw.Format() != "raw" {
		t.Errorf("item format must be raw")
	}
}
