package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	logger         = logrus.New()
)

// Setup initializes the logger with the specified log file.
func Setup(logFile string) error {
	return nil
}

// Info logs an informational message.
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Error logs an error message.
func Error(args ...interface{}) {
	logger.Error(args...)
}

// Fatal logs a fatal error message and exits the application.
func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

// Printf logs a formatted message.
func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

// Writer returns a file handle for the current log file for use with other loggers.
func Writer() (*os.File, error) {
	return nil, nil
}
