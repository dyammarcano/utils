package hash

import (
	"io"
)

// HashFile returns the hash of the file at the given path.
func HashFile(file io.Reader, hashType Hash) string {
	data := readFile(file)

	switch hashType {
	case MD5:
		return md5Bytes(data)
	case SHA1:
		return sha1String(data)
	case SHA256:
		return sha256String(data)
	case SHA512:
		return sha512String(data)
	case SHA3:
		return sha3String(data)
	case BLAKE2B:
		return blake2bString(data)
	case BLAKE2S:
		return blake2sString(data)
	}
	return ""
}

func readFile(file io.Reader) []byte {
	var data []byte
	buffer := make([]byte, 1024)
	for {
		bytesRead, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil
		}
		data = append(data, buffer[:bytesRead]...)
	}
	return data
}
