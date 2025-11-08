package lib

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

func GenerateJWT(name, email, role string, permissions []string) (string, error) {
	cfg := loadConfig()
	jwtKey := cfg.JWT_SECRET
	expiration := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Name:        name,
		Email:       email,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenStr string) (*Claims, error) {
	cfg := loadConfig()
	jwtKey := cfg.JWT_SECRET
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return claims, nil
}

func VerifyJWT(tokenStr string) (*Claims, error) {
	cfg := loadConfig()
	jwtKey := cfg.JWT_SECRET

	// Remove possible "Bearer " prefix (common in Authorization headers)
	if after, ok := strings.CutPrefix(tokenStr, "Bearer "); ok {
		tokenStr = after
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token is signed with HMAC and expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("malformed token")
		}
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token expired")
		}
		if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("invalid token signature")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
