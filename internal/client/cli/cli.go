package cli

import (
	"context"
	"fmt"

	"github.com/pinbrain/gophkeeper/internal/client/config"
	"github.com/pinbrain/gophkeeper/internal/model"
	"github.com/spf13/cobra"
)

// Service описывает структуру сервиса для работы с данными и сервером.
type Service interface {
	UserService
	VaultService
}

// UserService описывает методы для работы с регистрацией и аутентификацией.
type UserService interface {
	Register(ctx context.Context, login, password string) (token string, err error)
	Login(ctx context.Context, login, password string) (token string, err error)
}

// VaultService описывает методы для работы с данными.
type VaultService interface {
	AddPassword(ctx context.Context, data string, meta model.PasswordMeta) error
	AddText(ctx context.Context, data string, meta model.TextMeta) error
	AddBankCard(ctx context.Context, data model.BankCardData, meta model.BankCardMeta) error
	AddFile(ctx context.Context, file string, comment string) error
	GetData(ctx context.Context, id string) (model.DataType, any, error)
	GetAllByType(ctx context.Context, dataType model.DataType) ([]model.ItemInfo, error)
	DeleteData(ctx context.Context, id string) error
}

// CLI описывает структуру cli приложения.
type CLI struct {
	service Service

	rootCMD  *cobra.Command
	userCMD  *cobra.Command
	vaultCMD *cobra.Command
}

// NewCLI создает и возвращает новое cli приложение.
func NewCLI(ctx context.Context, service Service) *CLI {
	cli := &CLI{
		service: service,
		rootCMD: &cobra.Command{
			Use: "gophkeeper",
		},
		userCMD: &cobra.Command{
			Use:   "user",
			Short: "Команды аутентификации",
			Long:  "Команды создания нового пользователя, аутентификации по логину и паролю",
		},
		vaultCMD: &cobra.Command{
			Use:   "vault",
			Short: "Команды для работы с хранилищем",
			Long:  "Команды для работы с хранилищем - добавление, удаление, загрузка данных",
		},
	}

	aboutCMD := &cobra.Command{
		Use:   "about",
		Short: "О программе",
		Long:  "Информация о версии и дате сборки клиента",
		Run: func(_ *cobra.Command, _ []string) {
			version, date := config.GetBuildInfo()
			fmt.Printf("version=%s, build time=%s\n", version, date)
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
	cli.rootCMD.AddCommand(aboutCMD)

	return cli
}

// Execute запускает обработку команды.
func (c *CLI) Execute() error {
	return c.rootCMD.Execute()
}
