package sign

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"github.com/dyammarcano/utils/encoding/base58"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io"
	"log"
	"os"
)

type (
	Encoder interface {
		Encode([]byte) string
	}

	EncodeType    int
	Base64Encoder struct{}
	Base58Encoder struct{}
	HexEncoder    struct{}
	BynaryEncoder struct{}
)

const (
	Base64 EncodeType = iota
	Base58
	Hex
	Binary
)

func (e Base64Encoder) Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func (e Base58Encoder) Encode(data []byte) string {
	return base58.StdEncoding.Encode(data)
}

func (e HexEncoder) Encode(data []byte) string {
	return hex.EncodeToString(data)
}

func (e BynaryEncoder) Encode(data []byte) string {
	return string(data)
}

type Decoder interface {
	Decode(string) ([]byte, error)
}

func (e Base64Encoder) Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func (e Base58Encoder) Decode(data string) ([]byte, error) {
	return base58.StdEncoding.Decode(data)
}

func (e HexEncoder) Decode(data string) ([]byte, error) {
	return hex.DecodeString(data)
}

func (e BynaryEncoder) Decode(data string) ([]byte, error) {
	return []byte(data), nil
}

// SignFile signs a file using the HMAC-SHA256 algorithm
func SignFile(key []byte, file io.Reader, encodeType EncodeType) string {
	var buffer bytes.Buffer
	if _, err := buffer.ReadFrom(file); err != nil {
		return ""
	}

	return generateHMAC(key, buffer.Bytes(), encodeType)
}

// SignString signs a string using the HMAC-SHA256 algorithm
func SignString(key []byte, data []byte, encodeType EncodeType) string {
	return generateHMAC(key, data, encodeType)
}

// SignBytes signs a byte array using the HMAC-SHA256 algorithm
func SignBytes(key []byte, data []byte, encodeType EncodeType) string {
	return generateHMAC(key, data, encodeType)
}

// VerifyFile verifies the signature of a file using the HMAC-SHA256 algorithm
func VerifyFile(key []byte, file io.Reader, signature string, encodeType EncodeType) bool {
	return signature == SignFile(key, file, encodeType)
}

// VerifyString verifies the signature of a string using the HMAC-SHA256 algorithm
func VerifyString(key []byte, data []byte, signature string, encodeType EncodeType) bool {
	return signature == SignString(key, data, encodeType)
}

// SignFileGPG signs a file using the GPG algorithm
func SignFileGPG(privateKeyPath string, filePath string, signaturePath string) error {
	// Read the private key file
	privateKeyEntity, err := readPrivateKey(privateKeyPath)
	if err != nil {
		return err
	}

	// Create the output signature file
	signatureFile, err := os.Create(signaturePath)
	if err != nil {
		return err
	}
	defer func(signatureFile *os.File) {
		if err := signatureFile.Close(); err != nil {
			log.Println(err)
		}
	}(signatureFile)

	// Sign the file
	return signFile(privateKeyEntity, filePath, signatureFile)
}

//type CustomKeyRing struct {
//	publicKey *openpgp.Entity
//}
//
//func (kr *CustomKeyRing) KeysById(id uint64) []openpgp.Key {
//	if kr.publicKey.PrimaryKey.KeyId == id {
//		return []openpgp.Key{kr.publicKey.PrimaryKey}
//	}
//	return nil
//}
//
//func (kr *CustomKeyRing) KeysByIdUsage(id uint64, requiredUsage byte) []openpgp.Key {
//	if kr.publicKey.PrimaryKey.KeyId == id {
//		return []openpgp.Key{kr.publicKey.PrimaryKey}
//	}
//	return nil
//}
//
//func (kr *CustomKeyRing) DecryptionKeys() []openpgp.Key {
//	return nil
//}
//
//func VerifyFileGPG(publicKeyPath string, filePath string, signaturePath string) error {
//	// Read the public key file
//	publicKeyEntity, err := readPublicKey(publicKeyPath)
//	if err != nil {
//		return err
//	}
//
//	signature, err := os.Open(signaturePath)
//	if err != nil {
//		log.Println(err)
//	}
//	defer func(file *os.File) {
//		if err := file.Close(); err != nil {
//			log.Println(err)
//		}
//	}(signature)
//
//	// Create a buffered reader for the signature file
//	signatureReader := bufio.NewReader(signature)
//
//	file, err := os.Open(filePath)
//	if err != nil {
//		log.Println(err)
//	}
//	defer func(file *os.File) {
//		if err := file.Close(); err != nil {
//			log.Println(err)
//		}
//	}(file)
//
//	// Create a slice of keys with the public key
//	publicKeys := []openpgp.Entity{*publicKeyEntity}
//
//	// Create a new KeyRing with the public keys
//	publicKeyRing := openpgp.NewKeyRing(publicKeys)
//
//	// Verify the file
//	return verifyFile(publicKeyRing, signatureReader, file)
//}
//
//func verifyFile(publicKeyEntity openpgp.KeyRing, signed io.Reader, file io.Reader) error {
//	signature, err := armor.Decode(signed)
//	if err != nil {
//		return err
//	}
//
//	_, err = openpgp.CheckArmoredDetachedSignature(publicKeyEntity, file, signature.Body)
//	return err
//}
//
//func readPublicKey(publicKeyPath string) (*openpgp.Entity, error) {
//	publicKeyFile, err := os.Open(publicKeyPath)
//	if err != nil {
//		return nil, err
//	}
//	defer publicKeyFile.Close()
//
//	entityList, err := openpgp.ReadArmoredKeyRing(publicKeyFile)
//	if err != nil {
//		return nil, err
//	}
//	return entityList[0], nil
//}

func readPrivateKey(privateKeyPath string) (*openpgp.Entity, error) {
	privateKeyFile, err := os.Open(privateKeyPath)
	if err != nil {
		return nil, err
	}
	defer func(privateKeyFile *os.File) {
		if err := privateKeyFile.Close(); err != nil {
			log.Println(err)
		}
	}(privateKeyFile)

	entityList, err := openpgp.ReadArmoredKeyRing(privateKeyFile)
	if err != nil {
		return nil, err
	}
	return entityList[0], nil
}

func signFile(privateKeyEntity *openpgp.Entity, filePath string, signatureFile *os.File) error {
	armorWriter, err := armor.Encode(signatureFile, openpgp.SignatureType, nil)
	if err != nil {
		return err
	}
	defer func(armorWriter io.WriteCloser) {
		if err := armorWriter.Close(); err != nil {
			log.Println(err)
		}
	}(armorWriter)

	signer, err := openpgp.Sign(armorWriter, privateKeyEntity, nil, nil)
	if err != nil {
		return err
	}
	defer func(signer io.WriteCloser) {
		if err := signer.Close(); err != nil {
			log.Println(err)
		}
	}(signer)

	fileToSign, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(fileToSign *os.File) {
		if err := fileToSign.Close(); err != nil {
			log.Println(err)
		}
	}(fileToSign)

	if _, err = io.Copy(signer, fileToSign); err != nil {
		return err
	}

	return nil
}

func generateHMAC(key []byte, data []byte, encodeType EncodeType) string {
	// Create a new HMAC instance with the SHA-256 hash function
	h := hmac.New(sha256.New, key)

	// Write the data to the HMAC
	h.Write(data)

	// Calculate the HMAC signature
	signature := h.Sum(nil)

	var encodedSignature Encoder

	switch encodeType {
	case Base64:
		encodedSignature = Base64Encoder{}
	case Base58:
		encodedSignature = Base58Encoder{}
	case Hex:
		encodedSignature = HexEncoder{}
	default:
		encodedSignature = BynaryEncoder{}
	}

	// Encode the signature to a string
	return encodedSignature.Encode(signature)
}
