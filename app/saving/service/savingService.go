package service

import (
	"brimobile/app/helper"
	saving "brimobile/app/saving"
	"brimobile/app/saving/repository"
	"brimobile/graph/model"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
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
	wg := &sync.WaitGroup{}
	savingChan, errChan := make(chan saving.Saving, 1), make(chan error, 1)

	defer func() {
		close(savingChan)
		close(errChan)
	}()

	wg.Add(1)
	go s.SavingRepo.GetByAccountNumber(ctx, wg, savingChan, errChan, accountNumber)
	saving, err := <-savingChan, <-errChan
	wg.Wait()

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

func (s *SavingService) OverbookingLocal(ctx context.Context, overbookingInputParams model.OvbRequest) (*model.OvbResponse, error) {
	wg := &sync.WaitGroup{}

	// cek data rekening pengirim dan penerima
	reqChan, benefChan := make(chan saving.Saving, 1), make(chan saving.Saving, 1)
	errReqChan, errBenefChan := make(chan error, 1), make(chan error, 1)

	defer func() {
		close(reqChan)
		close(benefChan)
		close(errReqChan)
		close(errBenefChan)
	}()

	wg.Add(2)
	go s.SavingRepo.GetByAccountNumber(ctx, wg, reqChan, errReqChan, overbookingInputParams.AccountDebit)
	go s.SavingRepo.GetByAccountNumber(ctx, wg, benefChan, errBenefChan, overbookingInputParams.AccountCredit)

	req, benef := <-reqChan, <-benefChan
	errReq, errBenef := <-errReqChan, <-errBenefChan
	wg.Wait()

	fmt.Println(req)
	fmt.Println(benef)

	// cek reuestor dan benef
	if errReq != nil {
		return nil, errors.New("requestor not found")
	}

	if errBenef != nil {
		return nil, errors.New("beneficiary not found")
	}

	/*
		// set saldo setelah transaksi
		cbalReq, _ := strconv.ParseFloat(req.Cbal, 64)
		var reqCbal float64 = cbalReq - overbookingInputParams.AmountTrx

		cbalBenef, _ := strconv.ParseFloat(benef.Cbal, 64)
		var benefCBal float64 = cbalBenef + overbookingInputParams.AmountTrx
	*/

	// update

	return &model.OvbResponse{
		StatusCode: http.StatusOK,
		StatusDesc: "ok",
	}, nil // sengaja dulu
}
