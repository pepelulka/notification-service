package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type WorkerConfig struct {
	EmailWorkerConfig EmailSenderConfig `yaml:"email"`
	TgSenderConfig    TgSenderConfig    `yaml:"tg"`
}

type EmailSenderConfig struct {
	SmtpHost       string `yaml:"smtp_host"`
	SmtpPort       string `yaml:"smtp_port"`
	SenderAddress  string `yaml:"sender_address"`
	SenderPassword string `yaml:"sender_password"`
}

type EtcdConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type TgSenderConfig struct {
	Token string     `yaml:"token"`
	Etcd  EtcdConfig `yaml:"etcd"`
}

func LoadWorkerConfig(filename string) (WorkerConfig, error) {
	// Open our yaml file
	yamlFile, err := os.Open(filename)
	// if we os.Open returns an error then handle it
	if err != nil {
		return WorkerConfig{}, err
	}
	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)
	var config WorkerConfig
	yaml.Unmarshal(byteValue, &config)
	return config, nil
}
