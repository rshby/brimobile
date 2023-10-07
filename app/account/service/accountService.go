package service

import (
	"brimobile/app/account"
	"brimobile/app/account/repository"
	"brimobile/app/helper"
	"brimobile/graph/model"
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"
)

type AccountService struct {
	AccountRepo repository.IAccountRepository
}

// function provider
func NewAccountService(accRepo repository.IAccountRepository) *AccountService {
	return &AccountService{
		AccountRepo: accRepo,
	}
}

func (a *AccountService) CreateAccount(ctx context.Context, uname string, pass string) (*model.CreateAccountResponse, error) {
	// proses hash password

	entity := account.Account{
		Uname: uname,
		Pass:  pass,
		AccessToken: sql.NullString{
			Valid: false,
		},
		RefreshToken: sql.NullString{
			Valid: false,
		},
	}

	account, err := a.AccountRepo.Insert(ctx, &entity)
	if err != nil {
		return nil, err
	}

	// success
	response := model.CreateAccountResponse{
		ID:    account.Id,
		Uname: account.Uname,
		Pass:  account.Pass,
	}

	return &response, nil
}

func (a *AccountService) Login(ctx context.Context, uname string, pass string, idNum string, deviceID string) (*model.LoginResponse, error) {
	// cek apakah uname ada di database
	account, err := a.AccountRepo.GetByUname(ctx, uname)
	if err != nil {
		return nil, errors.New("record not found")
	}

	// cek apakah password benar
	if !helper.CheckPassword(account.Pass, pass) {
		return nil, errors.New("password not match")
	}

	// cek apakah sudah login
	if account.AccessToken.Valid {
		return nil, errors.New("user has been logged in")
	}

	// lolos -> generate access_oken dan refresh_token
	wg := &sync.WaitGroup{}
	accessToken := make(chan string, 1)
	refreshToken := make(chan string, 1)

	wg.Add(2)
	go helper.GenerateToken(uname, time.Duration(1*time.Hour), wg, accessToken)
	go helper.GenerateToken(uname, time.Duration(5*time.Hour), wg, refreshToken)
	wg.Wait()

	// update access_token dan refresh_token by uname
	accToken, refreshTkn := <-accessToken, <-refreshToken
	err = a.AccountRepo.UpdateToken(ctx, uname, accToken, refreshTkn)
	if err != nil {
		return nil, err
	}

	// return response
	return &model.LoginResponse{
		accToken, refreshTkn,
	}, nil
}

func (a *AccountService) Logout(ctx context.Context, refreshToken string) (string, error) {
	err := a.AccountRepo.DeleteToken(ctx, refreshToken)
	if err != nil {
		return "", errors.New("you are not logged in")
	}

	return "ok", nil
}

func (a *AccountService) Account(ctx context.Context, uname string) (*model.AccountResponse, error) {
	account, err := a.AccountRepo.GetByUname(ctx, uname)
	if err != nil {
		return nil, errors.New("record not found")
	}

	return &model.AccountResponse{
		ID:    account.Id,
		Uname: account.Uname,
		Pass:  account.Pass,
	}, nil
}
