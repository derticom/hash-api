package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	LogLevel string `yaml:"log_level"`

	RedisURL string        `yaml:"redis_url"`
	RedisTTL time.Duration `yaml:"redis_ttl"`

	Server HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address string        `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
}

func NewConfig() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config: %s", err)
	}

	return &cfg
}
