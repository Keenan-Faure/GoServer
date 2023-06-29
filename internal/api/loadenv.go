package api

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() string {
	// by default, godotenv will look for a file named .env in the current directory
	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	return jwtSecret
}
