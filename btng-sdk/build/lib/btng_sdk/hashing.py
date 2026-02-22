"""
BTNG Hashing Module
===================

Provides SHA256, RIPEMD160, Hash160, and Hash256 functions.
Combines functionality from sha256.py, ripemd160.py, hash160.py, hash256.py.
"""

import hashlib
import hmac


def sha256(data):
    """
    SHA256 hash

    Args:
        data (bytes): Data to hash

    Returns:
        bytes: SHA256 digest
    """
    return hashlib.sha256(data).digest()


def ripemd160(data):
    """
    RIPEMD160 hash

    Args:
        data (bytes): Data to hash

    Returns:
        bytes: RIPEMD160 digest
    """
    h = hashlib.new('ripemd160')
    h.update(data)
    return h.digest()


def hash160(data):
    """
    Hash160 = RIPEMD160(SHA256(data))

    Args:
        data (bytes): Data to hash

    Returns:
        bytes: 20-byte hash160
    """
    return ripemd160(sha256(data))


def hash256(data):
    """
    Hash256 = SHA256(SHA256(data))

    Args:
        data (bytes): Data to hash

    Returns:
        bytes: 32-byte double SHA256
    """
    return sha256(sha256(data))


def hmac_sha512(key, data):
    """
    HMAC-SHA512

    Args:
        key (bytes): HMAC key
        data (bytes): Data to hash

    Returns:
        bytes: HMAC-SHA512 digest
    """
    return hmac.new(key, data, hashlib.sha512).digest()