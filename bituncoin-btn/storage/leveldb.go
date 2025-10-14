package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// WalletData represents wallet information stored in database
type WalletData struct {
	Address         string             `json:"address"`
	Balances        map[string]float64 `json:"balances"`
	Transactions    []string           `json:"transactions"`
	EncryptedKeys   string             `json:"encrypted_keys"`
	AuthConfig      map[string]string  `json:"auth_config"`
}

// LevelDB simulates a LevelDB key-value store
type LevelDB struct {
	dataDir string
	cache   map[string][]byte
	mu      sync.RWMutex
}

// NewLevelDB creates a new LevelDB instance
func NewLevelDB(dataDir string) (*LevelDB, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	db := &LevelDB{
		dataDir: dataDir,
		cache:   make(map[string][]byte),
	}

	// Load existing data
	if err := db.loadCache(); err != nil {
		return nil, err
	}

	return db, nil
}

// Put stores a key-value pair
func (db *LevelDB) Put(key string, value []byte) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.cache[key] = value
	
	// Persist to disk
	return db.persistKey(key, value)
}

// Get retrieves a value by key
func (db *LevelDB) Get(key string) ([]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	value, exists := db.cache[key]
	if !exists {
		return nil, errors.New("key not found")
	}

	return value, nil
}

// Delete removes a key-value pair
func (db *LevelDB) Delete(key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.cache, key)
	
	// Remove from disk
	filePath := filepath.Join(db.dataDir, key+".json")
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// Has checks if a key exists
func (db *LevelDB) Has(key string) bool {
	db.mu.RLock()
	defer db.mu.RUnlock()

	_, exists := db.cache[key]
	return exists
}

// GetAll returns all keys and values
func (db *LevelDB) GetAll() (map[string][]byte, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	result := make(map[string][]byte)
	for k, v := range db.cache {
		result[k] = v
	}

	return result, nil
}

// persistKey writes a key-value pair to disk
func (db *LevelDB) persistKey(key string, value []byte) error {
	filePath := filepath.Join(db.dataDir, key+".json")
	return os.WriteFile(filePath, value, 0600)
}

// loadCache loads all data from disk into memory
func (db *LevelDB) loadCache() error {
	files, err := os.ReadDir(db.dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		key := file.Name()[:len(file.Name())-5] // Remove .json extension
		filePath := filepath.Join(db.dataDir, file.Name())
		
		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		db.cache[key] = data
	}

	return nil
}

// Close closes the database
func (db *LevelDB) Close() error {
	// In a real implementation, this would flush any pending writes
	return nil
}

// SaveWallet saves wallet data to database
func (db *LevelDB) SaveWallet(address string, wallet *WalletData) error {
	data, err := json.Marshal(wallet)
	if err != nil {
		return err
	}

	return db.Put("wallet:"+address, data)
}

// LoadWallet loads wallet data from database
func (db *LevelDB) LoadWallet(address string) (*WalletData, error) {
	data, err := db.Get("wallet:" + address)
	if err != nil {
		return nil, err
	}

	var wallet WalletData
	if err := json.Unmarshal(data, &wallet); err != nil {
		return nil, err
	}

	return &wallet, nil
}

// SaveTransaction saves a transaction to database
func (db *LevelDB) SaveTransaction(txID string, txData []byte) error {
	return db.Put("tx:"+txID, txData)
}

// LoadTransaction loads a transaction from database
func (db *LevelDB) LoadTransaction(txID string) ([]byte, error) {
	return db.Get("tx:" + txID)
}

// GetAllWallets returns all wallet addresses
func (db *LevelDB) GetAllWallets() ([]string, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	var addresses []string
	for key := range db.cache {
		if len(key) > 7 && key[:7] == "wallet:" {
			addresses = append(addresses, key[7:])
		}
	}

	return addresses, nil
}
