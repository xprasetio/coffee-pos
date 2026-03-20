package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/xprasetio/coffee-pos/internal/entity"
	"github.com/xprasetio/coffee-pos/internal/service"
	"github.com/xprasetio/coffee-pos/pkg/response"
	"github.com/xprasetio/coffee-pos/pkg/validator"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authService *service.AuthService
	validator   *validator.Validator
}

// NewAuthHandler creates a new AuthHandler instance
func NewAuthHandler(authService *service.AuthService, v *validator.Validator) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   v,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var req entity.RegisterRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.ValidationError(c, errors)
		return
	}

	// Call service
	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, "Registrasi berhasil", result)
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Format request tidak valid")
		return
	}

	// Validate request
	if errors := h.validator.Validate(&req); errors != nil {
		response.ValidationError(c, errors)
		return
	}

	// Call service
	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, "Login berhasil", result)
}
