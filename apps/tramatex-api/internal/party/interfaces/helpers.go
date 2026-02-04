package interfaces

import "github.com/gin-gonic/gin"

// ListResponse is a generic list payload for paginated endpoints.
type ListResponse struct {
	Data       interface{} `json:"data"`
	PageNumber int         `json:"page_number,omitempty"`
	PageSize   int         `json:"page_size,omitempty"`
	Total      int         `json:"total"`
}

// getUserIDFromContext returns the authenticated user ID or a fallback value.
func getUserIDFromContext(c *gin.Context) string {
	if c == nil {
		return "system"
	}
	if value, exists := c.Get("userID"); exists {
		if userID, ok := value.(string); ok && userID != "" {
			return userID
		}
	}
	return "system"
}
