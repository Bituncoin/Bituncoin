"""
BTNG Transaction Module
=======================

Handles transaction building, signing, and serialization.
Provides low-level transaction construction for BTNG blockchain.
"""

from .hashing import hash256
from .encoding import int_to_bytes, bytes_to_int, reverse_endian


class Transaction:
    """BTNG transaction builder"""

    def __init__(self, version=1):
        self.version = version
        self.inputs = []
        self.outputs = []
        self.locktime = 0

    def add_input(self, txid, vout, script_sig=b'', sequence=0xFFFFFFFF):
        """
        Add transaction input

        Args:
            txid (str): Previous transaction ID (hex)
            vout (int): Output index
            script_sig (bytes): Unlocking script
            sequence (int): Sequence number
        """
        self.inputs.append({
            'txid': txid,
            'vout': vout,
            'script_sig': script_sig,
            'sequence': sequence
        })

    def add_output(self, amount, script_pubkey):
        """
        Add transaction output

        Args:
            amount (int): Amount in satoshis
            script_pubkey (bytes): Locking script
        """
        self.outputs.append({
            'amount': amount,
            'script_pubkey': script_pubkey
        })

    def serialize(self):
        """
        Serialize transaction to bytes

        Returns:
            bytes: Serialized transaction
        """
        result = b''

        # Version
        result += int_to_bytes(self.version, 4, 'little')

        # Input count
        result += self._varint(len(self.inputs))

        # Inputs
        for inp in self.inputs:
            result += reverse_endian(bytes.fromhex(inp['txid']))
            result += int_to_bytes(inp['vout'], 4, 'little')
            result += self._varint(len(inp['script_sig']))
            result += inp['script_sig']
            result += int_to_bytes(inp['sequence'], 4, 'little')

        # Output count
        result += self._varint(len(self.outputs))

        # Outputs
        for out in self.outputs:
            result += int_to_bytes(out['amount'], 8, 'little')
            result += self._varint(len(out['script_pubkey']))
            result += out['script_pubkey']

        # Locktime
        result += int_to_bytes(self.locktime, 4, 'little')

        return result

    def txid(self):
        """
        Calculate transaction ID

        Returns:
            str: Transaction ID (hex)
        """
        serialized = self.serialize()
        return reverse_endian(hash256(serialized)).hex()

    def sign(self, keypair, input_index=0, sighash_type=0x01, script_pubkey=None):
        """
        Sign transaction input (P2PKH)

        Args:
            keypair (KeyPair): Key pair for signing
            input_index (int): Index of input to sign
            sighash_type (int): Sighash type (default: SIGHASH_ALL)
            script_pubkey (bytes): Script pubkey for the input

        Returns:
            bytes: Signed transaction
        """
        if input_index >= len(self.inputs):
            raise ValueError("Input index out of range")

        # Create sighash
        sighash = self._sighash(input_index, sighash_type, script_pubkey)

        # Sign the sighash
        signature = keypair.sign(sighash)

        # Add sighash type to signature
        signature += bytes([sighash_type])

        # Create script sig: <sig> <pubkey>
        script_sig = self._varint(len(signature)) + signature
        script_sig += self._varint(len(keypair.public_key)) + keypair.public_key

        self.inputs[input_index]['script_sig'] = script_sig

        return self.serialize()

    def _sighash(self, input_index, sighash_type, script_pubkey):
        """
        Create sighash for signing (BIP143 style for simplicity)

        Args:
            input_index (int): Index of input being signed
            sighash_type (int): Sighash type
            script_pubkey (bytes): Script pubkey for the input

        Returns:
            bytes: Sighash to sign
        """
        # For SIGHASH_ALL, we hash the transaction with all inputs and outputs
        # This is simplified - real implementation needs proper sighash handling

        # Create a copy of the transaction for sighash
        sighash_tx = Transaction()
        sighash_tx.version = self.version
        sighash_tx.locktime = self.locktime

        # Add all inputs
        for i, inp in enumerate(self.inputs):
            sighash_input = inp.copy()
            if i == input_index:
                # For the input being signed, use the script_pubkey
                sighash_input['script_sig'] = script_pubkey or b''
            else:
                # For other inputs, script_sig is empty for SIGHASH_ALL
                sighash_input['script_sig'] = b''
            sighash_tx.inputs.append(sighash_input)

        # Add all outputs
        sighash_tx.outputs = self.outputs.copy()

        # Serialize and double hash
        preimage = sighash_tx.serialize() + int_to_bytes(sighash_type, 4, 'little')
        return hash256(preimage)

    def _varint(self, n):
        """Encode variable length integer"""
        if n < 0xFD:
            return bytes([n])
        elif n <= 0xFFFF:
            return b'\xFD' + int_to_bytes(n, 2, 'little')
        elif n <= 0xFFFFFFFF:
            return b'\xFE' + int_to_bytes(n, 4, 'little')
        else:
            return b'\xFF' + int_to_bytes(n, 8, 'little')

    @classmethod
    def from_hex(cls, hex_string):
        """
        Parse transaction from hex string

        Args:
            hex_string (str): Hex-encoded transaction

        Returns:
            Transaction: Parsed transaction
        """
        # Simplified parser - would need full implementation
        raise NotImplementedError("Transaction parsing not implemented")


def create_p2pkh_script_pubkey(address):
    """
    Create P2PKH script pubkey from address

    Args:
        address (str): P2PKH address

    Returns:
        bytes: Script pubkey
    """
    # Decode address to get hash160
    from .encoding import base58check_decode
    decoded = base58check_decode(address)
    pubkey_hash = decoded[1:]  # Remove version byte

    # Create script: OP_DUP OP_HASH160 <pubkey_hash> OP_EQUALVERIFY OP_CHECKSIG
    script = b'\x76\xA9\x14' + pubkey_hash + b'\x88\xAC'
    return script