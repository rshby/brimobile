package service

import (
	"brimobile/graph/model"
	"context"
)

type ISavingService interface {
	Insert(ctx context.Context, input model.InsertSavingRequest) (*model.InqAccountSaving, error)
	InqAccountSaving(ctx context.Context, accountNumber string) (*model.InqAccountSaving, error)
	OverbookingLocal(ctx context.Context, overbookingInputParams model.OvbRequest) (*model.OvbResponse, error)
}
