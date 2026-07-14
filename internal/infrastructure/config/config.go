package config

import "os"

type Config struct {
	// DatabaseURL may be a libpq keyword DSN or a postgres:// URL; pgx
	// accepts both.
	DatabaseURL string
	// Port is the HTTP listen port without colon.
	Port string
}

// Load reads configuration from the environment. Defaults live here — next
// to the code that enforces them — not in the database or deployment files.
func Load() Config {
	return Config{
		DatabaseURL: getEnv("DATABASE_URL", "host=localhost user=marketplace password=marketplace dbname=marketplace port=5432 sslmode=disable"),
		Port:        getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
