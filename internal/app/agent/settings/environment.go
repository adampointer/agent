package settings

import (
	"os"
	"time"
)

const (
	defaultScanInterval = time.Second
	defaultOutputFile   = "./metrics.log"
	defaultAlgodAddress = "http://localhost:8080"
)

var (
	ScanInterval time.Duration
	OutputFile   string
	AlgodAddress string
)

func SetupEnvironment() {
	ScanInterval = durationFromEnvOrDefault("SCAN_INTERVAL", defaultScanInterval)
	OutputFile = stringFromEnvOrDefault("OUTPUT_FILE", defaultOutputFile)
	AlgodAddress = stringFromEnvOrDefault("ALGOD_ADDRESS", defaultAlgodAddress)
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
