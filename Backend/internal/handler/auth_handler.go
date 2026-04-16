package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/service"
	"github.com/petrushandika/one-log/pkg/utils"
)

// AuthHandler handles the admin authentication endpoints.
// Business logic (credential validation, JWT signing) lives in AuthService.
type AuthHandler struct {
	authSvc service.AuthService
	logSvc  service.LogService
}

func NewAuthHandler(authSvc service.AuthService, logSvc service.LogService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc, logSvc: logSvc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusUnprocessableEntity, "Validation failed", []utils.ErrorDetail{
			{Field: "body", Message: err.Error()},
		})
		return
	}

	user, accessToken, refreshToken, err := h.authSvc.Login(req.Email, req.Password)
	if err != nil {
		h.logFailedAttempt(c, req.Email)
		utils.Error(c, http.StatusUnauthorized, "Invalid credentials", nil)
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

	// Set httpOnly cookies.
	// NOTE: secure flag should be true in production behind HTTPS.
	c.SetCookie("ulam_access", accessToken, 3600*24, "/", "", false, true)
	c.SetCookie("ulam_refresh", refreshToken, 3600*24*7, "/api/auth/refresh", "", false, true)
	c.SetCookie("token", accessToken, 3600*24, "/", "", false, true) // legacy

	utils.Success(c, http.StatusOK, "Login successful", gin.H{
		// Keep returning token for current frontend compatibility (will be removed once frontend switches fully to cookies).
		"token": accessToken,
		"email": user.Email,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshTokenString, err := c.Cookie("ulam_refresh")
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Missing refresh token", nil)
		return
	}

	newAccessToken, err := h.authSvc.RefreshAccessToken(refreshTokenString)
	if err != nil {
		utils.Error(c, http.StatusUnauthorized, "Invalid or expired refresh token", nil)
		return
	}

	c.SetCookie("ulam_access", newAccessToken, 3600*24, "/", "", false, true)
	c.SetCookie("token", newAccessToken, 3600*24, "/", "", false, true) // legacy

	utils.Success(c, http.StatusOK, "Token refreshed", nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear all auth cookies
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
