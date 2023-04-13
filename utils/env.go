package utils

import (
	"os"
	"strconv"
)

func GetIntFromEnv(key string) int {
	str := os.Getenv(key)
	i, _ := strconv.Atoi(str)
	return i
}

func GetStringEnvWithDefault(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}
