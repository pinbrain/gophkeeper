// Package client инициализирует запуск клиента для работы с хранилищем.
package client

import (
	"context"

	"github.com/pinbrain/gophkeeper/internal/client/cli"
	"github.com/pinbrain/gophkeeper/internal/client/config"
	"github.com/pinbrain/gophkeeper/internal/client/grpc"
	"github.com/pinbrain/gophkeeper/internal/client/service"
)

// Client описывает структуру клиента.
type Client struct {
	cli *cli.CLI
}

// NewClient создает и возвращает новый клиент.
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

// Execute запускает обработку команды.
func (c *Client) Execute() error {
	return c.cli.Execute()
}
