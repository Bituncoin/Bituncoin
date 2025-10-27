package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Role represents user roles in the system
type Role string

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleMerchant  Role = "merchant"
	RoleValidator Role = "validator"
)

// Permission represents system permissions
type Permission string

const (
	PermissionRead           Permission = "read"
	PermissionWrite          Permission = "write"
	PermissionDelete         Permission = "delete"
	PermissionManageUsers    Permission = "manage_users"
	PermissionManageTokens   Permission = "manage_tokens"
	PermissionSystemConfig   Permission = "system_config"
	PermissionViewDashboard  Permission = "view_dashboard"
	PermissionManageMerchant Permission = "manage_merchant"
)

// User represents a user account in the system
type User struct {
	ID              string       `json:"id"`
	Username        string       `json:"username"`
	Email           string       `json:"email"`
	PasswordHash    string       `json:"-"` // Never expose in JSON
	Role            Role         `json:"role"`
	Permissions     []Permission `json:"permissions"`
	WalletAddresses []string     `json:"walletAddresses"`
	CreatedAt       time.Time    `json:"createdAt"`
	LastLogin       time.Time    `json:"lastLogin"`
	IsActive        bool         `json:"isActive"`
	TwoFactorSecret string       `json:"-"`
	TwoFactorEnabled bool        `json:"twoFactorEnabled"`
}

// AccountManager manages user accounts and authentication
type AccountManager struct {
	users       map[string]*User // username -> user
	usersByID   map[string]*User // id -> user
	sessions    map[string]*Session
	rolePerms   map[Role][]Permission
	mutex       sync.RWMutex
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	IPAddress string    `json:"ipAddress"`
}

// NewAccountManager creates a new account manager
func NewAccountManager() *AccountManager {
	am := &AccountManager{
		users:     make(map[string]*User),
		usersByID: make(map[string]*User),
		sessions:  make(map[string]*Session),
		rolePerms: make(map[Role][]Permission),
	}
	
	// Initialize default role permissions
	am.initializeRolePermissions()
	
	return am
}

// initializeRolePermissions sets up default permissions for each role
func (am *AccountManager) initializeRolePermissions() {
	// User permissions
	am.rolePerms[RoleUser] = []Permission{
		PermissionRead,
		PermissionWrite,
		PermissionViewDashboard,
	}
	
	// Merchant permissions
	am.rolePerms[RoleMerchant] = []Permission{
		PermissionRead,
		PermissionWrite,
		PermissionViewDashboard,
		PermissionManageMerchant,
	}
	
	// Validator permissions
	am.rolePerms[RoleValidator] = []Permission{
		PermissionRead,
		PermissionWrite,
		PermissionViewDashboard,
	}
	
	// Admin permissions (all permissions)
	am.rolePerms[RoleAdmin] = []Permission{
		PermissionRead,
		PermissionWrite,
		PermissionDelete,
		PermissionManageUsers,
		PermissionManageTokens,
		PermissionSystemConfig,
		PermissionViewDashboard,
		PermissionManageMerchant,
	}
}

// CreateUser creates a new user account
func (am *AccountManager) CreateUser(username, email, password string, role Role) (*User, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if username == "" || email == "" || password == "" {
		return nil, errors.New("username, email, and password are required")
	}
	
	// Check if username already exists
	if _, exists := am.users[username]; exists {
		return nil, errors.New("username already exists")
	}
	
	// Generate user ID
	userID, err := generateID()
	if err != nil {
		return nil, err
	}
	
	// Hash password
	passwordHash := hashPassword(password)
	
	// Get role permissions
	permissions, ok := am.rolePerms[role]
	if !ok {
		permissions = am.rolePerms[RoleUser] // Default to user permissions
		role = RoleUser
	}
	
	user := &User{
		ID:               userID,
		Username:         username,
		Email:            email,
		PasswordHash:     passwordHash,
		Role:             role,
		Permissions:      permissions,
		WalletAddresses:  make([]string, 0),
		CreatedAt:        time.Now(),
		LastLogin:        time.Time{},
		IsActive:         true,
		TwoFactorEnabled: false,
	}
	
	am.users[username] = user
	am.usersByID[userID] = user
	
	return user, nil
}

// Authenticate verifies user credentials and creates a session
func (am *AccountManager) Authenticate(username, password string) (*Session, error) {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	user, exists := am.users[username]
	if !exists {
		return nil, errors.New("invalid credentials")
	}
	
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}
	
	// Verify password
	if !verifyPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}
	
	// Create session
	sessionID, err := generateID()
	if err != nil {
		return nil, err
	}
	
	session := &Session{
		ID:        sessionID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour session
		IPAddress: "",
	}
	
	am.sessions[sessionID] = session
	user.LastLogin = time.Now()
	
	return session, nil
}

// ValidateSession checks if a session is valid
func (am *AccountManager) ValidateSession(sessionID string) (*User, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	session, exists := am.sessions[sessionID]
	if !exists {
		return nil, errors.New("invalid session")
	}
	
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}
	
	user, exists := am.usersByID[session.UserID]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}
	
	return user, nil
}

// HasPermission checks if a user has a specific permission
func (am *AccountManager) HasPermission(userID string, permission Permission) bool {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	user, exists := am.usersByID[userID]
	if !exists || !user.IsActive {
		return false
	}
	
	for _, perm := range user.Permissions {
		if perm == permission {
			return true
		}
	}
	
	return false
}

// UpdateUserRole updates a user's role and permissions
func (am *AccountManager) UpdateUserRole(userID string, newRole Role) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	user, exists := am.usersByID[userID]
	if !exists {
		return errors.New("user not found")
	}
	
	permissions, ok := am.rolePerms[newRole]
	if !ok {
		return errors.New("invalid role")
	}
	
	user.Role = newRole
	user.Permissions = permissions
	
	return nil
}

// DeactivateUser deactivates a user account
func (am *AccountManager) DeactivateUser(userID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	user, exists := am.usersByID[userID]
	if !exists {
		return errors.New("user not found")
	}
	
	user.IsActive = false
	
	return nil
}

// ListUsers returns all users (admin only)
func (am *AccountManager) ListUsers() []*User {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	users := make([]*User, 0, len(am.usersByID))
	for _, user := range am.usersByID {
		users = append(users, user)
	}
	
	return users
}

// GetUser returns a user by ID
func (am *AccountManager) GetUser(userID string) (*User, error) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	
	user, exists := am.usersByID[userID]
	if !exists {
		return nil, errors.New("user not found")
	}
	
	return user, nil
}

// AddWalletAddress adds a wallet address to a user
func (am *AccountManager) AddWalletAddress(userID, address string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	user, exists := am.usersByID[userID]
	if !exists {
		return errors.New("user not found")
	}
	
	user.WalletAddresses = append(user.WalletAddresses, address)
	
	return nil
}

// Logout invalidates a session
func (am *AccountManager) Logout(sessionID string) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()
	
	if _, exists := am.sessions[sessionID]; !exists {
		return errors.New("session not found")
	}
	
	delete(am.sessions, sessionID)
	
	return nil
}

// hashPassword creates a secure hash of the password using bcrypt
func hashPassword(password string) string {
	// Use bcrypt with default cost (10)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Fallback to a simple hash if bcrypt fails (should never happen in practice)
		return password
	}
	return string(hash)
}

// verifyPassword checks if a password matches the hash
func verifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generateID generates a random unique ID
func generateID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
