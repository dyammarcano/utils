package auto_crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/dyammarcano/utils/base58"
)

// generateKeys generates a 32 byte master key
func generateKeys() (masterKey []byte) {
	masterKey = make([]byte, 32)

	if _, err := rand.Read(masterKey); err != nil {
		panic(err)
	}

	return masterKey
}

// xorBytes is a list of 32 hex strings
func xorBytes(key, version []byte) []byte {
	versionInt := binary.BigEndian.Uint16(version)
	bytes, err := hex.DecodeString(array[int(versionInt)%32])
	if err != nil {
		panic(err)
	}

	r := make([]byte, len(bytes))
	for i := 0; i < len(bytes); i++ {
		r[i] = bytes[i] ^ key[i%len(key)]
	}

	return r
}

// extractKeys extracts the iv, nonce and secret from the master key
func extractKeys(key []byte) ([]byte, []byte) {
	iv := xorBytes(key[:16], key[16:18])
	nonce := make([]byte, 12)
	secret := xorBytes(key[18:], key[16:18])
	copy(nonce, iv)
	return secret, nonce
}

// splitResult splits the result into iv, key and cypherText
func splitResult(result []byte) ([]byte, []byte, []byte) {
	iv, key := extractKeys(result[:32])
	return iv, key, result[32:]
}

// AutoEncrypt encrypts a message using AES-256-GCM
func AutoEncrypt(message string) (string, error) {
	if message == "" {
		return "", errors.New("message to be encrypted cannot be null or empty")
	}

	response := make([]byte, 0)
	response = append(response, generateKeys()...)
	key, nonce := extractKeys(response)

	cc, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		panic(err)
	}

	response = append(response, gcm.Seal(nil, nonce, []byte(message), nil)...)
	return base58.StdEncoding.EncodeToString(response), nil
}

// AutoDecrypt decrypts a message using AES-256-GCM
func AutoDecrypt(message string) (string, error) {
	if message == "" {
		return "", errors.New("message to be decrypted cannot be null or empty")
	}

	decoded, err := base58.StdEncoding.DecodeString(message)
	if err != nil {
		return "", err
	}

	key, nonce, cypherText := splitResult(decoded)

	cc, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		panic(err)
	}

	decrypted, err := gcm.Open(nil, nonce, cypherText, nil)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}

func AutoEncryptBytes(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("message to be encrypted cannot be null or empty")
	}

	response := make([]byte, 0)
	response = append(response, generateKeys()...)
	key, nonce := extractKeys(response)

	cc, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		panic(err)
	}

	response = append(response, gcm.Seal(nil, nonce, data, nil)...)
	return response, nil
}

func AutoDecryptBytes(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("message to be decrypted cannot be null or empty")
	}

	key, nonce, cypherText := splitResult(data)

	cc, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(cc)
	if err != nil {
		panic(err)
	}

	decrypted, err := gcm.Open(nil, nonce, cypherText, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
