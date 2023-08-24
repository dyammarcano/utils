package uuid

import "github.com/google/uuid"

func UUID() uuid.UUID {
	return uuid.New()
}

func ToString() string {
	return UUID().String()
}

func ToBytes() []byte {
	return UUID().NodeID()
}

func FromString(str string) (uuid.UUID, error) {
	return uuid.Parse(str)
}

func FromBytes(b []byte) (uuid.UUID, error) {
	return uuid.ParseBytes(b)
}
