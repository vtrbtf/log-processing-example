package main

import (
	"encoding/json"
	"flag"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"

	gonx "github.com/satyrius/gonx"
)

const format = "$remote_addr - - [$time_local] \"$request\" $http_status $bytes \"-\" \"userid=$userid\""

var logFile1 string
var logFile2 string
var logFile3 string
var logFile4 string
var counter int
var wg sync.WaitGroup

func init() {
	counter = 0
	flag.StringVar(&logFile1, "log1", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
	flag.StringVar(&logFile2, "log2", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
	flag.StringVar(&logFile3, "log3", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
	flag.StringVar(&logFile4, "log4", "dummy", "Log file name to read. Read from STDIN if file name is '-'")
}

type ApacheLogEntry struct {
	Timestamp string `json:"t"`
	Line      int    `json:"l"`
}

func parseLog(logFile string) {

	indexes := make(map[string][]ApacheLogEntry)
	counter++
	var logReader io.Reader

	file, err := os.Open(logFile)
	if err != nil {
		panic(err)
	}
	logReader = file
	defer file.Close()

	reader := gonx.NewReader(logReader, format)
	line := 0
	for {
		line = line + 1
		rec, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		userid, _ := rec.Field("userid")
		timestamp, _ := rec.Field("time_local")
		if _, ok := indexes[userid]; ok {
			indexes[userid] = append(indexes[userid], ApacheLogEntry{Timestamp: timestamp, Line: line})
		} else {
			indexes[userid] = append(make([]ApacheLogEntry, 0), ApacheLogEntry{Timestamp: timestamp, Line: line})
		}
	}

	jsonString, _ := json.Marshal(indexes)
	jsonFile, err := os.Create("./" + strconv.Itoa(rand.Intn(100)) + "results.json")

	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	jsonFile.Write(jsonString)
	defer wg.Done()
}

func main() {
	flag.Parse()

	wg.Add(1)
	go parseLog(logFile1)
	wg.Add(1)
	go parseLog(logFile2)
	wg.Add(1)
	go parseLog(logFile3)
	wg.Add(1)
	go parseLog(logFile4)
	wg.Wait()
}
