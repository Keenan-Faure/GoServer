package db

import (
	"fmt"
	"objects"
	"os"
	"testing"
)

const dbPath = "../../database.json"

func TestNewDB(t *testing.T) {
	fmt.Println("Test case 1 - DB does not exist")
	_, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	fmt.Println("Test case 1 - DB does exist")
	_, err = NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	os.Remove(dbPath)
}

func TestCreateChirp(t *testing.T) {
	fmt.Println("Test case 1 - DB does not exist")
	Db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	chirp, err := Db.CreateChirp("Test", 1)
	if err != nil {
		t.Errorf("Expected nil but found error")
	}
	if chirp.ID != 1 {
		t.Errorf("Expected 1 but found %d", chirp.ID)
	}
	if chirp.Body != "Test" {
		t.Errorf("Expected 'Test' but found %s", chirp.Body)
	}

	fmt.Println("Test case 2 - DB already exists")
	Db, err = NewDB(dbPath)
	if err != nil {
		t.Errorf("Unexpected error")
	}
	chirp, err = Db.CreateChirp("Test", 1)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	if chirp.ID != 2 {
		t.Errorf("Expected 2 but found %d", chirp.ID)
	}
	if chirp.Body != "Test" {
		t.Errorf("Expected 'Test' but found %s", chirp.Body)
	}

	fmt.Println("Test case 3 - Appending new Chirp and load from disk")
	Db, err = NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	chirp, err = Db.CreateChirp("Test String", 2)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	chirps, err := Db.GetChirps("1", "")
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	if len(chirps) != 3 {
		t.Errorf("Expected '0' but found %d", len(chirps))
	}
	os.Remove(dbPath)
}

func TestEnsureDb(t *testing.T) {
	fmt.Println("Test case 1 - DB does not exist")
	Db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	err = Db.ensureDB()
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	if !CheckFileExists(dbPath) {
		t.Errorf("Expected 'true' but found %v", !CheckFileExists(dbPath))
	}

	fmt.Println("Test case 2 - DB does exist")
	Db, err = NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	err = Db.ensureDB()
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	if !CheckFileExists(dbPath) {
		t.Errorf("Expected 'true' but found %v", !CheckFileExists(dbPath))
	}
	os.Remove(dbPath)
}

func TestGetChirps(t *testing.T) {
	fmt.Println("Test case 1 - DB does not exist")
	Db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	chirps, err := Db.GetChirps("", "")
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	if len(chirps) > 0 {
		t.Errorf("Expected '0' but found %d", len(chirps))
	}

	fmt.Println("Test case 2 - DB does exist && added a chirp")
	Db, err = NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	Db.CreateChirp("Test String", 1)
	chirps, err = Db.GetChirps("", "")
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	if len(chirps) != 1 {
		t.Errorf("Expected '0' but found %d", len(chirps))
	}
	os.Remove(dbPath)
}

func TestLoadDb(t *testing.T) {
	fmt.Println("Test case 1 - DB does not exist")
	Db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	data, err := Db.LoadDB()
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	if len(data.Chirps) != 0 {
		t.Errorf("Expected '0' but found %d", len(data.Chirps))
	}
	os.Remove(dbPath)
}

func TestWriteData(t *testing.T) {
	fmt.Println("Test case 1 - Writing data to new database")
	Db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	dbStruct := objects.DBStructure{
		Chirps: make(map[int]objects.Chirp),
	}
	err = Db.writeDB(dbStruct)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	data, err := Db.LoadDB()
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	if len(data.Chirps) != 0 {
		t.Errorf("Expected '0' but found %d", len(data.Chirps))
	}
	os.Remove(dbPath)
}

func TestSortChirps(t *testing.T) {
	fmt.Println("Test case 1 - Test sorting two chirps")
	Db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	Db.CreateChirp("Test String 1", 1)
	Db.CreateChirp("Test String 2", 1)
	Db.CreateChirp("Test String 3", 1)
	Db.CreateChirp("Test String 4", 2)
	Db.CreateChirp("Test String 5", 2)
	Db.CreateChirp("Test String 6", 3)

	data, err := Db.LoadDB()
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	result := SortChirpsAsc(data.Chirps)
	if len(result) != 6 {
		t.Errorf("Expected 6 but found %d", len(result))
	}

	fmt.Println("Test case 2 - Manually creating a new map and sorting")
	chirps := map[int]objects.Chirp{
		3: {
			ID:   3,
			Body: "Text 3",
		},
		5: {
			ID:   5,
			Body: "Text 5",
		},
		1: {
			ID:   1,
			Body: "Text 1",
		},
	}
	result = SortChirpsAsc(chirps)
	if len(result) != 3 {
		t.Errorf("Expected 3 but found %d", len(result))
	}
	os.Remove(dbPath)
}

func TestRetrieveChirp(t *testing.T) {
	fmt.Println("Test case 1 - Fetch an ID that does exist")
	Db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	data, err := Db.LoadDB()
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	Db.CreateChirp("Test String 1", 1)
	data, _ = Db.LoadDB()
	_, exist := RetrieveChirp(1, data.Chirps)
	if !exist {
		t.Errorf("Expected 'true', but found %v", exist)
	}
	fmt.Println("Test case 1 - Fetch an ID that does not")
	_, exist = RetrieveChirp(10, data.Chirps)
	if exist {
		t.Errorf("Expected 'false', but found %v", exist)
	}
	os.Remove(dbPath)
}

func TestValidateUser(t *testing.T) {
	fmt.Println("Test case 1 - Email does not exist")
	Db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	data, err := Db.LoadDB()
	if err != nil {
		t.Errorf("Expected error to be nil but found %s", err.Error())
	}
	exist, _ := ValidateUser(data, "abc123gmail.com")
	if exist {
		t.Errorf("Expected 'false' but found 'true'")
	}
	_, err = Db.CreateUser("abc123gmail.com", []byte("passwordTest"))
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}

	fmt.Println("Test case 2 - Email does exist (error)")
	data, _ = Db.LoadDB()
	_, err = ValidateUser(data, "abc123gmail.com")
	if err == nil {
		t.Errorf("Expected 'error' but found 'nil'")
	}
	os.Remove(dbPath)
}

func TestValidatePassword(t *testing.T) {
	fmt.Println("Test case 1 - Password found")
	Db, err := NewDB(dbPath)
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	data, err := Db.LoadDB()
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	exist, _ := ValidateUser(data, "abc123gmail.com")
	if exist {
		t.Errorf("Expected 'false' but found 'true'")
	}
	_, err = Db.CreateUser("abc123gmail.com", []byte("passwordTest"))
	if err != nil {
		t.Errorf("Expected error to be 'nil' but found: %s", err.Error())
	}
	data, err = Db.LoadDB()
	if err != nil {
		t.Errorf("Expected error to be nil but found: %s", err.Error())
	}
	_, err = ValidateLogin(data, "abc123gmail.com", []byte("passwordTest"))
	if err != nil {
		t.Errorf("Expected 'nil', but found %s", err.Error())
	}
	fmt.Println("Test case 2 - Password not found")
	_, err = ValidateLogin(data, "abc123gmail.com", []byte("abc12345"))
	if err == nil {
		t.Errorf("Expected 'error', but found %v", err)
	}
	os.Remove(dbPath)
}
