package catalog

import (
	"context"
	"errors"
	"github.com/segmentio/ksuid"
)

type Service interface {
	PutProduct(ctx context.Context, name, description string, price float64) (*Product, error)
	GetProductById(ctx context.Context, id string) (*Product, error)
	ListProducts(ctx context.Context, skip, take uint64) ([]Product, error)
	ListProductsWithIds(ctx context.Context, ids []string) ([]Product, error)
	SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error)
}

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type ServiceCatalog struct {
	repository Repository
}

func NewService(repository Repository) *ServiceCatalog {
	return &ServiceCatalog{repository: repository}
}

func (p2 ServiceCatalog) PutProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	product := Product{
		ID:          ksuid.New().String(),
		Name:        name,
		Description: description,
		Price:       price,
	}
	err := p2.repository.PutProduct(ctx, product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (p2 ServiceCatalog) GetProductById(ctx context.Context, id string) (*Product, error) {
	return p2.repository.GetProductById(ctx, id)
}

func (p2 ServiceCatalog) ListProducts(ctx context.Context, skip, take uint64) ([]Product, error) {
	if take > 100 || (skip <= 0 && take <= 0) {
		take = 100
	}
	return p2.repository.ListProducts(ctx, skip, take)
}

func (p2 ServiceCatalog) ListProductsWithIds(ctx context.Context, ids []string) ([]Product, error) {
	if ids == nil || len(ids) == 0 {
		return []Product{}, errors.New("no product id provided")
	}
	return p2.repository.ListProductsWithIds(ctx, ids)
}

func (p2 ServiceCatalog) SearchProducts(ctx context.Context, query string, skip, take uint64) ([]Product, error) {
	if take > 100 || (skip <= 0 && take <= 0) {
		take = 100
	}
	return p2.repository.SearchProducts(ctx, query, skip, take)
}
