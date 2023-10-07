package repository

import (
	"brimobile/app/account"
	"context"
)

type IAccountRepository interface {
	Insert(ctx context.Context, entity *account.Account) (*account.Account, error)
	GetByUname(ctx context.Context, uname string) (*account.Account, error)
	UpdateToken(ctx context.Context, uname string, accessToken string, refreshToken string) error
	DeleteToken(ctx context.Context, refreshToken string) error
}
