package service

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/petrushandika/one-log/internal/domain"
	"github.com/petrushandika/one-log/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// AuthService encapsulates all authentication business logic.
type AuthService interface {
	Login(email, password string) (*domain.User, string, string, error)
	RefreshAccessToken(refreshTokenString string) (string, error)
}

type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService constructs an AuthService with the given UserRepository.
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Login validates credentials and returns the user plus signed JWT token strings.
// Returns (user, accessTokenString, refreshTokenString, error).
func (s *authService) Login(email, password string) (*domain.User, string, string, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	secret := []byte(os.Getenv("JWT_SECRET"))

	// Access token — valid 24 h
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(secret)
	if err != nil {
		return nil, "", "", err
	}

	// Refresh token — valid 7 days
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"typ":     "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(secret)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessTokenString, refreshTokenString, nil
}

// RefreshAccessToken validates a refresh token and issues a new access token string.
func (s *authService) RefreshAccessToken(refreshTokenString string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(refreshTokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil || claims["email"] == nil {
		return "", errors.New("invalid refresh token claims")
	}
	if typ, _ := claims["typ"].(string); typ != "refresh" {
		return "", errors.New("invalid refresh token type")
	}

	// Re-issue access token (24 h)
	newAccessClaims := jwt.MapClaims{
		"user_id": claims["user_id"],
		"email":   claims["email"],
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newAccessClaims)
	return newToken.SignedString(secret)
}
