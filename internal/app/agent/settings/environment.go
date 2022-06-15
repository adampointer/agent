package settings

import (
	"os"
	"time"
)

const defaultScanInterval = time.Second

var ScanInterval time.Duration

func SetupEnvironment() {
	ScanInterval = durationFromEnvOrDefault("SCAN_INTERVAL", defaultScanInterval)
}

func durationFromEnvOrDefault(key string, defVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defVal
}
