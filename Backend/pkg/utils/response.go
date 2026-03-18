package utils

import "github.com/gin-gonic/gin"

type APIResponse struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type ErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Status:  "success",
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Status:  "error",
		Code:    code,
		Message: message,
		Data:    data,
	})
}
