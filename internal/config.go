package internal

import "os"

/* Code to pull config values from .env environmental variables */

type Config struct {
	Env        string
	HTTPListen string
}

func ParseConfigFromEnvironment(prefix string) (*Config, error) {
	return &Config{
		Env:        os.Getenv("ENV"),
		HTTPListen: os.Getenv("HTTP_ADDRESS"),
	}, nil

}
