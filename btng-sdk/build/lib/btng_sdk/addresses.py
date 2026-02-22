"""
BTNG Address Generation Module
==============================

Handles P2PKH, P2SH-P2WPKH, and other address formats.
Combines functionality from p2pkh.py, p2shp2wpkh.py, and address utilities.
"""

import hashlib
from .hashing import hash160, sha256
from .encoding import base58check_encode


class Address:
    """BTNG address generation and validation"""

    @staticmethod
    def from_public_key(public_key, address_type='p2pkh', network='mainnet'):
        """
        Generate address from public key

        Args:
            public_key (bytes): Public key bytes
            address_type (str): 'p2pkh', 'p2wpkh', 'p2sh-p2wpkh'
            network (str): 'mainnet' or 'testnet'

        Returns:
            str: Base58Check encoded address
        """
        if address_type == 'p2pkh':
            return Address._p2pkh_address(public_key, network)
        elif address_type == 'p2wpkh':
            return Address._p2wpkh_address(public_key, network)
        elif address_type == 'p2sh-p2wpkh':
            return Address._p2sh_p2wpkh_address(public_key, network)
        else:
            raise ValueError(f"Unsupported address type: {address_type}")

    @staticmethod
    def _p2pkh_address(public_key, network='mainnet'):
        """Generate P2PKH address"""
        # Hash160 of public key
        pubkey_hash = hash160(public_key)

        # Version byte (0x00 for mainnet BTC, using 0x1B for BTNG)
        version = b'\x1B' if network == 'mainnet' else b'\x6F'

        return base58check_encode(version + pubkey_hash)

    @staticmethod
    def _p2wpkh_address(public_key, network='mainnet'):
        """Generate P2WPKH (SegWit v0) address"""
        # Hash160 of public key
        pubkey_hash = hash160(public_key)

        # Bech32 encoding would be needed here, but using simplified version
        # For now, return hex representation
        return 'bc1' + pubkey_hash.hex()[:40]  # Simplified

    @staticmethod
    def _p2sh_p2wpkh_address(public_key, network='mainnet'):
        """Generate P2SH-P2WPKH address (SegWit wrapped in P2SH)"""
        # Create the redeem script
        pubkey_hash = hash160(public_key)
        redeem_script = b'\x00\x14' + pubkey_hash  # OP_0 <20-byte hash>

        # Hash160 of redeem script
        script_hash = hash160(redeem_script)

        # Version byte for P2SH (0x05 for mainnet BTC, using 0x3A for BTNG)
        version = b'\x3A' if network == 'mainnet' else b'\xC4'

        return base58check_encode(version + script_hash)

    @staticmethod
    def validate(address):
        """
        Validate BTNG address format

        Args:
            address (str): Address to validate

        Returns:
            bool: True if valid
        """
        try:
            decoded = base58check_encode.decode(address)
            return len(decoded) == 21  # Version + 20-byte hash
        except:
            return False

    @staticmethod
    def from_hash160(hash160_bytes, address_type='p2pkh', network='mainnet'):
        """
        Generate address from hash160

        Args:
            hash160_bytes (bytes): 20-byte hash160
            address_type (str): Address type
            network (str): Network

        Returns:
            str: Address
        """
        if address_type == 'p2pkh':
            version = b'\x1B' if network == 'mainnet' else b'\x6F'
        elif address_type == 'p2sh':
            version = b'\x3A' if network == 'mainnet' else b'\xC4'
        else:
            raise ValueError(f"Unsupported address type: {address_type}")

        return base58check_encode(version + hash160_bytes)