package utils

import (
	"errors"
	"strings"

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

// profane replaced certain words with asterisks
// which are defined in a map
func Profane(sentence string) (string, error) {
	if sentence == "" || len(sentence) == 0 {
		return "", errors.New("undefined sentence")
	}
	result := []string{}
	words := strings.Split(sentence, " ")
	damena_kotoba := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}
	for _, value := range words {
		if entry, ok := damena_kotoba[strings.ToLower(value)]; ok {
			result = append(result, entry)
			continue
		}
		result = append(result, value)
	}
	return strings.Join(result, " "), nil
}
