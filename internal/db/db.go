package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
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
	data, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	newChirp := Chirp{
		ID:   len(data.Chirps),
		Body: body,
	}
	return newChirp, nil
}

func (db *DB) ensureDB() error {
	rawData := DBStructure{}
	data, _ := json.MarshalIndent(rawData, "", " ")
	err := ioutil.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	chirps := []Chirp{}
	data, err := db.loadDB()
	if err != nil {
		return chirps, err
	}
	for _, value := range data.Chirps {
		chirps = append(chirps, value)
	}
	return chirps, nil
}

// checks if a file exists
func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

func (db *DB) loadDB() (DBStructure, error) {
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
	err := ioutil.WriteFile(db.path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}
