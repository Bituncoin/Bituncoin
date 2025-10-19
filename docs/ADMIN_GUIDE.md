# Admin Guide

## Overview
This guide provides administrators with instructions for managing the Bituncoin wallet system, including user management, system configuration, and add-on module administration.

## Table of Contents
1. [User Management](#user-management)
2. [Role-Based Access Control](#role-based-access-control)
3. [Add-On Module Management](#add-on-module-management)
4. [System Monitoring](#system-monitoring)
5. [Security Best Practices](#security-best-practices)

## User Management

### Creating User Accounts

Administrators can create user accounts with different roles:

**API Endpoint:** `POST /api/auth/register`

**Request:**
```json
{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "securepassword123",
  "role": "user"
}
```

**Available Roles:**
- `user` - Standard user with basic wallet functionality
- `admin` - Administrator with full system access
- `merchant` - Merchant with payment processing capabilities
- `validator` - Network validator with staking privileges

### Listing All Users

**API Endpoint:** `GET /api/users/list`

**Headers:**
```
X-Session-ID: <admin-session-id>
```

**Response:**
```json
[
  {
    "id": "user-id-123",
    "username": "alice",
    "email": "alice@example.com",
    "role": "user",
    "permissions": ["read", "write", "view_dashboard"],
    "walletAddresses": ["GLD123..."],
    "createdAt": "2025-10-19T00:00:00Z",
    "lastLogin": "2025-10-19T12:00:00Z",
    "isActive": true
  }
]
```

### Updating User Roles

**API Endpoint:** `POST /api/users/update-role`

**Headers:**
```
X-Session-ID: <admin-session-id>
```

**Request:**
```json
{
  "userId": "user-id-123",
  "newRole": "merchant"
}
```

### Deactivating Users

**API Endpoint:** `POST /api/users/deactivate`

**Headers:**
```
X-Session-ID: <admin-session-id>
```

**Request:**
```json
{
  "userId": "user-id-123"
}
```

## Role-Based Access Control

### Permission System

Each role has specific permissions:

#### User Permissions
- `read` - Read wallet and transaction data
- `write` - Create transactions
- `view_dashboard` - Access user dashboard

#### Merchant Permissions
- All user permissions
- `manage_merchant` - Manage merchant services and payment requests

#### Validator Permissions
- All user permissions
- Validate blocks and participate in consensus

#### Admin Permissions
- All permissions including:
  - `delete` - Delete data
  - `manage_users` - Create, modify, and deactivate users
  - `manage_tokens` - Manage token configurations
  - `system_config` - Modify system settings

### Custom Permission Assignment

While roles have default permissions, administrators can customize permissions for specific use cases by modifying the role permission mappings in the system configuration.

## Add-On Module Management

### Listing Available Modules

**API Endpoint:** `GET /api/addons/list`

**Response:**
```json
[
  {
    "name": "Advanced Staking",
    "version": "1.0.0",
    "category": "staking",
    "description": "Advanced staking with multiple pools and auto-compounding",
    "status": "disabled",
    "author": "Bituncoin Team"
  },
  {
    "name": "DeFi Lending",
    "version": "1.0.0",
    "category": "lending",
    "description": "Decentralized lending and borrowing",
    "status": "enabled",
    "author": "Bituncoin Team"
  }
]
```

### Enabling a Module

**API Endpoint:** `POST /api/addons/enable`

**Request:**
```json
{
  "name": "Advanced Staking",
  "config": {
    "maxPools": 10,
    "defaultAPY": 5.0
  }
}
```

### Disabling a Module

**API Endpoint:** `POST /api/addons/disable`

**Request:**
```json
{
  "name": "Advanced Staking"
}
```

### Module Categories

Modules are organized into categories:

- **DeFi** - Decentralized finance features
- **Staking** - Staking and rewards
- **Lending** - Lending and borrowing
- **Trading** - Trading and exchange features
- **Payment** - Payment processing
- **Analytics** - Data analysis and insights
- **Security** - Security enhancements
- **Utility** - General utility features

## System Monitoring

### Dashboard Access

The admin dashboard is accessible at `/api/dashboard/status` and provides:

- System health status
- Component status (wallet, exchange, merchant services, etc.)
- Blockchain network connectivity
- Real-time metrics (users, transactions, volume)
- Performance indicators
- Active alerts

### Health Checks

**API Endpoint:** `GET /api/health`

Regular health checks ensure system availability.

### Metrics

Key metrics tracked:
- Total users
- Active wallets
- Total transactions
- Transaction volume (USD)
- Average response time
- Error rate
- System uptime

## Security Best Practices

### 1. Session Management
- Sessions expire after 24 hours
- Implement regular session cleanup
- Monitor active sessions for suspicious activity

### 2. Password Policies
- Enforce strong password requirements
- Implement password rotation policies
- Use secure password hashing (SHA-256)

### 3. Two-Factor Authentication
- Encourage all users to enable 2FA
- Require 2FA for admin accounts
- Provide backup codes for account recovery

### 4. Audit Logging
- Log all admin actions
- Monitor user authentication attempts
- Track permission changes
- Review logs regularly

### 5. Network Security
- Use HTTPS for all API communications
- Implement rate limiting
- Enable DDoS protection
- Whitelist trusted IP addresses for admin access

### 6. Regular Updates
- Keep system software up to date
- Apply security patches promptly
- Monitor security advisories

### 7. Backup and Recovery
- Regular automated backups
- Test recovery procedures
- Secure backup storage
- Encrypted backups

## Troubleshooting

### Common Issues

#### Users Cannot Log In
1. Check if account is active
2. Verify session hasn't expired
3. Check network connectivity
4. Review authentication logs

#### Module Won't Enable
1. Verify module is registered
2. Check configuration parameters
3. Review module dependencies
4. Check system resources

#### Performance Issues
1. Check system metrics
2. Review active connections
3. Analyze database performance
4. Monitor resource usage

## Support

For additional support:
- GitHub Issues: https://github.com/Bituncoin/Bituncoin/issues
- Documentation: /docs directory
- Community: Discord channel

## Appendix

### API Authentication Flow

1. User registers or logs in
2. System creates session with unique ID
3. Client includes session ID in subsequent requests
4. System validates session and permissions
5. Request processed if authorized

### Module Development

See the [Module Developer Guide](MODULE_DEVELOPER_GUIDE.md) for creating custom add-on modules.
