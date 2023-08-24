package base58

import "testing"

func TestEncoding_EncodeToString(t *testing.T) {
	msg := []byte("Hello world")

	encoded := StdEncoding.Encode(msg)

	if encoded != "JxF12TrwXzT5jvT" {
		t.Errorf("Encode() got = %v, want %v", encoded, "JxF12TrwXzT5jvT")
	}

	encoded = StdEncoding.Encode(nil)

	if encoded != "" {
		t.Errorf("Encode() got = %v, want %v", encoded, "")
	}
}

// benchmark:
func BenchmarkEncoding_EncodeToString(b *testing.B) {
	msg := []byte("Hello world")

	for i := 0; i < b.N; i++ {
		StdEncoding.Encode(msg)
	}
}

func TestEncoding_DecodeString(t *testing.T) {
	decoded, err := StdEncoding.Decode("JxF12TrwXzT5jvT")

	if err != nil {
		t.Errorf("Decode() error = %v, wantErr %v", err, false)
	}

	if string(decoded) != "Hello world" {
		t.Errorf("Decode() got = %v, want %v", string(decoded), "Hello world")
	}

	decoded, err = StdEncoding.Decode("")

	if err == nil {
		t.Errorf("Decode() error = %v, wantErr %v", err, true)
	}
}

// benchmark:
func BenchmarkEncoding_DecodeString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := StdEncoding.Decode("JxF12TrwXzT5jvT")
		if err != nil {
			return
		}
	}
}
