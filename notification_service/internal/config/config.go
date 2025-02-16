package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DbName   string `yaml:"db_name"`
	} `yaml:"db"`
}

func LoadConfig(filename string) (Config, error) {
	// Open our yaml file
	yamlFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		return Config{}, err
	}
	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)
	var config Config
	yaml.Unmarshal(byteValue, &config)
	return config, nil
}
