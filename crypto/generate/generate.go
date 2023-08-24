package generate

import (
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"log"
	"os"
)

type (
	Entity struct {
		KeyPath    string
		Name       string
		Comment    string
		Email      string
		Passphrase []byte
	}
)

func GenerateKey(entity Entity) error {
	privateKeyFile, err := os.Create(entity.KeyPath)
	if err != nil {
		return err
	}
	defer func(privateKeyFile *os.File) {
		if err := privateKeyFile.Close(); err != nil {
			log.Printf("Error: %v", err)
		}
	}(privateKeyFile)

	newEntity, err := openpgp.NewEntity(entity.Name, entity.Comment, entity.Email, nil)
	if err != nil {
		return err
	}

	w, err := armor.Encode(privateKeyFile, openpgp.PrivateKeyType, nil)
	if err != nil {
		return err
	}
	defer func(w io.WriteCloser) {
		if err := w.Close(); err != nil {
			log.Printf("Error: %v", err)
		}
	}(w)

	config := &packet.Config{
		DefaultCipher: packet.CipherAES256,
	}

	if err := newEntity.SerializePrivate(w, config); err != nil {
		return err
	}

	if err := newEntity.Serialize(os.Stdout); err != nil {
		return err
	}

	return nil
}
