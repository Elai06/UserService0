package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file")

		return err
	}
	return nil
}

func GetInt64(key string) (int64, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s is empty", key)
	}

	base, err := strconv.Atoi(os.Getenv("INT_PARSE_BASE"))
	if err != nil {
		return 0, err
	}

	size, err := strconv.Atoi(os.Getenv("INT_PARSE_SIZE"))
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(value, base, size)
}

func GetTimeDuration(key string) (time.Duration, error) {
	value := os.Getenv(key)
	if value == "" {
		return 0, fmt.Errorf("environment variable %s is empty", key)
	}
	return time.ParseDuration(value)
}
