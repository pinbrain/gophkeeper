package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func (c *CLI) RegisterCmd(ctx context.Context) *cobra.Command {
	var login, password string
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Регистрация",
		Long:  "Регистрация нового пользователя в gophkeeper",
		RunE: func(_ *cobra.Command, _ []string) error {
			token, err := c.service.Register(ctx, login, password)
			if err != nil {
				return err
			}
			fmt.Println("Пользователь успешно зарегистрирован! JWT: ", token)
			return nil
		},
	}
	cmd.Flags().StringVarP(&login, "login", "l", "", "логин")
	_ = cmd.MarkFlagRequired("login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "пароль")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}

func (c *CLI) LoginCmd(ctx context.Context) *cobra.Command {
	var login, password string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Вход",
		Long:  "Аутентификация по логину и паролю",
		RunE: func(_ *cobra.Command, _ []string) error {
			token, err := c.service.Login(ctx, login, password)
			if err != nil {
				return err
			}
			fmt.Println("Вход успешно выполнен! JWT: ", token)
			return nil
		},
	}
	cmd.Flags().StringVarP(&login, "login", "l", "", "логин")
	_ = cmd.MarkFlagRequired("login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "пароль")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}
