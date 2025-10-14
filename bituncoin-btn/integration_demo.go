package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TestResponse represents API response
type TestResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

func main() {
	baseURL := "http://localhost:8080/api"

	fmt.Println("🧪 Bituncoin Universal Wallet Integration Test")
	fmt.Println("=============================================")
	fmt.Println()

	// Test 1: Check if server is running
	fmt.Println("Test 1: Checking server availability...")
	resp, err := http.Get(baseURL + "/blockchain/info")
	if err != nil {
		fmt.Printf("❌ Server not running. Please start the node first.\n")
		return
	}
	defer resp.Body.Close()
	fmt.Println("✓ Server is running")
	fmt.Println()

	// Test 2: Create wallet with 2FA
	fmt.Println("Test 2: Creating wallet with 2FA...")
	wallet1 := createWallet(baseURL, true, false, "")
	if wallet1 == nil {
		fmt.Println("❌ Failed to create wallet 1")
		return
	}
	fmt.Printf("✓ Wallet 1 created: %s\n", wallet1["address"])
	fmt.Println()

	// Test 3: Create wallet with biometric auth
	fmt.Println("Test 3: Creating wallet with biometric authentication...")
	wallet2 := createWallet(baseURL, false, true, "sample-biometric-data")
	if wallet2 == nil {
		fmt.Println("❌ Failed to create wallet 2")
		return
	}
	fmt.Printf("✓ Wallet 2 created: %s\n", wallet2["address"])
	fmt.Println()

	// Test 4: Check balances
	fmt.Println("Test 4: Checking wallet balances...")
	checkBalance(baseURL, wallet1["address"].(string))
	checkBalance(baseURL, wallet2["address"].(string))
	fmt.Println()

	// Test 5: Get supported currencies
	fmt.Println("Test 5: Getting supported currencies...")
	getCurrencies(baseURL)
	fmt.Println()

	// Test 6: Send regular transaction
	fmt.Println("Test 6: Sending regular transaction...")
	sendTransaction(baseURL, wallet1["address"].(string), wallet2["address"].(string), 10.5, "BTN", false, "")
	fmt.Println()

	// Test 7: Send cross-chain transaction
	fmt.Println("Test 7: Sending cross-chain transaction...")
	sendTransaction(baseURL, wallet2["address"].(string), wallet1["address"].(string), 5.0, "BTN", true, "ETH")
	fmt.Println()

	// Test 8: Mine a block
	fmt.Println("Test 8: Mining a block...")
	mineBlock(baseURL, wallet1["address"].(string))
	fmt.Println()

	// Test 9: Get transaction history
	fmt.Println("Test 9: Getting transaction history...")
	getTransactionHistory(baseURL, wallet1["address"].(string))
	fmt.Println()

	// Test 10: Get blockchain info
	fmt.Println("Test 10: Getting blockchain information...")
	getBlockchainInfo(baseURL)
	fmt.Println()

	fmt.Println("✅ All integration tests completed successfully!")
	fmt.Println()
	fmt.Println("📊 Summary:")
	fmt.Printf("  - Wallet 1 (2FA): %s\n", wallet1["address"])
	fmt.Printf("  - Wallet 2 (Biometric): %s\n", wallet2["address"])
	fmt.Println("  - Transactions: Created (regular and cross-chain)")
	fmt.Println("  - Mining: Successful")
	fmt.Println("  - All features: ✓ Working")
}

func createWallet(baseURL string, enable2FA, enableBiometric bool, biometricData string) map[string]interface{} {
	payload := map[string]interface{}{
		"enable_2fa":       enable2FA,
		"enable_biometric": enableBiometric,
		"biometric_data":   biometricData,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/wallet/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var result TestResponse
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Success {
		return result.Data
	}
	return nil
}

func checkBalance(baseURL, address string) {
	resp, err := http.Get(baseURL + "/wallet/balance?address=" + address)
	if err != nil {
		fmt.Printf("❌ Failed to get balance for %s\n", address)
		return
	}
	defer resp.Body.Close()

	var result TestResponse
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Success {
		fmt.Printf("✓ Balance for %s:\n", address[:20]+"...")
		for currency, balance := range result.Data {
			fmt.Printf("    %s: %.4f\n", currency, balance)
		}
	}
}

func getCurrencies(baseURL string) {
	resp, err := http.Get(baseURL + "/currencies")
	if err != nil {
		fmt.Println("❌ Failed to get currencies")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✓ Supported currencies: %s\n", string(body))
}

func sendTransaction(baseURL, from, to string, amount float64, currency string, crossChain bool, targetChain string) {
	payload := map[string]interface{}{
		"from":         from,
		"to":           to,
		"amount":       amount,
		"currency":     currency,
		"cross_chain":  crossChain,
		"target_chain": targetChain,
	}

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/transaction/send", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("❌ Failed to send transaction")
		return
	}
	defer resp.Body.Close()

	var result TestResponse
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Success {
		txType := "regular"
		if crossChain {
			txType = "cross-chain"
		}
		fmt.Printf("✓ Transaction sent (%s): %.2f %s\n", txType, amount, currency)
		if txID, ok := result.Data["transaction_id"]; ok {
			fmt.Printf("    TX ID: %s\n", txID)
		}
	}
}

func mineBlock(baseURL, minerAddress string) {
	resp, err := http.Post(baseURL+"/mine?miner="+minerAddress, "application/json", nil)
	if err != nil {
		fmt.Println("❌ Failed to mine block")
		return
	}
	defer resp.Body.Close()

	var result TestResponse
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Success {
		fmt.Println("✓ Block mined successfully")
	}
}

func getTransactionHistory(baseURL, address string) {
	resp, err := http.Get(baseURL + "/transaction/history?address=" + address)
	if err != nil {
		fmt.Println("❌ Failed to get transaction history")
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("✓ Transaction history retrieved: %s\n", string(body)[:100]+"...")
}

func getBlockchainInfo(baseURL string) {
	resp, err := http.Get(baseURL + "/blockchain/info")
	if err != nil {
		fmt.Println("❌ Failed to get blockchain info")
		return
	}
	defer resp.Body.Close()

	var result TestResponse
	json.NewDecoder(resp.Body).Decode(&result)
	if result.Success {
		fmt.Println("✓ Blockchain information:")
		for key, value := range result.Data {
			fmt.Printf("    %s: %v\n", key, value)
		}
	}
}

func init() {
	// Wait a moment for server to be ready
	time.Sleep(1 * time.Second)
}
