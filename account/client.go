package account

import (
	"context"
	"github.com/mefourr/go-grpc-graphql-microservice/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.AccountServiceClient
}

func NewClient(url string) (*Client, error) {
	client, err := grpc.NewClient(url, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")))
	if err != nil {
		return nil, err
	}
	serviceClient := pb.NewAccountServiceClient(client)
	return &Client{conn: client, service: serviceClient}, nil
}

func (c *Client) Close() {
	err := c.conn.Close()
	if err != nil {
		return
	}
}

func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	account, err := c.service.PostAccount(
		ctx,
		&pb.PostAccountRequest{Name: name},
	)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   account.Account.Id,
		Name: account.Account.Name,
	}, nil
}

func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	account, err := c.service.GetAccount(
		ctx,
		&pb.GetAccountRequest{Id: id},
	)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   account.Account.Id,
		Name: account.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip, take uint64) ([]Account, error) {
	accountsReq, err := c.service.GetAccounts(
		ctx,
		&pb.GetAccountsRequest{Skip: skip, Take: take},
	)
	if err != nil {
		return nil, err
	}
	var accounts []Account
	for _, account := range accountsReq.Accounts {
		accounts = append(accounts, Account{
			ID:   account.Id,
			Name: account.Name,
		})
	}
	return accounts, nil
}
