package config

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
