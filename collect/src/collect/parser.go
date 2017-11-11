package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"
)

const format = "$remote_addr - - [$time_local] \"$request\" $http_status $bytes \"-\" \"userid=$userid\""

//ApacheLogEntry is a entry
type ApacheLogEntry struct {
	Timestamp int64
	File      *string
	Line      int64
}

func (i *ApacheLogEntry) String() string {
	var buffer bytes.Buffer
	buffer.WriteString(*i.File)
	buffer.WriteString(";")
	buffer.WriteString(strconv.FormatInt(i.Line, 10))
	return buffer.String()
}

//LogIndex of userid to log entry
type LogIndex map[string][]ApacheLogEntry

func openReader(logFile string) *os.File {
	logReader, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}
	return logReader
}

func parseLog(logFile string) LogIndex {
	file, err := os.Open(logFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	r := regexp.MustCompile(`.*- - \[([\+\w:\s\d-\/]+)\].*userid\=([\w-]+)`)

	indexes := make(LogIndex)
	var linec int64
	for scanner.Scan() {
		linec++
		line := scanner.Text()
		rs := r.FindStringSubmatch(line)
		userid := rs[2]
		rawtimestamp := rs[1]
		timestamp, _ := time.Parse("02/Jan/2006:15:04:05 -0700", rawtimestamp)
		log := ApacheLogEntry{Timestamp: timestamp.Unix(), Line: linec, File: &logFile}

		indexes[userid] = append(indexes[userid], log)
	}
	return indexes
}

func loadLogFiles(logfiles LogOption, bar *pb.ProgressBar) chan LogIndex {
	var wg sync.WaitGroup
	logIndexChannel := make(chan LogIndex, len(logfiles))
	for _, logfile := range logfiles {
		wg.Add(1)
		go func(log string, c chan LogIndex, wg *sync.WaitGroup) {
			defer wg.Done()
			defer bar.Increment()
			c <- parseLog(log)
		}(logfile, logIndexChannel, &wg)
	}
	wg.Wait()
	close(logIndexChannel)
	return logIndexChannel
}

//ParseLogs parses a access_log file
func ParseLogs(logfiles LogOption) []LogIndex {
	fmt.Println("--- Parsing log files (async)")
	for _, log := range logfiles {
		fmt.Println("- " + log)
	}

	bar := pb.StartNew(len(logfiles))
	logs := loadLogFiles(logfiles, bar)
	bar.FinishPrint("--- All logs were parsed")

	var logIndex []LogIndex
	for index := 0; index < len(logfiles); index++ {
		logIndex = append(logIndex, <-logs)
	}

	return logIndex
}
