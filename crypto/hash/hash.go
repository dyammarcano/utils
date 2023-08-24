package hash

type (
	HashType int
)

const (
	MD5 HashType = iota
	SHA1
	SHA256
	SHA512
	SHA3
	BLAKE2B
	BLAKE2S
)
