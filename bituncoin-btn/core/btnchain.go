package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// Block represents a block in the blockchain
type Block struct {
	Index        int64          `json:"index"`
	Timestamp    int64          `json:"timestamp"`
	Transactions []Transaction  `json:"transactions"`
	PrevHash     string         `json:"prev_hash"`
	Hash         string         `json:"hash"`
	Nonce        int            `json:"nonce"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	ID            string  `json:"id"`
	From          string  `json:"from"`
	To            string  `json:"to"`
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	Timestamp     int64   `json:"timestamp"`
	Signature     string  `json:"signature"`
	CrossChain    bool    `json:"cross_chain"`
	TargetChain   string  `json:"target_chain,omitempty"`
}

// Blockchain represents the main blockchain structure
type Blockchain struct {
	Blocks          []Block
	PendingTxs      []Transaction
	Difficulty      int
	MiningReward    float64
	SupportedChains []string
}

// NewBlockchain creates a new blockchain with genesis block
func NewBlockchain() *Blockchain {
	bc := &Blockchain{
		Blocks:       []Block{},
		PendingTxs:   []Transaction{},
		Difficulty:   4,
		MiningReward: 50.0,
		SupportedChains: []string{"BTN", "BTC", "ETH", "USDT", "BNB"},
	}
	bc.createGenesisBlock()
	return bc
}

// createGenesisBlock creates the first block in the chain
func (bc *Blockchain) createGenesisBlock() {
	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Transactions: []Transaction{},
		PrevHash:     "0",
		Nonce:        0,
	}
	genesis.Hash = bc.calculateHash(genesis)
	bc.Blocks = append(bc.Blocks, genesis)
}

// calculateHash calculates the hash of a block
func (bc *Blockchain) calculateHash(block Block) string {
	record := fmt.Sprintf("%d%d%s%d", block.Index, block.Timestamp, block.PrevHash, block.Nonce)
	for _, tx := range block.Transactions {
		txData, _ := json.Marshal(tx)
		record += string(txData)
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// AddTransaction adds a new transaction to pending transactions
func (bc *Blockchain) AddTransaction(tx Transaction) {
	bc.PendingTxs = append(bc.PendingTxs, tx)
}

// MineBlock mines a new block with pending transactions
func (bc *Blockchain) MineBlock(minerAddress string) Block {
	newBlock := Block{
		Index:        int64(len(bc.Blocks)),
		Timestamp:    time.Now().Unix(),
		Transactions: bc.PendingTxs,
		PrevHash:     bc.Blocks[len(bc.Blocks)-1].Hash,
		Nonce:        0,
	}

	// Proof of work
	for {
		newBlock.Hash = bc.calculateHash(newBlock)
		if bc.isValidHash(newBlock.Hash) {
			break
		}
		newBlock.Nonce++
	}

	bc.Blocks = append(bc.Blocks, newBlock)
	
	// Reset pending transactions and add mining reward
	bc.PendingTxs = []Transaction{
		{
			ID:        generateTxID(),
			From:      "SYSTEM",
			To:        minerAddress,
			Amount:    bc.MiningReward,
			Currency:  "BTN",
			Timestamp: time.Now().Unix(),
		},
	}

	return newBlock
}

// isValidHash checks if hash meets difficulty requirement
func (bc *Blockchain) isValidHash(hash string) bool {
	prefix := ""
	for i := 0; i < bc.Difficulty; i++ {
		prefix += "0"
	}
	return len(hash) >= bc.Difficulty && hash[:bc.Difficulty] == prefix
}

// GetBalance returns the balance for an address in a specific currency
func (bc *Blockchain) GetBalance(address string, currency string) float64 {
	balance := 0.0
	for _, block := range bc.Blocks {
		for _, tx := range block.Transactions {
			if tx.Currency == currency {
				if tx.From == address {
					balance -= tx.Amount
				}
				if tx.To == address {
					balance += tx.Amount
				}
			}
		}
	}
	return balance
}

// IsChainValid validates the entire blockchain
func (bc *Blockchain) IsChainValid() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		if currentBlock.Hash != bc.calculateHash(currentBlock) {
			return false
		}

		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}
	return true
}

// generateTxID generates a unique transaction ID
func generateTxID() string {
	return hex.EncodeToString([]byte(time.Now().String()))
}
