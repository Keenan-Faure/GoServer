package api

import (
	"fmt"
	"objects"
	"testing"
)

func createUser() objects.User {
	return objects.User{
		ID:       1,
		Email:    "test@gmail.com",
		Password: []byte("abc123"),
	}
}

func TestCreateJWT(t *testing.T) {
	jwtSecret := []byte("secret")
	expectedToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHkiLCJzdWIiOiIxIiwiZXhwIjoxNjg4MDI5MTQzLCJpYXQiOjE2ODgwMjkxMzN9.ETfOOWytt0JSnHOLj3lF-uauFZN8Oe-ae27Guf5z28g"
	fmt.Println("Test case 1 - Tokens not equal")
	expiredIn := 10
	user := createUser()
	token, err := CreateJWT(jwtSecret, expiredIn, user)
	if err != nil {
		t.Errorf("Unexpected error found: %s", err.Error())
	}
	if token == expectedToken {
		t.Errorf("Expected %s, but found %s", token, expectedToken)
	}
}
