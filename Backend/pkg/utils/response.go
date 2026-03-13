package utils

import (
	"github.com/gin-gonic/gin"
)

// Response structure for all API calls
type APIResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ErrorDetail represents specific field errors
type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Success sends a standard success response
func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Error sends a standard error response
func Error(c *gin.Context, code int, message string, errors interface{}) {
	c.JSON(code, APIResponse{
		Status:  "error",
		Code:    code,
		Message: message,
		Errors:  errors,
	})
}
