package service

import (
	"brimobile/graph/model"
	"context"
)

type IAccountService interface {
	CreateAccount(ctx context.Context, uname string, pass string) (*model.CreateAccountResponse, error)
	Login(ctx context.Context, uname string, pass string, idNum string, deviceID string) (*model.LoginResponse, error)
	Logout(ctx context.Context, refreshToken string) (string, error)
	Account(ctx context.Context, uname string) (*model.AccountResponse, error)
	InqAccountSaving(ctx context.Context, accountNumber string) (*model.InqAccountSaving, error)
}
