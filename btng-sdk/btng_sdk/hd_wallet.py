#!/usr/bin/env python3
"""
BTNG HD Wallet Module
=====================

Implements BIP32/BIP44-style hierarchical deterministic wallets.
Provides master seed generation, extended keys, and derivation paths.
"""

import os
import hashlib
import hmac
import ecdsa
from ecdsa import SECP256k1
from typing import Optional, Tuple, List
import btng_sdk

# BTNG coin type for BIP44
BTNG_COIN_TYPE = 0x80000000 | 1234  # 1234 is placeholder, should be assigned by SLIP44

class ExtendedKey:
    """BIP32 extended key (xpriv/xpub)"""

    def __init__(self, key: bytes, chain_code: bytes, depth: int = 0,
                 parent_fingerprint: bytes = b'\x00\x00\x00\x00',
                 child_index: int = 0, is_private: bool = True):
        self.key = key
        self.chain_code = chain_code
        self.depth = depth
        self.parent_fingerprint = parent_fingerprint
        self.child_index = child_index
        self.is_private = is_private

        if is_private:
            # Generate public key from private key
            sk = ecdsa.SigningKey.from_secret_exponent(int.from_bytes(key, 'big'), curve=SECP256k1)
            vk = sk.verifying_key
            self.public_key = b'\x04' + vk.to_string()  # Uncompressed
        else:
            self.public_key = key

    @classmethod
    def from_seed(cls, seed: bytes) -> 'ExtendedKey':
        """Create master extended key from seed"""
        # HMAC-SHA512 with "Bitcoin seed" as key
        hmac_key = b"Bitcoin seed"
        I = hmac.new(hmac_key, seed, hashlib.sha512).digest()

        master_key = I[:32]
        master_chain_code = I[32:]

        return cls(master_key, master_chain_code, is_private=True)

    def derive(self, index: int) -> 'ExtendedKey':
        """Derive child extended key"""
        if index < 0:
            raise ValueError("Index must be non-negative")

        # Serialize key for HMAC
        if self.is_private:
            key_bytes = b'\x00' + self.key  # Add 0x00 prefix for private keys
        else:
            key_bytes = self.public_key

        # Create data for HMAC
        data = key_bytes + index.to_bytes(4, 'big')

        # HMAC-SHA512
        I = hmac.new(self.chain_code, data, hashlib.sha512).digest()
        child_key = I[:32]
        child_chain_code = I[32:]

        # Calculate child key
        if self.is_private:
            # Add parent private key to child key (mod n)
            parent_int = int.from_bytes(self.key, 'big')
            child_int = int.from_bytes(child_key, 'big')
            combined = (parent_int + child_int) % SECP256k1.order
            final_key = combined.to_bytes(32, 'big')
        else:
            # Add to public key (EC point addition)
            # This is simplified - real implementation needs proper EC math
            final_key = child_key  # Placeholder

        # Calculate parent fingerprint
        parent_fingerprint = hashlib.new('ripemd160', hashlib.sha256(self.public_key).digest()).digest()[:4]

        return ExtendedKey(
            final_key, child_chain_code,
            depth=self.depth + 1,
            parent_fingerprint=parent_fingerprint,
            child_index=index,
            is_private=self.is_private
        )

    def derive_path(self, path: str) -> 'ExtendedKey':
        """Derive key from BIP32 path"""
        current = self

        for part in path.split('/'):
            if part == 'm':
                continue
            elif part.endswith("'") or part.endswith("h"):
                # Hardened derivation
                index = int(part[:-1]) + 0x80000000
            else:
                # Normal derivation
                index = int(part)

            current = current.derive(index)

        return current

    def to_xpriv(self) -> str:
        """Export as extended private key (xpriv)"""
        if not self.is_private:
            raise ValueError("Cannot export private key from public key")

        version = b'\x04\x88\xAD\xE4'  # Mainnet xpriv
        depth = bytes([self.depth])
        fingerprint = self.parent_fingerprint
        child_index = self.child_index.to_bytes(4, 'big')
        chain_code = self.chain_code
        key = b'\x00' + self.key  # Private key prefix

        data = version + depth + fingerprint + child_index + chain_code + key

        # Base58Check encoding
        from btng_sdk.encoding import base58check_encode
        return base58check_encode(data)

    def to_xpub(self) -> str:
        """Export as extended public key (xpub)"""
        version = b'\x04\x88\xB2\x1E'  # Mainnet xpub
        depth = bytes([self.depth])
        fingerprint = self.parent_fingerprint
        child_index = self.child_index.to_bytes(4, 'big')
        chain_code = self.chain_code
        key = self.public_key

        data = version + depth + fingerprint + child_index + chain_code + key

        # Base58Check encoding
        from btng_sdk.encoding import base58check_encode
        return base58check_encode(data)

    @classmethod
    def from_xpriv(cls, xpriv: str) -> 'ExtendedKey':
        """Import from extended private key"""
        from btng_sdk.encoding import base58check_decode
        data = base58check_decode(xpriv)

        if len(data) != 78:
            raise ValueError("Invalid xpriv length")

        version = data[:4]
        depth = data[4]
        fingerprint = data[5:9]
        child_index = int.from_bytes(data[9:13], 'big')
        chain_code = data[13:45]
        key = data[46:78]  # Skip the 0x00 prefix

        return cls(key, chain_code, depth, fingerprint, child_index, is_private=True)

    def to_keypair(self) -> btng_sdk.KeyPair:
        """Convert to BTNG KeyPair"""
        if not self.is_private:
            raise ValueError("Cannot create KeyPair from public key")

        return btng_sdk.KeyPair.from_private_key(self.key)

class HDWallet:
    """BIP44 HD wallet for BTNG"""

    def __init__(self, seed: Optional[bytes] = None):
        if seed is None:
            seed = os.urandom(64)  # 512-bit seed

        self.master_key = ExtendedKey.from_seed(seed)
        self.seed = seed

    @classmethod
    def from_mnemonic(cls, mnemonic: str, passphrase: str = "") -> 'HDWallet':
        """Create wallet from BIP39 mnemonic"""
        # This would implement BIP39 - simplified for now
        seed = hashlib.pbkdf2_hmac('sha512', mnemonic.encode(), b'mnemonic' + passphrase.encode(), 2048)
        return cls(seed[:64])

    @classmethod
    def from_seed_hex(cls, seed_hex: str) -> 'HDWallet':
        """Create wallet from hex seed"""
        seed = bytes.fromhex(seed_hex)
        return cls(seed)

    def get_account_key(self, account: int = 0, is_change: bool = False) -> ExtendedKey:
        """Get account extended key (BIP44 path: m/44'/coin'/account')"""
        path = f"m/44'/{BTNG_COIN_TYPE >> 1}'/{account}'"
        return self.master_key.derive_path(path)

    def get_address_key(self, account: int = 0, is_change: bool = False, index: int = 0) -> ExtendedKey:
        """Get address extended key (BIP44 path: m/44'/coin'/account'/change/index)"""
        change = 1 if is_change else 0
        path = f"m/44'/{BTNG_COIN_TYPE >> 1}'/{account}'/{change}/{index}"
        return self.master_key.derive_path(path)

    def get_keypair(self, account: int = 0, is_change: bool = False, index: int = 0) -> btng_sdk.KeyPair:
        """Get KeyPair for specific address"""
        ext_key = self.get_address_key(account, is_change, index)
        return ext_key.to_keypair()

    def get_address(self, account: int = 0, is_change: bool = False, index: int = 0) -> str:
        """Get BTNG address for specific derivation path"""
        keypair = self.get_keypair(account, is_change, index)
        return btng_sdk.Address.from_public_key(keypair.public_key, 'p2pkh')

    def derive_range(self, account: int = 0, is_change: bool = False, start: int = 0, count: int = 20) -> List[Tuple[str, btng_sdk.KeyPair]]:
        """Derive a range of addresses and keypairs"""
        addresses = []
        for i in range(start, start + count):
            keypair = self.get_keypair(account, is_change, i)
            address = self.get_address(account, is_change, i)
            addresses.append((address, keypair))
        return addresses

    def export_master_xpriv(self) -> str:
        """Export master extended private key"""
        return self.master_key.to_xpriv()

    def export_account_xpub(self, account: int = 0) -> str:
        """Export account extended public key (for watch-only wallets)"""
        account_key = self.get_account_key(account)
        # Convert to public key
        pub_key = ExtendedKey(
            account_key.public_key, account_key.chain_code,
            account_key.depth, account_key.parent_fingerprint,
            account_key.child_index, is_private=False
        )
        return pub_key.to_xpub()

# Utility functions
def generate_mnemonic() -> str:
    """Generate BIP39 mnemonic (simplified)"""
    # This is a simplified implementation
    # Real BIP39 uses a wordlist and proper entropy
    words = ["abandon", "ability", "able", "about", "above", "absent", "absorb", "abstract",
             "absurd", "abuse", "access", "accident", "account", "accuse", "achieve", "acid"]
    return " ".join(words[:12])

def mnemonic_to_seed(mnemonic: str, passphrase: str = "") -> bytes:
    """Convert mnemonic to seed (simplified BIP39)"""
    return hashlib.pbkdf2_hmac('sha512', mnemonic.encode(), b'mnemonic' + passphrase.encode(), 2048)[:64]

if __name__ == "__main__":
    # Test HD wallet
    wallet = HDWallet()

    # Get first address
    address = wallet.get_address(0, False, 0)
    print(f"First address: {address}")

    # Get account xpub
    xpub = wallet.export_account_xpub(0)
    print(f"Account xpub: {xpub}")

    # Derive multiple addresses
    addresses = wallet.derive_range(0, False, 0, 5)
    print("First 5 addresses:")
    for addr, _ in addresses:
        print(f"  {addr}")