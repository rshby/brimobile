package testing

import (
	mck "brimobile/app/mock"
	"brimobile/app/saving"
	"brimobile/app/saving/service"
	"brimobile/graph/model"
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// test insert saving
func TestInsertSaving(t *testing.T) {
	entity := saving.Saving{
		AccountNumber: "001",
		ShortName:     "reo sahobby",
		Cbal:          "100.00",
	}

	t.Run("test insert success", func(t *testing.T) {
		var savingRepo = mck.NewSavingRepositoryMock()
		var brinjournalRepo = mck.NewBrinJournalMock()
		var savingService = service.NewSavingService(savingRepo, brinjournalRepo)

		savingRepo.Mock.On("Insert", mock.Anything, &entity).Return(&entity, nil)

		res, err := savingService.Insert(context.Background(), model.InsertSavingRequest{
			AccountNumber: entity.AccountNumber,
			ShortName:     entity.ShortName,
			Cbal:          entity.Cbal,
		})

		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, entity.AccountNumber, res.AccountNumber)
	})

	// test insert saving error
	t.Run("test insert saving failed", func(t *testing.T) {
		savingRepo := mck.NewSavingRepositoryMock()
		brinjournalRepo := mck.NewBrinJournalMock()
		savingService := service.NewSavingService(savingRepo, brinjournalRepo)

		// mock set method insert in SavingRepository -> return error
		errorMessage := "gagal insert"
		savingRepo.Mock.On("Insert", mock.Anything, mock.Anything).Return(nil, errors.New(errorMessage)).Times(1)

		// run method Insert in SavingService
		res, err := savingService.Insert(context.Background(), model.InsertSavingRequest{
			AccountNumber: "0452001",
			ShortName:     "Reo Sahobby",
			Cbal:          "100.0",
		})

		// validate test
		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, errorMessage, err.Error())
		savingRepo.Mock.AssertExpectations(t)
	})
}

// test get by account_number
func TestInqSaving(t *testing.T) {
	accNum := "001"
	t.Run("test inq saving success", func(t *testing.T) {
		var savingRepo = mck.NewSavingRepositoryMock()
		var brinjournalRepo = mck.NewBrinJournalMock()
		var savingService = service.NewSavingService(savingRepo, brinjournalRepo)

		savingRepo.Mock.On("GetByAccountNumber", mock.Anything, mock.Anything, mock.Anything, mock.Anything, accNum).Return(saving.Saving{
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
		var savingRepo = mck.NewSavingRepositoryMock()
		var brinjournalRepo = mck.NewBrinJournalMock()
		var savingService = service.NewSavingService(savingRepo, brinjournalRepo)

		accNum = "002"
		savingRepo.Mock.On("GetByAccountNumber", mock.Anything, mock.Anything, mock.Anything, mock.Anything, accNum).Return(saving.Saving{}, errors.New("record not found"))

		resInq, err := savingService.InqAccountSaving(context.Background(), accNum)

		assert.Nil(t, resInq)
		assert.NotNil(t, err)
		assert.Equal(t, "record not found", err.Error())
	})
}
