package main

import (
	"fmt"
	"log"

	"github.com/Bituncoin/Bituncoin/goldcoin"
	"github.com/Bituncoin/Bituncoin/consensus"
	"github.com/Bituncoin/Bituncoin/identity"
)

func main() {
	fmt.Println("🪙 Gold-Coin Cryptocurrency Demo")
	fmt.Println("================================")

	// 1. Initialize Gold-Coin
	fmt.Println("1. Initializing Gold-Coin...")
	gc := goldcoin.NewGoldCoin()
	tokenomics := gc.GetTokenomics()
	fmt.Printf("   Name: %s (%s)\n", tokenomics["name"], tokenomics["symbol"])
	fmt.Printf("   Max Supply: %d\n", tokenomics["maxSupply"])
	fmt.Printf("   Staking Reward: %.1f%%\n", tokenomics["stakingReward"])
	fmt.Printf("   Transaction Fee: %.1f%%\n\n", tokenomics["transactionFee"].(float64)*100)

	// 2. Generate addresses
	fmt.Println("2. Generating addresses...")
	addrManager := identity.NewAddressManager()
	
	addr1, err := addrManager.GenerateAddress("Alice")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Alice's address: %s\n", addr1.Address[:20]+"...")

	addr2, err := addrManager.GenerateAddress("Bob")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Bob's address: %s\n\n", addr2.Address[:20]+"...")

	// 3. Create a transaction
	fmt.Println("3. Creating a transaction...")
	tx, err := gc.CreateTransaction(addr1.Address, addr2.Address, 100.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Transaction ID: %s\n", tx.ID[:16]+"...")
	fmt.Printf("   From: %s\n", tx.From[:20]+"...")
	fmt.Printf("   To: %s\n", tx.To[:20]+"...")
	fmt.Printf("   Amount: %.2f GLD\n", tx.Amount)
	fmt.Printf("   Fee: %.4f GLD\n\n", tx.Fee)

	// 4. Initialize Proof-of-Stake
	fmt.Println("4. Initializing Proof-of-Stake consensus...")
	pos := consensus.NewProofOfStake()
	fmt.Printf("   Min Validator Stake: %.0f GLD\n", pos.MinStake)
	fmt.Printf("   Block Time: %d seconds\n", pos.BlockTime)
	fmt.Printf("   Reward per Block: %.1f GLD\n\n", pos.RewardPerBlock)

	// 5. Register validators
	fmt.Println("5. Registering validators...")
	err = pos.RegisterValidator(addr1.Address, 2000.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Validator 1 registered: %s (2000 GLD)\n", addr1.Address[:20]+"...")

	err = pos.RegisterValidator(addr2.Address, 3000.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Validator 2 registered: %s (3000 GLD)\n\n", addr2.Address[:20]+"...")

	// 6. Create a block
	fmt.Println("6. Creating a block...")
	transactions := []string{tx.ID}
	block, err := pos.CreateBlock(transactions, "0000000000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Block #%d created\n", block.Index)
	fmt.Printf("   Validator: %s\n", block.Validator[:20]+"...")
	fmt.Printf("   Transactions: %d\n", len(block.Transactions))
	fmt.Printf("   Hash: %s\n\n", block.Hash[:16]+"...")

	// 7. Initialize staking pool
	fmt.Println("7. Creating staking pool...")
	stakingPool := goldcoin.NewStakingPool()
	poolInfo := stakingPool.GetPoolInfo()
	fmt.Printf("   Annual Reward: %.1f%%\n", poolInfo["annualReward"])
	fmt.Printf("   Min Stake: %.0f GLD\n", poolInfo["minStake"])
	fmt.Printf("   Lock Period: %d days\n\n", poolInfo["lockPeriod"].(int64)/(24*60*60))

	// 8. Create a stake
	fmt.Println("8. Creating a stake...")
	err = stakingPool.CreateStake(addr1.Address, 1000.0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Staked 1000 GLD from %s\n", addr1.Address[:20]+"...")
	
	stakeInfo, _ := stakingPool.GetStakeInfo(addr1.Address)
	fmt.Printf("   Active: %v\n", stakeInfo.IsActive)
	fmt.Printf("   Amount: %.2f GLD\n\n", stakeInfo.Amount)

	// 9. Calculate rewards
	fmt.Println("9. Calculating staking rewards...")
	rewards, err := stakingPool.CalculateRewards(addr1.Address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("   Current rewards: %.6f GLD\n\n", rewards)

	// Summary
	fmt.Println("✅ Gold-Coin Demo Complete!")
	fmt.Println("\nFeatures demonstrated:")
	fmt.Println("  ✓ Token initialization with tokenomics")
	fmt.Println("  ✓ Address generation and management")
	fmt.Println("  ✓ Transaction creation and validation")
	fmt.Println("  ✓ Proof-of-Stake consensus mechanism")
	fmt.Println("  ✓ Validator registration and block creation")
	fmt.Println("  ✓ Staking pool with rewards calculation")
	fmt.Println("\n🚀 Ready for deployment!")
}
