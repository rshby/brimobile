package repository

import (
	"brimobile/app/brinjournal"
	"context"
	"database/sql"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"sync"
)

type BrinJournalRepository struct {
	DB *sql.DB
}

func NewBrinJournalRepository(db *sql.DB) *BrinJournalRepository {
	return &BrinJournalRepository{db}
}

// method get data journalseq
func (b *BrinJournalRepository) GetByBranch(ctx context.Context, wg *sync.WaitGroup, resChan chan brinjournal.BrinJournal, errChan chan error, branchNo string) {
	defer wg.Done()

	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "BrinJournalRepository GetByBranch")
	defer span.Finish()

	query := "SELECT id, branch_code, joirnalseq FROM brinjournalseq Where branch_code=$1"

	row := b.DB.QueryRowContext(ctxTracing, query, branchNo)
	if row.Err() != nil {
		resChan <- brinjournal.BrinJournal{}
		errChan <- row.Err()
		return
	}

	var resultBrinjournal brinjournal.BrinJournal
	err := row.Scan(&resultBrinjournal.Id, &resultBrinjournal.BranchCode, &resultBrinjournal.JournalSeq)
	if err != nil {
		resChan <- brinjournal.BrinJournal{}
		errChan <- err
		return
	}

	// success
	span.LogFields(
		log.String("request-branchNo", branchNo),
		log.Object("response-brinJournal", resultBrinjournal),
		log.Object("response-error", nil),
	)
	resChan <- resultBrinjournal
	errChan <- nil
}

// method update journalseq + 1
func (b *BrinJournalRepository) UpdateOne(ctx context.Context, branchNo string) error {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "BrinJournalRepository UpdateOne")
	defer span.Finish()

	query := "UPDATE brinjournalseq set joirnalseq=joirnalseq+1 WHERE branch_code=$1"

	_, err := b.DB.ExecContext(ctxTracing, query, branchNo)

	span.LogFields(
		log.String("request-branchNo", branchNo),
		log.Object("response-error", nil),
	)
	return err
}
