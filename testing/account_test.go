package testing

import (
	"brimobile/app/account"
	"brimobile/app/account/service"
	mck "brimobile/app/mock"
	"context"
	"database/sql"
	"errors"
	"github.com/stretchr/testify/assert"
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
		accRepo.Mock.On("GetByUname", context.Background(), uname).Return(&account.Account{
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
		accRepo.Mock.On("GetByUname", context.Background(), uname).Return(nil, errors.New("uname not set"))

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
		accRepo.Mock.On("GetByUname", context.Background(), uname).Return(&account.Account{
			Id:    1,
			Uname: uname,
			Pass:  pass,
		}, nil)

		result, err := accService.Login(context.Background(), uname, pass, "1", "mb")

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("test login password not match", func(t *testing.T) {
		accRepo.Mock.On("GetByUname", context.Background(), uname).Return(&account.Account{
			Id:    1,
			Uname: uname,
			Pass:  "P@ssw0rd",
		}, nil)

		result, err := accService.Login(context.Background(), uname, "123", "!", "mb")

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "password not match", err.Error())
	})

}

// test Logour
func TestLogout(t *testing.T) {
	t.Run("test logout success", func(t *testing.T) {
		refreshToken := "qwerty123"
		accRepo.Mock.On("DeleteToken", context.Background(), refreshToken).Return(nil)

		res, err := accService.Logout(context.Background(), refreshToken)

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "ok", res)
	})
}
