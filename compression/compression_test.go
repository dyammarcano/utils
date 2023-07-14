package compression

import "testing"

func Test_GzipCompressStruct(t *testing.T) {
	type test struct {
		Name string
		Age  int
	}

	b, err := GzipCompressStruct(test{
		Name: "John",
		Age:  30,
	})

	if err != nil {
		t.Error("error should be nil")
	}

	if len(b) == 0 {
		t.Error("file is empty")
	}

	u := test{}

	err = GzipUncompressStruct(b, &u)

	if err != nil {
		t.Error("error should be nil")
	}

	if u.Name != "John" {
		t.Error("name should be John")
	}

	if u.Age != 30 {
		t.Error("age should be 30")
	}
}
