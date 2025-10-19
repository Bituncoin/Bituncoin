import React, { useState, useEffect } from 'react';
import './Wallet.css';

const Wallet = () => {
  const [activeTab, setActiveTab] = useState('overview');
  const [balance, setBalance] = useState({
    bituncoin: 0,  // BTN - Primary currency
    goldcoin: 0,   // GLD - Secondary currency
    bitcoin: 0,
    ethereum: 0,
  });
  const [stakingInfo, setStakingInfo] = useState({
    stakedAmount: 0,
    rewards: 0,
    apy: 5.0,
    currency: 'BTN',  // Default to BTN
  });
  const [transactions, setTransactions] = useState([]);
  const [securitySettings, setSecuritySettings] = useState({
    twoFactorEnabled: false,
    biometricEnabled: false,
  });

  // Mock data for demonstration
  useEffect(() => {
    // Simulate fetching wallet data
    setBalance({
      bituncoin: 2500.75,  // BTN balance
      goldcoin: 1250.5,    // GLD balance
      bitcoin: 0.05,
      ethereum: 2.3,
    });

    setTransactions([
      { id: '1', type: 'Received', amount: 500, coin: 'BTN', date: '2025-10-19', status: 'Completed' },
      { id: '2', type: 'Sent', amount: 100, coin: 'BTN', date: '2025-10-18', status: 'Completed' },
      { id: '3', type: 'Received', amount: 100, coin: 'GLD', date: '2025-10-14', status: 'Completed' },
      { id: '4', type: 'Sent', amount: 50, coin: 'GLD', date: '2025-10-13', status: 'Completed' },
      { id: '5', type: 'Staked', amount: 1000, coin: 'BTN', date: '2025-10-12', status: 'Active' },
    ]);

    setStakingInfo({
      stakedAmount: 1000,
      rewards: 5.2,
      apy: 5.0,
      currency: 'BTN',
    });
  }, []);

  const handleStake = () => {
    alert('Staking functionality - Connect to backend API');
  };

  const handleUnstake = () => {
    alert('Unstaking functionality - Connect to backend API');
  };

  const handleClaimRewards = () => {
    alert('Claim rewards functionality - Connect to backend API');
  };

  const handleSend = () => {
    alert('Send transaction functionality - Connect to backend API');
  };

  const handleReceive = () => {
    alert('Receive functionality - Show QR code with address');
  };

  const handleCrossChainSwap = () => {
    alert('Cross-chain swap functionality - Connect to bridge API');
  };

  const toggleTwoFactor = () => {
    setSecuritySettings({
      ...securitySettings,
      twoFactorEnabled: !securitySettings.twoFactorEnabled,
    });
  };

  const toggleBiometric = () => {
    setSecuritySettings({
      ...securitySettings,
      biometricEnabled: !securitySettings.biometricEnabled,
    });
  };

  const handleBackup = () => {
    alert('Backup wallet - Generate encrypted backup file');
  };

  const handleRestore = () => {
    alert('Restore wallet - Upload backup file');
  };

  return (
    <div className="wallet-container">
      <header className="wallet-header">
        <h1>Universal Wallet</h1>
        <p className="wallet-subtitle">Multi-Currency Digital Wallet with Bituncoin (BTN) as Primary</p>
      </header>

      <nav className="wallet-nav">
        <button
          className={activeTab === 'overview' ? 'active' : ''}
          onClick={() => setActiveTab('overview')}
        >
          Overview
        </button>
        <button
          className={activeTab === 'staking' ? 'active' : ''}
          onClick={() => setActiveTab('staking')}
        >
          Staking
        </button>
        <button
          className={activeTab === 'transactions' ? 'active' : ''}
          onClick={() => setActiveTab('transactions')}
        >
          Transactions
        </button>
        <button
          className={activeTab === 'security' ? 'active' : ''}
          onClick={() => setActiveTab('security')}
        >
          Security
        </button>
      </nav>

      <div className="wallet-content">
        {activeTab === 'overview' && (
          <div className="overview-tab">
            <div className="balance-section">
              <h2>Your Balances</h2>
              <div className="balance-cards">
                <div className="balance-card bituncoin primary">
                  <div className="coin-icon">🪙</div>
                  <h3>Bituncoin (BTN)</h3>
                  <p className="balance-amount">{balance.bituncoin.toFixed(2)} BTN</p>
                  <p className="balance-usd">≈ ${(balance.bituncoin * 15).toFixed(2)} USD</p>
                  <span className="primary-badge">PRIMARY</span>
                </div>
                <div className="balance-card goldcoin">
                  <div className="coin-icon">🥇</div>
                  <h3>Gold-Coin (GLD)</h3>
                  <p className="balance-amount">{balance.goldcoin.toFixed(2)} GLD</p>
                  <p className="balance-usd">≈ ${(balance.goldcoin * 10).toFixed(2)} USD</p>
                </div>
                <div className="balance-card bitcoin">
                  <div className="coin-icon">₿</div>
                  <h3>Bitcoin (BTC)</h3>
                  <p className="balance-amount">{balance.bitcoin.toFixed(4)} BTC</p>
                  <p className="balance-usd">≈ ${(balance.bitcoin * 45000).toFixed(2)} USD</p>
                </div>
                <div className="balance-card ethereum">
                  <div className="coin-icon">Ξ</div>
                  <h3>Ethereum (ETH)</h3>
                  <p className="balance-amount">{balance.ethereum.toFixed(2)} ETH</p>
                  <p className="balance-usd">≈ ${(balance.ethereum * 3000).toFixed(2)} USD</p>
                </div>
              </div>
            </div>

            <div className="actions-section">
              <h2>Quick Actions</h2>
              <div className="action-buttons">
                <button className="action-btn send" onClick={handleSend}>
                  <span>📤</span>
                  Send
                </button>
                <button className="action-btn receive" onClick={handleReceive}>
                  <span>📥</span>
                  Receive (QR Code)
                </button>
                <button className="action-btn swap" onClick={handleCrossChainSwap}>
                  <span>🔄</span>
                  Cross-Chain Swap
                </button>
                <button className="action-btn stake" onClick={handleStake}>
                  <span>💎</span>
                  Stake BTN
                </button>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'staking' && (
          <div className="staking-tab">
            <div className="staking-info">
              <h2>Bituncoin (BTN) Staking</h2>
              <div className="staking-cards">
                <div className="staking-card">
                  <h3>Staked Amount</h3>
                  <p className="staking-value">{stakingInfo.stakedAmount} BTN</p>
                </div>
                <div className="staking-card">
                  <h3>Rewards Earned</h3>
                  <p className="staking-value">{stakingInfo.rewards} BTN</p>
                </div>
                <div className="staking-card">
                  <h3>Annual Yield</h3>
                  <p className="staking-value">{stakingInfo.apy}%</p>
                </div>
              </div>
            </div>

            <div className="staking-actions">
              <h3>Manage Staking</h3>
              <div className="staking-buttons">
                <button className="btn-primary" onClick={handleStake}>
                  Stake More GLD
                </button>
                <button className="btn-secondary" onClick={handleClaimRewards}>
                  Claim Rewards
                </button>
                <button className="btn-warning" onClick={handleUnstake}>
                  Unstake
                </button>
              </div>
              <div className="staking-info-text">
                <p>• Minimum stake: 100 GLD</p>
                <p>• Lock period: 30 days</p>
                <p>• Rewards calculated daily</p>
                <p>• Annual percentage yield: 5%</p>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'transactions' && (
          <div className="transactions-tab">
            <h2>Transaction History</h2>
            <div className="transactions-list">
              {transactions.map((tx) => (
                <div key={tx.id} className="transaction-item">
                  <div className="tx-type">
                    <span className={`tx-icon ${tx.type.toLowerCase()}`}>
                      {tx.type === 'Received' ? '📥' : tx.type === 'Sent' ? '📤' : '💎'}
                    </span>
                    <div className="tx-details">
                      <p className="tx-type-text">{tx.type}</p>
                      <p className="tx-date">{tx.date}</p>
                    </div>
                  </div>
                  <div className="tx-amount">
                    <p className={`tx-value ${tx.type === 'Sent' ? 'negative' : 'positive'}`}>
                      {tx.type === 'Sent' ? '-' : '+'}{tx.amount} {tx.coin}
                    </p>
                    <p className="tx-status">{tx.status}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'security' && (
          <div className="security-tab">
            <h2>Security Settings</h2>
            
            <div className="security-section">
              <h3>Authentication</h3>
              <div className="security-option">
                <div className="security-info">
                  <h4>🔐 Two-Factor Authentication (2FA)</h4>
                  <p>Add an extra layer of security with 2FA</p>
                </div>
                <label className="toggle-switch">
                  <input
                    type="checkbox"
                    checked={securitySettings.twoFactorEnabled}
                    onChange={toggleTwoFactor}
                  />
                  <span className="toggle-slider"></span>
                </label>
              </div>

              <div className="security-option">
                <div className="security-info">
                  <h4>👆 Biometric Login</h4>
                  <p>Use fingerprint or face recognition</p>
                </div>
                <label className="toggle-switch">
                  <input
                    type="checkbox"
                    checked={securitySettings.biometricEnabled}
                    onChange={toggleBiometric}
                  />
                  <span className="toggle-slider"></span>
                </label>
              </div>
            </div>

            <div className="security-section">
              <h3>Backup & Recovery</h3>
              <div className="backup-buttons">
                <button className="btn-primary" onClick={handleBackup}>
                  💾 Backup Wallet
                </button>
                <button className="btn-secondary" onClick={handleRestore}>
                  🔄 Restore Wallet
                </button>
              </div>
              <div className="backup-info">
                <p>⚠️ Always keep your backup file secure and encrypted</p>
                <p>✓ Store backups in multiple secure locations</p>
                <p>✓ Never share your recovery phrase with anyone</p>
              </div>
            </div>

            <div className="security-section">
              <h3>🔒 Encryption Status</h3>
              <div className="encryption-status">
                <p>✓ Wallet encrypted with AES-256</p>
                <p>✓ Private keys stored securely</p>
                <p>✓ End-to-end encryption enabled</p>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Wallet;
