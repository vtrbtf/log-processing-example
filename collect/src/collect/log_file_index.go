package main

import "os"

type LogFileIndex struct {
	index map[string]*os.File
}

func NewLogFileIndex(logfiles LogOption) *LogFileIndex {
	logFileIndex := LogFileIndex{index: make(map[string]*os.File)}
	for _, log := range logfiles {
		file, err := os.Open(log)
		if err != nil {
			panic(err)
		}
		logFileIndex.index[log] = file
	}
	return &logFileIndex
}

func (i *LogFileIndex) Close() {
	for _, file := range i.index {
		file.Close()
	}
}
