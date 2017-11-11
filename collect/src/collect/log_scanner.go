package main

import (
	"bufio"
)

func LoadContent(userid string, parseIndex []LogIndex, fileIndex LogFileIndex) map[string][]byte {
	contentIndex := make(map[string][]byte)

	for _, parsedLog := range parseIndex {
		currentIndex, ispresent := parsedLog[userid]
		if ispresent && len(currentIndex) > 0 {
			filename := *currentIndex[0].File
			file := fileIndex.index[filename]
			sc := bufio.NewScanner(file)
			var lastLine int64
			for _, entry := range currentIndex {
				var lineContent []byte
				for sc.Scan() {
					lastLine++
					if entry.Line-lastLine == 0 {
						lineContent = []byte(sc.Text())
						break
					}
				}
				contentIndex[entry.String()] = lineContent
			}
			file.Seek(0, 0)
		}
	}

	return contentIndex
}
