package main

import (
	"fmt"
	"log"

	"github.com/Bituncoin/Bituncoin/identity"
)

func main() {
	fmt.Println("🔐 Bituncoin Address Generation Demo")
	fmt.Println("====================================")
	fmt.Println()

	// Initialize address manager
	manager := identity.NewAddressManager()

	// 1. Generate Bitcoin-Style Addresses
	fmt.Println("1. Bitcoin-Style Addresses (Btu prefix + Base58Check)")
	fmt.Println("   Features: Checksum validation, error detection")
	fmt.Println()

	for i := 1; i <= 3; i++ {
		label := fmt.Sprintf("BTC Wallet %d", i)
		addr, err := manager.GenerateAddressOfType(label, identity.BitcoinStyle)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("   %s:\n", label)
		fmt.Printf("   Address: %s\n", addr.Address)
		fmt.Printf("   Valid:   %v\n", identity.ValidateBitcoinStyleAddress(addr.Address) == nil)
		fmt.Println()
	}

	// 2. Generate Ethereum-Style Addresses
	fmt.Println("2. Ethereum-Style Addresses (0x + 40 hex chars)")
	fmt.Println("   Features: EVM compatible, smart contract support")
	fmt.Println()

	for i := 1; i <= 3; i++ {
		label := fmt.Sprintf("ETH Wallet %d", i)
		addr, err := manager.GenerateAddressOfType(label, identity.EthereumStyle)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("   %s:\n", label)
		fmt.Printf("   Address: %s\n", addr.Address)
		fmt.Printf("   Valid:   %v\n", identity.ValidateEthereumStyleAddress(addr.Address) == nil)
		fmt.Println()
	}

	// 3. Generate Legacy GLD Addresses
	fmt.Println("3. GLD Addresses (Legacy format)")
	fmt.Println("   Features: Gold-Coin native format")
	fmt.Println()

	gldAddr, err := manager.GenerateAddress("GLD Main Wallet")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   GLD Main Wallet:\n")
	fmt.Printf("   Address: %s\n", gldAddr.Address)
	fmt.Printf("   Valid:   %v\n", identity.ValidateAddress(gldAddr.Address) == nil)
	fmt.Println()

	// 4. Demonstrate Address Validation
	fmt.Println("4. Address Validation Examples")
	fmt.Println()

	// Use a freshly generated address for the Bitcoin example
	validBtcAddr, _, _ := identity.GenerateAddressWithType(identity.BitcoinStyle)

	testAddresses := []struct {
		address string
		name    string
	}{
		{validBtcAddr, "Valid Bitcoin-style"},
		{"0x1234567890abcdef1234567890abcdef12345678", "Valid Ethereum-style"},
		{"Btu123", "Invalid (too short)"},
		{"0x123", "Invalid (too short)"},
		{"invalid", "Invalid (unknown format)"},
	}

	for _, test := range testAddresses {
		addrType, err := identity.ValidateAddressType(test.address)
		if err != nil {
			fmt.Printf("   ❌ %s: %s\n", test.name, err)
		} else {
			fmt.Printf("   ✅ %s: Type = %s\n", test.name, addrType)
		}
	}
	fmt.Println()

	// 5. Demonstrate Direct Address Generation
	fmt.Println("5. Direct Address Generation (without manager)")
	fmt.Println()

	btcAddr, privKey1, err := identity.GenerateAddressWithType(identity.BitcoinStyle)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Bitcoin-style Address: %s\n", btcAddr)
	fmt.Printf("   Private Key Length:    %d bytes\n", len(privKey1))
	fmt.Println()

	ethAddr, privKey2, err := identity.GenerateAddressWithType(identity.EthereumStyle)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Ethereum-style Address: %s\n", ethAddr)
	fmt.Printf("   Private Key Length:     %d bytes\n", len(privKey2))
	fmt.Println()

	// 6. List all generated addresses
	fmt.Println("6. Address Summary")
	fmt.Println()

	allAddresses := manager.ListAddresses()
	fmt.Printf("   Total addresses generated: %d\n", len(allAddresses))
	fmt.Println()

	// Group by type
	var btcCount, ethCount, gldCount int
	for _, addr := range allAddresses {
		switch addr.AddressType {
		case identity.BitcoinStyle:
			btcCount++
		case identity.EthereumStyle:
			ethCount++
		case "gld":
			gldCount++
		}
	}

	fmt.Printf("   Bitcoin-style addresses:  %d\n", btcCount)
	fmt.Printf("   Ethereum-style addresses: %d\n", ethCount)
	fmt.Printf("   GLD addresses (legacy):   %d\n", gldCount)
	fmt.Println()

	// Summary
	fmt.Println("✅ Address Generation Demo Complete!")
	fmt.Println()
	fmt.Println("Key Features Demonstrated:")
	fmt.Println("  ✓ Bitcoin-style addresses with Base58Check encoding")
	fmt.Println("  ✓ Ethereum-style addresses with hex encoding")
	fmt.Println("  ✓ Automatic address validation and type detection")
	fmt.Println("  ✓ Checksum verification for Bitcoin-style addresses")
	fmt.Println("  ✓ Legacy GLD address support")
	fmt.Println("  ✓ Secure private key generation")
	fmt.Println()
	fmt.Println("📚 See docs/ADDRESS_FORMATS.md for detailed documentation")
	fmt.Println("🚀 Ready for universal wallet integration!")
}
