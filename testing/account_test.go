package testing

import (
	"brimobile/app/account"
	"brimobile/app/account/service"
	mck "brimobile/app/mock"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var accRepo = mck.NewAccountRepoMock()
var accService = service.AccountService{accRepo}

// func test insert
func TestInsertAccount(t *testing.T) {
	uname := "reo"
	pass := "123"

	t.Run("test insert account success", func(t *testing.T) {
		accRepo.Mock.On("Insert", context.Background(), &account.Account{
			Uname: uname,
			Pass:  pass,
		}).Return(&account.Account{
			Id:    1,
			Uname: uname,
			Pass:  pass,
		}, nil)

		account, err := accService.CreateAccount(context.Background(), uname, pass)

		assert.NotNil(t, account)
		assert.Nil(t, err)
		assert.Equal(t, 1, account.ID)
		assert.Equal(t, uname, account.Uname)
	})
	t.Run("test insert account failed", func(t *testing.T) {
		accRepo.Mock.On("Insert", context.Background(), &account.Account{
			Uname: uname,
			Pass:  "",
		}).Return(nil, errors.New("password not set"))
		account, err := accService.CreateAccount(context.Background(), uname, "")

		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Equal(t, "password not set", err.Error())
	})
}

// test get by username
func TestGetByUname(t *testing.T) {
	uname := "reo"

	t.Run("test get by uname success", func(t *testing.T) {
		accRepo.Mock.On("GetByUname", mock.Anything, uname).Return(&account.Account{
			Id:    1,
			Uname: uname,
			Pass:  "123",
			AccessToken: sql.NullString{
				Valid:  true,
				String: "qwerty123",
			},
			RefreshToken: sql.NullString{
				Valid:  true,
				String: "1234567890",
			},
		}, nil)

		account, err := accService.Account(context.Background(), uname)

		assert.Nil(t, err)
		assert.NotNil(t, account)
		assert.Equal(t, uname, account.Uname)
	})

	t.Run("test gey by uname not found", func(t *testing.T) {
		uname = ""
		accRepo.Mock.On("GetByUname", mock.Anything, uname).Return(nil, errors.New("uname not set"))

		account, err := accService.Account(context.Background(), uname)

		assert.Nil(t, account)
		assert.NotNil(t, err)
		assert.Equal(t, "record not found", err.Error())
	})
}

// test login
func TestLogin(t *testing.T) {
	uname := "reo"
	pass := "123"

	t.Run("test login success", func(t *testing.T) {
		var accRepo = mck.NewAccountRepoMock()
		var accService = service.AccountService{accRepo}
		accRepo.Mock.On("GetByUname", mock.Anything, uname).Return(&account.Account{
			Id:    1,
			Uname: uname,
			Pass:  pass,
		}, nil)

		accRepo.Mock.On("UpdateToken", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

		result, err := accService.Login(context.Background(), uname, pass, "1", "mb")

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})
	t.Run("test login password not match", func(t *testing.T) {
		var accRepo = mck.NewAccountRepoMock()
		var accService = service.AccountService{accRepo}

		accRepo.Mock.On("GetByUname", mock.Anything, uname).Return(&account.Account{
			Id:    1,
			Uname: uname,
			Pass:  "P@ssw0rd",
		}, nil)

		result, err := accService.Login(context.Background(), uname, "123", "!", "mb")

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "password not match", err.Error())
	})
	t.Run("test login sudah ada access dan refresh Token", func(t *testing.T) {
		accRepo = mck.NewAccountRepoMock()
		accService := service.NewAccountService(accRepo)

		// mock method GetByUname
		accRepo.Mock.On("GetByUname", mock.Anything, uname).Return(&account.Account{
			Id:           2,
			Uname:        uname,
			Pass:         pass,
			AccessToken:  sql.NullString{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVbmFtZSI6InJlbyIsIlJlZ2lzdGVyZWRDbGFpbXMiOnsiaXNzIjoiYnJtb2JpbGUiLCJleHAiOjIwNTc0ODc2OTJ9fQ.98VUJnHMRzy_TK6XzqbgRq4wogTwYcoQ9vUuTATr6Tw", true},
			RefreshToken: sql.NullString{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVbmFtZSI6InJlbyIsIlJlZ2lzdGVyZWRDbGFpbXMiOnsiaXNzIjoiYnJtb2JpbGUiLCJleHAiOjUyOTc0ODc2OTJ9fQ.ESs3y4IdLRVj9EO9o0b3ySOOJtNWUuL_n0nGL5YP5_I", true},
		}, nil).Times(1)

		res, err := accService.Login(context.Background(), uname, pass, "1", "iPhone")

		// validate test
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, "user sudah login", err.Error())
	})
	t.Run("test login errors uname not exist", func(t *testing.T) {
		accountRepo := mck.NewAccountRepoMock()
		accountService := service.NewAccountService(accountRepo)

		// mock method GetByUname
		errorMessage := "record not found"
		accountRepo.Mock.On("GetByUname", mock.Anything, mock.Anything).Return(nil, errors.New(errorMessage)).Times(1)

		result, err := accountService.Login(context.Background(), uname, pass, "1", "iPhone")

		// validate test result
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errorMessage, err.Error())
		accountRepo.Mock.AssertExpectations(t)
	})
	t.Run("test login token sudah ada tapi expired, gagal update token", func(t *testing.T) {
		accountRepo := mck.NewAccountRepoMock()
		accountService := service.NewAccountService(accountRepo)

		// mock method GetByUname
		accountRepo.Mock.On("GetByUname", mock.Anything, uname).Return(&account.Account{
			Id:           1,
			Uname:        uname,
			Pass:         pass,
			AccessToken:  sql.NullString{"", false},
			RefreshToken: sql.NullString{"", false},
		}, nil).Times(1)

		// mock method UpdateToken
		errorMessage := "gagal update token"
		accountRepo.Mock.On("UpdateToken", mock.Anything, uname, mock.Anything, mock.Anything).Return(errors.New(errorMessage)).Times(1)

		res, err := accountService.Login(context.Background(), uname, pass, "1", "iPhone")

		// validate
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errorMessage, err.Error())
		accountRepo.Mock.AssertExpectations(t)
	})
}

// test Logout
func TestLogout(t *testing.T) {
	t.Run("test logout success", func(t *testing.T) {
		accRepo = mck.NewAccountRepoMock()
		accService := service.NewAccountService(accRepo)
		accRepo.Mock.On("DeleteToken", mock.Anything, mock.Anything).Return(nil)
		res, err := accService.Logout(context.Background(), "refreshToken")
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "ok", res)
	})
	t.Run("test login you are not logged in", func(t *testing.T) {
		accountRepo := mck.NewAccountRepoMock()
		accountService := service.NewAccountService(accountRepo)

		// mock method DeleteToken
		errorMessgage := "you are not logged in"
		accountRepo.Mock.On("DeleteToken", mock.Anything, mock.Anything).Return(errors.New(errorMessgage)).Times(1)

		res, err := accountService.Logout(context.Background(), "refresh-token")

		// validate
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, "", res)
		assert.Equal(t, errorMessgage, err.Error())
		accountRepo.Mock.AssertExpectations(t)
	})
}
