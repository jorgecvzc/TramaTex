package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/iam/application/usecase"
	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
)

// IAMHandler handles authentication and user management endpoints.
type IAMHandler struct {
	loginUseCase *usecase.LoginUseCase
}

// NewIAMHandler creates a new IAM handler.
func NewIAMHandler(loginUseCase *usecase.LoginUseCase) *IAMHandler {
	return &IAMHandler{
		loginUseCase: loginUseCase,
	}
}

// Login handles POST /api/iam/login
func (h *IAMHandler) Login(c *gin.Context) {
	var req usecase.LoginInput

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	// Execute use case
	output, err := h.loginUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		// Handle specific domain errors
		if errors.Is(err, model.ErrUserNotFound) || errors.Is(err, model.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}

		// Handle other errors
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error: " + err.Error(),
		})
		return
	}

	// Success response
	c.JSON(http.StatusOK, output)
}

// Health returns API health status (placeholder)
func (h *IAMHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "tramatex-api-iam",
	})
}
