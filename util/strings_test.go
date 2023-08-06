package util

import "testing"

func TestGenerateRandomString2(t *testing.T) {
	str := GenerateRandomString(10, -1)
	if len(str) != 10 {
		t.Errorf("GenerateRandomString(10) = %s; want length 10", str)
	}
}
