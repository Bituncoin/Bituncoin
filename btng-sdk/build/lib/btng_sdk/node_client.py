#!/usr/bin/env python3
"""
BTNG Node RPC Client
====================

Connects the BTNG SDK to the live BTNG blockchain node for real transactions.
"""

import json
import time
from typing import Dict, List, Optional, Any
import btng_sdk

class BTNGNodeClient:
    """RPC client for BTNG blockchain node"""

    def __init__(self, host: str = "74.118.126.72", port: int = 54328, timeout: int = 30):
        self.host = host
        self.port = port
        self.timeout = timeout
        self.base_url = f"http://{host}:{port}"

    async def _make_request(self, endpoint: str, method: str = "GET", data: Optional[Dict] = None) -> Dict:
        """Make HTTP request to node"""
        import aiohttp

        url = f"{self.base_url}{endpoint}"

        async with aiohttp.ClientSession(timeout=aiohttp.ClientTimeout(total=self.timeout)) as session:
            if method == "GET":
                async with session.get(url) as response:
                    return await response.json()
            elif method == "POST":
                async with session.post(url, json=data) as response:
                    return await response.json()

    def _make_sync_request(self, endpoint: str, method: str = "GET", data: Optional[Dict] = None) -> Dict:
        """Synchronous request fallback"""
        import requests

        url = f"{self.base_url}{endpoint}"

        if method == "GET":
            response = requests.get(url, timeout=self.timeout)
        elif method == "POST":
            response = requests.post(url, json=data, timeout=self.timeout)

        response.raise_for_status()
        return response.json()

    async def get_blockchain_info(self) -> Dict:
        """Get blockchain information"""
        return await self._make_request("/api/v1/blockchain/info")

    async def get_block_count(self) -> int:
        """Get current block count"""
        info = await self.get_blockchain_info()
        return info.get("blocks", 0)

    async def get_block_hash(self, height: int) -> str:
        """Get block hash by height"""
        return await self._make_request(f"/api/v1/block/{height}/hash")

    async def get_block(self, block_hash: str) -> Dict:
        """Get block by hash"""
        return await self._make_request(f"/api/v1/block/{block_hash}")

    async def get_raw_transaction(self, txid: str) -> str:
        """Get raw transaction hex"""
        return await self._make_request(f"/api/v1/tx/{txid}/raw")

    async def get_transaction(self, txid: str) -> Dict:
        """Get transaction details"""
        return await self._make_request(f"/api/v1/tx/{txid}")

    async def get_utxos(self, address: str) -> List[Dict]:
        """Get UTXOs for address"""
        return await self._make_request(f"/api/v1/address/{address}/utxos")

    async def estimate_fee(self, blocks: int = 6) -> float:
        """Estimate transaction fee"""
        result = await self._make_request(f"/api/v1/fee/estimate/{blocks}")
        return result.get("fee_per_kb", 0.001) / 1000  # Convert to BTNG per byte

    async def broadcast_transaction(self, hex_tx: str) -> str:
        """Broadcast raw transaction"""
        result = await self._make_request("/api/v1/tx/broadcast", "POST", {"hex": hex_tx})
        return result.get("txid", "")

    async def get_mempool_info(self) -> Dict:
        """Get mempool information"""
        return await self._make_request("/api/v1/mempool/info")

    async def get_mempool_transactions(self) -> List[str]:
        """Get transactions in mempool"""
        return await self._make_request("/api/v1/mempool/txs")

    # Synchronous versions for backward compatibility
    def sync_get_blockchain_info(self) -> Dict:
        return self._make_sync_request("/api/v1/blockchain/info")

    def sync_get_block_count(self) -> int:
        info = self.sync_get_blockchain_info()
        return info.get("blocks", 0)

    def sync_get_utxos(self, address: str) -> List[Dict]:
        return self._make_sync_request(f"/api/v1/address/{address}/utxos")

    def sync_estimate_fee(self, blocks: int = 6) -> float:
        result = self._make_sync_request(f"/api/v1/fee/estimate/{blocks}")
        return result.get("fee_per_kb", 0.001) / 1000

    def sync_broadcast_transaction(self, hex_tx: str) -> str:
        result = self._make_sync_request("/api/v1/tx/broadcast", "POST", {"hex": hex_tx})
        return result.get("txid", "")

class BTNGWallet(btng_sdk.Wallet):
    """Enhanced wallet with live blockchain integration"""

    def __init__(self, node_client: Optional[BTNGNodeClient] = None):
        super().__init__()
        self.node = node_client or BTNGNodeClient()

    async def get_balance(self, address: str) -> float:
        """Get real balance from blockchain"""
        try:
            utxos = await self.node.get_utxos(address)
            balance = sum(utxo.get("value", 0) for utxo in utxos)
            return balance
        except Exception as e:
            print(f"Error getting balance: {e}")
            return 0.0

    async def create_transaction(self, to_address: str, amount: float, fee: Optional[float] = None) -> Dict:
        """Create transaction using real UTXOs"""
        if not self.keypairs:
            raise ValueError("No keypairs in wallet")

        # Get sender address
        sender_keypair = self.keypairs[0]
        sender_address = sender_keypair.get_address()

        # Get UTXOs
        utxos = await self.node.get_utxos(sender_address)
        if not utxos:
            raise ValueError("No UTXOs available")

        # Estimate fee if not provided
        if fee is None:
            fee = await self.node.estimate_fee()

        # Select UTXOs (simplified - just use first one)
        selected_utxo = utxos[0]
        input_amount = selected_utxo["value"]

        if input_amount < amount + fee:
            raise ValueError("Insufficient funds")

        # Calculate change
        change = input_amount - amount - fee

        # Build transaction
        tx = btng_sdk.Transaction()
        tx.version = 1

        # Add input
        tx.add_input(
            prev_tx_hash=selected_utxo["txid"],
            prev_output_index=selected_utxo["vout"],
            script_sig=b"",  # Will be filled during signing
            sequence=0xFFFFFFFF
        )

        # Add outputs
        tx.add_output(amount, to_address)
        if change > 0:
            change_address = sender_address  # Send change back to sender
            tx.add_output(change, change_address)

        # Sign transaction
        signed_tx = self.sign_transaction(tx, sender_keypair)

        return {
            "txid": signed_tx.get_txid(),
            "hex": signed_tx.serialize(),
            "fee": fee,
            "amount": amount,
            "change": change
        }

    async def broadcast_transaction(self, tx_hex: str) -> str:
        """Broadcast transaction to network"""
        return await self.node.broadcast_transaction(tx_hex)

    async def send_transaction(self, to_address: str, amount: float) -> str:
        """Create and broadcast transaction in one call"""
        tx = await self.create_transaction(to_address, amount)
        txid = await self.broadcast_transaction(tx["hex"])
        return txid

    def sign_transaction(self, tx: btng_sdk.Transaction, keypair: btng_sdk.KeyPair) -> btng_sdk.Transaction:
        """Sign transaction (simplified implementation)"""
        # This is a placeholder - real implementation would need proper script building
        # For now, return the transaction as-is
        return tx

# Synchronous versions for CLI compatibility
class SyncBTNGWallet(BTNGWallet):
    """Synchronous wallet for CLI usage"""

    def get_balance(self, address: str) -> float:
        """Get balance synchronously"""
        import asyncio
        return asyncio.run(super().get_balance(address))

    def create_transaction(self, to_address: str, amount: float, fee: Optional[float] = None) -> Dict:
        """Create transaction synchronously"""
        import asyncio
        return asyncio.run(super().create_transaction(to_address, amount, fee))

    def broadcast_transaction(self, tx_hex: str) -> str:
        """Broadcast transaction synchronously"""
        import asyncio
        return asyncio.run(super().broadcast_transaction(tx_hex))

    def send_transaction(self, to_address: str, amount: float) -> str:
        """Send transaction synchronously"""
        import asyncio
        return asyncio.run(super().send_transaction(to_address, amount))

if __name__ == "__main__":
    # Test the node client
    import asyncio

    async def test():
        client = BTNGNodeClient()
        try:
            info = await client.get_blockchain_info()
            print(f"Blockchain info: {json.dumps(info, indent=2)}")
        except Exception as e:
            print(f"Error: {e}")

    asyncio.run(test())