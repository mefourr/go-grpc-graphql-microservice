package order

import (
	"github.com/mefourr/go-grpc-graphql-microservice/catalog/pb"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.CatalogServiceClient
}

func NewClient(url string) *Client {
}
