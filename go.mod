module goserver

replace internal/api => ./internal/api

replace internal/db => ./internal/db

go 1.20

require github.com/go-chi/chi/v5 v5.0.8
