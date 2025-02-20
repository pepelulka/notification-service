package config

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
