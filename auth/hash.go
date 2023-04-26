package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type hashEngine struct{}

func (h hashEngine) GetNewHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed generate hash, error: %w", err)
	}
	return string(hash), nil
}

func (h hashEngine) CompareHashWithPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func NewHashEngine() *hashEngine {
	return &hashEngine{}
}
