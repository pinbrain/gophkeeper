package model

// User описывает структуру данных пользователя.
type User struct {
	ID              string
	Login           string
	PasswordHash    string
	EncryptedSecret string
}
