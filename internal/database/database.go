package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// Chirp represents a single chirp
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

// DBStructure represents the structure of the database file
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// DB encapsulates the path to the database file and a mutex
type DB struct {
	path string
	mux  *sync.RWMutex
}

func (db *DB) ensureDB() error {
	db.mux.Lock()
	defer db.mux.Unlock()

	if _, err := os.Stat(db.path); errors.Is(err, os.ErrNotExist) {
		// Create an empty database file if it doesn't exist
		initialData := DBStructure{
			Chirps: make(map[int]Chirp),
		}
		return db.writeDB(initialData)
	}
	return nil
}

func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	// Ensure the database file exists
	if err := db.ensureDB(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	file, err := os.Open(db.path)
	if err != nil {
		return DBStructure{}, err
	}
	defer file.Close()

	var dbStructure DBStructure
	if err := json.NewDecoder(file).Decode(&dbStructure); err != nil {
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	file, err := os.Create(db.path)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(dbStructure); err != nil {
		return err
	}

	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Find the next available ID
	newID := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   newID,
		Body: body,
	}

	// Save the chirp to the map
	dbStructure.Chirps[newID] = chirp

	// Write the updated structure to disk
	if err := db.writeDB(dbStructure); err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}
