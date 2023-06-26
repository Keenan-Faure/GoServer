module goserver

replace internal/validateChirp => ./internal/validateChirp

replace internal/db => ./internal/db

go 1.20

require github.com/go-chi/chi/v5 v5.0.8
