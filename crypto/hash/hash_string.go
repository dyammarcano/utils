package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/sha3"
)

// HashString returns the hash of the data.
func HashString(str string, hashType HashType) string {
	data := []byte(str)

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

func md5Bytes(data []byte) string {
	hasher := md5.New()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func sha1String(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func sha256String(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func sha512String(data []byte) string {
	hasher := sha512.New()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func sha3String(data []byte) string {
	hasher := sha3.New384()
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func blake2bString(data []byte) string {
	hasher, err := blake2b.New256(nil)
	if err != nil {
		return ""
	}
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func blake2sString(data []byte) string {
	hasher, err := blake2s.New256(nil)
	if err != nil {
		return ""
	}
	hasher.Write(data)
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

//func md5File(data []byte) string {
//	file, err := os.Open(path)
//	if err != nil {
//		return ""
//	}
//	defer file.Close()
//
//	hasher := md5.New()
//
//	// Copy the file contents to the hasher
//	_, err = io.Copy(hasher, file)
//	if err != nil {
//		return ""
//	}
//
//	return hex.EncodeToString(hasher.Sum(nil))
//}
//
//func md5ByteArray(data []byte) string {
//	return md5Bytes(data)
//}
//
//// sha1 string
//func Sha1String(str string) string {
//	return ""
//}

//sha1 file

//sha1 byte array

//sha256 string

//sha256 file

//sha256 byte array

//sha512 string

//sha512 file

//sha512 byte array

//sha3 string

//sha3 file

//sha3 byte array

//blake2b string

//blake2b file

//blake2b byte array

//blake2s string

//blake2s file

//blake2s byte array
