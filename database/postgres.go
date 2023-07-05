package database

import (
	"fmt"
	"errors"

	"gorm.io/driver/postgres"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/config"
	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/logger"
)

// NewPostgresDB creates a new PostgreSQL database connection using the provided configuration.
// It initializes a GORM database instance with the PostgreSQL driver and configuration options,
// and returns the database connection.
// Returns an error if there is a failure during the connection establishment or configuration.
func NewPostgresDB(cfg *config.RDBMS, logger *logger.Logger) (*gorm.DB, error) {
	// Create the connection string using the provided configuration
	dbString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.Username, cfg.Password, cfg.DBName, cfg.Port)

	// Determine the log level based on the DEBUG flag in the configuration
	var logLevel gormlogger.LogLevel
	if cfg.Debug {
		logLevel = gormlogger.Info
	} else {
		logLevel = gormlogger.Error
	}

	// Open a new GORM database connection with the PostgreSQL driver and logmode configuration
	db, err := gorm.Open(postgres.Open(dbString), &gorm.Config{Logger: gormlogger.Default.LogMode(logLevel)})
	if err != nil {
		return nil, errors.New("DB connection start failure")
	}

	// Retrieve the underlying SQL database connection
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(cfg.MaxConnectionPool)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(cfg.ConnectionsMaxLifeTime)

	return db, nil
}
