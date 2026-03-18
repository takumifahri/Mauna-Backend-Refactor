package config

import (
	"REFACTORING_MAUNA/pkg/database"
	"time"
)

func NewDatabaseConfig() database.Config {
	return database.Config{
		Host:            "localhost",
		Port:            "5432",
		User:            "postgres",
		Password:        "postgres",
		Database:        "mauna_db",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}
