package config

import(
	"github.com/joho/godotenv"
	"os"
)
type Config struct {
	DatabaseURL string
	Port        string
}

func Load() Config {
	godotenv.Load()

	return Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}
}