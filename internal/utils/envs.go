package utils

import (
	"os"
	"strconv"
)

type Environment struct{}

func GetEnvironment() *Environment {
	return &Environment{}
}

func (e *Environment) GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func (e *Environment) GetEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}
