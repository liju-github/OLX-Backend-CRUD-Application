package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte("olxsecret")

type Claims struct {
	UserEmail string `json:"useremail"`
	jwt.StandardClaims
}

func GenerateJWT(UserEmail string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserEmail: UserEmail,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
