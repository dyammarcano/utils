package mocks

import (
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
	testfile, err := NewLoremTestFile(t, 10)
	if err != nil {
		t.Error(err)
		return
	}

	defer testfile.CleanUp()
}
