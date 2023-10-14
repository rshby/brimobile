package helper

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPassword(ctx context.Context, passwordDb, inputPassword string) bool {
	span, _ := opentracing.StartSpanFromContext(ctx, "CheckPassword")
	defer span.Finish()

	isOk := passwordDb == inputPassword
	span.LogFields(
		log.String("request-password-db", passwordDb),
		log.String("request-input-password", inputPassword),
		log.Bool("response-match", isOk),
	)

	return isOk
}
