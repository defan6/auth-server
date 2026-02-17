package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string        `yaml:"env" env-default:"local"`
	Server  HttpConfig    `yaml:"server"`
	Logging LoggingConfig `yaml:"logging"`
}

type HttpConfig struct {
	Host         string        `yaml:"host" env-required:"true"`
	Port         int           `yaml:"port" env-required:"true"`
	ReadTimeout  time.Duration `yaml:"read_timeout" env-default:"10s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"10s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type LoggingConfig struct {
	Level string `yaml:"level" env-default:"info"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		panic("config file is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file is empty:" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config:" + path)
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string
	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
