package binary

import (
	"github.com/dyammarcano/utils/mocks"
	"path/filepath"
	"reflect"
	"testing"
)

type (
	MyStruct struct {
		Field1      int32
		Field2      string
		OtherStruct OtherStruct
	}

	OtherStruct struct {
		FieldA float64
		FieldB bool
	}
)

var mockData MyStruct

func init() {
	mockData = MyStruct{
		Field1: 42,
		Field2: "Hello, World!",
		OtherStruct: OtherStruct{
			FieldA: 3.14159,
			FieldB: true,
		},
	}
}

func TestMarshal(t *testing.T) {
	serializedData, err := Marshal(mockData)
	if err != nil {
		t.Error(err)
		return
	}

	var deserializedData MyStruct
	if err := Unmarshal(serializedData, &deserializedData); err != nil {
		t.Error(err)
		return
	}

	if reflect.DeepEqual(mockData, deserializedData) == false {
		t.Errorf("Expected %v, got %v", mockData, deserializedData)
	}
}

func TestUnmarshalToFile(t *testing.T) {
	serializedData, err := Marshal(mockData)
	if err != nil {
		t.Fatal("Couldn't serialize data:", err)
	}

	dir, err := mocks.CreateTempDir(t)
	if err != nil {
		return
	}

	defer mocks.CleanUpTmpDir(t, dir)

	testfile := filepath.Join(dir, "data.bin")

	if err := mocks.WriteToTestFile(t, testfile, serializedData); err != nil {
		return
	}

	readData, err := mocks.ReadFromTestFile(t, testfile)
	if err != nil {
		return
	}

	var deserializedData MyStruct
	if err := Unmarshal(readData, &deserializedData); err != nil {
		t.Fatal("Failed to unmarshal to struct:", err)
	}

	if reflect.DeepEqual(mockData, deserializedData) == false {
		t.Errorf("Expected %v, got %v", mockData, deserializedData)
	}
}
