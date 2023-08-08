package hash

import "testing"

func TestGenerateFixedHash(t *testing.T) {
	hash := GenerateFixedHash("test")
	if hash != "9F86D-08188-4C7D-659" {
		t.Errorf("GenerateFixedHash failed")
	}
}
