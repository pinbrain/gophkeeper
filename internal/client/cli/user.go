package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// RegisterCmd возвращает команду cobra для регистрации пользователя.
func (c *CLI) RegisterCmd(ctx context.Context) *cobra.Command {
	var login, password string
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register new user",
		Long:  "Register new user in gophkeeper",
		RunE: func(_ *cobra.Command, _ []string) error {
			token, err := c.service.Register(ctx, login, password)
			if err != nil {
				return err
			}
			fmt.Println("User successfully registered! JWT: ", token)
			return nil
		},
	}
	cmd.Flags().StringVarP(&login, "login", "l", "", "user login")
	_ = cmd.MarkFlagRequired("login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "user password")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}

// LoginCmd возвращает команду cobra для регистрации аутентификации пользователя.
func (c *CLI) LoginCmd(ctx context.Context) *cobra.Command {
	var login, password string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login",
		Long:  "Login with login and password",
		RunE: func(_ *cobra.Command, _ []string) error {
			token, err := c.service.Login(ctx, login, password)
			if err != nil {
				return err
			}
			fmt.Println("User successfully logged in! JWT: ", token)
			return nil
		},
	}
	cmd.Flags().StringVarP(&login, "login", "l", "", "user login")
	_ = cmd.MarkFlagRequired("login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "user password")
	_ = cmd.MarkFlagRequired("password")
	return cmd
}
