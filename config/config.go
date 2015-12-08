package config

import (
	"log"
	"os"
)

const (
	LogPrefix = "[app] "
	LogFlags  = log.LstdFlags
)

var defaultLogger *log.Logger

func DefaultLogger() *log.Logger {
	if defaultLogger == nil {
		defaultLogger = log.New(os.Stdout, LogPrefix, LogFlags)
	}
	return defaultLogger
}

func DatabaseURL() string {
	return "postgresql://localhost/gocard_dev?sslmode=disable"
}

// DatabaseTestUrl returns the connection string for the test database
func DatabaseTestURL() string {
	return "postgresql://localhost/gocard_test?sslmode=disable"
}

func CookieSecret() string {
	return "something-very-secret"
}
