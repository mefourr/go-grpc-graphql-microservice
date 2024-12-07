package main

import "context"

type mutationResolver struct {
	server *Server
}

func (r *mutationResolver)  CreateAccount(ctx context.Context, account AccountInput) (*Account, error) {
	return &Account{}, nil
}

func (r *mutationResolver) CreateProduct(ctx context.Context, product ProductIntput) (*Product, error) {
	return &Product{}, nil
}

func (r *mutationResolver) CreateOrder(ctx context.Context, order OrderIntput) (*Order, error) {
	return &Order{}, nil
}
