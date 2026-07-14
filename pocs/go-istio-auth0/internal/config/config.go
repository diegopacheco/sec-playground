package config

import "os"

type Config struct {
	Port       string
	DBPath     string
	AuthHeader string
}

func Load() Config {
	return Config{
		Port:       getenv("PORT", "8080"),
		DBPath:     getenv("DB_PATH", "books.db"),
		AuthHeader: getenv("AUTH_HEADER", "X-Auth-User"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
