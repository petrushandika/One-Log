package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthHandler handles the admin login process.
type AuthHandler struct {
	db     *gorm.DB
	logSvc service.LogService
}

func NewAuthHandler(db *gorm.DB, logSvc service.LogService) *AuthHandler {
	return &AuthHandler{db: db, logSvc: logSvc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "Validation failed", []utils.ErrorDetail{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	var user domain.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		h.logFailedAttempt(c, req.Email)
		utils.Error(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		h.logFailedAttempt(c, req.Email)
		utils.Error(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Create token (Valid for 24 hours)
	secret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   req.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
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

	// Set HTTP Only Cookie
	// name, value, maxAge (sec), path, domain, secure, httpOnly
	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)

	utils.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": tokenString,
		"email": req.Email,
	})
}

func (h *AuthHandler) logFailedAttempt(c *gin.Context, email string) {
	clientIP := c.ClientIP()
	_ = h.logSvc.IngestLog(domain.IngestLogRequest{
		Category:  "AUTH_EVENT",
		Level:     "WARN",
		Message:   "Failed login attempt for admin panel",
		IPAddress: clientIP,
		Context:   map[string]interface{}{"attempted_email": email},
	}, "00000000-0000-0000-0000-000000000001")

	// Brute Force Detection
	isBruteForce, _ := h.logSvc.CheckBruteForce(clientIP)
	if isBruteForce {
		_ = h.logSvc.IngestLog(domain.IngestLogRequest{
			Category:  "SECURITY",
			Level:     "CRITICAL",
			Message:   "Brute force attempt detected from IP: " + clientIP,
			IPAddress: clientIP,
			Context:   map[string]interface{}{"note": "Exceeded 5 failed attempts in 10 minutes"},
		}, "00000000-0000-0000-0000-000000000001")
	}
}
