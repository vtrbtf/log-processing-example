package main

import (
	"flag"
	"fmt"
	"sort"
	"sync"

	pb "gopkg.in/cheggaaa/pb.v1"
)

var logfiles LogOption

func init() {
	flag.Var(&logfiles, "logfile", "Log file name to read.")
	flag.Parse()
}

func main() {
	parseIndex := ParseLogs(logfiles)
	useridIndex := *NewuseridIndex(parseIndex)
	fileIndex := *NewLogFileIndex(logfiles)
	defer fileIndex.Close()

	var wg sync.WaitGroup

	fmt.Printf("--- Building log file for each userid (total: %d)\n", len(useridIndex))
	bar := pb.StartNew(len(useridIndex))

	for userid, entries := range useridIndex {
		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp < entries[j].Timestamp
		})

		content := LoadContent(userid, parseIndex, fileIndex)
		wg.Add(1)
		go func(userid string, entries []ApacheLogEntry, content map[string][]byte) {
			defer wg.Done()
			defer bar.Increment()
			WriteContentToFile(userid, entries, content)
		}(userid, entries, content)
	}

	wg.Wait()
	bar.FinishPrint("All Done ;)")
}
