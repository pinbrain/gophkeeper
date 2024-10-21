package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type ClientConfig struct {
	ServerAddress string // Адрес gRPC сервера.
	JWT           string // JWT токен
	JWTMetaKey    string // Название ключа в мета gRPC запроса.
}

func InitConfig() (*ClientConfig, error) {
	// Файл с конфигурацией
	viper.SetConfigFile("clientConfig.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	// Переменные окружения
	viper.AutomaticEnv()
	_ = viper.BindEnv("ServerAddress", "SERVER_ADDRESS")
	_ = viper.BindEnv("jwt", "JWT")
	_ = viper.BindEnv("JWTMetaKey", "JWT_META_KEY")

	// Дефолтные значения
	viper.SetDefault("ServerAddress", ":8080")
	viper.SetDefault("JWTMetaKey", "jwt")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	clientConfig := &ClientConfig{}
	if err := viper.Unmarshal(clientConfig); err != nil {
		return nil, err
	}
	return clientConfig, nil
}

func GetJWTMetaKey() string {
	return viper.GetString("JWTMetaKey")
}

func GetJWT() string {
	return viper.GetString("jwt")
}

func SaveJWT(jwt string) error {
	viper.Set("jwt", jwt)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("не удалось сохранить JWT токен: %w", err)
	}
	return nil
}
