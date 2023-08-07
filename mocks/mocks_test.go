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
