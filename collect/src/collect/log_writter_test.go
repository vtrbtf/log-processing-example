package main

import (
	"bufio"
	"os"
	"testing"
)

func TestWriteResultLogFile(t *testing.T) {
	userid := "user123"
	inputFile := "/tmp/log/abc.log"
	entry1 := ApacheLogEntry{Line: 1, Timestamp: 88888, File: &inputFile}
	entry2 := ApacheLogEntry{Line: 2, Timestamp: 88880, File: &inputFile}
	logIndex := []ApacheLogEntry{entry1, entry2}
	content := make(map[string][]byte)
	content[entry1.String()] = []byte("line1: abc")
	content[entry2.String()] = []byte("line2: abc")
	WriteContentToFile(userid, logIndex, content)

	file, err := os.Open("/tmp/user123.log")
	defer os.Remove(file.Name())

	if err != nil {
		t.Error("Result file was not created")
	}

	scanner := bufio.NewScanner(file)

	if scanner.Scan(); scanner.Text() != "line1: abc" {
		t.Error("First line in result file is not correct")
	}

	if scanner.Scan(); scanner.Text() != "line2: abc" {
		t.Error("Second line in result file is not correct")
	}
}
