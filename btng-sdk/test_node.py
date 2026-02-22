#!/usr/bin/env python3
"""
Test BTNG Node Client Connection
"""

import asyncio
import btng_sdk

async def test_node_connection():
    """Test connection to BTNG node"""
    print("🔗 Testing BTNG Node Connection...")
    print("==================================")

    client = btng_sdk.BTNGNodeClient()

    try:
        # Test health check
        print("\n1. Testing health endpoint...")
        health = await client._make_request("/api/btng/health")
        print(f"   ✅ Health: {health}")

        # Test oracle price
        print("\n2. Getting oracle price...")
        price = await client.get_oracle_price()
        print(f"   ✅ Oracle price: {price}")

        # Test mining hashrate
        print("\n3. Getting mining hashrate...")
        hashrate = await client.get_mining_hashrate()
        print(f"   ✅ Mining hashrate: {hashrate}")

        # Test mining info
        print("\n4. Getting mining info...")
        mining_info = await client.get_mempool_info()
        print(f"   ✅ Mining info: {mining_info}")

        print("\n🎉 Node connection successful!")

    except Exception as e:
        print(f"\n❌ Connection failed: {e}")
        print("   Note: This may be expected if the node endpoints aren't fully implemented yet")

if __name__ == "__main__":
    asyncio.run(test_node_connection())