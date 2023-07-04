package config

import (
	"time"
)

type ServerConfig struct {
	Version      string        `env:"RELEASE_VERSION"`
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
	SecretKey	 string        `env:"API_SECRET,required"`
}