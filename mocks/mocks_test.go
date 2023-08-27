package mocks

import (
	"os"
	"reflect"
	"testing"
)

func TestNewMockFile(t *testing.T) {
	data := []byte(`Hello, world!
This is a test file.`)

	f := NewMockFile(data)
	b := make([]byte, len(data))

	if _, err := f.Read(b); err != nil {
		t.Error(err)
	}

	if reflect.DeepEqual(b, data) == false {
		t.Error("content is not equal")
	}
}

func TestLorem_LoremTestFile(t *testing.T) {
	testfile, err, cleanup := NewLoremTestFile(t, 10)
	if err != nil {
		t.Error(err)
		return
	}
	defer cleanup()

	if testfile.File == nil {
		t.Error("File is nil")
		return
	}
}

func TestCreateTempDirCleanUp(t *testing.T) {
	dir, err, cleanup := CreateTempDirCleanUp(t)
	if err != nil {
		t.Error(err)
		return
	}

	cleanup()

	if _, err := os.Stat(dir); os.IsExist(err) {
		t.Error("Directory still exists")
	}
}
