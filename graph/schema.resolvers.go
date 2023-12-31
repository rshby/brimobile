package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.39

import (
	"brimobile/graph/model"
	"context"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// CreateAccount is the resolver for the createAccount field.
func (r *mutationResolver) CreateAccount(ctx context.Context, uname string, pass string) (*model.CreateAccountResponse, error) {
	return r.AccService.CreateAccount(ctx, uname, pass)
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, uname string, pass string, idNum string, deviceID string) (*model.LoginResponse, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Resolver Login")
	defer span.Finish()

	span.LogFields(
		log.String("uname", uname),
		log.String("idNum", idNum),
		log.String("deviceId", deviceID),
	)

	return r.AccService.Login(ctxTracing, uname, pass, idNum, deviceID)
}

// Logout is the resolver for the logout field.
func (r *mutationResolver) Logout(ctx context.Context, refreshToken string) (string, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "Resolver Logout")
	defer span.Finish()
	span.LogFields(
		log.String("refreshToken", refreshToken),
	)

	return r.AccService.Logout(ctxTracing, refreshToken)
}

// InsertSaving is the resolver for the insertSaving field.
func (r *mutationResolver) InsertSaving(ctx context.Context, input model.InsertSavingRequest) (*model.InqAccountSaving, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "resolver InsertSaving")
	defer span.Finish()

	span.SetTag("request", input)

	return r.SavingService.Insert(ctxTracing, input)
}

// OverbookingLocal is the resolver for the overbookingLocal field.
func (r *mutationResolver) OverbookingLocal(ctx context.Context, overbookingInputParams model.OvbRequest) (*model.OvbResponse, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "resolver OverbookingLocal")
	defer span.Finish()

	span.SetTag("request", overbookingInputParams)

	return r.SavingService.OverbookingLocal(ctxTracing, overbookingInputParams)
}

// Account is the resolver for the account field.
func (r *queryResolver) Account(ctx context.Context, uname string) (*model.AccountResponse, error) {
	return r.AccService.Account(ctx, uname)
}

// InqAccountSaving is the resolver for the inqAccountSaving field.
func (r *queryResolver) InqAccountSaving(ctx context.Context, accountNumber string) (*model.InqAccountSaving, error) {
	span, ctxTracing := opentracing.StartSpanFromContext(ctx, "resolver InqAccountSaving")
	span.SetTag("account_number", accountNumber)
	span.LogFields(
		log.String("acc", accountNumber),
		log.String("time", time.Now().Format("15:04:05")),
	)
	defer span.Finish()

	return r.SavingService.InqAccountSaving(ctxTracing, accountNumber)
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
