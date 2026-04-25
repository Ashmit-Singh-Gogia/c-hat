package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint
	jwt.RegisteredClaims
}

func CreateToken(userid uint, secret string) (string, error) {
	claims := Claims{
		UserID: userid,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 10)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(token string, secret string) (jwt.Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	return claims, nil
}
