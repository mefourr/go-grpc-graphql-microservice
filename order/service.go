package order

import (
	"context"
	"github.com/segmentio/ksuid"
	"time"
)

type Service interface {
	PostOrder(ctx context.Context, accountId string, product []OrderedProduct) (*Order, error)
	GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error)
}

type Order struct {
	ID         string
	CreatedAt  time.Time
	AccountId  string
	TotalPrice float64
	Products   []OrderedProduct
}

type OrderedProduct struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Quantity    int
}

type ServiceOrder struct {
	repository Repository
}

func (s ServiceOrder) PostOrder(ctx context.Context, accountId string, product []OrderedProduct) (*Order, error) {
	order := &Order{
		ID:         ksuid.New().String(),
		CreatedAt:  time.Now().UTC(),
		AccountId:  accountId,
		TotalPrice: 0.0,
		Products:   product,
	}

	for _, p := range product {
		order.TotalPrice += p.Price * float64(p.Quantity)
	}

	err := s.repository.PutOrder(ctx, *order)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s ServiceOrder) GetOrdersForAccount(ctx context.Context, accountId string) ([]Order, error) {
	return s.repository.GetOrdersForAccount(ctx, accountId)
}

func NewService(repository Repository) *ServiceOrder {
	return &ServiceOrder{repository: repository}
}
