package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// BTNGWalletContract implements on-chain mint/melt for BTNG Gold
type BTNGWalletContract struct {
	contractapi.Contract
}

// WalletState represents a wallet's on-chain state
type WalletState struct {
	WalletID  string  `json:"wallet_id"`
	Balance   float64 `json:"balance"`
	UpdatedAt string  `json:"updated_at"`
}

// Transaction represents a mint/melt transaction record
type Transaction struct {
	TxID         string  `json:"tx_id"`
	Type         string  `json:"type"` // "mint" or "melt"
	WalletID     string  `json:"wallet_id"`
	Amount       float64 `json:"amount"`
	BalanceAfter float64 `json:"balance_after"`
	Timestamp    string  `json:"timestamp"`
	GoldPriceUSD float64 `json:"gold_price_usd,omitempty"`
}

// Mint creates new BTNG tokens and adds them to a wallet
func (c *BTNGWalletContract) Mint(ctx contractapi.TransactionContextInterface, walletID string, amountStr string, goldPriceStr string) (*Transaction, error) {
	if walletID == "" {
		return nil, fmt.Errorf("wallet ID is required")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		return nil, fmt.Errorf("mint amount must be a positive number")
	}

	goldPrice, _ := strconv.ParseFloat(goldPriceStr, 64)

	// Get or create wallet state
	wallet, err := c.getWalletState(ctx, walletID)
	if err != nil {
		wallet = &WalletState{WalletID: walletID, Balance: 0}
	}

	wallet.Balance += amount
	wallet.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	// Save wallet state
	if err := c.putWalletState(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to save wallet state: %v", err)
	}

	// Create and store transaction record
	tx := &Transaction{
		TxID:         ctx.GetStub().GetTxID(),
		Type:         "mint",
		WalletID:     walletID,
		Amount:       amount,
		BalanceAfter: wallet.Balance,
		Timestamp:    wallet.UpdatedAt,
		GoldPriceUSD: goldPrice,
	}

	txKey := fmt.Sprintf("TX_%s", tx.TxID)
	txBytes, _ := json.Marshal(tx)
	if err := ctx.GetStub().PutState(txKey, txBytes); err != nil {
		return nil, fmt.Errorf("failed to record transaction: %v", err)
	}

	// Update total supply
	if err := c.updateTotalSupply(ctx, amount); err != nil {
		return nil, fmt.Errorf("failed to update total supply: %v", err)
	}

	return tx, nil
}

// Melt destroys BTNG tokens from a wallet (gold redemption)
func (c *BTNGWalletContract) Melt(ctx contractapi.TransactionContextInterface, walletID string, amountStr string, goldPriceStr string) (*Transaction, error) {
	if walletID == "" {
		return nil, fmt.Errorf("wallet ID is required")
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		return nil, fmt.Errorf("melt amount must be a positive number")
	}

	goldPrice, _ := strconv.ParseFloat(goldPriceStr, 64)

	wallet, err := c.getWalletState(ctx, walletID)
	if err != nil {
		return nil, fmt.Errorf("wallet %s not found", walletID)
	}

	if wallet.Balance < amount {
		return nil, fmt.Errorf("insufficient BTNG: wallet %s has %.4f, requested %.4f", walletID, wallet.Balance, amount)
	}

	wallet.Balance -= amount
	wallet.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := c.putWalletState(ctx, wallet); err != nil {
		return nil, fmt.Errorf("failed to save wallet state: %v", err)
	}

	tx := &Transaction{
		TxID:         ctx.GetStub().GetTxID(),
		Type:         "melt",
		WalletID:     walletID,
		Amount:       amount,
		BalanceAfter: wallet.Balance,
		Timestamp:    wallet.UpdatedAt,
		GoldPriceUSD: goldPrice,
	}

	txKey := fmt.Sprintf("TX_%s", tx.TxID)
	txBytes, _ := json.Marshal(tx)
	if err := ctx.GetStub().PutState(txKey, txBytes); err != nil {
		return nil, fmt.Errorf("failed to record transaction: %v", err)
	}

	if err := c.updateTotalSupply(ctx, -amount); err != nil {
		return nil, fmt.Errorf("failed to update total supply: %v", err)
	}

	return tx, nil
}

// Transfer moves BTNG tokens between wallets
func (c *BTNGWalletContract) Transfer(ctx contractapi.TransactionContextInterface, fromWallet string, toWallet string, amountStr string) (*Transaction, error) {
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		return nil, fmt.Errorf("transfer amount must be a positive number")
	}

	from, err := c.getWalletState(ctx, fromWallet)
	if err != nil {
		return nil, fmt.Errorf("source wallet %s not found", fromWallet)
	}
	if from.Balance < amount {
		return nil, fmt.Errorf("insufficient BTNG in wallet %s: has %.4f, requested %.4f", fromWallet, from.Balance, amount)
	}

	to, err := c.getWalletState(ctx, toWallet)
	if err != nil {
		to = &WalletState{WalletID: toWallet, Balance: 0}
	}

	now := time.Now().UTC().Format(time.RFC3339)
	from.Balance -= amount
	from.UpdatedAt = now
	to.Balance += amount
	to.UpdatedAt = now

	if err := c.putWalletState(ctx, from); err != nil {
		return nil, err
	}
	if err := c.putWalletState(ctx, to); err != nil {
		return nil, err
	}

	tx := &Transaction{
		TxID:         ctx.GetStub().GetTxID(),
		Type:         "transfer",
		WalletID:     fmt.Sprintf("%s->%s", fromWallet, toWallet),
		Amount:       amount,
		BalanceAfter: from.Balance,
		Timestamp:    now,
	}

	txKey := fmt.Sprintf("TX_%s", tx.TxID)
	txBytes, _ := json.Marshal(tx)
	ctx.GetStub().PutState(txKey, txBytes)

	return tx, nil
}

// GetBalance returns the BTNG balance for a wallet
func (c *BTNGWalletContract) GetBalance(ctx contractapi.TransactionContextInterface, walletID string) (*WalletState, error) {
	return c.getWalletState(ctx, walletID)
}

// GetTotalSupply returns the total BTNG supply across all wallets
func (c *BTNGWalletContract) GetTotalSupply(ctx contractapi.TransactionContextInterface) (string, error) {
	supplyBytes, err := ctx.GetStub().GetState("BTNG_TOTAL_SUPPLY")
	if err != nil {
		return "0", err
	}
	if supplyBytes == nil {
		return "0", nil
	}
	return string(supplyBytes), nil
}

// GetTransactionHistory returns all transactions for a wallet using composite key range query
func (c *BTNGWalletContract) GetTransactionHistory(ctx contractapi.TransactionContextInterface, walletID string) ([]*Transaction, error) {
	historyIterator, err := ctx.GetStub().GetHistoryForKey(fmt.Sprintf("WALLET_%s", walletID))
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %v", err)
	}
	defer historyIterator.Close()

	var transactions []*Transaction
	for historyIterator.HasNext() {
		response, err := historyIterator.Next()
		if err != nil {
			return nil, err
		}
		var wallet WalletState
		if err := json.Unmarshal(response.Value, &wallet); err != nil {
			continue
		}
		tx := &Transaction{
			TxID:         response.TxId,
			WalletID:     walletID,
			BalanceAfter: wallet.Balance,
			Timestamp:    time.Unix(response.Timestamp.Seconds, 0).UTC().Format(time.RFC3339),
		}
		transactions = append(transactions, tx)
	}
	return transactions, nil
}

// --- Internal helpers ---

func (c *BTNGWalletContract) getWalletState(ctx contractapi.TransactionContextInterface, walletID string) (*WalletState, error) {
	key := fmt.Sprintf("WALLET_%s", walletID)
	data, err := ctx.GetStub().GetState(key)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("wallet %s not found", walletID)
	}
	var wallet WalletState
	if err := json.Unmarshal(data, &wallet); err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (c *BTNGWalletContract) putWalletState(ctx contractapi.TransactionContextInterface, wallet *WalletState) error {
	key := fmt.Sprintf("WALLET_%s", wallet.WalletID)
	data, err := json.Marshal(wallet)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(key, data)
}

func (c *BTNGWalletContract) updateTotalSupply(ctx contractapi.TransactionContextInterface, delta float64) error {
	supplyBytes, _ := ctx.GetStub().GetState("BTNG_TOTAL_SUPPLY")
	currentSupply := 0.0
	if supplyBytes != nil {
		currentSupply, _ = strconv.ParseFloat(string(supplyBytes), 64)
	}
	newSupply := currentSupply + delta
	return ctx.GetStub().PutState("BTNG_TOTAL_SUPPLY", []byte(fmt.Sprintf("%.4f", newSupply)))
}

func main() {
	chaincode, err := contractapi.NewChaincode(&BTNGWalletContract{})
	if err != nil {
		log.Panicf("Error creating btng-wallet chaincode: %v", err)
	}

	if err := chaincode.Start(); err != nil {
		log.Panicf("Error starting btng-wallet chaincode: %v", err)
	}
}
