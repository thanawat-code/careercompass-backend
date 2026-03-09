package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordMismatch   = errors.New("passwords do not match")
)

type AuthService struct {
	jwtSecret     string
	jwtExpiration time.Duration
}

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(jwtSecret string, jwtExpiration string) (*AuthService, error) {
	duration, err := time.ParseDuration(jwtExpiration)
	if err != nil {
		return nil, err
	}

	return &AuthService{
		jwtSecret:     jwtSecret,
		jwtExpiration: duration,
	}, nil
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword compares a password with a hash
func (s *AuthService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ValidatePasswordMatch checks if password and confirm password match
func (s *AuthService) ValidatePasswordMatch(password, confirmPassword string) error {
	if password != confirmPassword {
		return ErrPasswordMismatch
	}
	return nil
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(userID, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidCredentials
}
