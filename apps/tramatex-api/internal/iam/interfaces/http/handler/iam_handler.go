package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joran-cortez/tramatex/internal/iam/application/usecase"
	"github.com/joran-cortez/tramatex/internal/iam/domain/model"
	"github.com/joran-cortez/tramatex/internal/shared/infrastructure/logging"
)

// IAMHandler handles authentication and user management endpoints.
type IAMHandler struct {
	loginUseCase      *usecase.LoginUseCase
	registerUseCase   *usecase.RegisterUserUseCase
	createUserUseCase *usecase.CreateUserUseCase
	refreshUseCase    *usecase.RefreshTokenUseCase
	logoutUseCase     *usecase.LogoutUserUseCase
	assignRoleUseCase *usecase.AssignRoleUseCase
	checkAuthUseCase  *usecase.CheckAuthorizationUseCase
	listUsersUseCase  *usecase.ListUsersUseCase
	deleteUserUseCase *usecase.DeleteUserUseCase
	updateUserUseCase *usecase.UpdateUserUseCase
}

// NewIAMHandler creates a new IAM handler.
func NewIAMHandler(
	loginUseCase *usecase.LoginUseCase,
	registerUseCase *usecase.RegisterUserUseCase,
	createUserUseCase *usecase.CreateUserUseCase,
	refreshUseCase *usecase.RefreshTokenUseCase,
	logoutUseCase *usecase.LogoutUserUseCase,
	assignRoleUseCase *usecase.AssignRoleUseCase,
	checkAuthUseCase *usecase.CheckAuthorizationUseCase,
	listUsersUseCase *usecase.ListUsersUseCase,
	deleteUserUseCase *usecase.DeleteUserUseCase,
	updateUserUseCase *usecase.UpdateUserUseCase,
) *IAMHandler {
	return &IAMHandler{
		loginUseCase:      loginUseCase,
		registerUseCase:   registerUseCase,
		createUserUseCase: createUserUseCase,
		refreshUseCase:    refreshUseCase,
		logoutUseCase:     logoutUseCase,
		assignRoleUseCase: assignRoleUseCase,
		checkAuthUseCase:  checkAuthUseCase,
		listUsersUseCase:  listUsersUseCase,
		deleteUserUseCase: deleteUserUseCase,
		updateUserUseCase: updateUserUseCase,
	}
}

// Login handles POST /auth/login
func (h *IAMHandler) Login(c *gin.Context) {
	var req usecase.LoginInput

	// Get request ID for logging
	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	// Validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		logging.WithRequestID(reqID).WithField("error", err.Error()).Warn("Login failed: invalid request")
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
			// Log failed login attempt (security event)
			logging.WithRequestID(reqID).WithFields(map[string]interface{}{
				"email":  logging.MaskEmail(req.Email),
				"reason": "invalid_credentials",
				"ip":     c.ClientIP(),
			}).Warn("Login failed: invalid credentials")

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}

		// Handle other errors
		logging.WithRequestID(reqID).WithError(err).Error("Login failed: internal error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error: " + err.Error(),
		})
		return
	}

	// Log successful login (security event)
	logging.WithRequestID(reqID).WithFields(map[string]interface{}{
		"userID":    output.User.ID,
		"userEmail": logging.MaskEmail(output.User.Email),
		"userRole":  output.User.Role,
		"ip":        c.ClientIP(),
	}).Info("Login successful")

	// Success response
	c.JSON(http.StatusOK, output)
}

// Register handles POST /auth/register
func (h *IAMHandler) Register(c *gin.Context) {
	var req usecase.RegisterInput

	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.WithRequestID(reqID).WithField("error", err.Error()).Warn("Register failed: invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	output, err := h.registerUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "user already exists",
			})
			return
		}

		logging.WithRequestID(reqID).WithError(err).Error("Register failed: internal error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, output)
}

// CreateUser handles POST /auth/users
func (h *IAMHandler) CreateUser(c *gin.Context) {
	var req usecase.CreateUserInput

	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.WithRequestID(reqID).WithField("error", err.Error()).Warn("CreateUser failed: invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	output, err := h.createUserUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "user already exists",
			})
			return
		}
		logging.WithRequestID(reqID).WithError(err).Warn("CreateUser failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, output)
}

// Refresh handles POST /auth/refresh
func (h *IAMHandler) Refresh(c *gin.Context) {
	var req usecase.RefreshInput

	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.WithRequestID(reqID).WithField("error", err.Error()).Warn("Refresh failed: invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	output, err := h.refreshUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		logging.WithRequestID(reqID).WithError(err).Warn("Refresh failed: invalid token")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, output)
}

// Logout handles POST /auth/logout
func (h *IAMHandler) Logout(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
		return
	}

	var req usecase.LogoutInput
	_ = c.ShouldBindJSON(&req)

	if err := h.logoutUseCase.Execute(c.Request.Context(), parts[1], req); err != nil {
		logging.WithRequestID(reqID).WithError(err).Warn("Logout failed")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		return
	}

	c.Status(http.StatusNoContent)
}

// AssignRole handles POST /auth/assign-role
func (h *IAMHandler) AssignRole(c *gin.Context) {
	var req usecase.AssignRoleInput

	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.WithRequestID(reqID).WithField("error", err.Error()).Warn("AssignRole failed: invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	output, err := h.assignRoleUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		logging.WithRequestID(reqID).WithError(err).Warn("AssignRole failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, output)
}

// CheckAuthorization handles POST /auth/authorize
func (h *IAMHandler) CheckAuthorization(c *gin.Context) {
	var req usecase.CheckAuthorizationInput

	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.WithRequestID(reqID).WithField("error", err.Error()).Warn("Authorize failed: invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}

	// If user_id not provided, use authenticated user from context
	if req.UserID == "" {
		if userID, ok := c.Get("userID"); ok {
			if idStr, ok := userID.(string); ok {
				req.UserID = idStr
			}
		}
	}

	output, err := h.checkAuthUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		logging.WithRequestID(reqID).WithError(err).Warn("Authorize failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, output)
}

// UpdateUser handles PUT /auth/users/:id
func (h *IAMHandler) UpdateUser(c *gin.Context) {
	var req usecase.UpdateUserInput

	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logging.WithRequestID(reqID).WithField("error", err.Error()).Warn("UpdateUser failed: invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request: " + err.Error(),
		})
		return
	}
	req.ID = userID

	output, err := h.updateUserUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if errors.Is(err, model.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "email already in use by another user"})
			return
		}
		logging.WithRequestID(reqID).WithError(err).Warn("UpdateUser failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, output)
}

// DeleteUser handles DELETE /auth/users/:id
func (h *IAMHandler) DeleteUser(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id is required"})
		return
	}

	if err := h.deleteUserUseCase.Execute(c.Request.Context(), userID); err != nil {
		logging.WithRequestID(reqID).WithError(err).Warn("DeleteUser failed")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ListUsers handles GET /auth/users
func (h *IAMHandler) ListUsers(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	reqID, _ := requestID.(string)

	users, err := h.listUsersUseCase.Execute(c.Request.Context())
	if err != nil {
		logging.WithRequestID(reqID).WithError(err).Error("ListUsers failed")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// Health returns API health status (placeholder)
func (h *IAMHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "tramatex-api-iam",
	})
}
