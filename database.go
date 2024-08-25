package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	lock *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBstructure struct {
	Chirps []Chirp `json:"chirps"`
}

// Creating a new Database
func newDataBase(filePath string) *DB {
	db := &DB{
		path: filePath,
		lock: &sync.RWMutex{},
	}

	DBstructure := DBstructure{}

	// Check if file exists; if not, create it
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {

		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Can't create the file: %v", err)
		}

		data, err := json.Marshal(DBstructure)
		if err != nil {
			log.Fatalf("Can't create the file: %v", err)
		}

		err = os.WriteFile(db.path, data, 0777)
		if err != nil {
			log.Fatalf("Can't create the file: %v", err)
		}

		file.Close()
	}

	return db
}

// Ensures whether the file already exists or not
func (db *DB) ensureDB() {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		newDB := newDataBase(db.path)
		db.path = newDB.path
		db.lock = newDB.lock
	}
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBstructure, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	data, err := os.ReadFile(db.path)
	if err != nil {
		return DBstructure{}, err
	}

	dbStructure := DBstructure{}
	err = json.Unmarshal(data, &dbStructure)
	if err != nil {
		return DBstructure{}, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBstructure) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	data, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0777)
	if err != nil {
		return err
	}

	return nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newChirp := Chirp{
		Id:   len(dbStructure.Chirps) + 1,
		Body: body,
	}

	dbStructure.Chirps = append(dbStructure.Chirps, newChirp)

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return newChirp, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	return dbStructure.Chirps, nil
}
