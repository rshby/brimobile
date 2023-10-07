package repository

import (
	"brimobile/app/account"
	"context"
	"database/sql"
)

type AccountRepository struct {
	DB *sql.DB
}

// function provider
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db,
	}
}

// method delete token
func (a *AccountRepository) DeleteToken(ctx context.Context, refreshToken string) error {
	query := "UPDATE accounts set access_token=null, refresh_token=null WHERE refresh_token=$1"

	_, err := a.DB.ExecContext(ctx, query, refreshToken)
	return err
}

// method create account
func (a *AccountRepository) Insert(ctx context.Context, entity *account.Account) (*account.Account, error) {
	query := "INSERT INTO accounts(uname, pass) VALUES($1, $2) RETURNING id"

	res := a.DB.QueryRowContext(ctx, query, entity.Uname, entity.Pass)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var id int
	if err := res.Scan(&id); err != nil {
		return nil, err
	}

	entity.Id = int(id)
	return entity, nil
}

// method get by username
func (a *AccountRepository) GetByUname(ctx context.Context, uname string) (*account.Account, error) {
	query := "SELECT id, uname, pass, access_token, refresh_token FROM accounts WHERE uname=$1"

	row := a.DB.QueryRowContext(ctx, query, uname)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var account account.Account
	if err := row.Scan(&account.Id, &account.Uname, &account.Pass, &account.AccessToken, &account.RefreshToken); err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountRepository) UpdateToken(ctx context.Context, uname string, accessToken string, refreshToken string) error {
	query := "UPDATE accounts set access_token=$1, refresh_token=$2 WHERE uname=$3"

	_, err := a.DB.ExecContext(ctx, query, accessToken, refreshToken, uname)
	return err
}
