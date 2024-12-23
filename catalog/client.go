package catalog

import (
	"context"
	"github.com/mefourr/go-grpc-graphql-microservice/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) *Client {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil
	}
	serviceClient := pb.NewCatalogServiceClient(client)
	return &Client{conn: client, service: serviceClient}
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}

func (c *Client) PostProduct(ctx context.Context, name, description string, price float64) (*Product, error) {
	res, err := c.service.PostProduct(
		ctx,
		&pb.PostProductRequest{
			Name:        name,
			Description: description,
			Price:       price,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	res, err := c.service.GetProduct(
		ctx,
		&pb.GetProductRequest{
			Id: id,
		},
	)
	if err != nil {
		return nil, err
	}
	return &Product{
		ID:          res.Product.Id,
		Name:        res.Product.Name,
		Description: res.Product.Description,
		Price:       res.Product.Price,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context, ids []string, query string, skip, take uint64) ([]Product, error) {
	res, err := c.service.GetProducts(
		ctx,
		&pb.GetProductsRequest{
			Ids:   ids,
			Query: query,
			Skip:  skip,
			Take:  take,
		},
	)
	if err != nil {
		return nil, err
	}

	var products []Product
	for _, prod := range res.Products {
		products = append(products, Product{
			ID:          prod.Id,
			Name:        prod.Name,
			Description: prod.Description,
			Price:       prod.Price,
		})
	}
	return products, nil
}
