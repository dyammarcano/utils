package util

import "testing"

func TestStructToStringFormat(t *testing.T) {
	type testStruct struct {
		Name string
		Age  int
	}

	test := testStruct{
		Name: "John",
		Age:  30,
	}

	result := StructToStringFormat(test)
	if result != "Name=John;Age=30" {
		t.Errorf("StructToStringFormat(test) = %s; want Name=John;Age=30", result)
	}
}
