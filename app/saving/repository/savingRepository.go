package repository

import (
	"brimobile/app/saving"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"sync"
	"time"
)

type SavingRepository struct {
	DB *sql.DB
}

// function provider
func NewSavingRepository(db *sql.DB) *SavingRepository {
	return &SavingRepository{
		db,
	}
}

func (s *SavingRepository) Insert(ctx context.Context, entity *saving.Saving) (*saving.Saving, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "repository Insert Saving")
	defer span.Finish()

	query := "INSERT INTO saving (account_number, account_type, branch_code, short_name, currency, cbal, hold, opening_date, product_group, product_name, status) VALUES ($1, 'S', '0999', $2, 'USD', $3, '0.00', $4, '10001', 'Britama Saving', '1');"

	todayDate := time.Now().Local().Format("2006-01-02 15:04:05")

	_, err := s.DB.ExecContext(ctxTracing, query, entity.AccountNumber, entity.ShortName, entity.Cbal, todayDate)
	if err != nil {
		return nil, err
	}

	wg := &sync.WaitGroup{}
	savingChan, errChan := make(chan saving.Saving, 1), make(chan error, 1)

	wg.Add(1)
	go s.GetByAccountNumber(ctx, wg, savingChan, errChan, entity.AccountNumber)
	wg.Wait()

	err = <-errChan
	saving := <-savingChan
	if err != nil {
		return nil, err
	}

	// success insert
	span.SetTag("result_data", saving)
	return &saving, nil
}

func (s *SavingRepository) GetByAccountNumber(ctx context.Context, wg *sync.WaitGroup, resChan chan saving.Saving, errChan chan error, accountNumber string) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "repository GetByAccountNumber")
	defer span.Finish()

	defer wg.Done()
	query := "SELECT account_number, account_type, branch_code, short_name, currency, cbal, hold, opening_date, product_group, product_name, status FROM saving WHERE account_number=$1"

	row := s.DB.QueryRowContext(ctxTracing, query, accountNumber)
	if row.Err() != nil {
		resChan <- saving.Saving{}
		errChan <- errors.New("record not found")
		return
	}

	var sv saving.Saving
	if err := row.Scan(&sv.AccountNumber, &sv.AccountType, &sv.BranchCode, &sv.ShortName, &sv.Currency, &sv.Cbal, &sv.Hold, &sv.OpeningDate, &sv.ProductGroup, &sv.ProductName, &sv.Status); err != nil {
		resChan <- saving.Saving{}
		errChan <- errors.New("record not found")
		return
	}

	span.SetTag("result_data", sv)

	// success
	resChan <- sv
	errChan <- nil
}

func (s *SavingRepository) UpdateCbal(ctx context.Context, wg *sync.WaitGroup, errChan chan error, accountNumber string, cbal float64) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "repository UpdateCbal")
	defer func() {
		span.Finish()
		wg.Done()
	}()

	span.SetTag("account_number", accountNumber)
	span.SetTag("cbal", cbal)

	query := "UPDATE saving set cbal=$1 where account_number=$2"

	_, err := s.DB.ExecContext(ctxTracing, query, fmt.Sprintf("%.7f", cbal), accountNumber)
	if err != nil {
		errChan <- err
		return
	}

	// success update
	errChan <- nil
}
