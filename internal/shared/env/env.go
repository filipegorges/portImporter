package env

import (
	"log"
	"os"
)

func MustEnv(varName string) string {
	value := os.Getenv(varName)
	if value != "" {
		return value
	}

	log.Fatalf("%s: must not be empty", varName)
	return ""
}
