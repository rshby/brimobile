package mock

import (
	"brimobile/app/saving"
	"context"
	"github.com/stretchr/testify/mock"
)

type SavingRepositoryMock struct {
	Mock mock.Mock
}

// function provider
func NewSavingRepositoryMock() *SavingRepositoryMock {
	return &SavingRepositoryMock{
		mock.Mock{},
	}
}

func (s *SavingRepositoryMock) Insert(ctx context.Context, entity *saving.Saving) (*saving.Saving, error) {
	args := s.Mock.Called(ctx, entity)

	// jika error
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}

	// jika success
	return args.Get(0).(*saving.Saving), nil
}

func (s *SavingRepositoryMock) GetByAccountNumber(ctx context.Context, accountNumber string) (*saving.Saving, error) {
	args := s.Mock.Called(ctx, accountNumber)

	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}

	// success
	return args.Get(0).(*saving.Saving), nil
}
