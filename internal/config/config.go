package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration settings.
type Config struct {
	// DatabaseURL is the URL for the database connection.
	DatabaseURL string
	// Port is the port on which the server will run.
	Port string
	// JWTSecret is the secret key used for JWT token signing.
	JWTSecret string
}

// Load loads the configuration from environment variables.
func Load() Config {
	godotenv.Load()

	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
	}
}
