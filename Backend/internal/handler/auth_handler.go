package handler

import (
	"crypto/subtle"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/pkg/utils"
)

// AuthHandler handles the admin login process.
// In this MVP version, credentials are only taken from environment variables for simplicity.
type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
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

	// Default values if .env is missing
	if adminEmail == "" {
		adminEmail = "admin@ulam.io"
		adminPassword = "adminpassword"
	}

	// Use constant time comparison to mitigate timing attacks on the login route
	emailMatch := subtle.ConstantTimeCompare([]byte(req.Email), []byte(adminEmail)) == 1
	passMatch := subtle.ConstantTimeCompare([]byte(req.Password), []byte(adminPassword)) == 1

	if !emailMatch || !passMatch {
		utils.Error(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Create token (Valid for 24 hours)
	secret := []byte(os.Getenv("JWT_SECRET"))
	if len(secret) == 0 {
		secret = []byte("mvp-super-secret-key-123") // Fallback
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": req.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	utils.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": tokenString,
		"email": req.Email,
	})
}
