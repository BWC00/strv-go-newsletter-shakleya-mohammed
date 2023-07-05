package config

import (
	"time"
)

type ServerConfig struct {
	Major      	 int           `env:"SERVER_MAJOR,required"`
	Minor      	 int           `env:"SERVER_MINOR,required"`
	Port         int           `env:"SERVER_PORT,required"`
	Debug        bool          `env:"SERVER_DEBUG,required"`
	SecretKey	 string        `env:"API_SECRET,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}