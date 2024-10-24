// Package utils предоставляет различные универсальные функции, используемые в разных пакетах сервера.
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// GenerateRandomBytes генерирует рандомный массив байт заданной длины.
func GenerateRandomBytes(size int) ([]byte, error) {
	// генерируем криптостойкие случайные байты в b
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GeneratePasswordHash генерирует hash из пароля.
func GeneratePasswordHash(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate pwd hash: %w", err)
	}
	return string(hashedBytes), nil
}

// ComparePwdAndHash сравнивает хэш пароля и паролем (проверка валидности пароля).
func ComparePwdAndHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GenerateUserKey генерирует ключ пользователя для шифрования данных.
func GenerateUserKey() ([]byte, error) {
	key, err := GenerateRandomBytes(2 * aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt шифрует данные с помощью переданного ключа.
func Encrypt(data []byte, key string) ([]byte, error) {
	keyB, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	aesblock, err := aes.NewCipher(keyB)
	if err != nil {
		return nil, err
	}
	// NewGCM возвращает заданный 128-битный блочный шифр
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}
	// создаём вектор инициализации
	nonce, err := GenerateRandomBytes(aesgcm.NonceSize())
	if err != nil {
		return nil, err
	}
	dst := aesgcm.Seal(nonce, nonce, data, nil) // зашифровываем
	return dst, nil
}

// Decrypt расшифровывает данные с помощью ключа.
func Decrypt(data []byte, key string) ([]byte, error) {
	keyB, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}
	aesblock, err := aes.NewCipher(keyB)
	if err != nil {
		return nil, err
	}
	// NewGCM возвращает заданный 128-битный блочный шифр
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}
	nonceSize := aesgcm.NonceSize()
	// получаем вектор инициализации
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plain, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
