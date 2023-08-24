package ulid

import "github.com/oklog/ulid/v2"

func ULID() ulid.ULID {
	return ulid.Make()
}

func ToString() string {
	return ULID().String()
}

func ToBytes() []byte {
	return ULID().Entropy()
}

func FromString(str string) (ulid.ULID, error) {
	return ulid.Parse(str)
}

func FromBytes(b []byte) (ulid.ULID, error) {
	return ulid.Parse(string(b))
}
