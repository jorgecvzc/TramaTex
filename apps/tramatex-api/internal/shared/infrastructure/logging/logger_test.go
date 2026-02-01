package logging_test

import (
	"testing"

	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/logging"
	"github.com/stretchr/testify/assert"
)

func TestMaskEmail(t *testing.T) {
	// Initialize logger for testing
	logging.InitLogger("test")

	tests := []struct {
		name     string
		email    string
		expected string
	}{
		{
			name:     "normal email",
			email:    "jorge@tramatex.com",
			expected: "j***@tramatex.com",
		},
		{
			name:     "short email",
			email:    "a@b.com",
			expected: "a***@b.com",
		},
		{
			name:     "empty email",
			email:    "",
			expected: "",
		},
		{
			name:     "no @ symbol",
			email:    "invalid",
			expected: "i***",
		},
		{
			name:     "@ at start",
			email:    "@domain.com",
			expected: "@***",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use WithUser to test maskEmail indirectly
			entry := logging.WithUser("user-123", tt.email, "admin")

			// Verify email field is masked in the entry
			maskedEmail, exists := entry.Data["userEmail"]
			assert.True(t, exists)
			assert.Equal(t, tt.expected, maskedEmail)
		})
	}
}

func TestWithRequestID(t *testing.T) {
	logging.InitLogger("test")

	entry := logging.WithRequestID("req-12345")

	requestID, exists := entry.Data["requestID"]
	assert.True(t, exists)
	assert.Equal(t, "req-12345", requestID)
}

func TestWithUser(t *testing.T) {
	logging.InitLogger("test")

	entry := logging.WithUser("user-123", "jorge@tramatex.com", "admin")

	assert.Equal(t, "user-123", entry.Data["userID"])
	assert.Equal(t, "j***@tramatex.com", entry.Data["userEmail"])
	assert.Equal(t, "admin", entry.Data["userRole"])
}

func TestWithUserAndRequest(t *testing.T) {
	logging.InitLogger("test")

	entry := logging.WithUserAndRequest("req-123", "user-456", "test@example.com", "commercial")

	assert.Equal(t, "req-123", entry.Data["requestID"])
	assert.Equal(t, "user-456", entry.Data["userID"])
	assert.Equal(t, "t***@example.com", entry.Data["userEmail"])
	assert.Equal(t, "commercial", entry.Data["userRole"])
}
