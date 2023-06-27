package db

import (
	"encoding/json"
	"errors"
	"os"
	"sort"
	"sync"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}
	if !checkFileExists(path) {
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
func (db *DB) CreateChirp(body string) (Chirp, error) {
	//get the ID of the new Chirp
	data, err := db.LoadDB()
	if err != nil {
		return Chirp{}, err
	}
	newChirp := Chirp{
		ID:   len(data.Chirps) + 1,
		Body: body,
	}
	data.Chirps[len(data.Chirps)+1] = newChirp
	err = db.writeDB(data)
	if err != nil {
		return Chirp{}, err
	}
	return newChirp, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	rawData := DBStructure{
		Chirps: make(map[int]Chirp),
	}
	data, _ := json.MarshalIndent(rawData, "", " ")
	err := os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	data, err := db.LoadDB()
	if err != nil {
		return []Chirp{}, err
	}
	chirps := SortChirps(data.Chirps)
	return chirps, nil
}

// checkFileExists checks if a file exists
func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

// loadDB reads the database file into memory
func (db *DB) LoadDB() (DBStructure, error) {
	rawData, err := os.ReadFile(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	var chirpdb DBStructure
	db.mux.Lock()
	defer db.mux.Unlock()
	json.Unmarshal(rawData, &chirpdb)
	return chirpdb, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	data, _ := json.MarshalIndent(dbStructure, "", " ")
	db.mux.Lock()
	defer db.mux.Unlock()
	err := os.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// sorts the chirps in ascending order
func SortChirps(chirps map[int]Chirp) []Chirp {
	result := []Chirp{}
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

// retrieves a chirp that has the respective `id`
func RetrieveChirp(id int, chirps map[int]Chirp) (Chirp, bool) {
	for key, value := range chirps {
		if key == id {
			return value, true
		}
		continue
	}
	return Chirp{}, false
}
