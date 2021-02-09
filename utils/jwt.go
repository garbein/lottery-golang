package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserId int
	jwt.StandardClaims
}

func GenerateAccessToken(userId int) (string, error) {
	var claims = Claims{UserId: userId, StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + 30*86400,
	}}
	claimsToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := claimsToken.SignedString([]byte("5432fin4303f00994qq0afgj44e400s"))
	return token, err
}

func ParseAccessToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("5432fin4303f00994qq0afgj44e400s"), nil
	})

	if err == nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
