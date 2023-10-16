package helper

import (
	"brimobile/app/auth"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"os"
	"sync"
	"time"
)

func GenerateToken(ctx context.Context, uname string, hour time.Duration, wg *sync.WaitGroup, token chan<- string) (string, error) {
	defer wg.Done()

	span, _ := opentracing.StartSpanFromContext(ctx, "GenerateToken")
	defer span.Finish()

	// create claims
	claims := auth.JwtClaims{
		Uname: uname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "brmobile",
			ExpiresAt: jwt.NewNumericDate(time.Now().Local().Add(hour)),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	tokenString, err := tokenClaims.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	// success
	token <- tokenString

	span.LogFields(
		log.Object("request-claims", claims),
		log.String("response-token", tokenString))
	return tokenString, nil
}

func GetClaims(ctx context.Context, tokenString string) (auth.JwtClaims, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "Helper GetClaims")
	defer span.Finish()

	span.LogFields(
		log.String("request-token", tokenString))
	claims := auth.JwtClaims{}
	_, _ = jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	return claims, nil
}
