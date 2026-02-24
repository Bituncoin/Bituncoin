package auth

import (
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	am := NewAccountManager()
	
	user, err := am.CreateUser("testuser", "test@example.com", "password123", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	if user.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", user.Username)
	}
	
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
	
	if user.Role != RoleUser {
		t.Errorf("Expected role 'user', got '%s'", user.Role)
	}
	
	if !user.IsActive {
		t.Error("Expected user to be active")
	}
}

func TestCreateAdminUser(t *testing.T) {
	am := NewAccountManager()
	
	admin, err := am.CreateUser("admin", "admin@example.com", "adminpass", RoleAdmin)
	if err != nil {
		t.Fatalf("Failed to create admin: %v", err)
	}
	
	if admin.Role != RoleAdmin {
		t.Errorf("Expected role 'admin', got '%s'", admin.Role)
	}
	
	// Verify admin has all permissions
	expectedPerms := []Permission{
		PermissionRead,
		PermissionWrite,
		PermissionDelete,
		PermissionManageUsers,
		PermissionManageTokens,
		PermissionSystemConfig,
		PermissionViewDashboard,
		PermissionManageMerchant,
	}
	
	if len(admin.Permissions) != len(expectedPerms) {
		t.Errorf("Expected %d permissions, got %d", len(expectedPerms), len(admin.Permissions))
	}
}

func TestDuplicateUsername(t *testing.T) {
	am := NewAccountManager()
	
	_, err := am.CreateUser("duplicate", "user1@example.com", "pass1", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create first user: %v", err)
	}
	
	_, err = am.CreateUser("duplicate", "user2@example.com", "pass2", RoleUser)
	if err == nil {
		t.Error("Expected error for duplicate username")
	}
}

func TestAuthenticate(t *testing.T) {
	am := NewAccountManager()
	
	_, err := am.CreateUser("authtest", "auth@example.com", "mypassword", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	session, err := am.Authenticate("authtest", "mypassword")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}
	
	if session.ID == "" {
		t.Error("Expected valid session ID")
	}
	
	if session.UserID == "" {
		t.Error("Expected valid user ID in session")
	}
}

func TestAuthenticateWrongPassword(t *testing.T) {
	am := NewAccountManager()
	
	_, err := am.CreateUser("wrongpass", "wrong@example.com", "correctpass", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	_, err = am.Authenticate("wrongpass", "wrongpass")
	if err == nil {
		t.Error("Expected authentication to fail with wrong password")
	}
}

func TestValidateSession(t *testing.T) {
	am := NewAccountManager()
	
	_, err := am.CreateUser("sessiontest", "session@example.com", "password", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	session, err := am.Authenticate("sessiontest", "password")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}
	
	user, err := am.ValidateSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to validate session: %v", err)
	}
	
	if user.Username != "sessiontest" {
		t.Errorf("Expected username 'sessiontest', got '%s'", user.Username)
	}
}

func TestExpiredSession(t *testing.T) {
	am := NewAccountManager()
	
	_, err := am.CreateUser("expiretest", "expire@example.com", "password", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	session, err := am.Authenticate("expiretest", "password")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}
	
	// Manually expire the session
	am.mutex.Lock()
	am.sessions[session.ID].ExpiresAt = time.Now().Add(-1 * time.Hour)
	am.mutex.Unlock()
	
	_, err = am.ValidateSession(session.ID)
	if err == nil {
		t.Error("Expected error for expired session")
	}
}

func TestHasPermission(t *testing.T) {
	am := NewAccountManager()
	
	user, err := am.CreateUser("permtest", "perm@example.com", "password", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	// User should have read permission
	if !am.HasPermission(user.ID, PermissionRead) {
		t.Error("Expected user to have read permission")
	}
	
	// User should not have admin permissions
	if am.HasPermission(user.ID, PermissionManageUsers) {
		t.Error("Expected user to not have manage users permission")
	}
}

func TestUpdateUserRole(t *testing.T) {
	am := NewAccountManager()
	
	user, err := am.CreateUser("roletest", "role@example.com", "password", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	err = am.UpdateUserRole(user.ID, RoleAdmin)
	if err != nil {
		t.Fatalf("Failed to update role: %v", err)
	}
	
	updatedUser, err := am.GetUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	
	if updatedUser.Role != RoleAdmin {
		t.Errorf("Expected role 'admin', got '%s'", updatedUser.Role)
	}
	
	// Should now have admin permissions
	if !am.HasPermission(user.ID, PermissionManageUsers) {
		t.Error("Expected admin to have manage users permission")
	}
}

func TestDeactivateUser(t *testing.T) {
	am := NewAccountManager()
	
	user, err := am.CreateUser("deactivate", "deactivate@example.com", "password", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	err = am.DeactivateUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to deactivate user: %v", err)
	}
	
	// Try to authenticate with deactivated account
	_, err = am.Authenticate("deactivate", "password")
	if err == nil {
		t.Error("Expected authentication to fail for deactivated account")
	}
}

func TestListUsers(t *testing.T) {
	am := NewAccountManager()
	
	am.CreateUser("user1", "user1@example.com", "pass1", RoleUser)
	am.CreateUser("user2", "user2@example.com", "pass2", RoleUser)
	am.CreateUser("admin1", "admin1@example.com", "pass3", RoleAdmin)
	
	users := am.ListUsers()
	
	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}

func TestAddWalletAddress(t *testing.T) {
	am := NewAccountManager()
	
	user, err := am.CreateUser("wallettest", "wallet@example.com", "password", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	err = am.AddWalletAddress(user.ID, "GLD1234567890")
	if err != nil {
		t.Fatalf("Failed to add wallet address: %v", err)
	}
	
	updatedUser, err := am.GetUser(user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	
	if len(updatedUser.WalletAddresses) != 1 {
		t.Errorf("Expected 1 wallet address, got %d", len(updatedUser.WalletAddresses))
	}
	
	if updatedUser.WalletAddresses[0] != "GLD1234567890" {
		t.Errorf("Expected wallet address 'GLD1234567890', got '%s'", updatedUser.WalletAddresses[0])
	}
}

func TestLogout(t *testing.T) {
	am := NewAccountManager()
	
	_, err := am.CreateUser("logouttest", "logout@example.com", "password", RoleUser)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	
	session, err := am.Authenticate("logouttest", "password")
	if err != nil {
		t.Fatalf("Failed to authenticate: %v", err)
	}
	
	err = am.Logout(session.ID)
	if err != nil {
		t.Fatalf("Failed to logout: %v", err)
	}
	
	_, err = am.ValidateSession(session.ID)
	if err == nil {
		t.Error("Expected session validation to fail after logout")
	}
}
