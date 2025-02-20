package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DbName   string `yaml:"db_name"`
}

type RabbitConfig struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	EmailQueue string `yaml:"email_queue"`
	TgQueue    string `yaml:"tg_queue"`
}

type Config struct {
	Database DatabaseConfig `yaml:"db"`
	Rabbit   RabbitConfig   `yaml:"rabbit"`
}

func LoadConfig[T any](filename string) (T, error) {
	yamlFile, err := os.Open(filename)
	if err != nil {
		var t T
		return t, err
	}
	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)
	var config T
	yaml.Unmarshal(byteValue, &config)
	return config, nil
}
