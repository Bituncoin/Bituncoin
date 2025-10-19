# Bituncoin Address Formats

This document describes the different address formats supported by the Bituncoin blockchain ecosystem for the universal wallet feature.

## Overview

The Bituncoin blockchain supports multiple address formats to provide compatibility with different blockchain ecosystems and use cases:

1. **GLD Address** - Legacy format for Gold-Coin (GLD) cryptocurrency
2. **Bitcoin-Style Address** - Compatible with Bitcoin-based blockchains with custom "Btu" prefix
3. **Ethereum-Style Address** - Compatible with Ethereum and EVM-based smart contracts

## Address Types

### 1. GLD Address (Legacy)

The original address format for Gold-Coin cryptocurrency.

**Format:**
- Prefix: `GLD`
- Encoding: Hexadecimal
- Length: 43 characters (3-char prefix + 40 hex characters)

**Example:**
```
GLD5a7d3e9f1c2b4a8e6d3f9c7b2e4a8d6c3f9b7e
```

**Usage:**
- Default format for Gold-Coin transactions
- Backward compatible with existing wallets

### 2. Bitcoin-Style Address

Bitcoin-compatible address format with custom "Btu" (Bituncoin) prefix.

**Format:**
- Prefix: `Btu`
- Encoding: Base58Check
- Checksum: 4-byte double SHA-256 hash
- Length: ~29-40 characters

**Example:**
```
BtuzP1eP5QGefi2DMPTfTL5SLmv7DivfNa
```

**Features:**
- **Base58Check Encoding**: Uses Bitcoin's Base58 alphabet (excludes 0, O, I, l to avoid confusion)
- **Built-in Checksum**: 4-byte checksum prevents typos and errors
- **Error Detection**: Invalid addresses are detected before transactions
- **Human Readable**: Easier to read and transcribe than hexadecimal

**Technical Details:**
```
Address = "Btu" + Base58Check(PublicKeyHash + Checksum)
Checksum = First4Bytes(SHA256(SHA256(PublicKeyHash)))
```

**Use Cases:**
- Cross-chain transactions with Bitcoin-based networks
- Bitcoin wallet compatibility
- Enhanced error detection for manual address entry

### 3. Ethereum-Style Address

Ethereum-compatible address format for smart contract interactions.

**Format:**
- Prefix: `0x`
- Encoding: Hexadecimal
- Length: 42 characters (2-char prefix + 40 hex characters)

**Example:**
```
0x1234567890abcdef1234567890abcdef12345678
```

**Features:**
- **EVM Compatible**: Works with Ethereum Virtual Machine
- **Smart Contract Support**: Can interact with deployed smart contracts
- **Standard Format**: Follows ERC-20 and other Ethereum standards
- **Case Insensitive**: Accepts both uppercase and lowercase hex characters

**Technical Details:**
```
Address = "0x" + Hex(Last20Bytes(SHA256(PublicKey)))
```

**Use Cases:**
- Smart contract deployment and interaction
- Cross-chain transactions with Ethereum and BSC
- DeFi applications and token swaps
- NFT minting and trading

## Address Generation

### Using the Go API

```go
package main

import (
    "fmt"
    "log"
    "github.com/Bituncoin/Bituncoin/identity"
)

func main() {
    // Initialize address manager
    manager := identity.NewAddressManager()
    
    // Generate Bitcoin-style address
    btcAddr, err := manager.GenerateAddressOfType("My BTC Wallet", identity.BitcoinStyle)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Bitcoin-style: %s\n", btcAddr.Address)
    
    // Generate Ethereum-style address
    ethAddr, err := manager.GenerateAddressOfType("My ETH Wallet", identity.EthereumStyle)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Ethereum-style: %s\n", ethAddr.Address)
    
    // Generate legacy GLD address
    gldAddr, err := manager.GenerateAddress("My GLD Wallet")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("GLD Address: %s\n", gldAddr.Address)
}
```

### Direct Address Generation

```go
// Generate Bitcoin-style address
address, privateKey, err := identity.GenerateAddressWithType(identity.BitcoinStyle)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Address: %s\n", address)

// Generate Ethereum-style address
address, privateKey, err := identity.GenerateAddressWithType(identity.EthereumStyle)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Address: %s\n", address)
```

## Address Validation

### Validate Bitcoin-Style Address

```go
err := identity.ValidateBitcoinStyleAddress("BtuzP1eP5QGefi2DMPTfTL5SLmv7DivfNa")
if err != nil {
    fmt.Printf("Invalid address: %v\n", err)
} else {
    fmt.Println("Valid Bitcoin-style address!")
}
```

**Validation Checks:**
- Prefix must be "Btu"
- Length must be between 29-40 characters
- Must contain only Base58 characters
- Checksum must be valid

### Validate Ethereum-Style Address

```go
err := identity.ValidateEthereumStyleAddress("0x1234567890abcdef1234567890abcdef12345678")
if err != nil {
    fmt.Printf("Invalid address: %v\n", err)
} else {
    fmt.Println("Valid Ethereum-style address!")
}
```

**Validation Checks:**
- Prefix must be "0x"
- Length must be exactly 42 characters
- Must contain only hexadecimal characters (0-9, a-f, A-F)

### Auto-Detect Address Type

```go
addressType, err := identity.ValidateAddressType(someAddress)
if err != nil {
    fmt.Printf("Unknown address type: %v\n", err)
} else {
    fmt.Printf("Address type: %s\n", addressType)
}
```

## Security Considerations

### Private Key Management

- **Never share private keys**: Keep private keys secure and encrypted
- **Use strong encryption**: Wallet uses AES-256 encryption for key storage
- **Backup securely**: Create encrypted backups of private keys
- **Use recovery phrases**: Store 12-word recovery phrases offline

### Address Validation

- **Always validate**: Check address format before sending transactions
- **Verify checksums**: Bitcoin-style addresses include checksums for error detection
- **Double-check recipients**: Confirm recipient address before large transactions
- **Use copy-paste**: Avoid manual typing to prevent errors

### Best Practices

1. **Generate new addresses** for each transaction to enhance privacy
2. **Use appropriate type** based on the blockchain you're interacting with
3. **Test with small amounts** before sending large transactions
4. **Keep backups** of all addresses and their labels
5. **Enable 2FA** on wallet for additional security

## Compatibility Matrix

| Address Type | Bitcoin | Ethereum | BSC | Gold-Coin | Smart Contracts |
|-------------|---------|----------|-----|-----------|-----------------|
| GLD Address | ❌ | ❌ | ❌ | ✅ | ❌ |
| Bitcoin-Style | ✅ | ❌ | ❌ | ✅ | ❌ |
| Ethereum-Style | ❌ | ✅ | ✅ | ✅ | ✅ |

## Error Handling

### Common Validation Errors

**Bitcoin-Style Addresses:**
- `"invalid Bitcoin-style address: must start with 'Btu'"` - Wrong prefix
- `"invalid Bitcoin-style address: too short"` - Address length < 29
- `"invalid Bitcoin-style address: too long"` - Address length > 40
- `"invalid Bitcoin-style address: invalid character"` - Contains 0, O, I, or l
- `"invalid Bitcoin-style address: checksum mismatch"` - Corrupted address

**Ethereum-Style Addresses:**
- `"invalid Ethereum-style address: must start with '0x'"` - Wrong prefix
- `"invalid Ethereum-style address: must be 42 characters"` - Wrong length
- `"invalid Ethereum-style address: must contain only hexadecimal characters"` - Invalid characters

### Error Recovery

```go
// Attempt to detect and fix common issues
address := strings.TrimSpace(userInput) // Remove whitespace
address = strings.ToLower(address)      // Normalize case for Ethereum

// Validate the cleaned address
if err := identity.ValidateEthereumStyleAddress(address); err != nil {
    return fmt.Errorf("could not validate address: %w", err)
}
```

## Future Enhancements

- Support for more blockchain address formats (Cardano, Solana, etc.)
- QR code generation for addresses
- Address book with contact management
- Multi-signature addresses
- Hierarchical Deterministic (HD) wallet support
- Address aliasing and ENS-style naming

## References

- [Bitcoin Base58Check Encoding](https://en.bitcoin.it/wiki/Base58Check_encoding)
- [Ethereum Address Format](https://ethereum.org/en/developers/docs/accounts/)
- [BIP-32: Hierarchical Deterministic Wallets](https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki)
- [EIP-55: Mixed-case checksum address encoding](https://eips.ethereum.org/EIPS/eip-55)

## Support

For issues or questions:
- GitHub Issues: https://github.com/Bituncoin/Bituncoin/issues
- Documentation: See README.md for more information

---

**Last Updated:** October 2025  
**Version:** 1.0.0
