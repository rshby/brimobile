package service

import (
	"brimobile/app/brinjournal"
	brinjournalRepo "brimobile/app/brinjournal/repository"
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
	"time"
)

type SavingService struct {
	SavingRepo      repository.ISavingRepository
	BrinJournalRepo brinjournalRepo.IBrinjournalRepo
}

// function provider
func NewSavingService(svRepo repository.ISavingRepository, brinjournalRepo brinjournalRepo.IBrinjournalRepo) *SavingService {
	return &SavingService{
		svRepo, brinjournalRepo,
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
		Status:           helper.StatusToString(saving.Status),
		CurrentBalance:   fmt.Sprintf("%.7f", cbal),
		ShortName:        saving.ShortName,
	}, nil
}

// method overbooking
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

	// cek reuestor dan benef
	if errReq != nil {
		return nil, errors.New("requestor not found")
	}
	if errBenef != nil {
		return nil, errors.New("beneficiary not found")
	}

	// cek saldo requestor apakah cukup
	reqcbal, _ := strconv.ParseFloat(req.Cbal, 64)
	reqhold, _ := strconv.ParseFloat(req.Hold, 64)
	if (reqcbal - reqhold) <= overbookingInputParams.AmountTrx {
		return nil, errors.New("requestor not enough balance")
	}

	// set saldo setelah transaksi
	cbalReq, _ := strconv.ParseFloat(req.Cbal, 64)
	var reqCbal float64 = cbalReq - overbookingInputParams.AmountTrx

	cbalBenef, _ := strconv.ParseFloat(benef.Cbal, 64)
	var benefCbal float64 = cbalBenef + overbookingInputParams.AmountTrx

	// update saldo(cbal)
	errChanUpdateReq, errChanUpdateBenef := make(chan error, 1), make(chan error, 1)
	defer func() {
		close(errChanUpdateReq)
		close(errChanUpdateBenef)
	}()

	wg.Add(2)
	go s.SavingRepo.UpdateCbal(ctx, wg, errChanUpdateReq, req.AccountNumber, reqCbal)
	go s.SavingRepo.UpdateCbal(ctx, wg, errChanUpdateBenef, benef.AccountNumber, benefCbal)

	errUpdateReq, errUpdateBenef := <-errChanUpdateReq, <-errChanUpdateBenef
	wg.Wait()

	if errUpdateReq != nil {
		return nil, errUpdateReq
	}
	if errUpdateBenef != nil {
		return nil, errUpdateBenef
	}

	// get brinjournal
	wgb := &sync.WaitGroup{}
	errChanBrinjournal, brinjournalChan := make(chan error, 1), make(chan brinjournal.BrinJournal, 1)
	defer func() {
		close(errChanBrinjournal)
		close(brinjournalChan)
	}()

	wgb.Add(1)
	go s.BrinJournalRepo.GetByBranch(ctx, wgb, brinjournalChan, errChanBrinjournal, "09999")
	brinJournalData, errBrinJournal := <-brinjournalChan, <-errChanBrinjournal
	wgb.Wait()
	if errBrinJournal != nil {
		return nil, errBrinJournal
	}

	// set trfReff
	tellerId := "0999999"
	dateToday := time.Now().Local().Format("020106")
	trRefn := fmt.Sprintf("%v%v%v0", helper.PadLeft(tellerId, 7), helper.PadLeft(dateToday, 6), helper.PadLeft(strconv.Itoa(brinJournalData.JournalSeq), 7))

	// update brinjournal
	if err := s.BrinJournalRepo.UpdateOne(ctx, "09999"); err != nil {
		return nil, err
	}

	// success transaction
	return &model.OvbResponse{
		StatusCode:     http.StatusOK,
		StatusDesc:     "ok",
		AccountDebit:   req.AccountNumber,
		NameDebit:      req.ShortName,
		StatusDebit:    helper.StatusToString(req.Status),
		AccountCredit:  benef.AccountNumber,
		NameCredit:     benef.ShortName,
		StatusCredit:   helper.StatusToString(benef.Status),
		AmountTrx:      fmt.Sprintf("%.7f", overbookingInputParams.AmountTrx),
		Remark:         overbookingInputParams.Remark,
		DateTrx:        time.Now().Local().Format("2006-01-02 15:04:05"),
		Trrefn:         trRefn,
		CurrencyDebit:  req.Currency,
		CurrencyCredit: req.Currency,
	}, nil // sengaja dulu
}
