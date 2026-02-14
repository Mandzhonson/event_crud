package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HTTPConfing struct {
		Host string `env:"SRV_HOST" env-default:"localhost"`
		Port string `env:"SRV_PORT" env-default:"8080"`
	}
	DatabaseConfig struct {
		Name string `env:"POSTGRES_DB" required: "true"`
		User string `env:"POSTGRES_USER" default:"postgres"`
		Pass string `env:"POSTGRES_PASSWORD" default:"postgres"`
		Host string `env:"POSTGRES_HOST" default:"postgres"`
		Port string `env:"POSTGRES_PORT" default:"postgres"`
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

func (cfg Config) GetDBString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Pass,
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Name,
	)
}
