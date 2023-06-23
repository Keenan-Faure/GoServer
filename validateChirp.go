package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type chirp struct {
	Body string `json:"body"`
}

type responseBody struct {
	Cleaned_body string `json:"cleaned_body"`
}

// decodes the data sent over the request body
func validate_chirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) < 140 {
		validated, err := profane(params.Body)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid request body")
		} else {
			respondWithJSON(w, http.StatusOK, responseBody{
				Cleaned_body: validated,
			})
		}
		return
	}
	respondWithError(w, 400, "Chirp is too long")
}

func encode_chirp(w http.ResponseWriter, r *http.Request, data chirp) {
	defer r.Body.Close()

	type responseBody struct {
		error string
	}

	dat, err := json.Marshal(data)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error marshalling JSON "+err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func profane(sentence string) (string, error) {
	if sentence == "" || len(sentence) == 0 {
		return "", errors.New("undefined sentence")
	}
	result := []string{}
	words := strings.Split(sentence, " ")
	damena_kotoba := map[string]string{
		"kerfuffle": "****",
		"sharbert":  "****",
		"fornax":    "****",
	}
	for _, value := range words {
		if entry, ok := damena_kotoba[strings.ToLower(value)]; ok {
			result = append(result, entry)
			continue
		}
		result = append(result, value)
	}
	return strings.Join(result, " "), nil
}
