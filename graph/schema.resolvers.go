package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.39

import (
	"brimobile/app/account/repository"
	"brimobile/app/account/service"
	"brimobile/db/connection"
	"brimobile/graph/model"
	"context"
)

// CreateAccount is the resolver for the createAccount field.
func (r *mutationResolver) CreateAccount(ctx context.Context, uname string, pass string) (*model.CreateAccountResponse, error) {
	return r.AccService.CreateAccount(ctx, uname, pass)
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, uname string, pass string, idNum string, deviceID string) (*model.LoginResponse, error) {
	return r.AccService.Login(ctx, uname, pass, idNum, deviceID)
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context, refreshToken string) (string, error) {
	return r.AccService.Logout(ctx, refreshToken)
}

// Account is the resolver for the account field.
func (r *queryResolver) Account(ctx context.Context, uname string) (*model.AccountResponse, error) {
	return r.AccService.Account(ctx, uname)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{
		Resolver:   r,
		AccService: service.NewAccountService(repository.NewAccountRepository(connection.DB)),
	}
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{
		Resolver:   r,
		AccService: service.NewAccountService(repository.NewAccountRepository(connection.DB)),
	}
}

type mutationResolver struct {
	Resolver   *Resolver
	AccService service.IAccountService
}
type queryResolver struct {
	Resolver   *Resolver
	AccService service.IAccountService
}