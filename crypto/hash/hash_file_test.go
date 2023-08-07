package hash

import (
	"github.com/dyammarcano/utils/mocks"
	"testing"
)

var data = []byte(`Hello, world!
This is a test file.`)

func TestHashFileMD5(t *testing.T) {
	mFile := mocks.NewMockFile(data)
	md5File := HashFile(mFile, MD5)
	if md5File != "81d7303d21b9791a38aa673d5c478e45" {
		t.Errorf("MD5 hash failed")
	}
}

func TestHashFileSHA1(t *testing.T) {
	mFile := mocks.NewMockFile(data)
	sha1File := HashFile(mFile, SHA1)
	if sha1File != "5ae781a34d2c0b450eb59138e5d6a366baf6187c" {
		t.Errorf("SHA1 hash failed")
	}
}

func TestHashFileHash256(t *testing.T) {
	mFile := mocks.NewMockFile(data)
	sha256File := HashFile(mFile, SHA256)
	if sha256File != "07e1a31f7433cb5eedb2d1f784cb1622139a5c85051d22452ab718d22739cabe" {
		t.Errorf("SHA256 hash failed")
	}
}
