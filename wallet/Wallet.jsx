import React, { useState, useEffect } from 'react';
import './Wallet.css';

const Wallet = () => {
  const [activeTab, setActiveTab] = useState('overview');
  const [balance, setBalance] = useState({
    btn: 0,
    bitcoin: 0,
    ethereum: 0,
    usdt: 0,
    bnb: 0,
    goldcoin: 0,
  });
  const [totalUSD, setTotalUSD] = useState(0);
  const [stakingInfo, setStakingInfo] = useState({
    stakedAmount: 0,
    rewards: 0,
    apy: 5.0,
  });
  const [transactions, setTransactions] = useState([]);
  const [cards, setCards] = useState([]);
  const [insights, setInsights] = useState([]);
  const [alerts, setAlerts] = useState([]);
  const [securitySettings, setSecuritySettings] = useState({
    twoFactorEnabled: false,
    biometricEnabled: false,
    fraudDetectionEnabled: true,
  });

  // Mock data for demonstration
  useEffect(() => {
    // Simulate fetching wallet data
    setBalance({
      btn: 5000.0,
      bitcoin: 0.05,
      ethereum: 2.3,
      usdt: 10000.0,
      bnb: 15.5,
      goldcoin: 1250.5,
    });

    setTotalUSD(68750.0);

    setTransactions([
      { id: '1', type: 'Received', amount: 1000, coin: 'BTN', date: '2025-10-19', status: 'Completed' },
      { id: '2', type: 'Sent', amount: 500, coin: 'BTN', date: '2025-10-18', status: 'Completed' },
      { id: '3', type: 'Staked', amount: 500, coin: 'BTN', date: '2025-10-17', status: 'Active' },
      { id: '4', type: 'Exchange', amount: 0.5, coin: 'BTC', date: '2025-10-16', status: 'Completed' },
      { id: '5', type: 'Card', amount: 125.50, coin: 'USD', date: '2025-10-15', status: 'Completed' },
    ]);

    setCards([
      { id: 'card_1', type: 'Virtual Visa', last4: '4532', status: 'Active', balance: 5000.0 },
      { id: 'card_2', type: 'Physical MasterCard', last4: '8765', status: 'Active', balance: 3000.0 },
    ]);

    setInsights([
      { id: 'ins_1', title: 'Portfolio Diversification', message: 'Consider diversifying into stable assets', priority: 'medium' },
      { id: 'ins_2', title: 'Staking Opportunity', message: 'Stake BTN to earn 5% APY', priority: 'high' },
    ]);

    setAlerts([
      { id: 'alert_1', type: 'Market Trend', message: 'BTN price increased 5% in 24h', severity: 'info' },
      { id: 'alert_2', type: 'Security', message: 'New device login detected', severity: 'warning' },
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

  const handleExchange = () => {
    alert('Exchange cryptocurrency - Convert between currencies');
  };

  const handleCreateCard = () => {
    alert('Create new BTN-Pay card - Virtual or Physical');
  };

  const handleMerchantPay = () => {
    alert('Merchant payment - QR/NFC/Mobile Money');
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
        <p className="wallet-subtitle">Multi-Currency Digital Wallet with Gold-Coin Support</p>
      </header>

      <nav className="wallet-nav">
        <button
          className={activeTab === 'overview' ? 'active' : ''}
          onClick={() => setActiveTab('overview')}
        >
          Overview
        </button>
        <button
          className={activeTab === 'cards' ? 'active' : ''}
          onClick={() => setActiveTab('cards')}
        >
          Cards
        </button>
        <button
          className={activeTab === 'exchange' ? 'active' : ''}
          onClick={() => setActiveTab('exchange')}
        >
          Exchange
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
          className={activeTab === 'insights' ? 'active' : ''}
          onClick={() => setActiveTab('insights')}
        >
          AI Insights
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
              <div className="total-portfolio">
                <h3>Total Portfolio Value</h3>
                <p className="total-usd">${totalUSD.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 })}</p>
              </div>
              <div className="balance-cards">
                <div className="balance-card btn">
                  <div className="coin-icon">₿</div>
                  <h3>Bituncoin (BTN)</h3>
                  <p className="balance-amount">{balance.btn.toFixed(2)} BTN</p>
                  <p className="balance-usd">≈ ${(balance.btn * 10).toFixed(2)} USD</p>
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
                  <div className="coin-icon">₮</div>
                  <h3>Tether (USDT)</h3>
                  <p className="balance-amount">{balance.usdt.toFixed(2)} USDT</p>
                  <p className="balance-usd">≈ ${balance.usdt.toFixed(2)} USD</p>
                </div>
                <div className="balance-card bnb">
                  <div className="coin-icon">🔸</div>
                  <h3>Binance Coin (BNB)</h3>
                  <p className="balance-amount">{balance.bnb.toFixed(2)} BNB</p>
                  <p className="balance-usd">≈ ${(balance.bnb * 600).toFixed(2)} USD</p>
                </div>
                <div className="balance-card goldcoin">
                  <div className="coin-icon">🪙</div>
                  <h3>Gold-Coin (GLD)</h3>
                  <p className="balance-amount">{balance.goldcoin.toFixed(2)} GLD</p>
                  <p className="balance-usd">≈ ${(balance.goldcoin * 10).toFixed(2)} USD</p>
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
                  Cross-Chain
                </button>
                <button className="action-btn exchange" onClick={handleExchange}>
                  <span>💱</span>
                  Exchange
                </button>
                <button className="action-btn stake" onClick={handleStake}>
                  <span>💎</span>
                  Stake
                </button>
                <button className="action-btn card" onClick={handleCreateCard}>
                  <span>💳</span>
                  Cards
                </button>
                <button className="action-btn merchant" onClick={handleMerchantPay}>
                  <span>🏪</span>
                  Pay Merchant
                </button>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'cards' && (
          <div className="cards-tab">
            <h2>BTN-Pay Cards</h2>
            <div className="cards-info">
              <p>Manage your virtual and physical payment cards backed by your cryptocurrency.</p>
              <button className="btn-primary" onClick={handleCreateCard}>
                Create New Card
              </button>
            </div>
            
            <div className="cards-list">
              {cards.map((card) => (
                <div key={card.id} className="card-item">
                  <div className="card-visual">
                    <div className="card-type">{card.type}</div>
                    <div className="card-number">•••• •••• •••• {card.last4}</div>
                    <div className="card-status">
                      <span className={`status-badge ${card.status.toLowerCase()}`}>
                        {card.status}
                      </span>
                    </div>
                  </div>
                  <div className="card-details">
                    <p className="card-balance">Available: ${card.balance.toFixed(2)}</p>
                    <div className="card-actions">
                      <button className="btn-small">Manage</button>
                      <button className="btn-small">Transactions</button>
                      <button className="btn-small">Freeze</button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {activeTab === 'exchange' && (
          <div className="exchange-tab">
            <h2>Cryptocurrency Exchange</h2>
            <div className="exchange-form">
              <div className="exchange-input">
                <label>From</label>
                <select className="currency-select">
                  <option>BTN</option>
                  <option>BTC</option>
                  <option>ETH</option>
                  <option>USDT</option>
                  <option>BNB</option>
                </select>
                <input type="number" placeholder="0.00" className="amount-input" />
              </div>
              
              <div className="exchange-arrow">
                <span>🔄</span>
              </div>
              
              <div className="exchange-input">
                <label>To</label>
                <select className="currency-select">
                  <option>USDT</option>
                  <option>BTN</option>
                  <option>BTC</option>
                  <option>ETH</option>
                  <option>BNB</option>
                </select>
                <input type="number" placeholder="0.00" className="amount-input" readOnly />
              </div>
            </div>
            
            <div className="exchange-info">
              <div className="rate-info">
                <span>Exchange Rate:</span>
                <span className="rate-value">1 BTN = 10 USDT</span>
              </div>
              <div className="fee-info">
                <span>Fee (0.1%):</span>
                <span className="fee-value">0.01 USDT</span>
              </div>
            </div>
            
            <button className="btn-primary exchange-btn" onClick={handleExchange}>
              Exchange Now
            </button>
            
            <div className="exchange-features">
              <h3>Features</h3>
              <ul>
                <li>✓ Live exchange rates updated every second</li>
                <li>✓ Low fees (0.1% per transaction)</li>
                <li>✓ Instant conversion</li>
                <li>✓ Support for all major cryptocurrencies</li>
              </ul>
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

        {activeTab === 'insights' && (
          <div className="insights-tab">
            <h2>AI-Driven Insights</h2>
            
            <div className="alerts-section">
              <h3>Recent Alerts</h3>
              <div className="alerts-list">
                {alerts.map((alert) => (
                  <div key={alert.id} className={`alert-item ${alert.severity}`}>
                    <div className="alert-header">
                      <span className="alert-type">{alert.type}</span>
                      <span className={`severity-badge ${alert.severity}`}>{alert.severity}</span>
                    </div>
                    <p className="alert-message">{alert.message}</p>
                  </div>
                ))}
              </div>
            </div>
            
            <div className="insights-section">
              <h3>Portfolio Insights</h3>
              <div className="insights-list">
                {insights.map((insight) => (
                  <div key={insight.id} className="insight-item">
                    <div className="insight-header">
                      <h4>{insight.title}</h4>
                      <span className={`priority-badge ${insight.priority}`}>{insight.priority}</span>
                    </div>
                    <p className="insight-message">{insight.message}</p>
                    <button className="btn-small">Learn More</button>
                  </div>
                ))}
              </div>
            </div>
            
            <div className="recommendations-section">
              <h3>AI Recommendations</h3>
              <div className="recommendations-list">
                <div className="recommendation-item">
                  <h4>💎 Optimize Staking Rewards</h4>
                  <p>Based on your portfolio, you could earn additional rewards by staking your BTN holdings.</p>
                  <div className="recommendation-stats">
                    <span>Potential APY: 5%</span>
                    <span>Confidence: 95%</span>
                  </div>
                  <button className="btn-primary">Stake Now</button>
                </div>
                
                <div className="recommendation-item">
                  <h4>📊 Portfolio Rebalancing</h4>
                  <p>Your portfolio is heavily weighted in volatile assets. Consider diversifying into USDT.</p>
                  <div className="recommendation-stats">
                    <span>Risk Level: Medium</span>
                    <span>Confidence: 85%</span>
                  </div>
                  <button className="btn-secondary">View Details</button>
                </div>
              </div>
            </div>
            
            <div className="gold-reserve-section">
              <h3>🏅 Gold Reserve Backing</h3>
              <div className="reserve-info">
                <p><strong>Total Reserve:</strong> $2.689 Trillion</p>
                <p><strong>Your BTN Backing:</strong> ${(balance.btn * 10).toFixed(2)}</p>
                <p><strong>Backing Ratio:</strong> 100%</p>
                <p><strong>Last Verified:</strong> 2 minutes ago</p>
                <span className="verified-badge">✓ Verified by Consensus</span>
              </div>
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
                  <h4>🛡️ Real-Time Fraud Detection</h4>
                  <p>AI-powered fraud detection and prevention</p>
                </div>
                <label className="toggle-switch">
                  <input
                    type="checkbox"
                    checked={securitySettings.fraudDetectionEnabled}
                    onChange={() => setSecuritySettings({...securitySettings, fraudDetectionEnabled: !securitySettings.fraudDetectionEnabled})}
                  />
                  <span className="toggle-slider"></span>
                </label>
              </div>
            </div>
            
            <div className="security-section">
              <h3>Device Management</h3>
              <div className="devices-list">
                <div className="device-item">
                  <div className="device-info">
                    <h4>💻 Current Device</h4>
                    <p>Chrome on Windows • Last used: Just now</p>
                  </div>
                  <span className="trusted-badge">✓ Trusted</span>
                </div>
                <div className="device-item">
                  <div className="device-info">
                    <h4>📱 iPhone 14</h4>
                    <p>Safari on iOS • Last used: 2 hours ago</p>
                  </div>
                  <button className="btn-small">Remove</button>
                </div>
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
