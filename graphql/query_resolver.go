package main

import "context"

type queryResolver struct {
	server *Server
}

func (r* queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string) ([]*Account, error) {
	return []*Account{}, nil
}

func (r* queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string) ([]*Product, error) {
	return []*Product{}, nil
}
