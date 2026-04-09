package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret []byte

// SetJWTSecret sets the JWT secret key used for signing tokens.
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

// Claims represents the JWT claims structure.
type Claims struct {
	// UserID is the ID of the user.
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for the given user ID.
func GenerateAccessToken(userID uuid.UUID) (string, error) {
	if jwtSecret == nil {
		return "", errors.New("jwt secret not set")
	}
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//use the variable, not the setter
	return token.SignedString(jwtSecret)
}

func GenerateRefreshToken() string{
	return uuid.New().String()
}

// ParseToken parses and validates a JWT token string.
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
