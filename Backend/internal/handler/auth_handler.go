package handler

import (
	"crypto/subtle"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

// AuthHandler handles the admin login process.
type AuthHandler struct {
	logSvc service.LogService
}

func NewAuthHandler(logSvc service.LogService) *AuthHandler {
	return &AuthHandler{logSvc: logSvc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "Validation failed", []utils.ErrorDetail{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	// Use constant time comparison to mitigate timing attacks on the login route
	emailMatch := subtle.ConstantTimeCompare([]byte(req.Email), []byte(adminEmail)) == 1
	passMatch := subtle.ConstantTimeCompare([]byte(req.Password), []byte(adminPassword)) == 1

	if !emailMatch || !passMatch {
		// Log failed attempt audit trail
		_ = h.logSvc.IngestLog(domain.IngestLogRequest{
			Category:  "AUTH_EVENT",
			Level:     "WARN",
			Message:   "Failed login attempt for admin panel",
			IPAddress: c.ClientIP(),
			Context:   map[string]interface{}{"attempted_email": req.Email},
		}, "00000000-0000-0000-0000-000000000001") // Mock ULAM Internal UUID

		utils.Error(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Create token (Valid for 24 hours)
	secret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	// Log successful login audit trail
	_ = h.logSvc.IngestLog(domain.IngestLogRequest{
		Category:  "AUTH_EVENT",
		Level:     "INFO",
		Message:   "Admin logged in successfully",
		IPAddress: c.ClientIP(),
		Context:   map[string]interface{}{"email": req.Email},
	}, "00000000-0000-0000-0000-000000000001")

	utils.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": tokenString,
		"email": req.Email,
	})
}
