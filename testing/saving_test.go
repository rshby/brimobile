package testing

import (
	"brimobile/app/mock"
	"brimobile/app/saving"
	"brimobile/app/saving/service"
	"brimobile/graph/model"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var savingRepo = mock.NewSavingRepositoryMock()
var savingService = service.NewSavingService(savingRepo)

// test insert saving
func TestInsertSaving(t *testing.T) {
	entity := saving.Saving{
		AccountNumber: "001",
		ShortName:     "reo sahobby",
		Cbal:          "100.00",
	}

	t.Run("test insert success", func(t *testing.T) {
		savingRepo.Mock.On("Insert", context.Background(), &entity).Return(&saving.Saving{
			AccountNumber: entity.AccountNumber,
			ShortName:     entity.ShortName,
			Cbal:          entity.Cbal,
		}, nil)

		res, err := savingService.Insert(context.Background(), model.InsertSavingRequest{
			AccountNumber: entity.AccountNumber,
			ShortName:     entity.ShortName,
			Cbal:          entity.Cbal,
		})

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, entity.AccountNumber, res.AccountNumber)
	})
}

// test get by account_number
func TestInqSaving(t *testing.T) {
	accNum := "001"
	t.Run("test inq success", func(t *testing.T) {
		savingRepo.Mock.On("GetByAccountNumber", context.Background(), accNum).Return(&saving.Saving{
			AccountNumber: accNum,
			AccountType:   "S",
			ShortName:     "reo sahobby",
			Cbal:          "100.00",
			Hold:          "0.00",
		}, nil)

		resultInq, err := savingService.InqAccountSaving(context.Background(), accNum)

		assert.Nil(t, err)
		assert.NotNil(t, resultInq)
		fmt.Println(resultInq.AvailableBalance)
	})

	t.Run("test inq saving not found", func(t *testing.T) {
		accNum = "002"

		savingRepo.Mock.On("GetByAccountNumber", context.Background(), accNum).Return(nil, errors.New("record not found"))

		resInq, err := savingService.InqAccountSaving(context.Background(), accNum)

		assert.Nil(t, resInq)
		assert.NotNil(t, err)
		assert.Equal(t, "record not found", err.Error())
	})
}
