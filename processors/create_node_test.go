package processors

import (
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"testing"
)

func TestReadFiles(t *testing.T) {

	l, err := zap.NewDevelopment()
	if err != nil {
		panic("Failed to initialize Zap logger")
	}
	logger = l
	defer logger.Sync() // Flushes buffer, if any

	SetupProcessing("../output_temp", logger)
	currDir, err := os.Getwd()
	if err != nil {
		t.Errorf("ReadFiles error getting current directory: %v", err)
	}
	t.Log("Current Directory", currDir)

	// Call the function being tested
	err = ReadFiles("../output_temp/v1.21.0")
	if err != nil {
		t.Errorf("ReadFiles returned an error: %v", err)
	}

	// Add more positive and negative test cases as needed

	// Test case: Directory does not exist
	err = ReadFiles("nonexistent")
	if err == nil {
		t.Error("ReadFiles(nonexistent) did not return an error")
	}

	// Test case: Directory is empty
	emptyDir, err := ioutil.TempDir("", "empty")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(emptyDir)

	err = ReadFiles(emptyDir)
	if err != nil {
		t.Errorf("ReadFiles(%s) returned an error: %v", emptyDir, err)
	}
}
