package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileserverHits int
}

//stateful handlers
// For example, we might want to keep track of the number of requests we've received,
// or we may want to pass around an open connection to a database, or credentials to an API

// increments each time a request is made to the handler
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, cfg.fileserverHits)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

// middleware to log every request to server
func middleWareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		fmt.Println(r.Method)
		next.ServeHTTP(w, r)
	})
}
