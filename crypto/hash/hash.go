package hash

type (
	Hash int
)

const (
	MD5 Hash = iota
	SHA1
	SHA256
	SHA512
	SHA3
	BLAKE2B
	BLAKE2S
)
