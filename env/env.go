package env

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file %v", err)
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
		return 0, fmt.Errorf("environment variable %s is invalid", key)
	}

	size, err := strconv.Atoi(os.Getenv("INT_PARSE_SIZE"))
	if err != nil {
		return 0, fmt.Errorf("environment variable %s is empty", key)
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
