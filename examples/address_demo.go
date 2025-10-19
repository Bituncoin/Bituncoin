package main

import (
	"fmt"
	"github.com/Bituncoin/Bituncoin/identity"
)

func main() {
	fmt.Println("==============================================")
	fmt.Println("   Bituncoin Address Generation Demo")
	fmt.Println("==============================================\n")

	// Generate Bitcoin-Style Address
	fmt.Println("📍 Bitcoin-Style Address Generation")
	fmt.Println("--------------------------------------------")
	btcAddr, err := identity.GenerateBitcoinStyleAddress()
	if err != nil {
		fmt.Printf("Error generating Bitcoin-style address: %v\n", err)
		return
	}

	fmt.Printf("Address:     %s\n", btcAddr.Address)
	fmt.Printf("Public Key:  %s\n", btcAddr.PublicKey[:32]+"...")
	fmt.Printf("Private Key: %s\n", btcAddr.PrivateKey[:32]+"...")
	
	// Validate Bitcoin-Style Address
	if err := identity.ValidateBitcoinStyleAddress(btcAddr.Address); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Address is valid!")
	}

	fmt.Println()

	// Generate Ethereum-Style Address
	fmt.Println("📍 Ethereum-Style Address Generation")
	fmt.Println("--------------------------------------------")
	ethAddr, err := identity.GenerateEthereumStyleAddress()
	if err != nil {
		fmt.Printf("Error generating Ethereum-style address: %v\n", err)
		return
	}

	fmt.Printf("Address:     %s\n", ethAddr.Address)
	fmt.Printf("Public Key:  %s\n", ethAddr.PublicKey[:32]+"...")
	fmt.Printf("Private Key: %s\n", ethAddr.PrivateKey[:32]+"...")
	
	// Validate Ethereum-Style Address
	if err := identity.ValidateEthereumStyleAddress(ethAddr.Address); err != nil {
		fmt.Printf("❌ Validation failed: %v\n", err)
	} else {
		fmt.Println("✅ Address is valid!")
	}

	fmt.Println()

	// Generate Multiple Addresses
	fmt.Println("📍 Multiple Address Generation")
	fmt.Println("--------------------------------------------")
	
	fmt.Println("Bitcoin-style addresses:")
	for i := 0; i < 3; i++ {
		addr, err := identity.GenerateBitcoinStyleAddress()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("  %d. %s\n", i+1, addr.Address)
	}

	fmt.Println("\nEthereum-style addresses:")
	for i := 0; i < 3; i++ {
		addr, err := identity.GenerateEthereumStyleAddress()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("  %d. %s\n", i+1, addr.Address)
	}

	fmt.Println()

	// Test Validation with Invalid Addresses
	fmt.Println("📍 Validation Testing")
	fmt.Println("--------------------------------------------")
	
	invalidBtcAddresses := []string{
		"1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", // Missing Btu prefix
		"Btu123",                                // Too short
		"",                                      // Empty
	}

	fmt.Println("Testing invalid Bitcoin-style addresses:")
	for _, addr := range invalidBtcAddresses {
		err := identity.ValidateBitcoinStyleAddress(addr)
		if err != nil {
			fmt.Printf("  ❌ '%s' - %v\n", addr, err)
		} else {
			fmt.Printf("  ✅ '%s' - valid\n", addr)
		}
	}

	fmt.Println("\nTesting invalid Ethereum-style addresses:")
	invalidEthAddresses := []string{
		"1234567890abcdef1234567890abcdef12345678", // Missing 0x
		"0x123",                                     // Too short
		"0x1234567890abcdef1234567890abcdefghijklmn", // Invalid hex
		"",                                          // Empty
	}

	for _, addr := range invalidEthAddresses {
		err := identity.ValidateEthereumStyleAddress(addr)
		if err != nil {
			fmt.Printf("  ❌ '%s' - %v\n", addr, err)
		} else {
			fmt.Printf("  ✅ '%s' - valid\n", addr)
		}
	}

	fmt.Println("\n==============================================")
	fmt.Println("   Demo Complete!")
	fmt.Println("==============================================")
}
