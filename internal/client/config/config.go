// Package config формирует конфигурацию приложения.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// ClientConfig определяет структуру конфигурации клиента.
type ClientConfig struct {
	ServerAddress string // Адрес gRPC сервера.
	JWT           string // JWT токен
	JWTMetaKey    string // Название ключа в мета gRPC запроса.
}

// InitConfig формирует итоговую конфигурацию приложения.
func InitConfig(Version, BuildTime string) (*ClientConfig, error) {
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

	viper.Set("Version", Version)
	viper.Set("BuildTime", BuildTime)

	clientConfig := &ClientConfig{}
	if err := viper.Unmarshal(clientConfig); err != nil {
		return nil, err
	}
	return clientConfig, nil
}

// GetJWTMetaKey возвращает мета ключ jwt из конфигурации.
func GetJWTMetaKey() string {
	return viper.GetString("JWTMetaKey")
}

// GetJWT возвращает текущий jwt.
func GetJWT() string {
	return viper.GetString("jwt")
}

// SaveJWT сохраняет новый jwt (заменяет существующий и перезаписывает в файле конфигурации).
func SaveJWT(jwt string) error {
	viper.Set("jwt", jwt)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("не удалось сохранить JWT токен: %w", err)
	}
	return nil
}

// GetBuildInfo возвращает инфо о сборке - версию и дату
func GetBuildInfo() (string, string) {
	return viper.GetString("Version"), viper.GetString("BuildTime")
}
