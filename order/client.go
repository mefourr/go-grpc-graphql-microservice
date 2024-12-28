package order

import (
	"context"
	"github.com/mefourr/go-grpc-graphql-microservice/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"time"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.OrderServiceClient
}

func NewClient(url string) (*Client, error) {
	conn, err := grpc.NewClient(url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return &Client{conn, c}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}

func (c *Client) PostOrder(
	ctx context.Context,
	accountID string,
	products []OrderedProduct,
) (*Order, error) {
	var protoProducts []*pb.PostOrderRequest_OrderProduct
	for _, p := range products {
		protoProducts = append(protoProducts, &pb.PostOrderRequest_OrderProduct{
			ProductId: p.ID,
			Quantity:  uint32(p.Quantity),
		})
	}
	r, err := c.service.PostOrder(
		ctx,
		&pb.PostOrderRequest{
			AccountId: accountID,
			Products:  protoProducts,
		},
	)
	if err != nil {
		return nil, err
	}

	// Create response order
	newOrder := r.Order
	newOrderCreatedAt := time.Time{}
	err = newOrderCreatedAt.UnmarshalBinary(newOrder.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &Order{
		ID:         newOrder.Id,
		CreatedAt:  newOrderCreatedAt,
		TotalPrice: newOrder.TotalPrice,
		AccountId:  newOrder.AccountId,
		Products:   products,
	}, nil
}

func (c *Client) GetOrdersForAccount(ctx context.Context, accountID string) ([]Order, error) {
	r, err := c.service.GetOrdersForAccount(ctx, &pb.GetOrdersForAccountRequest{
		AccountId: accountID,
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Create response orders
	var orders []Order
	for _, orderProto := range r.Orders {
		newOrder := Order{
			ID:         orderProto.Id,
			TotalPrice: orderProto.TotalPrice,
			AccountId:  orderProto.AccountId,
		}
		newOrder.CreatedAt = time.Time{}
		err := newOrder.CreatedAt.UnmarshalBinary(orderProto.CreatedAt)
		if err != nil {
			return nil, err
		}

		var products []OrderedProduct
		for _, p := range orderProto.Products {
			products = append(products, OrderedProduct{
				ID:          p.Id,
				Quantity:    int(p.Quantity),
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			})
		}
		newOrder.Products = products

		orders = append(orders, newOrder)
	}
	return orders, nil
}
