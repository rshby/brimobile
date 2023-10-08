package repository

import (
	"brimobile/app/saving"
	"context"
	"sync"
)

type ISavingRepository interface {
	Insert(ctx context.Context, entity *saving.Saving) (*saving.Saving, error)
	GetByAccountNumber(ctx context.Context, wg *sync.WaitGroup, resChan chan saving.Saving, errChan chan error, accountNumber string)
	UpdateCbal(ctx context.Context, wg *sync.WaitGroup, errChan chan error, accountNumber string, cbal float64)
}
