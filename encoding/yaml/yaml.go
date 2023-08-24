package yaml

import (
	"gopkg.in/yaml.v3"
	"os"
)

func encodeAndSaveToFile(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	return encoder.Encode(data)
}

func decodeFromFile(filename string, out interface{}) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(out)
}

func encodeToString(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	return string(bytes), err
}

func decodeFromString(input string, out interface{}) error {
	return yaml.Unmarshal([]byte(input), out)
}
