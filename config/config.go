package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

// Configuration - server, databases and email service configuration variables
type Config struct {
	DB         DatabaseConfig
	Server     ServerConfig
	Email  	   EmailConfig
}

func New() *Config {
	var c Config

	// Unmarshel .env in Config struct fields
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}