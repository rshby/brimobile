package auth

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	Uname            string
	RegisteredClaims jwt.RegisteredClaims
}

func (j JwtClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (j JwtClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (j JwtClaims) GetNotBefore() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (j JwtClaims) GetIssuer() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (j JwtClaims) GetSubject() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (j JwtClaims) GetAudience() (jwt.ClaimStrings, error) {
	//TODO implement me
	panic("implement me")
}
