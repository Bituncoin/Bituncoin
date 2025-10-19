# Bituncoin Address Generation

This document describes the address generation feature in the Bituncoin blockchain network. Bituncoin supports two types of blockchain addresses: **Bitcoin-style addresses** and **Ethereum-style addresses**.

## Table of Contents

- [Overview](#overview)
- [Bitcoin-Style Addresses](#bitcoin-style-addresses)
- [Ethereum-Style Addresses](#ethereum-style-addresses)
- [Usage Examples](#usage-examples)
- [Validation](#validation)
- [Integration](#integration)

## Overview

The Bituncoin blockchain supports dual address generation formats to ensure compatibility with both Bitcoin-style Base58Check encoding and Ethereum-style hexadecimal addressing. This provides flexibility for users familiar with different blockchain ecosystems.

## Bitcoin-Style Addresses

### Format

Bitcoin-style addresses in Bituncoin use the following format:

- **Prefix**: `Btu` (Bituncoin identifier)
- **Encoding**: Base58Check
- **Checksum**: Included (first 4 bytes of double SHA-256 hash)
- **Example**: `Btu1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa`

### Generation Process

1. **Generate ECDSA Key Pair**: Uses the P-256 elliptic curve
2. **Hash Public Key**: 
   - Apply SHA-256 to the public key
   - Apply RIPEMD-160 to the SHA-256 hash
3. **Add Version Byte**: Prepend version byte (0x00) to the public key hash
4. **Calculate Checksum**:
   - Apply SHA-256 twice to the versioned payload
   - Take the first 4 bytes as the checksum
5. **Concatenate**: Combine versioned payload and checksum
6. **Encode**: Base58 encode the full payload
7. **Add Prefix**: Prepend "Btu" to create the final address

### Properties

- **Length**: 26-40 characters (including "Btu" prefix)
- **Error Detection**: Built-in checksum for typo prevention
- **Collision Resistant**: Uses cryptographic hashing (SHA-256 + RIPEMD-160)
- **Human Readable**: Base58 encoding avoids ambiguous characters (0, O, I, l)

### Code Example

```go
package main

import (
    "fmt"
    "github.com/Bituncoin/Bituncoin/identity"
)

func main() {
    // Generate a Bitcoin-style address
    addr, err := identity.GenerateBitcoinStyleAddress()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Address: %s\n", addr.Address)
    fmt.Printf("Public Key: %s\n", addr.PublicKey)
    fmt.Printf("Private Key: %s\n", addr.PrivateKey)

    // Validate the address
    err = identity.ValidateBitcoinStyleAddress(addr.Address)
    if err != nil {
        fmt.Printf("Invalid address: %v\n", err)
    } else {
        fmt.Println("Address is valid!")
    }
}
```

## Ethereum-Style Addresses

### Format

Ethereum-style addresses in Bituncoin use the following format:

- **Prefix**: `0x`
- **Encoding**: Hexadecimal (lowercase)
- **Length**: 40 hexadecimal characters (excluding 0x prefix)
- **Example**: `0x1234567890abcdef1234567890abcdef12345678`

### Generation Process

1. **Generate ECDSA Key Pair**: Uses the P-256 elliptic curve
2. **Hash Public Key**: Apply SHA-256 to the public key (in production, use Keccak-256)
3. **Take Last 20 Bytes**: Extract the last 20 bytes of the hash
4. **Encode**: Convert to hexadecimal string
5. **Add Prefix**: Prepend "0x" to create the final address

### Properties

- **Length**: Exactly 42 characters (including "0x" prefix)
- **Smart Contract Compatible**: Compatible with Ethereum smart contract standards
- **Case Insensitive**: Uses lowercase hexadecimal characters
- **Fixed Size**: Always 20 bytes (40 hex characters)

### Code Example

```go
package main

import (
    "fmt"
    "github.com/Bituncoin/Bituncoin/identity"
)

func main() {
    // Generate an Ethereum-style address
    addr, err := identity.GenerateEthereumStyleAddress()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Address: %s\n", addr.Address)
    fmt.Printf("Public Key: %s\n", addr.PublicKey)
    fmt.Printf("Private Key: %s\n", addr.PrivateKey)

    // Validate the address
    err = identity.ValidateEthereumStyleAddress(addr.Address)
    if err != nil {
        fmt.Printf("Invalid address: %v\n", err)
    } else {
        fmt.Println("Address is valid!")
    }
}
```

## Usage Examples

### Complete Example

```go
package main

import (
    "fmt"
    "github.com/Bituncoin/Bituncoin/identity"
)

func main() {
    fmt.Println("=== Bitcoin-Style Address ===")
    btcAddr, err := identity.GenerateBitcoinStyleAddress()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Address: %s\n", btcAddr.Address)
    
    if err := identity.ValidateBitcoinStyleAddress(btcAddr.Address); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("✓ Valid Bitcoin-style address")
    }

    fmt.Println("\n=== Ethereum-Style Address ===")
    ethAddr, err := identity.GenerateEthereumStyleAddress()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Address: %s\n", ethAddr.Address)
    
    if err := identity.ValidateEthereumStyleAddress(ethAddr.Address); err != nil {
        fmt.Printf("Validation failed: %v\n", err)
    } else {
        fmt.Println("✓ Valid Ethereum-style address")
    }
}
```

### Testing Multiple Addresses

```go
package main

import (
    "fmt"
    "github.com/Bituncoin/Bituncoin/identity"
)

func main() {
    // Generate multiple Bitcoin-style addresses
    fmt.Println("Generating 5 Bitcoin-style addresses:")
    for i := 0; i < 5; i++ {
        addr, err := identity.GenerateBitcoinStyleAddress()
        if err != nil {
            panic(err)
        }
        fmt.Printf("%d. %s\n", i+1, addr.Address)
    }

    // Generate multiple Ethereum-style addresses
    fmt.Println("\nGenerating 5 Ethereum-style addresses:")
    for i := 0; i < 5; i++ {
        addr, err := identity.GenerateEthereumStyleAddress()
        if err != nil {
            panic(err)
        }
        fmt.Printf("%d. %s\n", i+1, addr.Address)
    }
}
```

## Validation

Both address types include validation functions to ensure address integrity.

### Bitcoin-Style Validation

The validation checks:
- ✓ Address starts with "Btu" prefix
- ✓ Valid Base58 encoding (no invalid characters)
- ✓ Correct checksum (error detection)
- ✓ Valid version byte
- ✓ Minimum length requirements

### Ethereum-Style Validation

The validation checks:
- ✓ Address starts with "0x" prefix
- ✓ Exactly 40 hexadecimal characters
- ✓ Valid hexadecimal encoding
- ✓ Case sensitivity (lowercase)

### Validation Examples

```go
// Validate Bitcoin-style address
valid := "Btu1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"
if err := identity.ValidateBitcoinStyleAddress(valid); err == nil {
    fmt.Println("Valid Bitcoin-style address")
}

// Invalid examples
invalid1 := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa" // Missing Btu prefix
invalid2 := "Btu0OIl" // Invalid Base58 characters (0, O, I, l)
invalid3 := "Btu123" // Too short

// Validate Ethereum-style address
validEth := "0x1234567890abcdef1234567890abcdef12345678"
if err := identity.ValidateEthereumStyleAddress(validEth); err == nil {
    fmt.Println("Valid Ethereum-style address")
}

// Invalid examples
invalidEth1 := "1234567890abcdef1234567890abcdef12345678" // Missing 0x
invalidEth2 := "0x123" // Too short
invalidEth3 := "0x1234567890abcdef1234567890abcdefghijklmn" // Invalid hex
```

## Integration

### Wallet Integration

The address generation feature integrates seamlessly with the existing Bituncoin wallet system:

```go
import (
    "github.com/Bituncoin/Bituncoin/identity"
    "github.com/Bituncoin/Bituncoin/wallet"
)

// Generate address for wallet
func CreateWalletAddress(addressType string) error {
    var address string
    
    switch addressType {
    case "bitcoin":
        addr, err := identity.GenerateBitcoinStyleAddress()
        if err != nil {
            return err
        }
        address = addr.Address
        
    case "ethereum":
        addr, err := identity.GenerateEthereumStyleAddress()
        if err != nil {
            return err
        }
        address = addr.Address
        
    default:
        return errors.New("invalid address type")
    }
    
    // Store in wallet
    // ... wallet integration code ...
    
    return nil
}
```

### Transaction Integration

Addresses can be used in transactions:

```go
import (
    "github.com/Bituncoin/Bituncoin/identity"
    "github.com/Bituncoin/Bituncoin/goldcoin"
)

func SendTransaction(from, to string, amount float64) error {
    // Validate addresses before creating transaction
    if strings.HasPrefix(from, "Btu") {
        if err := identity.ValidateBitcoinStyleAddress(from); err != nil {
            return err
        }
    } else if strings.HasPrefix(from, "0x") {
        if err := identity.ValidateEthereumStyleAddress(from); err != nil {
            return err
        }
    }
    
    // Similar validation for 'to' address
    // ...
    
    // Create and send transaction
    gc := goldcoin.NewGoldCoin()
    tx, err := gc.CreateTransaction(from, to, amount)
    if err != nil {
        return err
    }
    
    // Process transaction...
    return nil
}
```

## Security Considerations

### Key Management

- **Private Keys**: Never share or expose private keys
- **Secure Storage**: Store private keys encrypted
- **Backup**: Keep secure backups of private keys
- **Recovery**: Use the private key to recover funds if needed

### Best Practices

1. **Always Validate**: Validate addresses before sending transactions
2. **Verify Checksums**: Bitcoin-style addresses include checksums - use them
3. **Use Unique Addresses**: Generate new addresses for enhanced privacy
4. **Test First**: Test with small amounts before large transactions
5. **Keep Software Updated**: Use the latest version for security fixes

## API Reference

### Functions

#### `GenerateBitcoinStyleAddress() (*BitcoinStyleAddress, error)`

Generates a new Bitcoin-style address with "Btu" prefix.

**Returns:**
- `BitcoinStyleAddress`: Contains Address, PublicKey, and PrivateKey
- `error`: Error if generation fails

#### `GenerateEthereumStyleAddress() (*EthereumStyleAddress, error)`

Generates a new Ethereum-style address with "0x" prefix.

**Returns:**
- `EthereumStyleAddress`: Contains Address, PublicKey, and PrivateKey
- `error`: Error if generation fails

#### `ValidateBitcoinStyleAddress(address string) error`

Validates a Bitcoin-style address.

**Parameters:**
- `address`: The address string to validate

**Returns:**
- `error`: nil if valid, error describing the issue if invalid

#### `ValidateEthereumStyleAddress(address string) error`

Validates an Ethereum-style address.

**Parameters:**
- `address`: The address string to validate

**Returns:**
- `error`: nil if valid, error describing the issue if invalid

### Types

```go
type BitcoinStyleAddress struct {
    PrivateKey string // Hexadecimal encoded private key
    PublicKey  string // Hexadecimal encoded public key
    Address    string // Base58Check encoded address with Btu prefix
}

type EthereumStyleAddress struct {
    PrivateKey string // Hexadecimal encoded private key
    PublicKey  string // Hexadecimal encoded public key
    Address    string // Hexadecimal address with 0x prefix
}
```

## Testing

Run the comprehensive test suite:

```bash
# Run all identity tests
go test ./identity -v

# Run specific test
go test ./identity -v -run TestGenerateBitcoinStyleAddress

# Run with coverage
go test ./identity -cover
```

## Future Enhancements

Potential improvements for future versions:

1. **Keccak-256 Hashing**: Use proper Keccak-256 for Ethereum addresses (currently using SHA-256)
2. **secp256k1 Curve**: Use secp256k1 curve for better Ethereum compatibility
3. **EIP-55 Checksums**: Implement EIP-55 checksum encoding for Ethereum addresses
4. **HD Wallets**: Support BIP-32/BIP-44 hierarchical deterministic wallets
5. **Multi-Signature**: Add support for multi-signature addresses
6. **Vanity Addresses**: Generate custom vanity addresses

## Support

For issues, questions, or contributions:
- **GitHub Issues**: https://github.com/Bituncoin/Bituncoin/issues
- **Documentation**: See README.md in the repository root
- **Community**: Join our Discord server

---

**Note**: This is a reference implementation. For production use, ensure proper security audits and use industry-standard cryptographic libraries.
