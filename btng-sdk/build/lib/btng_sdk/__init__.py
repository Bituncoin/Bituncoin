"""
BTNG Sovereign Blockchain SDK
============================

A unified Python SDK for BTNG blockchain operations, providing:
- Cryptographic primitives (keys, addresses, signatures)
- Transaction building and validation
- Wallet functionality
- Sovereign identity integration

This SDK combines the core blockchain toolkit modules into a cohesive package
for building BTNG applications, wallets, and node integrations.

Usage:
    from btng_sdk import KeyPair, Address, Transaction

    # Generate a new keypair
    keys = KeyPair.generate()

    # Create P2PKH address
    address = Address.from_public_key(keys.public_key, 'p2pkh')

    # Build a transaction
    tx = Transaction()
    tx.add_input(...)
    tx.add_output(...)
    signed_tx = tx.sign(keys)
"""

__version__ = "0.1.0"
__author__ = "BTNG Sovereign Platform"

from .keys import KeyPair
from .addresses import Address
from .transactions import Transaction
from .wallet import Wallet
from .sovereign import SovereignIdentity
from .node_client import BTNGNodeClient, BTNGWallet, SyncBTNGWallet

__all__ = [
    'KeyPair',
    'Address',
    'Transaction',
    'Wallet',
    'SovereignIdentity',
    'BTNGNodeClient',
    'BTNGWallet',
    'SyncBTNGWallet'
]