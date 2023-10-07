package graph

import "brimobile/app/account/service"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	AccService service.IAccountService
}
