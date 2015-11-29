package config

func DatabaseUrl() string {
	return "postgresql://localhost/gocard_dev?sslmode=disable"
}

func CookieSecret() string {
	return "something-very-secret"
}
