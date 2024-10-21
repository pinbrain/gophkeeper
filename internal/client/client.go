package client

import (
	"context"

	"github.com/pinbrain/gophkeeper/internal/client/cli"
	"github.com/pinbrain/gophkeeper/internal/client/config"
	"github.com/pinbrain/gophkeeper/internal/client/grpc"
	"github.com/pinbrain/gophkeeper/internal/client/service"
)

type Client struct {
	cli *cli.CLI
}

func NewClient(cfg *config.ClientConfig) (*Client, error) {
	ctx := context.Background()
	grpcClient, err := grpc.NewGRPCConnection(cfg)
	if err != nil {
		return nil, err
	}
	service := service.NewService(grpcClient)
	cliClient := cli.NewCLI(ctx, service)
	return &Client{
		cli: cliClient,
	}, nil
}

func (c *Client) Execute() error {
	return c.cli.Execute()
}
