package settings

import (
	"fmt"
	"os"
	"time"
)

const (
	defaultScanInterval = time.Second
	defaultOutputFile   = "./metrics.log"
	defaultAlgodAddress = "http://localhost:8080"
)

var (
	ScanInterval  time.Duration
	OutputFile    string
	AlgodAddress  string
	AlgodApiToken string
)

func SetupEnvironment() {
	ScanInterval = durationFromEnvOrDefault("SCAN_INTERVAL", defaultScanInterval)
	OutputFile = stringFromEnvOrDefault("OUTPUT_FILE", defaultOutputFile)
	AlgodAddress = stringFromEnvOrDefault("ALGOD_ADDRESS", defaultAlgodAddress)
	AlgodApiToken = mustGetStringFromEnvironment("ALGOD_API_TOKEN")
}

func durationFromEnvOrDefault(key string, defVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defVal
}

func stringFromEnvOrDefault(key string, defVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defVal
}

func mustGetStringFromEnvironment(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	panic(fmt.Sprintf("%s not set", key))
}
