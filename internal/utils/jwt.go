package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret []byte

// Setter function
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID) (string, error) {
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


func ParseToken(tokenStr string)(*Claims, error){
	token,err:=jwt.ParseWithClaims(tokenStr,&Claims{},func(token *jwt.Token) (interface{},error){
		return jwtSecret,nil
	})
	if err!=nil{
		return nil,err
	}
	claims,ok:=token.Claims.(*Claims)
	if !ok || !token.Valid{
		return nil, errors.New("invalid token")
	}
	return claims,nil
}

