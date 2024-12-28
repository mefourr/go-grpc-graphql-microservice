//go:generate protoc --go_out=./pb --go-grpc_out=./pb order.proto
package order

import (
	"context"
	"fmt"
	"github.com/mefourr/go-grpc-graphql-microservice/account"
	"github.com/mefourr/go-grpc-graphql-microservice/catalog"
	"github.com/mefourr/go-grpc-graphql-microservice/order/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type grpcServer struct {
	pb.UnimplementedOrderServiceServer
	service       Service
	accountClient *account.Client
	catalogClient *catalog.Client
}

func ListenGRPC(s Service, accountURL, catalogURL string, port int) error {
	accountClient, err := account.NewClient(accountURL)
	if err != nil {
		log.Printf("failed to create account client: %w", err)
		return err
	}

	catalogClient, err := catalog.NewClient(catalogURL)
	if err != nil {
		accountClient.Close()
		log.Printf("failed to create catalog client: %w", err)
		return err
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		accountClient.Close()
		catalogClient.Close()
		log.Printf("failed to listen on port %d: %w", port, err)
		return err
	}

	server := grpc.NewServer()
	pb.RegisterOrderServiceServer(server, &grpcServer{
		service:       s,
		accountClient: accountClient,
		catalogClient: catalogClient,
	})
	reflection.Register(server)

	log.Printf("gRPC server listening on port %d", port)
	return server.Serve(lis)
}

func (s *grpcServer) PostOrder(ctx context.Context, req *pb.PostOrderRequest) (*pb.PostOrderResponse, error) {
	_, err := s.accountClient.GetAccount(ctx, req.AccountId)
	if err != nil {
		return nil, err
	}

	var productIds []string
	orderedProducts, err := s.catalogClient.GetProducts(ctx, productIds, "", 0, 0)
	if err != nil {
		return nil, err
	}

	var products []OrderedProduct
	for _, orderedProduct := range orderedProducts {
		product := OrderedProduct{
			ID:          orderedProduct.ID,
			Name:        orderedProduct.Name,
			Description: orderedProduct.Description,
			Price:       orderedProduct.Price,
		}
		for _, rp := range req.Products {
			if rp.ProductId == orderedProduct.ID {
				product.Quantity = int(rp.Quantity)
				break
			}
		}
		if product.Quantity != 0 {
			products = append(products, product)
		}
	}

	order, err := s.service.PostOrder(ctx, req.AccountId, products)
	if err != nil {
		return nil, err
	}

	orderProto := pb.Order{
		Id:         order.ID,
		AccountId:  order.AccountId,
		TotalPrice: order.TotalPrice,
		Products:   []*pb.Order_OrderProduct{},
	}
	orderProto.CreatedAt, _ = order.CreatedAt.MarshalBinary()

	for _, product := range order.Products {
		orderProto.Products = append(orderProto.Products, &pb.Order_OrderProduct{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Quantity:    uint32(product.Quantity),
		})
	}

	return &pb.PostOrderResponse{Order: &orderProto}, nil
}

func (s *grpcServer) GetOrdersForAccount(ctx context.Context, req *pb.GetOrdersForAccountRequest) (*pb.GetOrdersForAccountResponse, error) {
	// Get orders for account
	accountOrders, err := s.service.GetOrdersForAccount(ctx, req.AccountId)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Get all ordered products
	productIDMap := map[string]bool{}
	for _, o := range accountOrders {
		for _, p := range o.Products {
			productIDMap[p.ID] = true
		}
	}
	var productIDs []string
	for id := range productIDMap {
		productIDs = append(productIDs, id)
	}
	products, err := s.catalogClient.GetProducts(ctx, productIDs, "", 0, 0)
	if err != nil {
		log.Println("Error getting account products: ", err)
		return nil, err
	}

	// Construct orders
	var orders []*pb.Order
	for _, o := range accountOrders {
		// Encode order
		op := &pb.Order{
			AccountId:  o.AccountId,
			Id:         o.ID,
			TotalPrice: o.TotalPrice,
			Products:   []*pb.Order_OrderProduct{},
		}
		op.CreatedAt, _ = o.CreatedAt.MarshalBinary()

		// Decorate orders with products
		for _, product := range o.Products {
			// Populate product fields
			for _, p := range products {
				if p.ID == product.ID {
					product.Name = p.Name
					product.Description = p.Description
					product.Price = p.Price
					break
				}
			}

			op.Products = append(op.Products, &pb.Order_OrderProduct{
				Id:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Quantity:    uint32(product.Quantity),
			})
		}

		orders = append(orders, op)
	}
	return &pb.GetOrdersForAccountResponse{Orders: orders}, nil
}
