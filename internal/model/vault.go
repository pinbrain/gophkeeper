package model

import "time"

// DataType enum типов данных.
type DataType string

// Типы данных.
const (
	Password DataType = "PASSWORD"
	Text     DataType = "TEXT"
	BankCard DataType = "BANK_CARD"
	File     DataType = "FILE"
)

// VaultItem описывает структуру хранящихся данных.
type VaultItem struct {
	ID          string
	UserID      string
	EncryptData []byte
	Meta        string
	Type        DataType
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PasswordMeta описывает структуру мета данных пароля.
type PasswordMeta struct {
	Resource string `json:"resource"`
	Login    string `json:"login"`
	Comment  string `json:"comment"`
}

// TextMeta описывает структуру мета данных текстовой информации.
type TextMeta struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

// BankCardMeta описывает структуру мета данных банковской карты.
type BankCardMeta struct {
	Bank    string `json:"bank"`
	Comment string `json:"comment"`
}

// FileMeta описывает структуру мета данных файла.
type FileMeta struct {
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Comment   string `json:"comment"`
}

// BankCardData описывает структуру данных банковской карты.
type BankCardData struct {
	Number     string `json:"number"`
	ValidMonth int    `json:"validMonth"`
	ValidYear  int    `json:"validYear"`
	Holder     string `json:"holder"`
	CSV        string `json:"csv"`
}

// PasswordItem описывает структуру пароля, которую возвращает сервис.
type PasswordItem struct {
	Type DataType
	Meta PasswordMeta
	Data string
}

// TextItem описывает структуру текстовой информации, которую возвращает сервис.
type TextItem struct {
	Type DataType
	Meta TextMeta
	Data string
}

// BankCardItem описывает структуру банковской карты, которую возвращает сервис.
type BankCardItem struct {
	Type DataType
	Meta BankCardMeta
	Data BankCardData
}

// FileItem описывает структуру файла, которую возвращает сервис.
type FileItem struct {
	Type DataType
	Meta FileMeta
}

// ItemInfo описывает структуру данных для вывода списка.
type ItemInfo struct {
	ID   string
	Meta any
}
