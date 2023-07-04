package api

import (
	"fmt"
	"testing"
	"utils"
)

func TestProfane(t *testing.T) {
	fmt.Println("Test case 1 - No profane words")
	sentence := "I am chuky want to play?"
	expected := "I am chuky want to play?"
	actual, _ := utils.Profane(sentence)
	if actual != expected {
		t.Errorf("Expected %s but found %s", expected, actual)
	}

	fmt.Println("Test case 2 - few profane words")
	sentence = "I am chuky sharbert want to play the game kerfuffle pokemon? Please do not fornax me"
	expected = "I am chuky **** want to play the game **** pokemon? Please do not **** me"
	actual, _ = utils.Profane(sentence)
	if actual != expected {
		t.Errorf("Expected %s but found %s", expected, actual)
	}

	fmt.Println("Test case 3 - empty sentence")
	sentence = ""
	expected = "error"
	_, err := utils.Profane(sentence)
	if err == nil {
		t.Errorf("Expected error but found nil")
	}
}
