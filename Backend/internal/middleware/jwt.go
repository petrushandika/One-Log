package middleware

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/petrushandika/one-log/pkg/utils"
)

// JWTAuth middleware validates the access token from admin login
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		var tokenString string
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			var err error
			// Prefer documented cookie name; fallback to legacy cookie for compatibility.
			tokenString, err = c.Cookie("ulam_access")
			if err != nil {
				tokenString, err = c.Cookie("token")
				if err != nil {
					utils.Error(c, http.StatusUnauthorized, "Missing Authorization token", nil)
					c.Abort()
					return
				}
			}
		}
		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Validate algo
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			utils.Error(c, http.StatusUnauthorized, "Invalid or expired token", nil)
			c.Abort()
			return
		}

		// Save payload
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("admin_email", claims["email"])
			if uid, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", uint(uid))
			}
		}

		c.Next()
	}
}
