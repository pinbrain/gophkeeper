// Package config формирует конфигурацию сервера.
package config

import (
	"errors"

	"github.com/spf13/viper"
)

// ServerConfig определяет структуру конфигурации сервера.
type ServerConfig struct {
	MasterKey     string    // Мастер ключ для шифрования.
	ServerAddress string    // Адрес gRPC сервера.
	LogLevel      string    // Уровень логирования.
	DSN           string    // Строка с адресом подключения к БД.
	JWT           JWTConfig // JWT конфигурация.
}

// JWTConfig определяет структуру конфигурации jwt.
type JWTConfig struct {
	LifeTime  int    // Время жизни токена в минутах.
	SecretKey string // Ключ для подписи jwt токена.
	MetaKey   string // Название ключа в мета gRPC запроса.
}

// InitConfig формирует итоговую конфигурацию сервера.
func InitConfig() (*ServerConfig, error) {
	// Файл с конфигурацией
	viper.SetConfigFile("serverConfig.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	// Переменные окружения
	viper.AutomaticEnv()
	_ = viper.BindEnv("MasterKey", "MASTER_KEY")
	_ = viper.BindEnv("ServerAddress", "SERVER_ADDRESS")
	_ = viper.BindEnv("LogLevel", "LOG_LEVEL")
	_ = viper.BindEnv("DSN", "DATABASE_DSN")
	_ = viper.BindEnv("JWT.LifeTime", "JWT_LIFE_TIME")
	_ = viper.BindEnv("JWT.SecretKey", "JWT_SECRET_KEY")
	_ = viper.BindEnv("JWT.MetaKey", "JWT_META_KEY")

	// Дефолтные значения
	viper.SetDefault("ServerAddress", ":8080")
	viper.SetDefault("LogLevel", "info")
	viper.SetDefault("JWT.LifeTime", "60")
	viper.SetDefault("JWT.SecretKey", "jwt_secret_key")
	viper.SetDefault("JWT.MetaKey", "jwt")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	severConfig := &ServerConfig{}
	if err := viper.Unmarshal(severConfig); err != nil {
		return nil, err
	}

	if severConfig.MasterKey == "" {
		return nil, errors.New("отсутствует мастер ключ")
	}

	return severConfig, nil
}
