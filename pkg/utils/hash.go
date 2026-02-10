package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates bcrypt hash from the given password.
func HashPassword(pass string, cost int) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(pass), cost)
	if err != nil {
		return "", fmt.Errorf("failed to generate bcrypt hash from password: %w", err)
	}

	return string(passwordHash), nil
}

// ComparePassword compare hashed password with given password.
func ComparePassword(hash, pass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass))
	if err != nil {
		return fmt.Errorf("failed to compare hash and bcrypted password: %w", err)
	}

	return nil
}
