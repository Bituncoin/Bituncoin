package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Bituncoin/Bituncoin/bituncoin-btn/api"
	"github.com/Bituncoin/Bituncoin/bituncoin-btn/consensus"
)

func main() {
	fmt.Println("🚀 Starting Bituncoin Universal Wallet Node...")
	fmt.Println("=" + string(make([]byte, 50)) + "=")

	// Get home directory for data storage
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get home directory:", err)
	}

	dataDir := filepath.Join(homeDir, ".bituncoin", "data")

	// Create data directory
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// Initialize consensus mechanism
	validators := consensus.NewValidatorEcho(3)
	fmt.Println("✓ Consensus mechanism initialized (Validator Echo)")

	// Register initial validators
	validators.RegisterValidator("BTN0000000000000000000000000000000000001", 5000.0)
	validators.RegisterValidator("BTN0000000000000000000000000000000000002", 3000.0)
	validators.RegisterValidator("BTN0000000000000000000000000000000000003", 2000.0)
	fmt.Println("✓ Initial validators registered")

	// Create and start node
	node, err := api.NewNode(8080, dataDir)
	if err != nil {
		log.Fatal("Failed to create node:", err)
	}

	fmt.Println("✓ Node initialized")
	fmt.Println("\n📦 Universal Wallet Features:")
	fmt.Println("  • Multi-Currency Support (BTN, BTC, ETH, USDT, BNB)")
	fmt.Println("  • Cross-Chain Transactions")
	fmt.Println("  • Two-Factor Authentication (2FA)")
	fmt.Println("  • Biometric Authentication")
	fmt.Println("  • Encrypted Storage")
	fmt.Println("  • Secure Identity Management")
	fmt.Println("\n🌐 API Endpoints:")
	fmt.Println("  POST   /api/wallet/create         - Create new wallet")
	fmt.Println("  GET    /api/wallet/balance        - Get wallet balance")
	fmt.Println("  POST   /api/transaction/send      - Send transaction")
	fmt.Println("  GET    /api/transaction/history   - Get transaction history")
	fmt.Println("  GET    /api/blockchain/info       - Get blockchain info")
	fmt.Println("  POST   /api/mine                  - Mine a new block")
	fmt.Println("  GET    /api/currencies            - Get supported currencies")
	fmt.Println("  POST   /api/crosschain/bridge     - Cross-chain bridge")
	fmt.Println("\n🔗 Starting HTTP server on port 8080...")
	fmt.Println("=" + string(make([]byte, 50)) + "=")

	// Start the node
	if err := node.Start(); err != nil {
		log.Fatal("Failed to start node:", err)
	}
}
