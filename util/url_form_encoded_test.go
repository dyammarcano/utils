package util

import "testing"

func TestUrlFormEncodedToObject(t *testing.T) {
	type testStruct struct {
		Name string
		Age  int
	}

	test := testStruct{
		Name: "John",
		Age:  30,
	}

	str := "Name=John&Age=30"

	err := UrlFormEncodedToObject(str, &test)
	if err != nil {
		t.Errorf("UrlFormEncodedToObject(StructToStringFormat(test)) = %v; want Name=John&Age=30", err)
	}

	if test.Name != "John" {
		t.Errorf("UrlFormEncodedToObject(StructToStringFormat(test)) = %s; want Name=John&Age=30", test.Name)
	}
}
