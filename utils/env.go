package utils

import (
	"bufio"
	"os"
	"strings"
	"sync"

	logger "positive-vibes-spotter/log"
)

var (
	loadEnvOnce sync.Once
	loadEnvErr  error
)

// LoadEnv loads the environment variables from a file.
func LoadEnv(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] != '#' {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				os.Setenv(parts[0], parts[1])
				logger.Info("Env variable loaded: ", parts[0])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// GetEnv retrieves the value of the environment variable named by the key.
// It ensures the .env file is loaded only once.
func GetEnv(key string) string {
	loadEnvOnce.Do(func() {
		loadEnvErr = LoadEnv(".env")
	})
	if loadEnvErr != nil {
		logger.Error("Error loading .env file: ", loadEnvErr)
		return ""
	}
	return os.Getenv(key)
}