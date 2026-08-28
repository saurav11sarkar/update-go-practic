package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrEmptyID    = errors.New("id is empty")
	ErrEmptyEmail = errors.New("email is empty")
	ErrEmptyRole  = errors.New("role is empty")
)

type JwtClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`

	jwt.RegisteredClaims
}

// JwtClime is kept as an alias for backward compatibility.
type JwtClime = JwtClaims

type JwtService interface {
	GenerateToken(id, email, role string) (string, error)
	ValidateToken(token string) (*JwtClaims, error)
}

type jwtService struct {
	secretKey     string
	tokenDuration time.Duration
}

func NewJwtService(secretKey string, tokenDuration time.Duration) JwtService {
	return &jwtService{
		secretKey:     secretKey,
		tokenDuration: tokenDuration,
	}
}

func (j *jwtService) GenerateToken(id, email, role string) (string, error) {
	if id == "" {
		return "", ErrEmptyID
	}
	if email == "" {
		return "", ErrEmptyEmail
	}
	if role == "" {
		return "", ErrEmptyRole
	}

	claims := JwtClaims{
		ID:    id,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-service",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *jwtService) ValidateToken(tokenString string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %s", t.Method.Alg())
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if claims, ok := token.Claims.(*JwtClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
