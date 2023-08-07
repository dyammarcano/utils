package generate

import (
	"os"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	entity := Entity{
		KeyPath:    "test.key",
		Name:       "John Doe",
		Comment:    "Test",
		Email:      "john.doe@mail.com",
		Passphrase: []byte("123456"),
	}

	// delete file after test
	defer func() {
		if err := os.Remove(entity.KeyPath); err != nil {
			t.Error(err)
		}
	}()

	if err := GenerateKey(entity); err != nil {
		t.Error(err)
	}
}
