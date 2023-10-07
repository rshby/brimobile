package graph

import (
	acc "brimobile/app/account/service"
	saving "brimobile/app/saving/service"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AccService    acc.IAccountService
	SavingService saving.ISavingService
}
