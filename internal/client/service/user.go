package service

import (
	"context"
	"fmt"

	"github.com/pinbrain/gophkeeper/internal/client/config"
	"github.com/pinbrain/gophkeeper/internal/proto"
	"google.golang.org/grpc/status"
)

// Register регистрирует нового пользователя.
func (s *Service) Register(ctx context.Context, login, password string) (string, error) {
	res, err := s.grpcClient.UserClient.Register(ctx, &proto.RegisterReq{
		Login: login, Password: password,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return "", fmt.Errorf("не удалось зарегистрировать пользователя: %s", s.Message())
		}
		return "", err
	}
	err = config.SaveJWT(res.GetToken())
	if err != nil {
		return res.GetToken(), fmt.Errorf("ошибка регистрации: %w", err)
	}
	return res.GetToken(), nil
}

// Login аутентифицирует пользователя по логину и паролю.
func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	res, err := s.grpcClient.UserClient.Login(ctx, &proto.LoginReq{
		Login: login, Password: password,
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return "", fmt.Errorf("не удалось залогиниться: %s", s.Message())
		}
		return "", err
	}
	err = config.SaveJWT(res.GetToken())
	if err != nil {
		return res.GetToken(), fmt.Errorf("ошибка аутентификации: %w", err)
	}
	return res.GetToken(), nil
}
