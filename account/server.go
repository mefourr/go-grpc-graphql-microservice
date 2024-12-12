//go:generate protoc --go_out=./pb --go-grpc_out=./pb account.proto
package account

import (
	"context"
	"fmt"
	"github.com/mefourr/go-grpc-graphql-microservice/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type grpcServer struct {
	pb.UnimplementedAccountServiceServer
	service Service
}

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	pb.RegisterAccountServiceServer(server, &grpcServer{service: s})
	reflection.Register(server)
	return server.Serve(lis)
}

func (s *grpcServer) PostAccount(ctx context.Context, req *pb.PostAccountRequest) (*pb.PostAccountResponse, error) {
	account, err := s.service.PostAccount(ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.PostAccountResponse{
		Account: &pb.Account{
			Id:   account.ID,
			Name: account.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccount(ctx context.Context, req *pb.GetAccountRequest) (*pb.GetAccountResponse, error) {
	account, err := s.service.GetAccount(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetAccountResponse{
		Account: &pb.Account{
			Id:   account.ID,
			Name: account.Name,
		},
	}, nil
}

func (s *grpcServer) GetAccounts(ctx context.Context, req *pb.GetAccountsRequest) (*pb.GetAccountsResponse, error) {
	accounts, err := s.service.GetAccounts(ctx, req.Skip, req.Take)
	if err != nil {
		return nil, err
	}

	var pbAccount []*pb.Account
	for _, account := range accounts {
		pbAccount = append(pbAccount, &pb.Account{
			Id:   account.ID,
			Name: account.Name,
		})
	}

	return &pb.GetAccountsResponse{Accounts: pbAccount}, nil
}
