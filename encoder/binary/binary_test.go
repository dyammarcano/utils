package binary

import (
	"fmt"
	"os"
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

func TestMarshal(t *testing.T) {
	data := MyStruct{
		Field1: 42,
		Field2: "Hello, World!",
		OtherStruct: OtherStruct{
			FieldA: 3.14159,
			FieldB: true,
		},
	}

	serializedData, err := Marshal(data)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(fmt.Sprintf("%v", serializedData))

	var deserializedData MyStruct
	if err := Unmarshal(serializedData, &deserializedData); err != nil {
		t.Error(err)
		return
	}

	if reflect.DeepEqual(data, deserializedData) == false {
		t.Errorf("Expected %v, got %v", data, deserializedData)
	}
}

func TestUnmarshalToFile(t *testing.T) {
	data := MyStruct{
		Field1: 42,
		Field2: "Hello, World!",
		OtherStruct: OtherStruct{
			FieldA: 3.14159,
			FieldB: true,
		},
	}

	serializedData, err := Marshal(data)
	if err != nil {
		t.Error(err)
		return
	}

	// write to file
	f, err := os.OpenFile("test.bin", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = f.Write(serializedData)
	if err != nil {
		t.Error(err)
		return
	}

	// read from file
	f, err = os.Open("test.bin")
	if err != nil {
		t.Error(err)
		return
	}

	var deserializedData MyStruct
	if err := Unmarshal(serializedData, &deserializedData); err != nil {
		t.Error(err)
		return
	}

	if reflect.DeepEqual(data, deserializedData) == false {
		t.Errorf("Expected %v, got %v", data, deserializedData)
	}
}
