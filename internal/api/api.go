package api

import (
	"db"
	"docs"
	"encoding/json"
	"net/http"
	"objects"
	"strconv"
	"utils"

	"github.com/go-chi/chi/v5"
)

const dbPath = "./database.json"

// displays all available endpoints
func Endpoints(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, http.StatusOK, docs.Endpoints())
}

// webhook that accept payment details from `Polka` for a user
func PostWebhook(w http.ResponseWriter, r *http.Request) {
	// _ = LoadAPIKey()
	apiKey := ExtractAPIKey(r.Header.Get("Authorization"))
	if apiKey == "" {
		RespondWithJSON(w, http.StatusUnauthorized, objects.BaseResponse{
			Body: "",
		})
		return
	}
	// if apiKey != apiSecret {
	// 	RespondWithJSON(w, http.StatusUnauthorized, objects.BaseResponse{
	// 		Body: "",
	// 	})
	// 	return
	// }
	decoder := json.NewDecoder(r.Body)
	params := objects.WebhookRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if params.Event != "user.upgraded" {
		RespondWithJSON(w, http.StatusOK, objects.BaseResponse{
			Body: "",
		})
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
	usr, err := Db.GetUserByID(dbstruct, params.Data.UserID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
	}
	_, err = Db.UpdateUserUpgrade(usr.ID, true, dbstruct)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	RespondWithJSON(w, http.StatusOK, objects.BaseResponse{
		Body: "",
	})
}

// revokes a refresh token present in the header
func RevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
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
	jwtToken := ExtractJWT(r.Header.Get("Authorization"))
	err = Db.RevokeToken(jwtToken, dbstruct)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.ResponseRefreshToken{})
}

// returns a new access token if the refresh token is valid
func RefreshToken(w http.ResponseWriter, r *http.Request) {
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
	jwtToken := ExtractJWT(r.Header.Get("Authorization"))
	usr, err := ValidateJWTRefresh(jwtToken, Db, dbstruct)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	jwtSecret := []byte(LoadJWTSecret())
	jwtNewTokenAccess, err := CreateJWTAccess(jwtSecret, usr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	response := objects.ResponseRefreshToken{
		Token: jwtNewTokenAccess,
	}
	RespondWithJSON(w, http.StatusOK, response)
}

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
	jwtToken := ExtractJWT(r.Header.Get("Authorization"))

	id, exist, err := ValidateJWTAccess(jwtToken)
	if err != nil {
		if exist {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	exist, err = db.ValidateUserByID(dbstruct, id)
	if !exist && err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	updatedUser, err := Db.UpdateUser(id, params.Email, params.Password, false, dbstruct)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	response := objects.ResponseUser{
		ID:    updatedUser.ID,
		Email: updatedUser.Email,
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
	jwtSecret := []byte(LoadJWTSecret())
	jwtTokenAccess, err := CreateJWTAccess(jwtSecret, usr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	jwtTokenRefresh, err := CreateJWTRefresh(jwtSecret, usr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
	}
	response := objects.ResponseUserLogon{
		ID:           usr.ID,
		Email:        usr.Email,
		IsChirpyRed:  usr.IsChirpyRed,
		Token:        jwtTokenAccess,
		RefreshToken: jwtTokenRefresh,
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

// deletes a chirp from the database
func DelChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(chi.URLParam(r, "chirpID"))
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "error retrieving id")
	}
	jwtToken := ExtractJWT(r.Header.Get("Authorization"))
	id, exist, err := ValidateJWTAccess(jwtToken)
	if err != nil {
		if exist {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	Db, err := db.NewDB(dbPath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	dbstruct, err := Db.LoadDB()
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	exist, err = db.ValidateUserByID(dbstruct, id)
	if !exist && err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	err = Db.DeleteUserChirp(chirpID, id, dbstruct)
	if err != nil {
		RespondWithError(w, http.StatusForbidden, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, objects.BaseResponse{
		Body: "",
	})
}

// posts data to the database to add a chirp
func PostChirp(w http.ResponseWriter, r *http.Request) {
	jwtToken := ExtractJWT(r.Header.Get("Authorization"))
	id, exist, err := ValidateJWTAccess(jwtToken)
	if err != nil {
		if exist {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	Db, err := db.NewDB(dbPath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	dbstruct, err := Db.LoadDB()
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	exist, err = db.ValidateUserByID(dbstruct, id)
	if !exist && err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := objects.RequestBodyChirp{}
	err = decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if len(params.Body) < 140 {
		validated, err := utils.Profane(params.Body)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "invalid request body")
		} else {
			chirp, err := Db.CreateChirp(validated, id)
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
	author_id := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	Db, err := db.NewDB(dbPath)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Error: "+err.Error())
		return
	}
	chirps, err := Db.GetChirps(author_id, sort)
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
