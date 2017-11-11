package main

import (
	"bytes"
	"os"
)

func WriteContentToFile(userid string, entries []ApacheLogEntry, contentIndex map[string][]byte) {
	f, _ := os.OpenFile("/tmp/"+userid+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	var buffer bytes.Buffer
	newline := '\n'
	for _, entry := range entries {
		buffer.Write(contentIndex[entry.String()])
		buffer.WriteRune(newline)
	}
	f.Write(buffer.Bytes())
}
