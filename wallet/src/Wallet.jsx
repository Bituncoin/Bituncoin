import React, { useState, useEffect } from 'react';
import './Wallet.css';

const Wallet = () => {
  const [activeTab, setActiveTab] = useState('overview');
  const [balance, setBalance] = useState({
    btn: 0,
    goldcoin: 0,
    bitcoin: 0,
    ethereum: 0,
    usdt: 0,
    bnb: 0,
  });
  const [stakingInfo, setStakingInfo] = useState({
    stakedAmount: 0,
    rewards: 0,
    apy: 5.0,
  });
  const [transactions, setTransactions] = useState([]);
  const [securitySettings, setSecuritySettings] = useState({
    twoFactorEnabled: false,
    biometricEnabled: false,
    fraudMonitoring: true,
  });
  const [theme, setTheme] = useState('light');
  const [language, setLanguage] = useState('en');

  // Mock data for demonstration
  useEffect(() => {
    // Simulate fetching wallet data
    setBalance({
      btn: 5000.0,
      goldcoin: 1250.5,
      bitcoin: 0.05,
      ethereum: 2.3,
      usdt: 1000.0,
      bnb: 10.5,
    });

    setTransactions([
      { id: '1', type: 'Received', amount: 100, coin: 'BTN', date: '2025-10-14', status: 'Completed' },
      { id: '2', type: 'Sent', amount: 50, coin: 'GLD', date: '2025-10-13', status: 'Completed' },
      { id: '3', type: 'Staked', amount: 500, coin: 'BTN', date: '2025-10-12', status: 'Active' },
      { id: '4', type: 'Swap', amount: 250, coin: 'ETH', date: '2025-10-11', status: 'Completed' },
      { id: '5', type: 'Received', amount: 1000, coin: 'USDT', date: '2025-10-10', status: 'Completed' },
    ]);

    setStakingInfo({
      stakedAmount: 500,
      rewards: 2.5,
      apy: 5.0,
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

  const handlePay = () => {
    alert('BTN-Pay - Merchant payment functionality');
  };

  const handleMobileMoney = () => {
    alert('Mobile Money integration - MTN, AirtelTigo, Vodafone Cash');
  };

  const toggleTheme = () => {
    setTheme(theme === 'light' ? 'dark' : 'light');
  };

  const toggleFraudMonitoring = () => {
    setSecuritySettings({
      ...securitySettings,
      fraudMonitoring: !securitySettings.fraudMonitoring,
    });
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
    <div className={`wallet-container ${theme}`}>
      <header className="wallet-header">
        <h1>BTNg Wallet - Bituncoin Blockchain Operating System</h1>
        <p className="wallet-subtitle">Multi-Currency Digital Wallet with Advanced Features</p>
        <div className="header-controls">
          <button className="theme-toggle" onClick={toggleTheme}>
            {theme === 'light' ? '🌙' : '☀️'}
          </button>
          <select 
            className="language-selector"
            value={language}
            onChange={(e) => setLanguage(e.target.value)}
          >
            <option value="en">English</option>
            <option value="fr">Français</option>
            <option value="es">Español</option>
            <option value="zh">中文</option>
          </select>
        </div>
      </header>

      <nav className="wallet-nav">
        <button
          className={activeTab === 'overview' ? 'active' : ''}
          onClick={() => setActiveTab('overview')}
        >
          📊 Overview
        </button>
        <button
          className={activeTab === 'pay' ? 'active' : ''}
          onClick={() => setActiveTab('pay')}
        >
          💳 Pay
        </button>
        <button
          className={activeTab === 'staking' ? 'active' : ''}
          onClick={() => setActiveTab('staking')}
        >
          💎 Staking
        </button>
        <button
          className={activeTab === 'transactions' ? 'active' : ''}
          onClick={() => setActiveTab('transactions')}
        >
          📜 Transactions
        </button>
        <button
          className={activeTab === 'security' ? 'active' : ''}
          onClick={() => setActiveTab('security')}
        >
          🔒 Security
        </button>
        <button
          className={activeTab === 'ai' ? 'active' : ''}
          onClick={() => setActiveTab('ai')}
        >
          🤖 AI Assistant
        </button>
      </nav>

      <div className="wallet-content">
        {activeTab === 'overview' && (
          <div className="overview-tab">
            <div className="balance-section">
              <h2>Your Balances</h2>
              <div className="balance-cards">
                <div className="balance-card btn">
                  <div className="coin-icon">🪙</div>
                  <h3>Bituncoin (BTN)</h3>
                  <p className="balance-amount">{balance.btn.toFixed(2)} BTN</p>
                  <p className="balance-usd">≈ ${(balance.btn * 15).toFixed(2)} USD</p>
                </div>
                <div className="balance-card goldcoin">
                  <div className="coin-icon">🏆</div>
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
                <div className="balance-card usdt">
                  <div className="coin-icon">💵</div>
                  <h3>Tether (USDT)</h3>
                  <p className="balance-amount">{balance.usdt.toFixed(2)} USDT</p>
                  <p className="balance-usd">≈ ${balance.usdt.toFixed(2)} USD</p>
                </div>
                <div className="balance-card bnb">
                  <div className="coin-icon">🔶</div>
                  <h3>Binance Coin (BNB)</h3>
                  <p className="balance-amount">{balance.bnb.toFixed(2)} BNB</p>
                  <p className="balance-usd">≈ ${(balance.bnb * 300).toFixed(2)} USD</p>
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
                  Receive
                </button>
                <button className="action-btn swap" onClick={handleCrossChainSwap}>
                  <span>🔄</span>
                  Swap
                </button>
                <button className="action-btn stake" onClick={handleStake}>
                  <span>💎</span>
                  Stake
                </button>
                <button className="action-btn pay" onClick={handlePay}>
                  <span>💳</span>
                  Pay
                </button>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'pay' && (
          <div className="pay-tab">
            <h2>BTN-Pay - Merchant Payments</h2>
            
            <div className="payment-section">
              <h3>Payment Methods</h3>
              <div className="payment-methods">
                <div className="payment-card">
                  <div className="payment-icon">💳</div>
                  <h4>BTN-Pay Card</h4>
                  <p>Virtual & Physical MasterCard/Visa</p>
                  <button className="btn-primary">Manage Cards</button>
                </div>
                <div className="payment-card">
                  <div className="payment-icon">📱</div>
                  <h4>Mobile Money</h4>
                  <p>MTN, AirtelTigo, Vodafone Cash</p>
                  <button className="btn-primary" onClick={handleMobileMoney}>Connect</button>
                </div>
                <div className="payment-card">
                  <div className="payment-icon">📷</div>
                  <h4>QR Code Payment</h4>
                  <p>Scan to pay merchants</p>
                  <button className="btn-primary">Generate QR</button>
                </div>
              </div>
            </div>

            <div className="invoice-section">
              <h3>Create Invoice (Merchant)</h3>
              <div className="invoice-form">
                <input type="number" placeholder="Amount" className="invoice-input" />
                <select className="invoice-select">
                  <option value="BTN">BTN</option>
                  <option value="GLD">GLD</option>
                  <option value="USDT">USDT</option>
                </select>
                <input type="text" placeholder="Description" className="invoice-input" />
                <button className="btn-primary">Create Invoice</button>
              </div>
            </div>

            <div className="payment-history">
              <h3>Recent Payments</h3>
              <div className="payment-list">
                <p className="placeholder-text">No payment history available</p>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'staking' && (
          <div className="staking-tab">
            <div className="staking-info">
              <h2>Gold-Coin Staking</h2>
              <div className="staking-cards">
                <div className="staking-card">
                  <h3>Staked Amount</h3>
                  <p className="staking-value">{stakingInfo.stakedAmount} GLD</p>
                </div>
                <div className="staking-card">
                  <h3>Rewards Earned</h3>
                  <p className="staking-value">{stakingInfo.rewards} GLD</p>
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

              <div className="security-option">
                <div className="security-info">
                  <h4>🛡️ Fraud Monitoring</h4>
                  <p>Real-time transaction monitoring and alerts</p>
                </div>
                <label className="toggle-switch">
                  <input
                    type="checkbox"
                    checked={securitySettings.fraudMonitoring}
                    onChange={toggleFraudMonitoring}
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
        {activeTab === 'ai' && (
          <div className="ai-tab">
            <h2>🤖 AI-Powered Assistant</h2>
            
            <div className="ai-section">
              <h3>Wallet Insights</h3>
              <div className="ai-insights">
                <div className="insight-card">
                  <span className="insight-icon">📊</span>
                  <div className="insight-content">
                    <h4>Portfolio Analysis</h4>
                    <p>Your portfolio is well-diversified. Consider increasing BTC allocation by 5%.</p>
                  </div>
                </div>
                <div className="insight-card">
                  <span className="insight-icon">💡</span>
                  <div className="insight-content">
                    <h4>Staking Recommendation</h4>
                    <p>You have 1,250 GLD available. Staking could earn you ~62.5 GLD annually (5% APY).</p>
                  </div>
                </div>
                <div className="insight-card">
                  <span className="insight-icon">📈</span>
                  <div className="insight-content">
                    <h4>Market Trend</h4>
                    <p>BTN is showing strong upward momentum. +12% in the last 7 days.</p>
                  </div>
                </div>
              </div>
            </div>

            <div className="ai-section">
              <h3>Security Alerts</h3>
              <div className="ai-alerts">
                <div className="alert-item success">
                  <span className="alert-icon">✓</span>
                  <p>No suspicious activity detected in the last 30 days</p>
                </div>
                <div className="alert-item info">
                  <span className="alert-icon">ℹ️</span>
                  <p>Your 2FA is enabled. Your wallet is secure.</p>
                </div>
                <div className="alert-item warning">
                  <span className="alert-icon">⚠️</span>
                  <p>Backup reminder: Last backup was 15 days ago. Consider creating a new backup.</p>
                </div>
              </div>
            </div>

            <div className="ai-section">
              <h3>Ask AI Assistant</h3>
              <div className="ai-chat">
                <div className="chat-input-container">
                  <input 
                    type="text" 
                    placeholder="Ask about your wallet, crypto markets, or best practices..."
                    className="chat-input"
                  />
                  <button className="btn-primary">Ask</button>
                </div>
                <div className="chat-suggestions">
                  <button className="suggestion-btn">How can I improve my portfolio?</button>
                  <button className="suggestion-btn">What are the best staking strategies?</button>
                  <button className="suggestion-btn">How do I secure my wallet?</button>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default Wallet;
