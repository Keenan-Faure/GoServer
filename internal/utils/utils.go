package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// hashes a password
func HashPassword(password []byte) ([]byte, error) {
	// Hashing the password with the default cost of 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return []byte(""), err
	}
	return hashedPassword, nil
}
