package mock

import (
	"brimobile/app/saving"
	"context"
	"github.com/stretchr/testify/mock"
	"sync"
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

func (s *SavingRepositoryMock) GetByAccountNumber(ctx context.Context, wg *sync.WaitGroup, resChan chan saving.Saving, errChan chan error, accountNumber string) {
	wg.Add(1)
	defer wg.Done()

	args := s.Mock.Called(ctx, wg, resChan, errChan, accountNumber)

	resChan <- args.Get(0).(saving.Saving)
	errChan <- args.Error(1)
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

func (s *SavingRepositoryMock) UpdateCbal(ctx context.Context, wg *sync.WaitGroup, errChan chan error, accountNumber string, cbal float64) {

}
