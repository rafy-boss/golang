package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("streamsyc.122526")
// var tokenBlacklist = make(map[string]bool)
type Claims struct {
	Username string `json:"username"`
	UserId uint `json:"user_id"`
	jwt.RegisteredClaims
}
func GenerateJWT(username string,userId uint) (string, error) {
	expiration := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		UserId: userId,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func VerifyJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return token.Claims.(*Claims), nil
}