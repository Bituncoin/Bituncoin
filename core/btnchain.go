package core

import (
	"errors"
	"sync"
	"time"
)

// Blockchain represents the Bituncoin blockchain
type Blockchain struct {
	Blocks     []*Block
	Difficulty int
	mutex      sync.RWMutex
}

// Block represents a block in the blockchain
type Block struct {
	Index        int
	Timestamp    int64
	Transactions []string
	PrevHash     string
	Hash         string
	Nonce        int
	Validator    string
}

// NewBlockchain creates a new blockchain instance
func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Blocks:     make([]*Block, 0),
		Difficulty: 2,
	}
	
	// Create genesis block
	genesisBlock := &Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []string{"Genesis Block"},
		PrevHash:     "0",
		Hash:         "0",
		Nonce:        0,
		Validator:    "system",
	}
	
	bc.Blocks = append(bc.Blocks, genesisBlock)
	return bc
}

// AddBlock adds a new block to the blockchain
func (bc *Blockchain) AddBlock(block *Block) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()

	if block == nil {
		return errors.New("block is nil")
	}

	// Validate block
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	
	if block.Index != lastBlock.Index+1 {
		return errors.New("invalid block index")
	}
	
	if block.PrevHash != lastBlock.Hash {
		return errors.New("invalid previous hash")
	}

	bc.Blocks = append(bc.Blocks, block)
	return nil
}

// GetLatestBlock returns the latest block
func (bc *Blockchain) GetLatestBlock() *Block {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if len(bc.Blocks) == 0 {
		return nil
	}
	
	return bc.Blocks[len(bc.Blocks)-1]
}

// GetBlock returns a block by index
func (bc *Blockchain) GetBlock(index int) (*Block, error) {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	if index < 0 || index >= len(bc.Blocks) {
		return nil, errors.New("block not found")
	}

	return bc.Blocks[index], nil
}

// GetBlockCount returns the number of blocks
func (bc *Blockchain) GetBlockCount() int {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return len(bc.Blocks)
}

// ValidateChain validates the entire blockchain
func (bc *Blockchain) ValidateChain() error {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		if currentBlock.PrevHash != prevBlock.Hash {
			return errors.New("invalid chain: hash mismatch")
		}

		if currentBlock.Index != prevBlock.Index+1 {
			return errors.New("invalid chain: index mismatch")
		}
	}

	return nil
}

// GetBlockchainInfo returns information about the blockchain
func (bc *Blockchain) GetBlockchainInfo() map[string]interface{} {
	bc.mutex.RLock()
	defer bc.mutex.RUnlock()

	return map[string]interface{}{
		"blocks":     len(bc.Blocks),
		"difficulty": bc.Difficulty,
		"latestHash": bc.Blocks[len(bc.Blocks)-1].Hash,
	}
}
