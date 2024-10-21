package cli

import (
	"context"

	"github.com/pinbrain/gophkeeper/internal/client/service"
	"github.com/spf13/cobra"
)

type CLI struct {
	service *service.Service

	rootCMD  *cobra.Command
	userCMD  *cobra.Command
	vaultCMD *cobra.Command
}

func NewCLI(ctx context.Context, service *service.Service) *CLI {
	cli := &CLI{
		service: service,
		rootCMD: &cobra.Command{
			Use: "gophkeeper",
		},
		userCMD: &cobra.Command{
			Use:   "user",
			Short: "auth commands",
			Long:  "Create new user, or login with existing one",
		},
		vaultCMD: &cobra.Command{
			Use:   "vault",
			Short: "Work with data",
			Long:  "Work with data - save, get, delete and update",
		},
	}

	cli.userCMD.AddCommand(
		cli.RegisterCmd(ctx),
		cli.LoginCmd(ctx),
	)

	cli.vaultCMD.AddCommand(
		cli.GetDataCmd(ctx),
		cli.GetAllByTypeCmd(ctx),
		cli.AddDataCmd(ctx),
		cli.DeleteDataCmd(ctx),
	)

	cli.rootCMD.AddCommand(cli.userCMD)
	cli.rootCMD.AddCommand(cli.vaultCMD)

	return cli
}

func (c *CLI) Execute() error {
	return c.rootCMD.Execute()
}
