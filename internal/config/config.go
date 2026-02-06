package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string        `yaml:"env" env-default:"local"`
	DbPath   string        `yaml:"db_path" env-reload:"true"`
	TokenTTL time.Duration `yaml:"token_ttl" env-required:"true"`
	Grpc     GrpcConfig    `yaml:"grpc"`
}

type GrpcConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPatch()

	if path == "" {
		panic("Config file path is required")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("Config file not found: " + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Config load error: " + err.Error())
	}

	return &cfg
}

func fetchConfigPatch() string {
	var res string

	flag.StringVar(&res, "config", "", "Config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
