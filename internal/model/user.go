package model

type User struct {
	ID           int
	Login        string
	Email        string
	PasswordHash string
}
