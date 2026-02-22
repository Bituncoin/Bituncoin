# BTNG Sovereign Blockchain SDK

A unified Python SDK for BTNG blockchain operations, providing cryptographic primitives, transaction building, wallet functionality, and sovereign identity integration.

## Features

- **Cryptographic Operations**: Private/public key generation, ECDSA signing, WIF encoding
- **Address Generation**: P2PKH, P2SH-P2WPKH, and SegWit address creation
- **Transaction Building**: Low-level transaction construction and signing
- **Wallet Management**: Key storage, address generation, transaction creation
- **Sovereign Identity**: Proof-of-value claims, trust union membership, sovereign signing
- **Hashing Functions**: SHA256, RIPEMD160, Hash160, Hash256
- **Encoding Utilities**: Base58, Base58Check, endian conversion

## Installation

```bash
pip install btng-sdk
```

Or from source:

```bash
git clone https://github.com/btng-sovereign/btng-sdk.git
cd btng-sdk
pip install -e .
```

## CLI Wallet

The SDK includes a command-line wallet for easy interaction:

```bash
# Create new wallet
python btng_wallet.py create -o my_wallet.json

# Check balance
python btng_wallet.py balance my_wallet.json

# Create sovereign identity
python btng_wallet.py identity my_wallet.json --name "Your Name" --country "Your Country"

# Sign a message
python btng_wallet.py sign my_wallet.json "Your message here"

# Verify a message signature
python btng_wallet.py verify <wif> "message" <signature>
```

## 🚀 Full Activation: Live Blockchain Integration

The SDK now supports **full activation** with live blockchain integration:

### Node RPC Client
```python
from btng_sdk import BTNGNodeClient, SyncBTNGWallet

# Connect to BTNG node
node = BTNGNodeClient(host="localhost", port=3002, mock_mode=True)

# Get blockchain info
info = await node.get_blockchain_info()

# Get balance
balance = await node.get_balance("your-address")

# Get transactions
txs = await node.get_transactions("your-address")
```

### Live Wallet Operations
```python
# Create live wallet
wallet = SyncBTNGWallet(node)

# Load your keys
keypair = btng_sdk.KeyPair.from_wif("your-wif-key")
wallet.keypairs.append(keypair)

# Send transaction
txid = wallet.send_transaction("recipient-address", 10.5)
print(f"Transaction sent: {txid}")
```

### CLI Wallet Commands
```bash
# Send BTNG
python btng_wallet.py send wallet.json RECIPIENT_ADDRESS 10.5

# Check balance (via node)
python btng_wallet.py balance wallet.json

# Create sovereign identity
python btng_wallet.py identity wallet.json --name "Your Name"
```

### CLI Commands

- `create`: Generate new wallet with keypair
- `balance`: Check wallet balance (mock implementation)
- `identity`: Create sovereign identity with claims
- `sign`: Sign messages with wallet keys
- `verify`: Verify message signatures
```

### Sovereign Identity

```python
from btng_sdk import SovereignIdentity

identity = SovereignIdentity()

# Create proof-of-value claim
value_proof = identity.create_value_proof(
    value_type="work",
    value_amount=100.0,
    metadata={"project": "BTNG Platform"}
)

# Create trust union membership
membership = identity.create_trust_union_membership(
    union_id="african-union",
    member_level="sovereign"
)

# Sign sovereign transaction
sovereign_tx = identity.sign_sovereign_transaction({
    "amount": 50.0,
    "recipient": "BTNG1RECIPIENTADDRESS"
})
```

## Architecture

The SDK is organized into the following modules:

- `keys.py`: Key pair generation and WIF encoding
- `addresses.py`: Address generation and validation
- `transactions.py`: Transaction building and signing
- `wallet.py`: Wallet management and operations
- `sovereign.py`: Sovereign identity and proof-of-value
- `hashing.py`: Cryptographic hash functions
- `encoding.py`: Base58, endian conversion, utilities

## Security

This SDK implements Bitcoin-standard cryptographic operations and is designed for production use. However, always:

- Keep private keys secure and never expose them
- Use hardware security modules for production keys
- Verify all operations in a test environment first
- Keep dependencies updated

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions welcome! Please see CONTRIBUTING.md for guidelines.

## Support

For support, please open an issue on GitHub or contact the BTNG development team.