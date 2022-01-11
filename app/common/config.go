package common

import (
	"os"
	"strconv"
)

type Configuration struct {
}

func (c *Configuration) GetIntEnv(key string, defaultValue int) int {
	svalue := os.Getenv(key)
	if svalue == "" {
		return defaultValue
	}
	val, _ := strconv.ParseUint(svalue, 10, 32)
	return int(val)
}

func (c *Configuration) GetBoolEnv(key string, defaultValue bool) bool {
	svalue := os.Getenv(key)
	if svalue == "" {
		return defaultValue
	}
	val, _ := strconv.ParseBool(svalue)
	return val
}

func (c *Configuration) GetStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
