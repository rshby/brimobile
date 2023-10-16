package mock

import (
	"brimobile/app/account"
	"context"
	"github.com/stretchr/testify/mock"
)

type AccountRepositoryMock struct {
	Mock mock.Mock
}

// create function provider
func NewAccountRepoMock() *AccountRepositoryMock {
	return &AccountRepositoryMock{mock.Mock{}}
}

func (a AccountRepositoryMock) DeleteToken(ctx context.Context, refreshToken string) error {
	return a.Mock.Called(ctx, refreshToken).Error(0)
}

func (a AccountRepositoryMock) UpdateToken(ctx context.Context, uname string, accessToken string, refreshToken string) error {
	args := a.Mock.Called(ctx, uname, accessToken, refreshToken)
	return args.Error(0)
}

// implemented mock insert
func (a AccountRepositoryMock) Insert(ctx context.Context, entity *account.Account) (*account.Account, error) {
	args := a.Mock.Called(ctx, entity)

	// if error
	if args.Get(1) != nil {
		err := args.Get(1).(error)
		return nil, err
	}

	// success
	account := args.Get(0).(*account.Account)
	return account, nil
}

// implemented mock get by uname
func (a AccountRepositoryMock) GetByUname(ctx context.Context, uname string) (*account.Account, error) {
	args := a.Mock.Called(ctx, uname)

	// if error
	if args.Get(1) != nil {
		err := args.Get(1).(error)
		return nil, err
	}

	// success
	account := args.Get(0).(*account.Account)
	return account, nil
}
