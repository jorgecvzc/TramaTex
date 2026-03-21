package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the application-wide structured logger
var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
}

// InitLogger initializes the global logger with configuration based on environment
func InitLogger(env string) {
	Logger = logrus.New()

	// Set output to stdout
	Logger.SetOutput(os.Stdout)

	// Set JSON formatter for structured logging
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Set log level based on environment
	if env == "production" {
		Logger.SetLevel(logrus.InfoLevel)
	} else {
		Logger.SetLevel(logrus.DebugLevel)
	}

	Logger.Info("Logger initialized successfully")
}

// WithRequestID adds request ID to log entry
func WithRequestID(requestID string) *logrus.Entry {
	return Logger.WithField("requestID", requestID)
}

// WithUser adds user context to log entry
func WithUser(userID, email, role string) *logrus.Entry {
	return Logger.WithFields(logrus.Fields{
		"userID":    userID,
		"userEmail": MaskEmail(email),
		"userRole":  role,
	})
}

// WithUserAndRequest combines user and request context
func WithUserAndRequest(requestID, userID, email, role string) *logrus.Entry {
	return Logger.WithFields(logrus.Fields{
		"requestID": requestID,
		"userID":    userID,
		"userEmail": MaskEmail(email),
		"userRole":  role,
	})
}

// MaskEmail masks email for privacy in logs (j***@tramatex.com)
func MaskEmail(email string) string {
	if len(email) == 0 {
		return ""
	}

	// Find @ position
	atPos := -1
	for i, ch := range email {
		if ch == '@' {
			atPos = i
			break
		}
	}

	if atPos <= 0 {
		// Invalid email or @ at start, return first char + ***
		if len(email) >= 1 {
			return string(email[0]) + "***"
		}
		return "***"
	}

	// Show first character + *** + domain
	return string(email[0]) + "***" + email[atPos:]
}
