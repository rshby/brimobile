package service

import (
	saving "brimobile/app/saving"
	"brimobile/app/saving/repository"
	"brimobile/graph/model"
	"context"
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

func (s *SavingService) Insert(ctx context.Context, input model.InsertSavingRequest) (*model.InqAccountSaving, error) {
	// create entity
	saving := saving.Saving{
		AccountNumber: input.AccountNumber,
		ShortName:     input.ShortName,
		Cbal:          input.Cbal,
	}

	// call repository insert
	res, err := s.SavingRepo.Insert(ctx, &saving)
	if err != nil {
		return nil, err
	}

	response := model.InqAccountSaving{
		AccountNumber: res.AccountNumber,
	}

	// success
	return &response, nil
}
