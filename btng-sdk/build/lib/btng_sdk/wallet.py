"""
BTNG Wallet Module
==================

Provides wallet functionality including key management, address generation,
and basic transaction operations.
"""

from .keys import KeyPair
from .addresses import Address
from .transactions import Transaction, create_p2pkh_script_pubkey


class Wallet:
    """BTNG wallet with key management and transaction capabilities"""

    def __init__(self):
        self.keypairs = []
        self.addresses = []

    def generate_keypair(self):
        """Generate a new keypair and add to wallet"""
        keypair = KeyPair.generate()
        self.keypairs.append(keypair)

        # Generate address
        address = Address.from_public_key(keypair.public_key, 'p2pkh')
        self.addresses.append(address)

        return {
            'address': address,
            'private_key': keypair.to_wif(),
            'public_key': keypair.public_key_hex
        }

    def import_private_key(self, wif_key):
        """Import private key from WIF format"""
        keypair = KeyPair.from_wif(wif_key)
        self.keypairs.append(keypair)

        address = Address.from_public_key(keypair.public_key, 'p2pkh')
        self.addresses.append(address)

        return address

    def get_balance(self, address):
        """
        Get balance for address (mock implementation)

        Args:
            address (str): Address to check

        Returns:
            float: Balance in BTNG
        """
        # Mock balance - in real implementation, query blockchain
        if address in self.addresses:
            return 125.487  # Mock balance
        return 0.0

    def create_transaction(self, to_address, amount, fee=0.001):
        """
        Create a signed transaction

        Args:
            to_address (str): Recipient address
            amount (float): Amount to send
            fee (float): Transaction fee

        Returns:
            dict: Transaction details
        """
        if not self.keypairs:
            raise ValueError("No keys in wallet")

        # Use first keypair
        keypair = self.keypairs[0]
        from_address = self.addresses[0]

        # Check balance
        balance = self.get_balance(from_address)
        if balance < amount + fee:
            raise ValueError("Insufficient balance")

        # Create transaction
        tx = Transaction()

        # Add input (mock - would need real UTXO)
        tx.add_input('00' * 32, 0)  # Mock previous tx

        # Add output to recipient
        script_pubkey = create_p2pkh_script_pubkey(to_address)
        tx.add_output(int(amount * 100000000), script_pubkey)  # Convert to satoshis

        # Add change output
        change_amount = balance - amount - fee
        if change_amount > 0:
            change_script = create_p2pkh_script_pubkey(from_address)
            tx.add_output(int(change_amount * 100000000), change_script)

        # Sign transaction
        signed_tx = tx.sign(keypair)

        return {
            'txid': tx.txid(),
            'hex': signed_tx.hex(),
            'from': from_address,
            'to': to_address,
            'amount': amount,
            'fee': fee
        }

    def sign_message(self, message, address=None):
        """
        Sign a message with the private key for an address

        Args:
            message (str): Message to sign
            address (str, optional): Address to sign with

        Returns:
            str: Base64 encoded signature
        """
        if not address:
            if not self.addresses:
                raise ValueError("No addresses in wallet")
            address = self.addresses[0]

        # Find keypair for address
        for i, addr in enumerate(self.addresses):
            if addr == address:
                keypair = self.keypairs[i]
                signature = keypair.sign(message.encode())
                return signature.hex()

        raise ValueError(f"Address {address} not found in wallet")

    def verify_message(self, message, signature, address):
        """
        Verify a signed message

        Args:
            message (str): Original message
            signature (str): Hex signature
            address (str): Signing address

        Returns:
            bool: True if signature is valid
        """
        # Find public key for address
        for i, addr in enumerate(self.addresses):
            if addr == address:
                keypair = self.keypairs[i]
                return keypair.verify(message.encode(), bytes.fromhex(signature))

        return False

    @property
    def addresses(self):
        """Get all wallet addresses"""
        return self._addresses

    @addresses.setter
    def addresses(self, value):
        self._addresses = value

    def export_wallet(self):
        """Export wallet data (for backup)"""
        return {
            'addresses': self.addresses,
            'keypairs': [kp.to_wif() for kp in self.keypairs]
        }

    @classmethod
    def import_wallet(cls, wallet_data):
        """Import wallet from exported data"""
        wallet = cls()
        for wif in wallet_data['keypairs']:
            wallet.import_private_key(wif)
        return wallet