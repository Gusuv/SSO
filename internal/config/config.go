package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env       string     `yaml:"env" env-default:"local"`
	TokenTTL  TokensTTL  `yaml:"token_ttl" env-required:"true"`
	Db        DbConfig   `yaml:"db"`
	Grpc      GRPCConfig `yaml:"grpc"`
	JWTSecret string     `env:"TOKEN_SECRET"`
}

type DbConfig struct {
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
	Port     int    `yaml:"port"`
	Password string `env:"DB_PASSWORD"`
	User     string `env:"DB_USERNAME"`
	SSLMode  string `yaml:"sslmode"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type TokensTTL struct {
	Refresh time.Duration `yaml:"refresh"`
	Access  time.Duration `yaml:"access"`
}

func MustLoad(configPath string) *Config {
	godotenv.Load(".env")

	if configPath == "" {
		panic("Config file path is required")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("Config file not found: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("Config load error: " + err.Error())
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(".ENV load error: " + err.Error())
	}

	return &cfg
}

func (c *Config) DSN() string {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Db.Host, c.Db.User, c.Db.Password, c.Db.Name, c.Db.Port, c.Db.SSLMode)
	return dsn
}

func (c *Config) DbUrl() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Db.User,
		c.Db.Password,
		c.Db.Host,
		c.Db.Port,
		c.Db.Name,
		c.Db.SSLMode,
	)
}

func FetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "Config file path")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
