package service

import (
	"brimobile/app/account"
	"brimobile/app/account/repository"
	"brimobile/app/helper"
	"brimobile/graph/model"
	"context"
	"database/sql"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
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
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountService Login")
	defer span.Finish()

	// cek apakah uname ada di database
	account, err := a.AccountRepo.GetByUname(ctxTracing, uname)
	if err != nil {
		return nil, errors.New("record not found")
	}

	// cek apakah password benar
	if !helper.CheckPassword(ctxTracing, account.Pass, pass) {
		return nil, errors.New("password not match")
	}

	// cek apakah sudah login
	if account.AccessToken.Valid {
		token := account.RefreshToken.String
		claims, _ := helper.GetClaims(ctxTracing, token)

		// jika belum expired tapi user sudah login -> error user sudah login
		if !time.Now().Local().After(time.UnixMicro(claims.RegisteredClaims.ExpiresAt.Unix() * 1000000)) {
			return nil, errors.New("user sudah login")
		}
	}

	// lolos -> generate access_token dan refresh_token
	wg := &sync.WaitGroup{}
	accessToken := make(chan string, 1)
	refreshToken := make(chan string, 1)

	wg.Add(2)
	go helper.GenerateToken(ctxTracing, uname, time.Duration(5*time.Minute), wg, accessToken)
	go helper.GenerateToken(ctxTracing, uname, time.Duration(10*time.Minute), wg, refreshToken)
	wg.Wait()

	// update access_token dan refresh_token by uname
	accToken, refreshTkn := <-accessToken, <-refreshToken
	err = a.AccountRepo.UpdateToken(ctxTracing, uname, accToken, refreshTkn)
	if err != nil {
		return nil, err
	}

	// return response
	response := model.LoginResponse{
		time.Now().Format("2006-01-02 15:04:05"), accToken, refreshTkn,
	}
	span.LogFields(
		log.String("request-uname", uname),
		log.String("request-pass", pass),
		log.String("request-idNum", idNum),
		log.String("request-deviceId", deviceID),
		log.Object("response-login", response))
	return &response, nil
}

func (a *AccountService) Logout(ctx context.Context, refreshToken string) (string, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountService Logout")
	defer span.Finish()

	err := a.AccountRepo.DeleteToken(ctxTracing, refreshToken)
	if err != nil {
		return "", errors.New("you are not logged in")
	}

	span.LogFields(
		log.String("request-refreshToken", refreshToken),
		log.String("response-message", "ok"),
		log.Object("response-error", err),
	)
	return "ok", nil
}

func (a *AccountService) Account(ctx context.Context, uname string) (*model.AccountResponse, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "AccountService Account")
	defer span.Finish()

	account, err := a.AccountRepo.GetByUname(ctxTracing, uname)
	if err != nil {
		return nil, errors.New("record not found")
	}

	response := model.AccountResponse{
		ID:    account.Id,
		Uname: account.Uname,
		Pass:  account.Pass,
	}

	span.LogFields(
		log.String("request-uname", uname),
		log.Object("response-account", response))
	return &response, nil
}
