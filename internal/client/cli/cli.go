package cli

import (
	"context"

	"github.com/pinbrain/gophkeeper/internal/model"
	"github.com/spf13/cobra"
)

type Service interface {
	UserService
	VaultService
}

type UserService interface {
	Register(ctx context.Context, login, password string) (token string, err error)
	Login(ctx context.Context, login, password string) (token string, err error)
}

type VaultService interface {
	AddPassword(ctx context.Context, data string, meta model.PasswordMeta) error
	AddText(ctx context.Context, data string, meta model.TextMeta) error
	AddBankCard(ctx context.Context, data model.BankCardData, meta model.BankCardMeta) error
	AddFile(ctx context.Context, file string, comment string) error
	GetData(ctx context.Context, id string) (model.DataType, any, error)
	GetAllByType(ctx context.Context, dataType model.DataType) ([]model.ItemInfo, error)
	DeleteData(ctx context.Context, id string) error
}

type CLI struct {
	service Service

	rootCMD  *cobra.Command
	userCMD  *cobra.Command
	vaultCMD *cobra.Command
}

func NewCLI(ctx context.Context, service Service) *CLI {
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
