package main

import (
	"net/http"
)

const baseAddress = ":8080"

func main() {
	// webServer()
	fileServer()
}

func fileServer() {
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./app/"))
	mux.Handle("/app/", http.StripPrefix("/app", fs))

	//HandleFunc registers the handler function for the given pattern.
	//endpoint
	mux.HandleFunc("/healthz", healthz)

	http.ListenAndServe(baseAddress, middlewareCors(mux))
}

func webServer() {
	mux := http.NewServeMux()

	//any software or service that enables the parts of a system to communicate and manage data
	corsMux := middlewareCors(mux)

	server := &http.Server{
		Addr:    baseAddress,
		Handler: corsMux,
	}
	server.ListenAndServe()
}

// Acts as the middleware that adds CORS headers to the response
func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Custom handlers

func healthz(writer http.ResponseWriter, reader *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.Write([]byte("OK"))
}
