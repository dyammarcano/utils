package sign

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func GenerateKeyPair(bits int) (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func SavePrivateKeyAsPEM(privateKey *rsa.PrivateKey, filePath string) error {
	derPKCS1Key := x509.MarshalPKCS1PrivateKey(privateKey)

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derPKCS1Key,
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}(file)

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}

	return nil
}

func SavePublicKeyAsPEM(publicKey *rsa.PublicKey, filePath string) error {
	derPKIXKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPKIXKey,
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}(file)

	if err := pem.Encode(file, block); err != nil {
		return err
	}

	return nil
}

func RSASignData(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hashed := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
}

func RSAVerifyData(data []byte, signature []byte, publicKey *rsa.PublicKey) error {
	hashed := sha256.Sum256(data)
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashed[:], signature)
}

func RSASignFile(filePath string, privateKey *rsa.PrivateKey) ([]byte, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return RSASignData(content, privateKey)
}

func RSAVerifyFile(filePath string, signature []byte, publicKey *rsa.PublicKey) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return RSAVerifyData(content, signature, publicKey)
}
