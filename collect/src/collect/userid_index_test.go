package main

import "testing"

func TestNewuseridIndex(t *testing.T) {
	firstIndex := make(LogIndex)
	firstIndex["abc123"] = []ApacheLogEntry{ApacheLogEntry{Line: 1}, ApacheLogEntry{Line: 2}}
	firstIndex["abc456"] = []ApacheLogEntry{ApacheLogEntry{Line: 3}}

	secondIndex := make(LogIndex)
	secondIndex["qqq123"] = []ApacheLogEntry{ApacheLogEntry{Line: 1}, ApacheLogEntry{Line: 2}}
	secondIndex["abc456"] = []ApacheLogEntry{ApacheLogEntry{Line: 6}}

	logs := []LogIndex{firstIndex, secondIndex}
	useridIndex := *NewuseridIndex(logs)

	if val, ok := useridIndex["abc123"]; !ok || len(val) != 2 {
		t.Errorf("Entry abc123 is not correct. Expected 3 but got %d", len(val))
	}

	if val, ok := useridIndex["abc456"]; !ok || len(val) != 2 {
		t.Errorf("Entry abc456 is not correct. Expected 3 but got %d", len(val))
	}

	if val, ok := useridIndex["qqq123"]; !ok || len(val) != 2 {
		t.Errorf("Entry qqq123 is not correct. Expected 3 but got %d", len(val))
	}
}
