package config

import (
	"time"
)

// DatabaseConfig - all database variables
type DatabaseConfig struct {
	
	// relational database
	RDBMS RDBMS

}

// RDBMS - relational database variables
type RDBMS struct {
	Driver   string 					 `env:"DB_DRIVER,required"`
	Host     string 					 `env:"DB_HOST,required"`
	Port     int    					 `env:"DB_PORT,required"`
	Username string 					 `env:"DB_USER,required"`
	Password string 					 `env:"DB_PASS,required"`
	DBName   string 					 `env:"DB_NAME,required"`
	Debug    bool   					 `env:"DB_DEBUG,required"`
	MaxConnectionPool      int           `env:"DB_MAX_CONNECTION_POOL",default:"4"`
	MaxIdleConnections     int           `env:"DB_MAX_IDLE_CONNECTIONS",default:"4"`
	ConnectionsMaxLifeTime time.Duration `env:"DB_CONNECTIONS_MAX_LIFETIME",default:"300s"`
}
