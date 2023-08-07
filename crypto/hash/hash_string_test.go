package hash

import "testing"

func TestHashString(t *testing.T) {
	md5Str := HashString("test", MD5)
	if md5Str != "098f6bcd4621d373cade4e832627b4f6" {
		t.Errorf("MD5 hash failed")
	}

	sha1Str := HashString("test", SHA1)
	if sha1Str != "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3" {
		t.Errorf("SHA1 hash failed")
	}

	sha256Str := HashString("test", SHA256)
	if sha256Str != "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08" {
		t.Errorf("SHA256 hash failed")
	}

	sha512Str := HashString("test", SHA512)
	if sha512Str != "ee26b0dd4af7e749aa1a8ee3c10ae9923f618980772e473f8819a5d4940e0db27ac185f8a0e1d5f84f88bc887fd67b143732c304cc5fa9ad8e6f57f50028a8ff" {
		t.Errorf("SHA512 hash failed")
	}

	sha3Str := HashString("test", SHA3)
	if sha3Str != "e516dabb23b6e30026863543282780a3ae0dccf05551cf0295178d7ff0f1b41eecb9db3ff219007c4e097260d58621bd" {
		t.Errorf("SHA3 hash failed")
	}

	blake2bStr := HashString("test", BLAKE2B)
	if blake2bStr != "928b20366943e2afd11ebc0eae2e53a93bf177a4fcf35bcc64d503704e65e202" {
		t.Errorf("BLAKE2B hash failed")
	}

	blake2sStr := HashString("test", BLAKE2S)
	if blake2sStr != "f308fc02ce9172ad02a7d75800ecfc027109bc67987ea32aba9b8dcc7b10150e" {
		t.Errorf("BLAKE2S hash failed")
	}
}
