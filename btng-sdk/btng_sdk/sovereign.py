"""
BTNG Sovereign Identity Module
==============================

Provides sovereign identity signing and verification for BTNG platform.
Integrates with the sovereign identity and value protocols.
"""

import ecdsa
from .keys import KeyPair
from .hashing import hash256
from .encoding import message_prefix


class SovereignIdentity:
    """BTNG sovereign identity management"""

    def __init__(self, keypair=None):
        """
        Initialize sovereign identity

        Args:
            keypair (KeyPair, optional): Key pair for identity
        """
        self.keypair = keypair or KeyPair.generate()
        self.claims = {}

    def create_identity_claim(self, claim_type, claim_data):
        """
        Create a signed identity claim

        Args:
            claim_type (str): Type of claim (e.g., 'citizenship', 'proof-of-value')
            claim_data (dict): Claim data

        Returns:
            dict: Signed claim
        """
        claim = {
            'type': claim_type,
            'data': claim_data,
            'timestamp': self._current_timestamp(),
            'public_key': self.keypair.public_key_hex
        }

        # Create claim hash
        claim_json = str(sorted(claim.items())).encode()
        claim_hash = hash256(claim_json)

        # Sign claim
        signature = self.keypair.sign(claim_hash)

        claim['signature'] = signature.hex()
        claim['claim_hash'] = claim_hash.hex()

        self.claims[claim_type] = claim
        return claim

    def verify_identity_claim(self, claim):
        """
        Verify an identity claim

        Args:
            claim (dict): Claim to verify

        Returns:
            bool: True if claim is valid
        """
        try:
            # Recreate claim hash
            claim_copy = claim.copy()
            signature = claim_copy.pop('signature')
            claim_hash = claim_copy.pop('claim_hash')

            claim_json = str(sorted(claim_copy.items())).encode()
            expected_hash = hash256(claim_json)

            if expected_hash.hex() != claim_hash:
                return False

            # Verify signature
            public_key = bytes.fromhex(claim['public_key'])
            signature_bytes = bytes.fromhex(signature)

            # Create a dummy keypair with the public key for verification
            # In practice, you'd store/verify against known identity keys
            return self._verify_signature(expected_hash, signature_bytes, public_key)

        except:
            return False

    def create_value_proof(self, value_type, value_amount, metadata=None):
        """
        Create a proof-of-value claim

        Args:
            value_type (str): Type of value (e.g., 'work', 'trade', 'trust')
            value_amount (float): Value amount
            metadata (dict, optional): Additional metadata

        Returns:
            dict: Signed value proof
        """
        proof_data = {
            'value_type': value_type,
            'amount': value_amount,
            'metadata': metadata or {},
            'timestamp': self._current_timestamp()
        }

        return self.create_identity_claim('proof-of-value', proof_data)

    def create_trust_union_membership(self, union_id, member_level):
        """
        Create trust union membership claim

        Args:
            union_id (str): Trust union identifier
            member_level (str): Membership level

        Returns:
            dict: Signed membership claim
        """
        membership_data = {
            'union_id': union_id,
            'member_level': member_level,
            'sovereign_id': self.keypair.public_key_hex[:16]  # Short ID
        }

        return self.create_identity_claim('trust-union-membership', membership_data)

    def sign_sovereign_transaction(self, transaction_data):
        """
        Sign a transaction with sovereign identity

        Args:
            transaction_data (dict): Transaction data

        Returns:
            dict: Signed sovereign transaction
        """
        # Add sovereign metadata
        sovereign_tx = {
            'transaction': transaction_data,
            'sovereign_signer': self.keypair.public_key_hex,
            'sovereign_timestamp': self._current_timestamp(),
            'sovereign_domain': 'btng.africa'
        }

        # Sign the entire package
        tx_json = str(sorted(sovereign_tx.items())).encode()
        tx_hash = hash256(tx_json)
        signature = self.keypair.sign(tx_hash)

        sovereign_tx['sovereign_signature'] = signature.hex()
        sovereign_tx['sovereign_hash'] = tx_hash.hex()

        return sovereign_tx

    def _verify_signature(self, message_hash, signature, public_key):
        """Verify signature with public key"""
        try:
            # This is a simplified verification
            # In practice, use proper ECDSA verification
            vk = ecdsa.VerifyingKey.from_string(public_key[1:], curve=ecdsa.SECP256k1)
            return vk.verify_digest(signature, message_hash, sigdecode=ecdsa.util.sigdecode_der)
        except:
            return False

    def _current_timestamp(self):
        """Get current timestamp"""
        import time
        return int(time.time())

    @property
    def identity_address(self):
        """Get identity address"""
        from .addresses import Address
        return Address.from_public_key(self.keypair.public_key, 'p2pkh')

    def export_identity(self):
        """Export identity data"""
        return {
            'public_key': self.keypair.public_key_hex,
            'identity_address': self.identity_address,
            'claims': self.claims
        }