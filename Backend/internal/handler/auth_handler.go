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

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   req.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	// Refresh token (Valid for 7 days)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   req.Email,
		"typ":     "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})
	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate refresh token", err.Error())
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

	// Set documented httpOnly cookies (hybrid: keep legacy cookie for compatibility).
	// NOTE: secure flag should be true in production behind HTTPS.
	c.SetCookie("ulam_access", accessTokenString, 3600*24, "/", "", false, true)
	c.SetCookie("ulam_refresh", refreshTokenString, 3600*24*7, "/api/auth/refresh", "", false, true)
	c.SetCookie("token", accessTokenString, 3600*24, "/", "", false, true) // legacy

	utils.Success(c, http.StatusOK, "Login successful", gin.H{
		// Keep returning token for current frontend compatibility (will be removed once frontend switches to cookies).
		"token": accessTokenString,
		"email": req.Email,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshTokenString, err := c.Cookie("ulam_refresh")
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Missing refresh token", nil)
		return
	}

	secret := []byte(os.Getenv("JWT_SECRET"))
	token, err := jwt.Parse(refreshTokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		utils.Error(c, http.StatusUnauthorized, "Invalid or expired refresh token", nil)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil || claims["email"] == nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid refresh token claims", nil)
		return
	}
	if typ, _ := claims["typ"].(string); typ != "refresh" {
		utils.Error(c, http.StatusUnauthorized, "Invalid refresh token type", nil)
		return
	}

	// Re-issue access token (24h)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": claims["user_id"],
		"email":   claims["email"],
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to generate token", err.Error())
		return
	}

	c.SetCookie("ulam_access", accessTokenString, 3600*24, "/", "", false, true)
	c.SetCookie("token", accessTokenString, 3600*24, "/", "", false, true) // legacy

	utils.Success(c, http.StatusOK, "Token refreshed", nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear cookies by setting MaxAge<0
	c.SetCookie("ulam_access", "", -1, "/", "", false, true)
	c.SetCookie("ulam_refresh", "", -1, "/api/auth/refresh", "", false, true)
	c.SetCookie("token", "", -1, "/", "", false, true) // legacy
	utils.Success(c, http.StatusOK, "Logged out successfully", nil)
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
