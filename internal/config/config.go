package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	PostgresConfig
	HTTPServerConfig
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST"`
	Port     int    `env:"POSTGRES_PORT"`
	DbName   string `env:"POSTGRES_DB_NAME"`
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
}

type HTTPServerConfig struct {
	Host string `env:"HTTP_SERVER_HOST"`
	Port int    `env:"HTTP_SERVER_PORT"`
}

func MustLoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	var cfg Config

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic("Error reading env: " + err.Error())
	}

	return cfg
}
