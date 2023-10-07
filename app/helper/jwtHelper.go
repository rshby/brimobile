package helper

import (
	"brimobile/app/auth"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"sync"
	"time"
)

func GenerateToken(uname string, hour time.Duration, wg *sync.WaitGroup, token chan<- string) (string, error) {
	defer wg.Done()

	// create claims
	claims := auth.JwtClaims{
		Uname: uname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "brmobile",
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(hour)),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := tokenClaims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	// success
	token <- tokenString
	return tokenString, nil
}
