package mock

import (
	"brimobile/app/brinjournal"
	"context"
	"github.com/stretchr/testify/mock"
	"sync"
)

type BrinJournalRepoMock struct {
	Mock mock.Mock
}

// func NewBrinJournalMock
func NewBrinJournalMock() *BrinJournalRepoMock {
	return &BrinJournalRepoMock{
		mock.Mock{},
	}
}

func (b *BrinJournalRepoMock) GetByBranch(ctx context.Context, wg *sync.WaitGroup, resChan chan brinjournal.BrinJournal, errChan chan error, branchNo string) {
	args := b.Mock.Called(ctx, wg, resChan, errChan, branchNo)
	resChan <- args.Get(0).(brinjournal.BrinJournal)
	errChan <- args.Error(1)
}

func (b *BrinJournalRepoMock) UpdateOne(ctx context.Context, branchNo string) error {
	// bisa error, bisa nil
	args := b.Mock.Called(ctx, branchNo)
	return args.Error(0)
}
