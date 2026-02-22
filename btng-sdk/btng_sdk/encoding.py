"""
BTNG Encoding Module
====================

Provides Base58, Base58Check encoding, endian conversion, and utilities.
Combines functionality from bs58.py, reverse_endian.py, msgprefixgen.py.
"""

import hashlib
from .hashing import sha256


# Base58 alphabet (Bitcoin standard)
BASE58_ALPHABET = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz'


def base58_encode(data):
    """
    Encode bytes to Base58 string

    Args:
        data (bytes): Data to encode

    Returns:
        str: Base58 encoded string
    """
    if not data:
        return ''

    # Count leading zeros
    leading_zeros = 0
    for byte in data:
        if byte == 0:
            leading_zeros += 1
        else:
            break

    # Convert to integer
    n = int.from_bytes(data, 'big')

    # Encode
    result = []
    while n > 0:
        n, remainder = divmod(n, 58)
        result.append(BASE58_ALPHABET[remainder])

    # Add leading '1's for leading zeros
    result.extend(['1'] * leading_zeros)
    result.reverse()

    return ''.join(result)


def base58_decode(encoded):
    """
    Decode Base58 string to bytes

    Args:
        encoded (str): Base58 encoded string

    Returns:
        bytes: Decoded data
    """
    if not encoded:
        return b''

    # Convert to integer
    n = 0
    for char in encoded:
        n = n * 58 + BASE58_ALPHABET.index(char)

    # Convert to bytes
    data = n.to_bytes((n.bit_length() + 7) // 8, 'big')

    # Add leading zeros
    leading_zeros = 0
    for char in encoded:
        if char == '1':
            leading_zeros += 1
        else:
            break

    return b'\x00' * leading_zeros + data


def base58check_encode(data):
    """
    Encode with Base58Check (data + checksum)

    Args:
        data (bytes): Data to encode

    Returns:
        str: Base58Check encoded string
    """
    # Calculate checksum
    checksum = sha256(sha256(data))[:4]

    # Append checksum
    data_with_checksum = data + checksum

    return base58_encode(data_with_checksum)


def base58check_decode(encoded):
    """
    Decode Base58Check string

    Args:
        encoded (str): Base58Check encoded string

    Returns:
        bytes: Decoded data (without checksum)

    Raises:
        ValueError: If checksum is invalid
    """
    decoded = base58_decode(encoded)

    if len(decoded) < 4:
        raise ValueError("Invalid Base58Check string")

    data = decoded[:-4]
    checksum = decoded[-4:]

    # Verify checksum
    expected_checksum = sha256(sha256(data))[:4]
    if checksum != expected_checksum:
        raise ValueError("Invalid checksum")

    return data


def reverse_endian(data):
    """
    Reverse byte order (endian conversion)

    Args:
        data (bytes): Data to reverse

    Returns:
        bytes: Reversed data
    """
    return data[::-1]


def int_to_bytes(n, length=None, byteorder='big'):
    """
    Convert integer to bytes

    Args:
        n (int): Integer to convert
        length (int, optional): Fixed length
        byteorder (str): 'big' or 'little'

    Returns:
        bytes: Byte representation
    """
    if length:
        return n.to_bytes(length, byteorder)
    else:
        return n.to_bytes((n.bit_length() + 7) // 8, byteorder)


def bytes_to_int(data, byteorder='big'):
    """
    Convert bytes to integer

    Args:
        data (bytes): Bytes to convert
        byteorder (str): 'big' or 'little'

    Returns:
        int: Integer representation
    """
    return int.from_bytes(data, byteorder)


def message_prefix(message):
    """
    Add Bitcoin-style message prefix

    Args:
        message (bytes): Message to prefix

    Returns:
        bytes: Prefixed message
    """
    prefix = b'\x18Bitcoin Signed Message:\n'
    length_byte = bytes([len(message)])
    return prefix + length_byte + message