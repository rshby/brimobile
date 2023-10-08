package repository

import (
	"brimobile/app/saving"
	"context"
	"database/sql"
	"errors"
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
	query := "INSERT INTO saving (account_number, account_type, branch_code, short_name, currency, cbal, hold, opening_date, product_group, product_name, status) VALUES ($1, 'S', '0999', $2, 'USD', $3, '0.00', $4, '10001', 'Britama Saving', '1');"

	todayDate := time.Now().Local().Format("2006-01-02 15:04:05")

	_, err := s.DB.ExecContext(ctx, query, entity.AccountNumber, entity.ShortName, entity.Cbal, todayDate)
	if err != nil {
		return nil, err
	}

	saving, err := s.GetByAccountNumber(ctx, entity.AccountNumber)
	if err != nil {
		return nil, err
	}

	// success insert
	return saving, nil
}

func (s *SavingRepository) GetByAccountNumber(ctx context.Context, accountNumber string) (*saving.Saving, error) {
	query := "SELECT account_number, account_type, branch_code, short_name, currency, cbal, hold, opening_date, product_group, product_name, status FROM saving WHERE account_number=$1"

	row := s.DB.QueryRowContext(ctx, query, accountNumber)
	if row.Err() != nil {
		return nil, errors.New("record not found")
	}

	var saving saving.Saving
	if err := row.Scan(&saving.AccountNumber, &saving.AccountType, &saving.BranchCode, &saving.ShortName, &saving.Currency, &saving.Cbal, &saving.Hold, &saving.OpeningDate, &saving.ProductGroup, &saving.ProductName, &saving.Status); err != nil {
		return nil, err
	}

	// success
	return &saving, nil
}
