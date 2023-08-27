package auto_crypto

import "testing"

func TestEncrypt(t *testing.T) {
	str := "Hello World"

	encrypted, err := AutoEncrypt(str)
	if err != nil {
		t.Errorf("Encrypt() error = %v", err)
		return
	}

	decrypted, err := AutoDecrypt(encrypted)
	if err != nil {
		t.Errorf("Encrypt() error = %v", err)
		return
	}

	if decrypted != str {
		t.Errorf("Encrypt() error = %v", err)
		return
	}
}

func TestDecrypt(t *testing.T) {
	str := "Hello World"
	encrypted := "9LELUzHgcqr5rYq9X57vGvzHiH5HpB8Lwsw5C2M85XY66JRRuinrjBMrbKjBgijPcHeZ1575RFjzfyUUv"

	decrypted, err := AutoDecrypt(encrypted)
	if decrypted != str {
		t.Errorf("Decrypt() error = %v", err)
		return
	}

	if decrypted != str {
		t.Errorf("Decrypt() error = %v", err)
		return
	}
}
