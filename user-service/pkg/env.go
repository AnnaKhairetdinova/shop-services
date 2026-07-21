package pkg

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

func MustEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("переменная окружения %s не задана", key)
	}
	return val
}

func EnvOr(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
