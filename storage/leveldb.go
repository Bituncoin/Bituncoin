package storage

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// LevelDB represents a simple key-value storage
type LevelDB struct {
	dataDir string
	cache   map[string][]byte
	mutex   sync.RWMutex
}

// NewLevelDB creates a new LevelDB instance
func NewLevelDB(dataDir string) (*LevelDB, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	db := &LevelDB{
		dataDir: dataDir,
		cache:   make(map[string][]byte),
	}

	// Load existing data
	if err := db.loadFromDisk(); err != nil {
		return nil, err
	}

	return db, nil
}

// Put stores a key-value pair
func (db *LevelDB) Put(key string, value []byte) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if key == "" {
		return errors.New("key cannot be empty")
	}

	db.cache[key] = value
	return db.saveToDisk(key, value)
}

// Get retrieves a value by key
func (db *LevelDB) Get(key string) ([]byte, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	value, exists := db.cache[key]
	if !exists {
		return nil, errors.New("key not found")
	}

	return value, nil
}

// Delete removes a key-value pair
func (db *LevelDB) Delete(key string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	if _, exists := db.cache[key]; !exists {
		return errors.New("key not found")
	}

	delete(db.cache, key)
	
	// Remove from disk
	filePath := filepath.Join(db.dataDir, key+".dat")
	return os.Remove(filePath)
}

// Has checks if a key exists
func (db *LevelDB) Has(key string) bool {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	_, exists := db.cache[key]
	return exists
}

// Keys returns all keys
func (db *LevelDB) Keys() []string {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	keys := make([]string, 0, len(db.cache))
	for key := range db.cache {
		keys = append(keys, key)
	}

	return keys
}

// saveToDisk saves a key-value pair to disk
func (db *LevelDB) saveToDisk(key string, value []byte) error {
	filePath := filepath.Join(db.dataDir, key+".dat")
	return ioutil.WriteFile(filePath, value, 0644)
}

// loadFromDisk loads all data from disk
func (db *LevelDB) loadFromDisk() error {
	files, err := ioutil.ReadDir(db.dataDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if filepath.Ext(fileName) != ".dat" {
			continue
		}

		key := fileName[:len(fileName)-4]
		filePath := filepath.Join(db.dataDir, fileName)
		
		value, err := ioutil.ReadFile(filePath)
		if err != nil {
			continue
		}

		db.cache[key] = value
	}

	return nil
}

// PutJSON stores a JSON-encoded value
func (db *LevelDB) PutJSON(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return db.Put(key, data)
}

// GetJSON retrieves and decodes a JSON value
func (db *LevelDB) GetJSON(key string, dest interface{}) error {
	data, err := db.Get(key)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// Close closes the database
func (db *LevelDB) Close() error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	// Flush any pending writes
	for key, value := range db.cache {
		if err := db.saveToDisk(key, value); err != nil {
			return err
		}
	}

	db.cache = make(map[string][]byte)
	return nil
}

// Size returns the number of key-value pairs
func (db *LevelDB) Size() int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	return len(db.cache)
}
