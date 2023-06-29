module goserver

replace internal/api => ./internal/api

replace internal/db => ./internal/db

replace internal/utils => ./internal/utils

replace internal/objects => ./internal/objects

go 1.20

require github.com/go-chi/chi/v5 v5.0.8

require (
	github.com/joho/godotenv v1.5.1
	golang.org/x/crypto v0.10.0
)

require github.com/golang-jwt/jwt/v5 v5.0.0 // indirect
