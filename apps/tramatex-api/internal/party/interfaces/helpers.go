package interfaces

import "github.com/gin-gonic/gin"

// ListResponse is a generic list payload for paginated endpoints.
type ListResponse struct {
	Data       interface{} `json:"data"`
	PageNumber int         `json:"page_number,omitempty"`
	PageSize   int         `json:"page_size,omitempty"`
	Total      int         `json:"total"`
}

func getUserIDFromContext(c *gin.Context) string {
	if c == nil {
		return "system"
	}
	userID, ok := c.Get("userID")
	if !ok {
		return "system"
	}
	value, ok := userID.(string)
	if !ok || value == "" {
		return "system"
	}
	return value
}
