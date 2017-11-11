package main

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestParseLogFiles(t *testing.T) {
	inputFile1 := "/tmp/abc.log"
	file1line1 := append([]byte(`240.62.40.107 - - [14/Nov/2017:11:55:06 +0000] "GET /wp-content HTTP/1.0" 200 5066 "-" "userid=2c5c8428-9fdb-4b2f-a32c-3d4e447acefe"`), '\n')
	file1line2 := append([]byte(`240.61.40.107 - - [14/Nov/2017:11:55:12 +0000] "GET /wp-content HTTP/1.0" 200 5066 "-" "userid=3c5c8429-a32c-9fdb-4b2f-3d4e447acefe"`), '\n')
	file1line3 := append([]byte(`242.62.40.110 - - [14/Nov/2017:08:55:14 -0300] "GET /wp-content HTTP/1.0" 200 5066 "-" "userid=2c5c8428-9fdb-4b2f-a32c-3d4e447acefe"`), '\n')

	/*
		1510660506
		1510660512
		1510660514
	*/

	inputFile2 := "/tmp/def.log"
	file2line1 := append([]byte(`240.62.40.107 - - [14/Jan/2021:21:41:46 +0000] "GET /wp-content HTTP/1.0" 200 5066 "-" "userid=2c5c8428-9fdb-4b2f-a32c-3d4e447acefe"`), '\n')
	file2line2 := append([]byte(`240.61.40.107 - - [14/Jan/2021:21:41:48 +0000] "GET /wp-content HTTP/1.0" 200 5066 "-" "userid=3c5c8429-a32c-9fdb-4b2f-3d4e447acefe"`), '\n')
	file2line3 := append([]byte(`242.62.40.110 - - [14/Jan/2021:19:41:50 -0200] "GET /wp-content HTTP/1.0" 200 5066 "-" "userid=d7bf80f0-66ba-478f-ac45-b8fd8481b2db"`), '\n')

	/*
		1610660506
		1610660508
		1610660510
	*/

	ioutil.WriteFile(inputFile1, append(file1line1, append(file1line2, file1line3...)...), 0644)
	ioutil.WriteFile(inputFile2, append(file2line1, append(file2line2, file2line3...)...), 0644)
	//defer os.Remove(inputFile1)
	//defer os.Remove(inputFile2)

	index := ParseLogs(LogOption{inputFile1, inputFile2})

	if len(index) != 2 {
		t.Errorf("Index size is not correct. Expected is 2 but got %d", len(index))
	}

	if (len(index[0]) != 2 && len(index[1]) != 3) && (len(index[0]) != 3 && len(index[1]) != 2) {
		t.Errorf("First/Second logIndex size is not correct")
	}

	entryAFile1 := []ApacheLogEntry{
		ApacheLogEntry{Line: 1, File: &inputFile1, Timestamp: 1510660506},
		ApacheLogEntry{Line: 3, File: &inputFile1, Timestamp: 1510660514},
	}
	useridCheck("2c5c8428-9fdb-4b2f-a32c-3d4e447acefe", entryAFile1, index, t)

	entryBFile1 := []ApacheLogEntry{
		ApacheLogEntry{Line: 2, File: &inputFile1, Timestamp: 1510660512},
	}
	useridCheck("3c5c8429-a32c-9fdb-4b2f-3d4e447acefe", entryBFile1, index, t)

	entryAFile2 := []ApacheLogEntry{
		ApacheLogEntry{Line: 1, File: &inputFile2, Timestamp: 1610660506},
	}
	useridCheck("3c5c8429-a32c-9fdb-4b2f-3d4e447acefe", entryAFile2, index, t)

	entryBFile2 := []ApacheLogEntry{
		ApacheLogEntry{Line: 2, File: &inputFile2, Timestamp: 1610660508},
	}
	useridCheck("2c5c8428-9fdb-4b2f-a32c-3d4e447acefe", entryBFile2, index, t)

	entryCFile2 := []ApacheLogEntry{
		ApacheLogEntry{Line: 3, File: &inputFile2, Timestamp: 1610660510},
	}
	useridCheck("d7bf80f0-66ba-478f-ac45-b8fd8481b2db", entryCFile2, index, t)
}

func useridCheck(userid string, expectedEntries []ApacheLogEntry, index []LogIndex, t *testing.T) {
	var selected LogIndex
	if len(index[0][userid]) > 0 && *index[0][userid][0].File == *expectedEntries[0].File {
		selected = index[0]
	} else if len(index[1][userid]) > 0 && *index[1][userid][0].File == *expectedEntries[0].File {
		selected = index[1]
	} else {
		t.Error("Could not find proper logIndex")
	}

	if len(selected[userid]) != len(expectedEntries) {
		t.Errorf("userid entry (%s) size is not correct. Expected is %d but got %d", userid, len(expectedEntries), len(selected[userid]))
	}

	reflect.DeepEqual(selected[userid], expectedEntries)
}
