package util

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func CheckHashFile(filepath, hash string) bool {
	f := OpenFileMetadata(filepath, true)
	if f.Error != nil {
		return false
	}

	defer DeferCloseFatal(f.File)

	return f.Hash == hash
}

func CalculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %s", err)
	}

	defer DeferCloseFatal(file)

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("error calculating MD5 hash: %s", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
