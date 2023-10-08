package repository

import (
	"brimobile/app/brinjournal"
	"context"
	"sync"
)

type IBrinjournalRepo interface {
	GetByBranch(ctx context.Context, wg *sync.WaitGroup, resChan chan brinjournal.BrinJournal, errChan chan error, branchNo string)
	UpdateOne(ctx context.Context, branchNo string) error
}
