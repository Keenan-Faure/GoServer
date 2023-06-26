package db

import (
	"fmt"
	"testing"
)

func TestNewDB(t *testing.T) {
	fmt.Println("Test case 1 - DB does not exist")
	_, err := NewDB("../../database.json")
	if err != nil {
		t.Errorf("Error")
	}
}
