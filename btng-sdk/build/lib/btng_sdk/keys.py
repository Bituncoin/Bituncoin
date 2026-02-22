"""
BTNG Key Management Module
==========================

Handles private/public key generation, WIF encoding, and cryptographic operations.
Combines functionality from pubkey.py, wif.py, and core key utilities.
"""

import os
import hashlib
import ecdsa
from ecdsa import SECP256k1
import base58


class KeyPair:
    """Represents a BTNG private/public key pair"""

    def __init__(self, private_key=None):
        """
        Initialize with existing private key or generate new one

        Args:
            private_key (bytes, optional): 32-byte private key
        """
        if private_key:
            if len(private_key) != 32:
                raise ValueError("Private key must be 32 bytes")
            self.private_key = private_key
        else:
            self.private_key = os.urandom(32)

        # Generate public key
        sk = ecdsa.SigningKey.from_secret_exponent(
            int.from_bytes(self.private_key, 'big'),
            curve=SECP256k1
        )
        vk = sk.verifying_key
        self.public_key = b'\x04' + vk.to_string()  # Uncompressed format

    @classmethod
    def generate(cls):
        """Generate a new random keypair"""
        return cls()

    @classmethod
    def from_private_key(cls, private_key):
        """Create KeyPair from existing private key"""
        return cls(private_key)

    @classmethod
    def from_wif(cls, wif_string):
        """Import from Wallet Import Format"""
        # Decode base58
        decoded = base58.b58decode(wif_string)

        # Remove checksum (last 4 bytes)
        if len(decoded) == 38:  # Compressed WIF
            key_data = decoded[1:-5]
            compressed = True
        elif len(decoded) == 37:  # Uncompressed WIF
            key_data = decoded[1:-4]
            compressed = False
        else:
            raise ValueError("Invalid WIF format")

        return cls(key_data)

    def to_wif(self, compressed=True):
        """Export to Wallet Import Format"""
        # Add version byte (0x80 for mainnet)
        versioned_key = b'\x80' + self.private_key

        if compressed:
            versioned_key += b'\x01'

        # Calculate checksum
        checksum = hashlib.sha256(hashlib.sha256(versioned_key).digest()).digest()[:4]

        return base58.b58encode(versioned_key + checksum).decode()

    def sign(self, message):
        """
        Sign a message using ECDSA

        Args:
            message (bytes): Message to sign

        Returns:
            bytes: DER-encoded signature
        """
        sk = ecdsa.SigningKey.from_secret_exponent(
            int.from_bytes(self.private_key, 'big'),
            curve=SECP256k1
        )

        # Add message prefix for Bitcoin-style signing
        prefixed_message = b'\x18Bitcoin Signed Message:\n' + bytes([len(message)]) + message

        # Double SHA256
        message_hash = hashlib.sha256(hashlib.sha256(prefixed_message).digest()).digest()

        signature = sk.sign_digest(message_hash, sigencode=ecdsa.util.sigencode_der)
        return signature

    def verify(self, message, signature):
        """
        Verify a signature

        Args:
            message (bytes): Original message
            signature (bytes): DER-encoded signature

        Returns:
            bool: True if signature is valid
        """
        try:
            sk = ecdsa.SigningKey.from_secret_exponent(
                int.from_bytes(self.private_key, 'big'),
                curve=SECP256k1
            )
            vk = sk.verifying_key

            prefixed_message = b'\x18Bitcoin Signed Message:\n' + bytes([len(message)]) + message
            message_hash = hashlib.sha256(hashlib.sha256(prefixed_message).digest()).digest()

            return vk.verify_digest(signature, message_hash, sigdecode=ecdsa.util.sigdecode_der)
        except:
            return False

    @property
    def public_key_hex(self):
        """Get public key as hex string"""
        return self.public_key.hex()

    @property
    def private_key_hex(self):
        """Get private key as hex string"""
        return self.private_key.hex()