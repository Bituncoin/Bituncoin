#!/usr/bin/env python3
"""
BTNG CLI Wallet
================

Command-line interface for BTNG Sovereign Blockchain wallet operations.
"""

import argparse
import sys
import json
import btng_sdk

Wallet = btng_sdk.Wallet
KeyPair = btng_sdk.KeyPair
SovereignIdentity = btng_sdk.SovereignIdentity

def create_wallet(args):
    """Create a new wallet"""
    print("🔑 Creating new BTNG wallet...")

    wallet = Wallet()
    key_info = wallet.generate_keypair()
    keypair = wallet.keypairs[0]  # Get the generated keypair

    # Save wallet data (in production, this should be encrypted)
    wallet_data = {
        'wif': key_info['private_key'],
        'address': key_info['address'],
        'public_key': key_info['public_key']
    }

    if args.output:
        with open(args.output, 'w') as f:
            json.dump(wallet_data, f, indent=2)
        print(f"💾 Wallet saved to {args.output}")

    print(f"🏠 Address: {wallet_data['address']}")
    print(f"🔐 Private Key (WIF): {wallet_data['wif']}")
    print("⚠️  WARNING: Keep your private key secure and never share it!")

def load_wallet(args):
    """Load wallet from file"""
    try:
        with open(args.file, 'r') as f:
            wallet_data = json.load(f)

        wallet = Wallet()

        # Check if this is an HD wallet
        if 'seed' in wallet_data:
            # HD wallet - recreate from seed
            hd_wallet = btng_sdk.HDWallet.from_seed_hex(wallet_data['seed'])
            # Use first keypair for now (in production, you'd derive the needed key)
            keypair = hd_wallet.get_keypair(0, False, 0)
            wallet.keypairs.append(keypair)
        else:
            # Regular wallet
            wif = wallet_data.get('wif') or wallet_data.get('first_wif')
            if not wif:
                raise ValueError("No private key found in wallet file")
            wallet.import_private_key(wif)

        address = wallet_data.get('address') or wallet_data.get('first_address', 'Unknown')
        print(f"✅ Wallet loaded: {address}")
        return wallet
    except FileNotFoundError:
        print(f"❌ Wallet file not found: {args.file}")
        sys.exit(1)
    except Exception as e:
        print(f"❌ Error loading wallet: {e}")
        sys.exit(1)

def show_balance(args):
    """Show wallet balance"""
    wallet = load_wallet(args)
    # In a real implementation, this would query the BTNG node
    print("🏦 Wallet Balance: 0.000 BTNG (mock - integrate with BTNG node)")

def create_identity(args):
    """Create sovereign identity"""
    wallet = load_wallet(args)

    print("🏛️  Creating Sovereign Identity...")
    keypair = wallet.keypairs[0]  # Use first keypair
    identity = SovereignIdentity(keypair)

    # Create identity claims
    identity.create_identity_claim("personal_info", {
        "name": args.name or "BTNG User",
        "country": args.country or "Digital Sovereign"
    })
    identity.create_identity_claim("trust_union", {
        "union_name": "BTNG Network",
        "member_level": "sovereign_citizen"
    })

    print(f"✅ Identity created for: {args.name or 'BTNG User'}")
    print(f"🏛️  Trust Union: BTNG Network")

    # Create value proof
    value_proof = identity.create_value_proof("sovereign_token", 100, {"purpose": "identity_verification"})
    print(f"💎 Value Proof: {value_proof['claim_hash'][:16]}...")

def sign_message(args):
    """Sign a message"""
    wallet = load_wallet(args)

    keypair = wallet.keypairs[0]
    signature = keypair.sign(args.message.encode())

    print(f"✍️  Message: {args.message}")
    print(f"🔏 Signature: {signature.hex()}")

def send_transaction(args):
    """Send transaction to address"""
    wallet = load_wallet(args)

    print(f"💸 Sending {args.amount} BTNG to {args.to_address}...")

    try:
        # Create live wallet with node client (synchronous mode)
        node_client = btng_sdk.BTNGNodeClient(mock_mode=True)
        live_wallet = btng_sdk.SyncBTNGWallet(node_client)

        # Copy keypairs to live wallet
        for keypair in wallet.keypairs:
            live_wallet.keypairs.append(keypair)

        # Create transaction
        tx = live_wallet.create_transaction(args.to_address, float(args.amount))
        print(f"📝 Transaction created: {tx['txid']}")

        # Broadcast transaction
        txid = live_wallet.broadcast_transaction(tx["hex"])

        print(f"✅ Transaction broadcast!")
        print(f"🔗 TXID: {txid}")
        print(f"💰 Amount: {args.amount} BTNG")
        print(f"🏠 To: {args.to_address}")
        print(f"💸 Fee: {tx['fee']} BTNG")

    except Exception as e:
        print(f"❌ Transaction failed: {e}")
        import traceback
        traceback.print_exc()

def create_hd_wallet(args):
    """Create a new HD wallet"""
    print("🌱 Creating new HD wallet...")

    # Create HD wallet
    hd_wallet = btng_sdk.HDWallet()

    # Get first address
    first_address = hd_wallet.get_address(0, False, 0)
    first_keypair = hd_wallet.get_keypair(0, False, 0)

    # Save wallet data (in production, encrypt this!)
    wallet_data = {
        'seed': hd_wallet.seed.hex(),
        'master_xpriv': hd_wallet.export_master_xpriv(),
        'first_address': first_address,
        'first_wif': first_keypair.to_wif()
    }

    if args.output:
        with open(args.output, 'w') as f:
            json.dump(wallet_data, f, indent=2)
        print(f"💾 HD Wallet saved to {args.output}")

    print(f"🏠 First address: {first_address}")
    print(f"🔐 First private key (WIF): {first_keypair.to_wif()}")
    print("⚠️  WARNING: Keep your seed secure! Never share it!")
    print(f"🌱 Master xpriv: {wallet_data['master_xpriv'][:20]}...")

def derive_addresses(args):
    """Derive addresses from HD wallet"""
    try:
        with open(args.file, 'r') as f:
            wallet_data = json.load(f)

        # Recreate HD wallet from seed
        hd_wallet = btng_sdk.HDWallet.from_seed_hex(wallet_data['seed'])

        print(f"📔 Deriving addresses for account {args.account}, change={args.change}")

        # Derive addresses
        addresses = hd_wallet.derive_range(args.account, args.change, args.start, args.count)

        for i, (address, keypair) in enumerate(addresses):
            index = args.start + i
            print(f"{index:3d}: {address} (WIF: {keypair.to_wif()})")

    except FileNotFoundError:
        print(f"❌ Wallet file not found: {args.file}")
        sys.exit(1)

def main():
    parser = argparse.ArgumentParser(description="BTNG Sovereign Blockchain CLI Wallet")
    subparsers = parser.add_subparsers(dest='command', help='Available commands')

    # Create wallet
    create_parser = subparsers.add_parser('create', help='Create new wallet')
    create_parser.add_argument('-o', '--output', help='Save wallet to file')
    create_parser.set_defaults(func=create_wallet)

    # Load wallet
    load_parser = subparsers.add_parser('load', help='Load wallet from file')
    load_parser.add_argument('file', help='Wallet file path')
    load_parser.set_defaults(func=lambda args: print("Wallet loaded - use with other commands"))

    # Balance
    balance_parser = subparsers.add_parser('balance', help='Show wallet balance')
    balance_parser.add_argument('file', help='Wallet file path')
    balance_parser.set_defaults(func=show_balance)

    # Identity
    identity_parser = subparsers.add_parser('identity', help='Create sovereign identity')
    identity_parser.add_argument('file', help='Wallet file path')
    identity_parser.add_argument('--name', help='Identity name')
    identity_parser.add_argument('--country', help='Country')
    identity_parser.set_defaults(func=create_identity)

    # Sign message
    sign_parser = subparsers.add_parser('sign', help='Sign a message')
    sign_parser.add_argument('file', help='Wallet file path')
    sign_parser.add_argument('message', help='Message to sign')
    sign_parser.set_defaults(func=sign_message)

    # Send transaction
    send_parser = subparsers.add_parser('send', help='Send BTNG to address')
    send_parser.add_argument('file', help='Wallet file path')
    send_parser.add_argument('to_address', help='Recipient address')
    send_parser.add_argument('amount', help='Amount to send')
    send_parser.set_defaults(func=send_transaction)

    # Create HD wallet
    hd_create_parser = subparsers.add_parser('hd-create', help='Create new HD wallet')
    hd_create_parser.add_argument('-o', '--output', help='Save HD wallet to file')
    hd_create_parser.set_defaults(func=create_hd_wallet)

    # Derive addresses
    derive_parser = subparsers.add_parser('hd-derive', help='Derive addresses from HD wallet')
    derive_parser.add_argument('file', help='HD wallet file path')
    derive_parser.add_argument('-a', '--account', type=int, default=0, help='Account number')
    derive_parser.add_argument('-c', '--change', action='store_true', help='Use change addresses')
    derive_parser.add_argument('-s', '--start', type=int, default=0, help='Start index')
    derive_parser.add_argument('-n', '--count', type=int, default=10, help='Number of addresses')
    derive_parser.set_defaults(func=derive_addresses)

    args = parser.parse_args()

    if not args.command:
        parser.print_help()
        return

    args.func(args)

if __name__ == "__main__":
    main()