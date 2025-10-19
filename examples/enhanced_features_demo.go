package main

import (
	"fmt"
	"log"

	"github.com/Bituncoin/Bituncoin/addons"
	"github.com/Bituncoin/Bituncoin/api"
	"github.com/Bituncoin/Bituncoin/auth"
)

// This example demonstrates the new features:
// 1. User authentication and role-based access control
// 2. Add-on module system
// 3. Enhanced API with authentication

func main() {
	fmt.Println("=== Bituncoin Enhanced Features Demo ===\n")

	// Demo 1: User Authentication and RBAC
	demoAuthentication()

	// Demo 2: Add-On Module System
	demoModuleSystem()

	// Demo 3: API Integration
	demoAPIIntegration()
}

func demoAuthentication() {
	fmt.Println("--- Demo 1: User Authentication & RBAC ---")

	// Create account manager
	am := auth.NewAccountManager()

	// Create regular user
	user, err := am.CreateUser("alice", "alice@example.com", "password123", auth.RoleUser)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created user: %s (Role: %s)\n", user.Username, user.Role)
	fmt.Printf("  Permissions: %v\n", user.Permissions)

	// Create admin user
	admin, err := am.CreateUser("admin", "admin@example.com", "adminpass", auth.RoleAdmin)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created admin: %s (Role: %s)\n", admin.Username, admin.Role)
	fmt.Printf("  Permissions: %v\n", admin.Permissions)

	// Authenticate user
	session, err := am.Authenticate("alice", "password123")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nAuthenticated user 'alice'\n")
	fmt.Printf("  Session ID: %s\n", session.ID)
	fmt.Printf("  Expires: %s\n", session.ExpiresAt.Format("2006-01-02 15:04:05"))

	// Validate session
	validatedUser, err := am.ValidateSession(session.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("  Validated session for: %s\n", validatedUser.Username)

	// Check permissions
	fmt.Println("\nPermission checks:")
	fmt.Printf("  Alice can read: %v\n", am.HasPermission(user.ID, auth.PermissionRead))
	fmt.Printf("  Alice can manage users: %v\n", am.HasPermission(user.ID, auth.PermissionManageUsers))
	fmt.Printf("  Admin can manage users: %v\n", am.HasPermission(admin.ID, auth.PermissionManageUsers))

	// Update user role
	fmt.Println("\nUpdating alice to merchant role...")
	err = am.UpdateUserRole(user.ID, auth.RoleMerchant)
	if err != nil {
		log.Fatal(err)
	}

	updatedUser, _ := am.GetUser(user.ID)
	fmt.Printf("  New role: %s\n", updatedUser.Role)
	fmt.Printf("  New permissions: %v\n", updatedUser.Permissions)
	fmt.Printf("  Can manage merchant: %v\n", am.HasPermission(user.ID, auth.PermissionManageMerchant))

	fmt.Println()
}

func demoModuleSystem() {
	fmt.Println("--- Demo 2: Add-On Module System ---")

	// Create module registry
	registry := addons.NewModuleRegistry()

	// Register staking module
	stakingModule := addons.NewStakingModule()
	err := registry.Register(stakingModule, "Bituncoin Team")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Registered module: %s v%s\n", stakingModule.GetName(), stakingModule.GetVersion())

	// Register lending module
	lendingModule := addons.NewLendingModule()
	err = registry.Register(lendingModule, "Bituncoin Team")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Registered module: %s v%s\n", lendingModule.GetName(), lendingModule.GetVersion())

	// List all modules
	fmt.Println("\nAvailable modules:")
	modules := registry.ListModules()
	for _, mod := range modules {
		fmt.Printf("  - %s (Category: %s, Status: %s)\n", mod.Name, mod.Category, mod.Status)
	}

	// Enable staking module
	fmt.Println("\nEnabling Advanced Staking module...")
	config := map[string]interface{}{
		"max_pools": 10,
		"default_apy": 5.0,
	}
	err = registry.Enable("Advanced Staking", config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Module enabled successfully")

	// List staking pools
	result, err := registry.Execute("Advanced Staking", "list_pools", nil)
	if err != nil {
		log.Fatal(err)
	}

	pools := result.([]addons.AdvancedStakePool)
	fmt.Printf("\nAvailable staking pools (%d):\n", len(pools))
	for _, pool := range pools {
		fmt.Printf("  - %s: APY %.1f%%, Min Stake: %.0f GLD\n", pool.Name, pool.APY, pool.MinStake)
	}

	// Enable lending module
	fmt.Println("\nEnabling DeFi Lending module...")
	err = registry.Enable("DeFi Lending", map[string]interface{}{})
	if err != nil {
		log.Fatal(err)
	}

	// Create a lending offer
	offerParams := map[string]interface{}{
		"lender":        "lender123",
		"amount":        1000.0,
		"interest_rate": 5.5,
	}

	offerResult, err := registry.Execute("DeFi Lending", "create_offer", offerParams)
	if err != nil {
		log.Fatal(err)
	}

	offer := offerResult.(*addons.LendingOffer)
	fmt.Printf("\nCreated lending offer:\n")
	fmt.Printf("  ID: %s\n", offer.ID)
	fmt.Printf("  Lender: %s\n", offer.Lender)
	fmt.Printf("  Amount: %.2f GLD\n", offer.Amount)
	fmt.Printf("  Interest Rate: %.2f%%\n", offer.InterestRate)
	fmt.Printf("  Status: %s\n", offer.Status)

	// Create a custom pool
	poolParams := map[string]interface{}{
		"id":   "demo-pool",
		"name": "Demo Staking Pool",
	}

	poolResult, err := registry.Execute("Advanced Staking", "create_pool", poolParams)
	if err != nil {
		log.Fatal(err)
	}

	pool := poolResult.(*addons.AdvancedStakePool)
	fmt.Printf("\nCreated custom staking pool:\n")
	fmt.Printf("  ID: %s\n", pool.ID)
	fmt.Printf("  Name: %s\n", pool.Name)
	fmt.Printf("  APY: %.1f%%\n", pool.APY)

	fmt.Println()
}

func demoAPIIntegration() {
	fmt.Println("--- Demo 3: Enhanced API Integration ---")

	// Create API node with new features
	node := api.NewNode("localhost", 8080)

	fmt.Println("API Node created with integrated features:")
	fmt.Println("  - Authentication & Authorization")
	fmt.Println("  - User Management")
	fmt.Println("  - Add-on Module System")

	fmt.Println("\nAvailable API Endpoints:")

	fmt.Println("\n  Authentication:")
	fmt.Println("    POST /api/auth/register - Register new user")
	fmt.Println("    POST /api/auth/login - User login")
	fmt.Println("    POST /api/auth/logout - User logout")
	fmt.Println("    GET /api/auth/validate - Validate session")

	fmt.Println("\n  User Management (Admin only):")
	fmt.Println("    GET /api/users/list - List all users")
	fmt.Println("    POST /api/users/update-role - Update user role")
	fmt.Println("    POST /api/users/deactivate - Deactivate user")

	fmt.Println("\n  Add-on Modules:")
	fmt.Println("    GET /api/addons/list - List available modules")
	fmt.Println("    POST /api/addons/enable - Enable a module")
	fmt.Println("    POST /api/addons/disable - Disable a module")
	fmt.Println("    POST /api/addons/execute - Execute module action")

	fmt.Println("\n  Gold-Coin & Wallet:")
	fmt.Println("    GET /api/goldcoin/balance - Get balance")
	fmt.Println("    POST /api/goldcoin/send - Send transaction")
	fmt.Println("    POST /api/goldcoin/stake - Stake tokens")
	fmt.Println("    GET /api/wallet/portfolio/:address - Get portfolio")

	info := node.GetNodeInfo()
	fmt.Printf("\nNode Information:\n")
	fmt.Printf("  Version: %s\n", info.Version)
	fmt.Printf("  Network: %s\n", info.Network)
	fmt.Printf("  Type: %s\n", info.NodeType)

	fmt.Println("\n=== Demo Complete ===")
	fmt.Println("\nThe Bituncoin wallet now includes:")
	fmt.Println("✓ User authentication with role-based access control")
	fmt.Println("✓ Extensible add-on module system")
	fmt.Println("✓ Enhanced API with security features")
	fmt.Println("✓ Multi-platform support (Web, Desktop, Mobile)")
	fmt.Println("✓ Automated CI/CD pipelines")
	fmt.Println("✓ Comprehensive documentation")
}
