package utils

import (
	"bufio"
	"os"
	logger "positive-vibes-spotter/log"
	"strings"
)

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