"""
BTNG Smart Contract
Bituncoin Gold (BTNG) - Mint & Melt Ledger Contract
"""

import json
from datetime import datetime, timezone


class BTNGContract:
    def __init__(self):
        self.ledger = {}
        self.transaction_log = []

    def mint(self, wallet: str, amount: float) -> dict:
        if amount <= 0:
            raise ValueError("Mint amount must be positive")
        self.ledger[wallet] = self.ledger.get(wallet, 0) + amount
        tx = {
            "type": "mint",
            "wallet": wallet,
            "amount": amount,
            "balance_after": self.ledger[wallet],
            "timestamp": datetime.now(timezone.utc).isoformat(),
        }
        self.transaction_log.append(tx)
        return tx

    def melt(self, wallet: str, amount: float) -> dict:
        if amount <= 0:
            raise ValueError("Melt amount must be positive")
        if self.ledger.get(wallet, 0) < amount:
            raise Exception(
                f"Insufficient BTNG: wallet {wallet} has {self.ledger.get(wallet, 0)}, requested {amount}"
            )
        self.ledger[wallet] -= amount
        tx = {
            "type": "melt",
            "wallet": wallet,
            "amount": amount,
            "balance_after": self.ledger[wallet],
            "timestamp": datetime.now(timezone.utc).isoformat(),
        }
        self.transaction_log.append(tx)
        return tx

    def balance(self, wallet: str) -> float:
        return self.ledger.get(wallet, 0)

    def get_ledger(self) -> dict:
        return dict(self.ledger)

    def get_transaction_log(self) -> list:
        return list(self.transaction_log)


if __name__ == "__main__":
    contract = BTNGContract()

    # Demo: Mint and Melt
    print(json.dumps(contract.mint("btng1qxyz_main", 1000.0), indent=2))
    print(json.dumps(contract.mint("btng1qxyz_main", 250.25), indent=2))
    print(json.dumps(contract.melt("btng1qxyz_main", 100.0), indent=2))

    print(f"\nFinal balance: {contract.balance('btng1qxyz_main')} BTNG")
    print(f"Ledger: {json.dumps(contract.get_ledger(), indent=2)}")
