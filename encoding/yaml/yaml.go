package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

func encodeAndSaveToFile(data any, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(file)

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(data)
}

func decodeFromFile(filename string, out any) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}(file)

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(out)
}

func encodeToString(data any) (string, error) {
	bytes, err := yaml.Marshal(data)
	return string(bytes), err
}

func decodeFromString(input string, out any) error {
	return yaml.Unmarshal([]byte(input), out)
}
