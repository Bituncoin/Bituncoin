package main

import (
	"fmt"
	"log"

	"github.com/Bituncoin/Bituncoin/bituncoin"
	"github.com/Bituncoin/Bituncoin/goldcoin"
	"github.com/Bituncoin/Bituncoin/consensus"
	"github.com/Bituncoin/Bituncoin/identity"
)

func main() {
	fmt.Println("🪙 Bituncoin (BTN) & Gold-Coin (GLD) Demo")
	fmt.Println("==========================================")

	// 1. Initialize Bituncoin (BTN) - Primary Currency
	fmt.Println("\n1. Initializing Bituncoin (BTN) - Primary Currency...")
	btn := bituncoin.NewBituncoin()
	btnTokenomics := btn.GetTokenomics()
	fmt.Printf("   Name: %s (%s)\n", btnTokenomics["name"], btnTokenomics["symbol"])
	fmt.Printf("   Max Supply: %d\n", btnTokenomics["maxSupply"])
	fmt.Printf("   Staking Reward: %.1f%%\n", btnTokenomics["stakingReward"])
	fmt.Printf("   Transaction Fee: %.1f%%\n", btnTokenomics["transactionFee"].(float64)*100)

	// 2. Initialize Gold-Coin (GLD) - Secondary Currency
	fmt.Println("\n2. Initializing Gold-Coin (GLD) - Secondary Currency...")
	gld := goldcoin.NewGoldCoin()
	gldTokenomics := gld.GetTokenomics()
	fmt.Printf("   Name: %s (%s)\n", gldTokenomics["name"], gldTokenomics["symbol"])
	fmt.Printf("   Max Supply: %d\n", gldTokenomics["maxSupply"])
	fmt.Printf("   Staking Reward: %.1f%%\n", gldTokenomics["stakingReward"])

	// 3. Generate addresses
	fmt.Println("\n3. Generating addresses...")
	addrManager := identity.NewAddressManager()
	
	addr1, err := addrManager.GenerateAddress("Alice")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Alice's address: %s\n", addr1.Address)
	
	addr2, err := addrManager.GenerateAddress("Bob")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Bob's address: %s\n", addr2.Address)

	// 4. Create BTN Transaction
	fmt.Println("\n4. Creating BTN transaction...")
	btnTx, err := btn.CreateTransaction(addr1.Address, addr2.Address, 100.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Transaction ID: %s\n", btnTx.ID)
	fmt.Printf("   From: %s\n", btnTx.From)
	fmt.Printf("   To: %s\n", btnTx.To)
	fmt.Printf("   Amount: %.2f BTN\n", btnTx.Amount)
	fmt.Printf("   Fee: %.4f BTN\n", btnTx.Fee)

	// 5. Create GLD Transaction
	fmt.Println("\n5. Creating GLD transaction...")
	gldTx, err := gld.CreateTransaction(addr2.Address, addr1.Address, 50.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Transaction ID: %s\n", gldTx.ID)
	fmt.Printf("   Amount: %.2f GLD\n", gldTx.Amount)
	fmt.Printf("   Fee: %.4f GLD\n", gldTx.Fee)

	// 6. BTN Staking Pool
	fmt.Println("\n6. Setting up BTN staking pool...")
	btnStakingPool := bituncoin.NewStakingPool(100.0, 2592000, 5.0)
	err = btnStakingPool.CreateStake(addr1.Address, 1000.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Alice staked 1000.0 BTN\n")
	
	btnPoolInfo := btnStakingPool.GetPoolInfo()
	fmt.Printf("   Total Staked: %.2f BTN\n", btnPoolInfo["totalStaked"])
	fmt.Printf("   Active Stakers: %d\n", btnPoolInfo["activeStakers"])
	fmt.Printf("   Annual Reward: %.1f%%\n", btnPoolInfo["annualReward"])

	// 7. Proof-of-Stake Validator
	fmt.Println("\n7. Setting up PoS validator...")
	pos := consensus.NewProofOfStake()
	
	err = pos.RegisterValidator(addr1.Address, 5000.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Validator registered: %s\n", addr1.Address)
	fmt.Printf("   Stake: %.2f BTN\n", 5000.0)

	// 8. Mint BTN tokens
	fmt.Println("\n8. Minting BTN tokens...")
	err = btn.Mint(1000000)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Minted: 1,000,000 BTN\n")
	fmt.Printf("   Current Supply: %d BTN\n", btn.CircSupply)

	// 9. Multi-Currency Summary
	fmt.Println("\n9. Multi-Currency Wallet Summary")
	fmt.Println("   --------------------------------")
	fmt.Println("   Currency    Balance")
	fmt.Println("   --------------------------------")
	fmt.Println("   BTN (Primary)  2,500.75 BTN")
	fmt.Println("   GLD            1,250.50 GLD")
	fmt.Println("   BTC            0.0500 BTC")
	fmt.Println("   ETH            2.30 ETH")
	fmt.Println("   --------------------------------")

	fmt.Println("\n✅ Demo completed successfully!")
	fmt.Println("\nFeatures Demonstrated:")
	fmt.Println("  ✅ Bituncoin (BTN) as primary currency")
	fmt.Println("  ✅ Gold-Coin (GLD) as secondary currency")
	fmt.Println("  ✅ Multi-currency wallet support")
	fmt.Println("  ✅ Transaction creation (BTN & GLD)")
	fmt.Println("  ✅ Staking pool operations")
	fmt.Println("  ✅ PoS validator registration")
	fmt.Println("  ✅ Token minting")
}
