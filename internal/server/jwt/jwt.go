// Package jwt содержит реализацию сервиса для работы с jwt.
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/pinbrain/gophkeeper/internal/model"
	"github.com/pinbrain/gophkeeper/internal/server/config"
)

// Service описывает структуру jwt сервиса.
type Service struct {
	lifeTime  time.Duration
	secretKey string
	mdJWTKey  string
}

// Claims описывает структуру данных jwt.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
	Login  string
}

// NewJWTService создает и возвращает новый jwt сервис.
func NewJWTService(cfg config.JWTConfig) *Service {
	return &Service{
		lifeTime:  time.Duration(cfg.LifeTime) * time.Minute,
		secretKey: cfg.SecretKey,
		mdJWTKey:  cfg.MetaKey,
	}
}

// BuildJWTSting формирует jwt с переданными данными.
func (j *Service) BuildJWTSting(user *model.User) (string, error) {
	if user.ID == "" || user.Login == "" {
		return "", errors.New("not valid user data")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: user.ID,
		Login:  user.Login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.lifeTime)),
		},
	})

	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to build jwt string: %w", err)
	}
	return tokenString, nil
}

// GetJWTClaims возвращает данные из jwt, проверяя его валидность.
func (j *Service) GetJWTClaims(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt token: %w", err)
	}
	if !token.Valid {
		return nil, errors.New("invalid jwt token")
	}
	return claims, nil
}

// GetMdJWTKey возвращает ключ jwt в метаданных.
func (j *Service) GetMdJWTKey() string {
	return j.mdJWTKey
}
