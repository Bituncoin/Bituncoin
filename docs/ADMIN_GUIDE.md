# Bituncoin Universal Wallet - Administrator Guide

## Overview

The Bituncoin Universal Wallet is a comprehensive financial platform supporting multiple cryptocurrencies with enterprise-grade security and AI-driven insights. This guide covers administrative operations, user management, and system configuration.

## System Architecture

### Core Components
- **Authentication System**: bcrypt password hashing, JWT tokens, role-based access control
- **AI Wallet Manager**: Automated balance verification, transaction analysis, personalized insights
- **Multi-Currency Support**: BTN, GLD, BTC, ETH, USDT, BNB with cross-chain capabilities
- **Security Layers**: AES-256 encryption, 2FA, biometric authentication, fraud detection
- **Add-on System**: Plug-and-play modules for staking, DeFi lending, and custom features

### Supported Platforms
- Web (React)
- Windows (Electron)
- macOS Intel/ARM (Electron)
- Linux (Electron)
- Android (React Native)
- iOS (React Native)

## User Management

### Roles and Permissions

The system supports 4 distinct roles with granular permissions:

1. **User** (Basic user access)
   - View own wallet
   - Send/receive crypto
   - View transaction history
   - Basic AI insights

2. **Merchant** (Business user)
   - All User permissions
   - Accept payments
   - Generate payment links
   - Business analytics

3. **Admin** (System administrator)
   - All Merchant permissions
   - User management
   - System configuration
   - Security monitoring

4. **Validator** (Network validator)
   - All Admin permissions
   - Network validation
   - Consensus participation
   - Block production

### Permission Matrix

| Permission | User | Merchant | Admin | Validator |
|------------|------|----------|-------|-----------|
| view_wallet | ✅ | ✅ | ✅ | ✅ |
| send_crypto | ✅ | ✅ | ✅ | ✅ |
| receive_crypto | ✅ | ✅ | ✅ | ✅ |
| view_history | ✅ | ✅ | ✅ | ✅ |
| ai_insights | ✅ | ✅ | ✅ | ✅ |
| accept_payments | ❌ | ✅ | ✅ | ✅ |
| payment_links | ❌ | ✅ | ✅ | ✅ |
| business_analytics | ❌ | ✅ | ✅ | ✅ |
| user_management | ❌ | ❌ | ✅ | ✅ |
| system_config | ❌ | ❌ | ✅ | ✅ |
| security_monitor | ❌ | ❌ | ✅ | ✅ |
| network_validation | ❌ | ❌ | ❌ | ✅ |
| consensus_participation | ❌ | ❌ | ❌ | ✅ |
| block_production | ❌ | ❌ | ❌ | ✅ |

## Security Configuration

### Password Policies
- Minimum length: 12 characters
- Must contain: uppercase, lowercase, numbers, special characters
- bcrypt hashing with 12 rounds
- Password history: last 5 passwords cannot be reused

### Session Management
- JWT tokens with 24-hour expiry
- Automatic logout on inactivity (30 minutes)
- Concurrent session limits: 3 per user
- Secure token storage with AES-256 encryption

### Two-Factor Authentication (2FA)
- TOTP (Time-based One-Time Password) support
- SMS backup codes
- Hardware security keys (FIDO2/WebAuthn)
- Biometric authentication on mobile devices

### Encryption Standards
- AES-256-GCM for data at rest
- TLS 1.3 for data in transit
- Hardware Security Module (HSM) integration
- Quantum-resistant signatures (future upgrade)

## System Monitoring

### Key Metrics
- Active users and sessions
- Transaction volume and success rates
- System uptime and performance
- Security incidents and alerts

### Alert Configuration
- Critical alerts: System downtime, security breaches
- Warning alerts: High CPU usage, failed transactions
- Info alerts: New user registrations, large transactions

### Audit Logging
- All administrative actions logged
- User activity tracking
- Security events monitoring
- 7-year retention policy
- Immutable blockchain-based audit trails

## Backup and Recovery

### Automated Backups
- Daily encrypted backups
- Off-site storage with geo-redundancy
- Point-in-time recovery capabilities
- Backup integrity verification

### Disaster Recovery
- Multi-region failover
- Automated recovery procedures
- Business continuity planning
- Regular DR testing

## Add-on Module Management

### Installing Modules
```bash
# Install from registry
wallet-cli module install staking-module

# Install from file
wallet-cli module install /path/to/module.zip

# List installed modules
wallet-cli module list
```

### Module Permissions
- Sandboxed execution environment
- Granular permission system
- Security scanning before installation
- Automatic updates and patches

### Custom Module Development
See MODULE_DEVELOPER_GUIDE.md for detailed instructions on creating custom add-on modules.

## Network Configuration

### Multi-Currency Support
- BTC: Bitcoin mainnet and testnet
- ETH: Ethereum mainnet and testnets
- BNB: Binance Smart Chain
- BTN: Bituncoin Gold mainnet
- GLD: Gold token network
- USDT: Tether on multiple chains

### Cross-Chain Operations
- Atomic swaps support
- Bridge protocols integration
- Interoperability standards compliance
- Multi-chain transaction monitoring

## Compliance and Regulation

### KYC/AML Integration
- Automated KYC verification
- AML transaction monitoring
- Regulatory reporting
- Geographic restrictions

### Data Privacy
- GDPR compliance
- Data minimization principles
- User consent management
- Right to erasure implementation

## Troubleshooting

### Common Issues
1. **Login failures**: Check password policies and 2FA settings
2. **Transaction delays**: Verify network connectivity and gas fees
3. **Module errors**: Check module permissions and compatibility
4. **Performance issues**: Monitor system resources and optimize queries

### Support Resources
- Technical documentation: docs.bituncoin.com
- Community forums: forum.bituncoin.com
- Emergency support: support@bituncoin.com
- Security incidents: security@bituncoin.com

## Maintenance Procedures

### Regular Tasks
- Weekly security updates
- Monthly performance optimization
- Quarterly security audits
- Annual disaster recovery testing

### Emergency Procedures
- Security breach response
- System outage recovery
- Data corruption handling
- Regulatory compliance updates

## API Reference

### REST API Endpoints
- `POST /api/auth/login` - User authentication
- `GET /api/users` - List users (admin only)
- `POST /api/users` - Create user (admin only)
- `PUT /api/users/{id}` - Update user (admin only)
- `DELETE /api/users/{id}` - Delete user (admin only)
- `GET /api/wallet/balance` - Get wallet balance
- `POST /api/wallet/send` - Send cryptocurrency
- `GET /api/transactions` - Get transaction history
- `POST /api/modules/install` - Install add-on module (admin only)

### WebSocket Events
- `wallet:update` - Wallet balance updates
- `transaction:new` - New transaction notifications
- `alert:security` - Security alerts
- `system:status` - System status updates

## Version Information

- Current Version: 1.0.0
- Release Date: February 24, 2026
- Supported Platforms: Web, Windows, macOS, Linux, Android, iOS
- Minimum Requirements: See PLATFORM_DEPLOYMENT.md

For additional support or questions, please contact the Bituncoin support team.
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
