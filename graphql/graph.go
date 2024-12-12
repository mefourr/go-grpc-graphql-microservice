package main

import (
	"github.com/99designs/gqlgen/graphql"
	"github.com/mefourr/go-grpc-graphql-microservice/account"
	"github.com/mefourr/go-grpc-graphql-microservice/catalog"
	"github.com/mefourr/go-grpc-graphql-microservice/order"
)

type Server struct {
	accountClient *account.Client
	catalogClient *catalog.Client
	orderClient   *order.Client
}

func NewGraphQlServer(accountUrl, catalogUrl, orderUrl string) (*Server, error) {
	accountClient := account.NewClient(accountUrl)
	catalogClient := catalog.NewClient(catalogUrl)
	orderClient := order.NewClient(orderUrl)

	return &Server{
		accountClient,
		catalogClient,
		orderClient,
	}, nil
}

func (s *Server) Mutation() MutationResolver {
	return &mutationResolver{
		server: s,
	}
}

func (s *Server) Query() QueryResolver {
	return &queryResolver{
		server: s,
	}
}

func (s *Server) Account() AccountResolver {
	return &accountResolver{
		server: s,
	}
}

func (s *Server) ToExecutableSchema() graphql.ExecutableSchema {
	return NewExecutableSchema(
		Config{Resolvers: s},
	)
}
