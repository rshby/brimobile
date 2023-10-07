package repository

import (
	"brimobile/app/saving"
	"context"
)

type ISavingRepository interface {
	Insert(ctx context.Context, entity *saving.Saving) (*saving.Saving, error)
	GetByAccountNumber(ctx context.Context, accountNumber string) (*saving.Saving, error)
}
