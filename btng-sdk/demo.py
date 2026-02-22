#!/usr/bin/env python3
"""
BTNG SDK Demo Script
====================

Demonstrates basic usage of the BTNG Sovereign Blockchain SDK.
"""

import sys
import os
sys.path.insert(0, os.path.dirname(__file__))

import keys
import addresses
import wallet
import sovereign

KeyPair = keys.KeyPair
Address = addresses.Address
Wallet = wallet.Wallet
SovereignIdentity = sovereign.SovereignIdentity

def main():
    print("🚀 BTNG Sovereign Blockchain SDK Demo")
    print("=" * 50)

    # 1. Key Generation
    print("\n1. 🔑 Key Generation")
    keys = KeyPair.generate()
    print(f"Private Key (WIF): {keys.to_wif()}")
    print(f"Public Key: {keys.public_key_hex}")

    # 2. Address Generation
    print("\n2. 🏠 Address Generation")
    p2pkh_address = Address.from_public_key(keys.public_key, 'p2pkh')
    print(f"P2PKH Address: {p2pkh_address}")

    p2sh_address = Address.from_public_key(keys.public_key, 'p2sh-p2wpkh')
    print(f"P2SH-P2WPKH Address: {p2sh_address}")

    # 3. Wallet Operations
    print("\n3. 👛 Wallet Operations")
    wallet = Wallet()

    # Generate new address
    new_addr = wallet.generate_keypair()
    print(f"New Wallet Address: {new_addr['address']}")

    # Mock balance check
    balance = wallet.get_balance(new_addr['address'])
    print(f"Wallet Balance: {balance} BTNG")

    # 4. Message Signing
    print("\n4. ✍️  Message Signing")
    message = "BTNG Sovereign Platform - Building Trust. Nurturing Growth."
    signature = wallet.sign_message(message, new_addr['address'])
    print(f"Message: {message}")
    print(f"Signature: {signature}")

    # Verify signature
    is_valid = wallet.verify_message(message, signature, new_addr['address'])
    print(f"Signature Valid: {is_valid}")

    # 5. Sovereign Identity
    print("\n5. 🏛️  Sovereign Identity")
    identity = SovereignIdentity()

    # Create value proof
    value_proof = identity.create_value_proof(
        value_type="trust",
        value_amount=1000.0,
        metadata={"platform": "BTNG", "verified": True}
    )
    print(f"Value Proof Created: {value_proof['claim_hash'][:16]}...")

    # Create trust union membership
    membership = identity.create_trust_union_membership(
        union_id="african-sovereign-union",
        member_level="founding-member"
    )
    print(f"Trust Union Membership: {membership['claim_hash'][:16]}...")

    # Sign sovereign transaction
    sovereign_tx = identity.sign_sovereign_transaction({
        "type": "sovereign-transfer",
        "amount": 100.0,
        "purpose": "infrastructure-development"
    })
    print(f"Sovereign Transaction: {sovereign_tx['sovereign_hash'][:16]}...")

    print("\n✅ BTNG SDK Demo Complete!")
    print("\nNext steps:")
    print("- Integrate with BTNG node API")
    print("- Build full wallet application")
    print("- Implement sovereign identity verification")
    print("- Deploy to production environment")

if __name__ == "__main__":
    main()