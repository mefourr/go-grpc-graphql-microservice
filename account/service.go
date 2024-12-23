package account

import (
	"context"
	"github.com/segmentio/ksuid"
)

type Service interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &accountService{repository: repository}
}

func (a accountService) PostAccount(ctx context.Context, name string) (*Account, error) {
	account := &Account{
		ID:   ksuid.New().String(),
		Name: name,
	}
	err := a.repository.PutAccount(ctx, *account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (a accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return a.repository.GetAccountById(ctx, id)
}

func (a accountService) GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	if take > 100 || (skip <= 0 && take <= 0) {
		take = 100
	}
	return a.repository.ListAccounts(ctx, skip, take)
}
