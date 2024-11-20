package monitoring

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// Define a whitelist of allowed fields that can be logged
var allowedLogFields = map[string]bool{
	"request_id": true, // Example: Include request IDs for debugging
	"status":     true, // Example: Include status codes
	"error":      true, // Example: Include error messages
	"message":    true, // Example: General messages
	"timestamp":  true, // Example: Include timestamps if necessary
	"query":      true, // Example: Include query strings
}

// NewLogger creates a new instance of a logrus.Logger with the specified log level and format.
// The log level can be any valid logrus log level (e.g., "debug", "info", "warn", "error").
// If the provided log level is invalid, it defaults to "info" level.
// The log format can be either "json" or "JSON" for JSON formatted logs, or any other value for text formatted logs.
// The logger output is set to os.Stdout.
//
// Parameters:
//   - logLevel: The desired log level as a string.
//   - logFormat: The desired log format as a string.
//
// Returns:
//   - *logrus.Logger: A pointer to the configured logrus.Logger instance.
func NewLogger(logLevel string, logFormat string) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel // Default to Info level if parsing fails
	}
	logger.SetLevel(level)

	switch logFormat {
	case "json", "JSON":
		logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.SetFormatter(&logrus.TextFormatter{})
	}

	// Open a file for writing logs
	// file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// if err == nil {
	// 		logger.SetOutput(file)
	// } else {
	logger.SetOutput(os.Stdout)
	// }

	return logger
}

func InitGlobalLogger(logLevel, logFormat string) {
	Logger = NewLogger(logLevel, logFormat)
}

// cleanFields sanitizes the provided log fields by retaining only the allowed fields
// and redacting the rest. The values of the allowed fields are sanitized using the
// SanitizeValue function.
//
// Parameters:
//
//	fields (logrus.Fields): The log fields to be sanitized.
//
// Returns:
//
//	logrus.Fields: The sanitized log fields with only allowed fields retained and
//	their values sanitized, while the rest are redacted.
func cleanFields(fields logrus.Fields) logrus.Fields {
	sanitizedFields := make(logrus.Fields)
	for key, value := range fields {
		if allowedLogFields[key] {
			sanitizedFields[key] = SanitizeValue(fmt.Sprintf("%v", value))
		} else {
			sanitizedFields[key] = "<REDACTED>"
		}
	}
	return sanitizedFields
}

// LogDebug logs a debug message with specific fields
func LogDebug(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Debug(message)
}

// LogInfo logs an informational message with specific fields
func LogInfo(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Info(message)
}

// LogError logs an error with a specific context
func LogError(err error, message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).WithError(err).Warn(message)
}

// LogWarn logs a warning message with specific fields
func LogWarn(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Warn(message)
}

// LogFatal logs a fatal message with specific fields
func LogFatal(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Fatal(message)
}

// LogPanic logs a panic message with specific fields
func LogPanic(message string, fields logrus.Fields) {
	Logger.WithFields(cleanFields(fields)).Panic(message)
}
