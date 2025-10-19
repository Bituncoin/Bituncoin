package main

import (
	"fmt"
	"github.com/Bituncoin/Bituncoin/identity"
)

func main() {
	fmt.Println("==============================================")
	fmt.Println("   Address Manager Integration Demo")
	fmt.Println("==============================================\n")

	// Create a new address manager
	manager := identity.NewAddressManager()

	// Generate different types of addresses
	fmt.Println("📍 Generating Addresses via AddressManager")
	fmt.Println("--------------------------------------------")

	// Generate legacy GLD address
	gldAddr, err := manager.GenerateAddress("My GLD Wallet")
	if err != nil {
		fmt.Printf("Error generating GLD address: %v\n", err)
		return
	}
	fmt.Printf("GLD Address:  %s (Label: %s)\n", gldAddr.Address, gldAddr.Label)

	// Generate Bitcoin-style address
	btcAddr, err := manager.GenerateBitcoinAddress("My Bitcoin-style Wallet")
	if err != nil {
		fmt.Printf("Error generating Bitcoin-style address: %v\n", err)
		return
	}
	fmt.Printf("Btu Address:  %s (Label: %s)\n", btcAddr.Address, btcAddr.Label)

	// Generate Ethereum-style address
	ethAddr, err := manager.GenerateEthereumAddress("My Ethereum-style Wallet")
	if err != nil {
		fmt.Printf("Error generating Ethereum-style address: %v\n", err)
		return
	}
	fmt.Printf("0x Address:   %s (Label: %s)\n", ethAddr.Address, ethAddr.Label)

	fmt.Println()

	// List all addresses
	fmt.Println("📍 Listing All Managed Addresses")
	fmt.Println("--------------------------------------------")
	allAddresses := manager.ListAddresses()
	for i, addr := range allAddresses {
		fmt.Printf("%d. %s (%s)\n", i+1, addr.Address, addr.Label)
	}

	fmt.Println()

	// Retrieve specific addresses
	fmt.Println("📍 Retrieving Specific Addresses")
	fmt.Println("--------------------------------------------")
	
	retrievedBtc, err := manager.GetAddress(btcAddr.Address)
	if err != nil {
		fmt.Printf("Error retrieving address: %v\n", err)
	} else {
		fmt.Printf("Retrieved: %s\n", retrievedBtc.Address)
		fmt.Printf("  Public Key: %s...\n", retrievedBtc.PublicKey[:32])
	}

	fmt.Println()

	// Validate addresses
	fmt.Println("📍 Universal Address Validation")
	fmt.Println("--------------------------------------------")
	
	testAddresses := []string{
		gldAddr.Address,
		btcAddr.Address,
		ethAddr.Address,
		"InvalidAddress",
	}

	for _, addr := range testAddresses {
		err := identity.ValidateAddress(addr)
		if err != nil {
			fmt.Printf("❌ %s - %v\n", addr, err)
		} else {
			fmt.Printf("✅ %s - valid\n", addr)
		}
	}

	fmt.Println()

	// Sign and verify message
	fmt.Println("📍 Message Signing (Bitcoin-style address)")
	fmt.Println("--------------------------------------------")
	message := "Hello, Bituncoin!"
	signature, err := manager.SignMessage(btcAddr.Address, message)
	if err != nil {
		fmt.Printf("Error signing message: %v\n", err)
	} else {
		fmt.Printf("Message: %s\n", message)
		fmt.Printf("Signature: %s...\n", signature[:32])
		
		// Verify signature
		valid := identity.VerifySignature(btcAddr.Address, message, signature)
		if valid {
			fmt.Println("✅ Signature verified!")
		} else {
			fmt.Println("❌ Signature verification failed!")
		}
	}

	fmt.Println()

	// Generate multiple addresses of each type
	fmt.Println("📍 Batch Address Generation")
	fmt.Println("--------------------------------------------")
	
	fmt.Println("Generating 3 Bitcoin-style addresses:")
	for i := 0; i < 3; i++ {
		addr, err := manager.GenerateBitcoinAddress(fmt.Sprintf("Btu Wallet %d", i+1))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("  %d. %s\n", i+1, addr.Address)
	}

	fmt.Println("\nGenerating 3 Ethereum-style addresses:")
	for i := 0; i < 3; i++ {
		addr, err := manager.GenerateEthereumAddress(fmt.Sprintf("ETH Wallet %d", i+1))
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("  %d. %s\n", i+1, addr.Address)
	}

	fmt.Println()

	// Final summary
	fmt.Println("📍 Summary")
	fmt.Println("--------------------------------------------")
	totalAddresses := len(manager.ListAddresses())
	fmt.Printf("Total addresses managed: %d\n", totalAddresses)
	
	// Count by type
	gldCount := 0
	btuCount := 0
	ethCount := 0
	
	for _, addr := range manager.ListAddresses() {
		if len(addr.Address) >= 3 {
			if addr.Address[:3] == "GLD" {
				gldCount++
			} else if addr.Address[:3] == "Btu" {
				btuCount++
			} else if len(addr.Address) >= 2 && addr.Address[:2] == "0x" {
				ethCount++
			}
		}
	}
	
	fmt.Printf("  - GLD addresses: %d\n", gldCount)
	fmt.Printf("  - Btu addresses: %d\n", btuCount)
	fmt.Printf("  - 0x addresses:  %d\n", ethCount)

	fmt.Println("\n==============================================")
	fmt.Println("   Integration Demo Complete!")
	fmt.Println("==============================================")
}
