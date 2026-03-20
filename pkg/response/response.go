package response

import (
	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// Meta represents pagination metadata
type Meta struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// Success sends a success response with custom status code
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Created sends a 201 Created response
func Created(c *gin.Context, message string, data interface{}) {
	Success(c, 201, message, data)
}

// OK sends a 200 OK response
func OK(c *gin.Context, message string, data interface{}) {
	Success(c, 200, message, data)
}

// Paginated sends a success response with pagination metadata
func Paginated(c *gin.Context, message string, data interface{}, meta Meta) {
	c.JSON(200, Response{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    &meta,
	})
}

// Error sends an error response with custom status code
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
	})
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFound sends a 404 Not Found response
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(c *gin.Context, message string) {
	Error(c, 500, message)
}

// ValidationError sends a 422 Unprocessable Entity response with validation errors
func ValidationError(c *gin.Context, errors interface{}) {
	c.JSON(422, Response{
		Success: false,
		Message: "Validasi gagal",
		Errors:  errors,
	})
}
