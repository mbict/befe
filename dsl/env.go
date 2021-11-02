package dsl

import (
	"os"
)

func FromEnv(name string) string {
	return FromEnvWithDefault(name, "")
}

func FromEnvWithDefault(name string, defaultValue string) string {
	value, found := os.LookupEnv(name)
	if found == false || value == "" {
		return defaultValue
	}
	return value
}
