package hasher

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type hasher struct{}

func (h hasher) GetNewHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed generate hash, error: %w", err)
	}
	return string(hash), nil
}

func (h hasher) CompareHashWithPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func New() *hasher {
	return &hasher{}
}
