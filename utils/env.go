package utils

import (
	"os"
	"strings"
)

// IsDebugMode returns true if the ENVIRONMENT env var is set to DEBUG
func IsDebugMode() bool {
	return strings.ToLower(os.Getenv("ENVIRONMENT")) == "debug"
}

func IsProductionMode() bool {
	return strings.ToLower(os.Getenv("ENVIRONMENT")) == "prod"
}
