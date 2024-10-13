package model

import "time"

type DataType string

const (
	Password DataType = "PASSWORD"
	Text     DataType = "TEXT"
	BankCard DataType = "BANK_CARD"
	File     DataType = "FILE"
)

type VaultItem struct {
	ID          string
	UserID      string
	EncryptData []byte
	Meta        string
	Type        DataType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type PasswordMeta struct {
	Resource string `json:"resource"`
	Login    string `json:"login"`
	Comment  string `json:"comment"`
}

type TextMeta struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type BankCardMeta struct {
	Bank    string `json:"bank"`
	Comment string `json:"comment"`
}

type FileMeta struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Comment   string `json:"comment"`
}

type BankCardData struct {
	Number     string `json:"number"`
	ValidMonth int    `json:"validMonth"`
	ValidYear  int    `json:"validYear"`
	Holder     string `json:"holder"`
	CSV        string `json:"csv"`
}
