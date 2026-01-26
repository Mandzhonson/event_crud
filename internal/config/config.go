package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPConfing struct {
		Host string `env:"SRV_HOST" env-default:"localhost"`
		Port string `env:"SRV_PORT" env-default:"8080"`
	}
	DatabaseConfig struct {
		ConnString string `env:"DB_CONN_STR" env-required:"true"`
	}
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
