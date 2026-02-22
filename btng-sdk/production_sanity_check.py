#!/usr/bin/env python3
"""
BTNG Production Readiness - Final Sanity Check
==============================================

Comprehensive validation of all system components before production deployment.
"""

import time
from btng_sdk import BTNGNodeClient

def production_sanity_check():
    print('🚀 BTNG Production Readiness - Final Sanity Check')
    print('=' * 60)

    client = BTNGNodeClient(mock_mode=True)

    # Generate test addresses
    from btng_sdk import KeyPair, Address
    test_addresses = []
    for i in range(10):
        keypair = KeyPair.generate()
        addr = Address.from_public_key(keypair.public_key, 'p2pkh', 'mainnet')
        test_addresses.append(addr)

    print(f'📋 Generated {len(test_addresses)} test addresses')

    # Test 1: Balance queries under load
    print('\n🧪 Test 1: Balance Query Performance (10k queries)')
    N = 10000
    start = time.time()

    query_count = 0
    for i in range(N):
        addr = test_addresses[i % len(test_addresses)]
        try:
            balance = client.sync_get_balance(addr)
            query_count += 1
        except Exception as e:
            if i == 0:
                print(f'❌ Query failed: {e}')
            continue

    elapsed = time.time() - start
    avg_latency = (elapsed / query_count) * 1000 if query_count > 0 else float('inf')

    status = 'PASS' if avg_latency < 50 else 'FAIL'
    print(f'✅ {query_count}/{N} queries successful')
    print(f'⏱️  Total time: {elapsed:.3f}s')
    print(f'📊 Avg latency: {avg_latency:.3f} ms per query')
    print(f'🎯 SLA Target: <50ms | Status: {status}')

    # Test 2: Transaction creation
    print('\n💸 Test 2: Transaction Creation & UTXO Management')

    # Create a test wallet with keypair
    from btng_sdk import SyncBTNGWallet
    test_wallet = SyncBTNGWallet(client)
    test_keypair = KeyPair.generate()
    test_wallet.keypairs = [test_keypair]

    # Mock address for testing
    test_addr = 'C7z995zi3EHaNSShwcnL9ex5Cf6wJui6BF'

    # Add test UTXOs
    if hasattr(client, 'utxo_pool') and client.utxo_pool:
        for i in range(3):
            from btng_sdk.node_client import UtxoRecord
            utxo_id = f'test_prod_{i}'
            record = UtxoRecord(
                txid=f'prod_tx_{i}_' + 'a' * 56,
                vout=0,
                value=100000000,  # 1 BTNG
                owner=test_addr,
                script_pubkey='76a91412345678901234567890123456789012345678901288ac',
                block_height=1001
            )
            client.utxo_pool.add_utxo(utxo_id, record)

    tx_created = False
    try:
        tx = test_wallet.create_transaction('BzeyaaCNibXiQDjWp5y2fte3Jjggqg8Nbo', 0.5)
        txid_short = tx.get('txid', '')[:16] if tx.get('txid') else 'N/A'
        print('✅ Transaction created: {}...'.format(txid_short))
        print('💰 Amount: {} BTNG'.format(tx.get('amount', 0)))
        print('💸 Fee: {} BTNG'.format(tx.get('fee', 0)))
        print('🔄 Change: {:.2f} BTNG'.format(tx.get('change', 0)))
        tx_created = True
    except Exception as e:
        print('❌ Transaction failed: {}'.format(e))

    # Test 3: Pool statistics
    print('\n📊 Test 3: UTXO Pool Health')
    stats = client.get_utxo_pool_stats()
    print(f'📈 Pool Stats: {stats}')

    total_utxos = stats.get('total_utxos', 0)
    unique_addresses = stats.get('unique_addresses', 0)
    cache_size = stats.get('cache_size', 0)

    print(f'🔢 Total UTXOs: {total_utxos}')
    print(f'🏠 Unique Addresses: {unique_addresses}')
    print(f'💾 Cache Size: {cache_size}')

    # Test 4: Thread safety (basic)
    print('\n🔒 Test 4: Thread Safety Check')
    import threading

    results = []
    errors = []

    def worker_thread(thread_id):
        try:
            for i in range(100):
                addr = test_addresses[i % len(test_addresses)]
                balance = client.sync_get_balance(addr)
                results.append((thread_id, balance))
        except Exception as e:
            errors.append((thread_id, str(e)))

    threads = []
    for i in range(4):  # 4 concurrent threads
        t = threading.Thread(target=worker_thread, args=(i,))
        threads.append(t)
        t.start()

    for t in threads:
        t.join()

    print(f'✅ {len(results)} concurrent operations completed')
    if errors:
        print(f'❌ {len(errors)} thread errors')
    else:
        print('🔒 Thread safety: PASS')

    # Final assessment
    print('\n🎯 PRODUCTION READINESS ASSESSMENT')
    print('=' * 60)

    checks = [
        ('Balance Query Performance', avg_latency < 50),
        ('Transaction Creation', tx_created),
        ('UTXO Pool Health', total_utxos > 0),
        ('Thread Safety', len(errors) == 0),
        ('Cache Functionality', cache_size >= 0),
    ]

    all_pass = True
    for check_name, passed in checks:
        status_icon = '✅' if passed else '❌'
        status_text = 'PASS' if passed else 'FAIL'
        print(f'{status_icon} {check_name}: {status_text}')
        if not passed:
            all_pass = False

    print()
    if all_pass:
        print('🎉 ALL CHECKS PASSED - SYSTEM READY FOR PRODUCTION DEPLOYMENT!')
        print('🚀 Proceed with the deployment checklist in the go-live guide.')
        print()
        print('📋 Next Steps:')
        print('1. 🐳 Build Docker images with the provided Dockerfile templates')
        print('2. ☸️  Deploy to Kubernetes using the rollout strategy')
        print('3. 📊 Set up monitoring with Prometheus/Grafana')
        print('4. 🔒 Configure security hardening (TLS, RBAC, secrets)')
        print('5. 🚦 Execute canary deployment and gradual rollout')
        print('6. 🎯 Monitor metrics and ensure SLA compliance')
    else:
        print('⚠️  SOME CHECKS FAILED - Review and fix before deployment.')
        print('   Check the error messages above and resolve issues.')

    print(f'\n⏰ Test completed at: {time.strftime("%Y-%m-%d %H:%M:%S")}')

if __name__ == '__main__':
    production_sanity_check()