package validateChirp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type responseBody struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// decodes the data sent over the request body
func Validate_chirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := Chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if len(params.Body) < 140 {
		validated, err := profane(params.Body)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "invalid request body")
		} else {
			RespondWithJSON(w, http.StatusOK, responseBody{
				ID:   0,
				Body: validated,
			})
		}
		return
	}
	RespondWithError(w, 400, "Chirp is too long")
}

// func encode_chirp(w http.ResponseWriter, r *http.Request, data Chirp) {
// 	defer r.Body.Close()

// 	type responseBody struct {
// 		error string
// 	}

// 	dat, err := json.Marshal(data)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "error marshalling JSON "+err.Error())
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	w.Write(dat)
// }

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
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

func RespondWithError(w http.ResponseWriter, code int, msg string) error {
	return RespondWithJSON(w, code, map[string]string{"error": msg})
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
