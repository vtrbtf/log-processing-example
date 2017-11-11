package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestOpenRequiredFiles(t *testing.T) {
	file1, _ := ioutil.TempFile(os.TempDir(), "log_file_index_test_1")
	filename1 := file1.Name()
	defer os.Remove(filename1)

	file2, _ := ioutil.TempFile(os.TempDir(), "log_file_index_test_2")
	filename2 := file2.Name()
	defer os.Remove(filename2)

	files := *NewLogFileIndex(LogOption{filename1, filename2})

	if len(files.index) != 2 {
		t.Errorf("File index len is not correct. Expected 2 but got %d", len(files.index))
	}
}
