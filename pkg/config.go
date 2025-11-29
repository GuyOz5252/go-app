package pkg

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

func LoadConfig[T any](path string) (T, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var config T
	if err := yaml.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	return config, nil
}