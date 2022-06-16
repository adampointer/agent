package settings

import (
	"os"
	"time"
)

const (
	defaultScanInterval = time.Second
	defaultOutputFile   = "./metrics.log"
)

var (
	ScanInterval time.Duration
	OutputFile   string
)

func SetupEnvironment() {
	ScanInterval = durationFromEnvOrDefault("SCAN_INTERVAL", defaultScanInterval)
	OutputFile = stringFromEnvOrDefault("OUTPUT_FILE", defaultOutputFile)
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
