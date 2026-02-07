package interfaces

import "github.com/gin-gonic/gin"

// ListResponse is a generic list payload for paginated endpoints.
type ListResponse struct {
	Data       interface{} `json:"data"`
	PageNumber int         `json:"page_number,omitempty"`
	PageSize   int         `json:"page_size,omitempty"`
	Total      int         `json:"total"`
}
