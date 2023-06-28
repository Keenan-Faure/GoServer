package utils

import (
	"fmt"
	"testing"
)

const dbPath = "./database.json"

func TestHashPassword(t *testing.T) {
	fmt.Println("Test Case 1 - empty password")

	password := []byte("")
	_, err := HashPassword(password)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	fmt.Println(password)

	fmt.Println("Test Case 2 - valid password")
	password = []byte("passwordString")
	_, err = HashPassword(password)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
}
