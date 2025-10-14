# Bituncoin Universal Wallet - Visual Guide

## Overview
The Bituncoin Universal Wallet is a modern, secure, multi-currency blockchain wallet with cross-chain capabilities.

## User Interface Sections

### 1. Create Wallet
The first step for new users is to create a wallet with optional security features:
- Two-Factor Authentication (2FA)
- Biometric Authentication
- Supports multiple currencies from creation

### 2. Dashboard
Main view showing:
- Wallet address
- Multi-currency balance cards (BTN, BTC, ETH, USDT, BNB)
- Refresh functionality
- Clean, modern design with gradient backgrounds

### 3. Send Transaction
Transaction interface featuring:
- Recipient address input
- Amount and currency selection
- Cross-chain transaction toggle
- Target chain selection for cross-chain transfers
- Form validation

### 4. Transaction History
Complete transaction log showing:
- Transaction ID
- From/To addresses
- Amount and currency
- Timestamp
- Cross-chain transaction indicators

### 5. Settings
Security configuration panel:
- Two-Factor Authentication toggle
- Biometric authentication toggle
- Wallet information display
- Security status indicators

## Color Scheme
- Primary: Purple gradient (#667eea to #764ba2)
- Background: Light gray (#f5f5f5)
- Cards: White with subtle shadows
- Accents: Various shades for different currencies

## Features Highlight

### Multi-Currency Support
- BTN (Bituncoin)
- BTC (Bitcoin)
- ETH (Ethereum)
- USDT (Tether)
- BNB (Binance Coin)

### Security Features
✓ Two-Factor Authentication (2FA)
✓ Biometric Authentication
✓ Encrypted Key Storage
✓ ECDSA P256 Key Generation
✓ SHA256 Address Hashing
✓ Transaction Signing

### Cross-Chain Capabilities
✓ Bridge transactions between chains
✓ Seamless currency conversion
✓ Multi-network support

## API Integration
All UI operations connect to the backend REST API:
- `POST /api/wallet/create` - Wallet creation
- `GET /api/wallet/balance` - Balance retrieval
- `POST /api/transaction/send` - Send transactions
- `GET /api/transaction/history` - Transaction history

## Responsive Design
The wallet interface is designed to be:
- Mobile-friendly
- Tablet-optimized
- Desktop-ready
- Accessible to all users

## Navigation
Simple tab-based navigation:
1. Create Wallet
2. Dashboard
3. Send
4. Transactions
5. Settings

Each tab provides focused functionality for specific wallet operations.
