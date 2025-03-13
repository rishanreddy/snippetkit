package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Global logger instance
var log = logrus.New()

// InitLogger initializes the logger with proper settings
func InitLogger() {
	// Ensure logs directory exists
	logDir := filepath.Join(os.Getenv("HOME"), ".config/snippetkit/logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Println("⚠️ Failed to create logs directory:", err)
		return
	}

	// Define log file with rotation based on date
	logFile := filepath.Join(logDir, fmt.Sprintf("snippetkit-%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("⚠️ Failed to open log file:", err)
		return
	}

	// Configure log output (console + file)
	log.SetOutput(io.MultiWriter(os.Stdout, file))

	// Set log format
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	// Load log level from config (default: INFO)
	logLevel := viper.GetString("log_level")
	SetLogLevel(logLevel)
}

// SetLogLevel dynamically updates the log level
func SetLogLevel(level string) {
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}
}

// Debug logs a debug message
func Debug(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Debug(msg)
}

// Info logs an info message
func Info(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Info(msg)
}

// Warn logs a warning message
func Warn(msg string, fields map[string]interface{}) {
	log.WithFields(fields).Warn(msg)
}

// Error logs an error message
func Error(msg string, err error, fields map[string]interface{}) {
	log.WithFields(fields).WithError(err).Error(msg)
}

// Fatal logs a fatal error and exits the program
func Fatal(msg string, err error, fields map[string]interface{}) {
	log.WithFields(fields).WithError(err).Fatal(msg)
	os.Exit(1)
}
