package db

import (
	"encoding/json"
	"errors"
	"objects"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
	"utils"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	if !CheckFileExists(path) {
		//creates db
		err := db.ensureDB()
		if err != nil {
			return &DB{}, err
		}
		return db, nil
	}
	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string, author int) (objects.Chirp, error) {
	//get the ID of the new Chirp
	data, err := db.LoadDB()
	if err != nil {
		return objects.Chirp{}, err
	}
	newChirp := objects.Chirp{
		ID:       len(data.Chirps) + 1,
		Body:     body,
		AuthorID: author,
	}
	db.AddChirp(data, newChirp)
	err = db.writeDB(data)
	if err != nil {
		return objects.Chirp{}, err
	}
	return newChirp, nil
}

// AddChirp to Database
func (db *DB) AddChirp(data objects.DBStructure, chirp objects.Chirp) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data.Chirps[len(data.Chirps)+1] = chirp
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	rawData := objects.DBStructure{
		Chirps:        make(map[int]objects.Chirp),
		Users:         make(map[int]objects.User),
		RevokedTokens: make(map[time.Time]string),
	}
	data, _ := json.MarshalIndent(rawData, "", " ")
	err := os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]objects.Chirp, error) {
	data, err := db.LoadDB()
	if err != nil {
		return []objects.Chirp{}, err
	}
	chirps := SortChirps(data.Chirps)
	return chirps, nil
}

// checkFileExists checks if a file exists
func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

// loadDB reads the database file into memory
func (db *DB) LoadDB() (objects.DBStructure, error) {
	rawData, err := os.ReadFile(db.path)
	if err != nil {
		return objects.DBStructure{}, err
	}
	var dbstruct objects.DBStructure
	db.mux.Lock()
	defer db.mux.Unlock()
	json.Unmarshal(rawData, &dbstruct)
	return dbstruct, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure objects.DBStructure) error {
	data, _ := json.MarshalIndent(dbStructure, "", " ")
	db.mux.Lock()
	defer db.mux.Unlock()
	err := os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// creates a new User and saves it to disk
func (db *DB) CreateUser(email string, password []byte) (objects.User, error) {
	data, err := db.LoadDB()
	if err != nil {
		return objects.User{}, err
	}
	exist, err := ValidateUser(data, email)
	if exist {
		return objects.User{}, err
	}
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return objects.User{}, err
	}
	newUser := objects.User{
		ID:       len(data.Users) + 1,
		Email:    email,
		Password: hashedPwd,
	}
	data.Users[len(data.Users)+1] = newUser
	err = db.writeDB(data)
	if err != nil {
		return objects.User{}, err
	}
	return newUser, nil
}

func (db *DB) GetUserByID(data objects.DBStructure, id int) (objects.User, error) {
	for _, value := range data.Users {
		if value.ID == id {
			return value, nil
		}
		continue
	}
	return objects.User{}, errors.New("unable to find user with id")
}

// updates a user in the database
func (db *DB) UpdateUser(id int, newEmail, newPassword string, database objects.DBStructure) (objects.User, error) {
	if id == 0 {
		return objects.User{}, errors.New("invalid ID")
	}
	db.mux.Lock()
	if entry, ok := database.Users[id]; ok {
		db.mux.Unlock()
		entry.Email = newEmail
		newPsw, err := utils.HashPassword([]byte(newPassword))
		if err != nil {
			return objects.User{}, err
		}
		entry.Password = newPsw
		database.Users[id] = entry
		defer db.writeDB(database)
		return entry, nil
	}
	return objects.User{}, errors.New("id not found in the database")
}

// revokes a JWT Token
func (db *DB) RevokeToken(token string, database objects.DBStructure) error {
	database.RevokedTokens[time.Now().UTC()] = token
	return db.writeDB(database)
}

// checks if a JWT Token has been revoked
func (db *DB) IsTokenRevoked(token string, database objects.DBStructure) bool {
	for _, value := range database.RevokedTokens {
		if value == token {
			return true
		}
		continue
	}
	return false
}

// removes a chrip from the specified database
func (db *DB) DeleteChirp(id int, database objects.DBStructure) error {
	for key, value := range database.Chirps {
		if value.ID == id {
			delete(database.Chirps, key)
			db.writeDB(database)
			return nil
		}
		continue
	}
	return errors.New("chirp with id '" + strconv.Itoa(id) + "' not found")
}

// helper functions

// sorts the chirps in ascending order
func SortChirps(chirps map[int]objects.Chirp) []objects.Chirp {
	result := []objects.Chirp{}
	if len(chirps) == 0 {
		return result
	}

	keys := make([]int, 0, len(chirps))
	for k := range chirps {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	for _, k := range keys {
		result = append(result, chirps[k])
	}
	return result
}

// retrieves a chirp that has the respective id
func RetrieveChirp(id int, chirps map[int]objects.Chirp) (objects.Chirp, bool) {
	for key, value := range chirps {
		if key == id {
			return value, true
		}
		continue
	}
	return objects.Chirp{}, false
}

// confirms if a user with the same email already exists
func ValidateUser(data objects.DBStructure, email string) (bool, error) {
	for _, value := range data.Users {
		if value.Email == email {
			return true, errors.New("user with email already exists")
		}
		continue
	}
	return false, nil
}

// returns the user with the respective email address and whether found
func GetUserByEmail(data objects.DBStructure, email string) (objects.User, bool) {
	for _, value := range data.Users {
		if value.Email == email {
			return value, true
		}
		continue
	}
	return objects.User{}, false
}

// determines if the password exists and returns the record of the user if found
func ValidateLogin(data objects.DBStructure, email string, password []byte) (objects.User, error) {
	usr, exists := GetUserByEmail(data, email)
	if !exists {
		return objects.User{}, errors.New("user with email " + email + " does not exist")
	}
	valid := bcrypt.CompareHashAndPassword(usr.Password, password)
	if valid == nil {
		return usr, nil
	}
	return objects.User{}, errors.New("password does not exist in database")
}

// determines if the user id exists in the database
func ValidateUserByID(data objects.DBStructure, id int) (bool, error) {
	for _, value := range data.Users {
		if value.ID == id {
			return true, nil
		}
		continue
	}
	return false, errors.New("unable to find user with id")
}

func (db *DB) ValidateUserChirp(chirpId, userID int, database objects.DBStructure) error {
	for _, value := range database.Chirps {
		if value.AuthorID == userID {
			db.DeleteChirp(chirpId, database)
			return nil
		}
		continue
	}
	return errors.New("user does not have any chirps")
}
