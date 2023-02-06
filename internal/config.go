package internal

import "os"

/* Code to pull config values from .env environmental variables */

type Config struct {
	Env          string
	HTTPListen   string
	PostgresHost string
	PostgresPort string
	PostgresDb   string
	PostgresUser string
	PostgresPw   string
}

func ParseConfigFromEnvironment(prefix string) (*Config, error) {
	return &Config{
		Env:          os.Getenv("ENV"),
		HTTPListen:   os.Getenv("HTTP_ADDRESS"),
		PostgresHost: os.Getenv("POSTGRES_HOST"),
		PostgresPort: os.Getenv("POSTGRES_PORT"),
		PostgresDb:   os.Getenv("POSTGRES_DB"),
		PostgresUser: os.Getenv("POSTGRES_USER"),
		PostgresPw:   os.Getenv("POSTGRES_PW"),
	}, nil

}
