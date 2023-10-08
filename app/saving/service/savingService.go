package service

import (
	"brimobile/app/helper"
	saving "brimobile/app/saving"
	"brimobile/app/saving/repository"
	"brimobile/graph/model"
	"context"
	"fmt"
	"strconv"
)

type SavingService struct {
	SavingRepo repository.ISavingRepository
}

// function provider
func NewSavingService(svRepo repository.ISavingRepository) *SavingService {
	return &SavingService{
		svRepo,
	}
}

// method insert
func (s *SavingService) Insert(ctx context.Context, input model.InsertSavingRequest) (*model.InqAccountSaving, error) {
	currentBal, _ := strconv.ParseFloat(input.Cbal, 64)

	// create entity
	saving := saving.Saving{
		AccountNumber: input.AccountNumber,
		ShortName:     input.ShortName,
		Cbal:          fmt.Sprintf("%.2f", currentBal),
	}

	// call repository insert
	res, err := s.SavingRepo.Insert(ctx, &saving)
	if err != nil {
		return nil, err
	}

	cbal, _ := strconv.ParseFloat(res.Cbal, 64)
	hold, _ := strconv.ParseFloat(res.Hold, 64)

	// success
	return &model.InqAccountSaving{
		AccountNumber:    res.AccountNumber,
		AvailableBalance: fmt.Sprintf("%.7f", cbal-hold),
		AccountType:      res.AccountType,
		BranchCode:       res.BranchCode,
		Currency:         res.Currency,
		OpeningDate:      res.OpeningDate,
		ProductGroup:     res.ProductGroup,
		ProductName:      res.ProductName,
		Status:           helper.StatusToString(res.Status),
		CurrentBalance:   fmt.Sprintf("%.7f", cbal),
		ShortName:        res.ShortName,
	}, nil
}

// method inquiry saving
func (s *SavingService) InqAccountSaving(ctx context.Context, accountNumber string) (*model.InqAccountSaving, error) {
	saving, err := s.SavingRepo.GetByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, err
	}

	cbal, _ := strconv.ParseFloat(saving.Cbal, 64)
	hold, _ := strconv.ParseFloat(saving.Hold, 64)

	// success
	return &model.InqAccountSaving{
		AccountNumber:    saving.AccountNumber,
		AvailableBalance: fmt.Sprintf("%.7f", (cbal - hold)),
		AccountType:      saving.AccountType,
		BranchCode:       saving.BranchCode,
		Currency:         saving.Currency,
		OpeningDate:      saving.OpeningDate,
		ProductGroup:     saving.ProductGroup,
		ProductName:      saving.ProductName,
		Status:           saving.Status,
		CurrentBalance:   fmt.Sprintf("%.7f", cbal),
		ShortName:        saving.ShortName,
	}, nil
}
