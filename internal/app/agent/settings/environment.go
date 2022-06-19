package settings

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	defaultScanInterval = time.Second
	defaultOutputFile   = "./metrics.log"
	defaultAlgodAddress = "http://localhost:8080"
	defaultProcFsPath   = "/proc"
)

var (
	ScanInterval     time.Duration
	OutputFile       string
	AlgodAddress     string
	AlgodApiToken    string
	ProcFSPath       string
	AlgodPidOverride int
)

// SetupEnvironment sets all the required variables from the environment or default values
func SetupEnvironment() {
	ScanInterval = durationFromEnvOrDefault("SCAN_INTERVAL", defaultScanInterval)
	OutputFile = stringFromEnvOrDefault("OUTPUT_FILE", defaultOutputFile)
	AlgodAddress = stringFromEnvOrDefault("ALGOD_ADDRESS", defaultAlgodAddress)
	AlgodApiToken = mustGetStringFromEnvironment("ALGOD_API_TOKEN")
	ProcFSPath = stringFromEnvOrDefault("PROCFS_PATH", defaultProcFsPath)
	AlgodPidOverride = intFromEnvOrDefault("ALGOD_PID_OVERRIDE", 0)
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

func intFromEnvOrDefault(key string, defVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return int(i)
		}
	}
	return defVal
}
