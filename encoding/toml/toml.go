package toml

import (
	"bytes"
	"github.com/BurntSushi/toml"
	"os"
)

func encodeAndSaveToFile(config any, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(file)

	if err := toml.NewEncoder(file).Encode(config); err != nil {
		return err
	}

	return nil
}

func decodeFromFile(filename string, config any) error {
	_, err := toml.DecodeFile(filename, &config)
	if err != nil {
		return err
	}

	return nil
}

func encodeToString(config any) (string, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(config); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func decodeFromString(input string, config any) error {
	_, err := toml.Decode(input, &config)
	if err != nil {
		return err
	}
	return nil
}
