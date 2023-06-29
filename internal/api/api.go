package api

import (
	"db"
	"encoding/json"
	"errors"
	"net/http"
	"objects"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

const dbPath = "./database.json"

// updates a user's email and password
// reads header for token
func UpdateUsers(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyUser{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	Db, err := db.NewDB(dbPath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	dbstruct, err := Db.LoadDB()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	//extract jwt from header and return it
	authHeader := ExtractJWT(r.Header.Get("Authorization"))

	//validate jwt
	id, exist, err := ValidateJWT(authHeader)
	if err != nil {
		if exist {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	//validate user id and update
	exist, err = db.ValidateUserByID(dbstruct, id)
	if !exist && err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	updatedUser, err := Db.UpdateUser(id, params.Email, params.Password, dbstruct)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	response := objects.RequestBodyUser{
		Password: string(updatedUser.Password),
		Email:    updatedUser.Email,
	}
	RespondWithJSON(w, http.StatusOK, response)
}

// login function, responds with a JWT token
func UserLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyLogin{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	Db, err := db.NewDB(dbPath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	dbstruct, err := Db.LoadDB()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	usr, err := db.ValidateLogin(dbstruct, params.Email, []byte(params.Password))
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	jwtSecret := []byte(LoadEnv())
	jwtToken, err := CreateJWT(jwtSecret, params.ExpiresSeconds, usr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	response := objects.ResponseUserLogon{
		ID:    usr.ID,
		Email: usr.Email,
		Token: jwtToken,
	}
	RespondWithJSON(w, http.StatusOK, response)
}

// posts data to add a new user
func PostUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyUser{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	Db, err := db.NewDB(dbPath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	user, err := Db.CreateUser(params.Email, []byte(params.Password))
	response := objects.ResponseUser{
		ID:    user.ID,
		Email: user.Email,
	}
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, response)
}

// posts data to the database to add a chirp
func PostValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyChirp{}
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
			Db, err := db.NewDB(dbPath)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
				return
			}
			chirp, err := Db.CreateChirp(validated)
			if err != nil {
				RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
				return
			}
			RespondWithJSON(w, http.StatusCreated, chirp)
		}
		return
	}
	RespondWithError(w, http.StatusBadRequest, "Chirp is too long")
}

// Gets all Chirps
func GetChirps(w http.ResponseWriter, r *http.Request) {
	Db, err := db.NewDB(dbPath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	chirps, err := Db.GetChirps()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, chirps)
}

// gets a specific chirp from the database
func GetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Invalid {chirpID}")
		return
	}
	Db, err := db.NewDB(dbPath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	dbstruct, err := Db.LoadDB()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	chirp, exist := db.RetrieveChirp(chirpID, dbstruct.Chirps)
	if !exist {
		RespondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	RespondWithJSON(w, http.StatusOK, chirp)
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

//helper methods

// profane replaced certain words with asterisks
// which are defined in a map
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
