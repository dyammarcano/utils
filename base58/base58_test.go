package base58

import "testing"

func TestEncoding_EncodeToString(t *testing.T) {
	msg := []byte("Hello world")

	encoded := StdEncoding.EncodeToString(msg)

	if encoded != "JxF12TrwXzT5jvT" {
		t.Errorf("EncodeToString() got = %v, want %v", encoded, "JxF12TrwXzT5jvT")
	}

	encoded = StdEncoding.EncodeToString(nil)

	if encoded != "" {
		t.Errorf("EncodeToString() got = %v, want %v", encoded, "")
	}
}

// benchmark:
func BenchmarkEncoding_EncodeToString(b *testing.B) {
	msg := []byte("Hello world")

	for i := 0; i < b.N; i++ {
		StdEncoding.EncodeToString(msg)
	}
}

func TestEncoding_DecodeString(t *testing.T) {
	decoded, err := StdEncoding.DecodeString("JxF12TrwXzT5jvT")

	if err != nil {
		t.Errorf("DecodeString() error = %v, wantErr %v", err, false)
	}

	if string(decoded) != "Hello world" {
		t.Errorf("DecodeString() got = %v, want %v", string(decoded), "Hello world")
	}

	decoded, err = StdEncoding.DecodeString("")

	if err == nil {
		t.Errorf("DecodeString() error = %v, wantErr %v", err, true)
	}
}

// benchmark:
func BenchmarkEncoding_DecodeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := StdEncoding.DecodeString("JxF12TrwXzT5jvT")
		if err != nil {
			return
		}
	}
}
