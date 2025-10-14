import React, { useState, useEffect } from 'react';

// Universal Wallet Component for Bituncoin Blockchain
const Wallet = () => {
  const [activeTab, setActiveTab] = useState('dashboard');
  const [walletAddress, setWalletAddress] = useState('');
  const [balances, setBalances] = useState({});
  const [transactions, setTransactions] = useState([]);
  const [selectedCurrency, setSelectedCurrency] = useState('BTN');
  const [enable2FA, setEnable2FA] = useState(false);
  const [enableBiometric, setEnableBiometric] = useState(false);
  const [isLoading, setIsLoading] = useState(false);

  // Supported currencies
  const supportedCurrencies = ['BTN', 'BTC', 'ETH', 'USDT', 'BNB'];

  // API endpoint configuration
  const API_BASE = 'http://localhost:8080/api';

  // Create new wallet
  const createWallet = async () => {
    setIsLoading(true);
    try {
      const response = await fetch(`${API_BASE}/wallet/create`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          enable_2fa: enable2FA,
          enable_biometric: enableBiometric,
          biometric_data: enableBiometric ? 'sample-biometric-data' : ''
        })
      });

      const data = await response.json();
      if (data.success) {
        setWalletAddress(data.data.address);
        alert('Wallet created successfully!\nAddress: ' + data.data.address);
        loadBalance(data.data.address);
      } else {
        alert('Failed to create wallet: ' + data.message);
      }
    } catch (error) {
      console.error('Error creating wallet:', error);
      alert('Error creating wallet');
    } finally {
      setIsLoading(false);
    }
  };

  // Load wallet balance
  const loadBalance = async (address) => {
    if (!address) return;
    
    try {
      const response = await fetch(`${API_BASE}/wallet/balance?address=${address}`);
      const data = await response.json();
      if (data.success) {
        setBalances(data.data);
      }
    } catch (error) {
      console.error('Error loading balance:', error);
    }
  };

  // Load transaction history
  const loadTransactions = async (address) => {
    if (!address) return;
    
    try {
      const response = await fetch(`${API_BASE}/transaction/history?address=${address}`);
      const data = await response.json();
      if (data.success) {
        setTransactions(data.data || []);
      }
    } catch (error) {
      console.error('Error loading transactions:', error);
    }
  };

  // Send transaction
  const sendTransaction = async (to, amount, currency, crossChain, targetChain) => {
    if (!walletAddress) {
      alert('Please create or load a wallet first');
      return;
    }

    setIsLoading(true);
    try {
      const response = await fetch(`${API_BASE}/transaction/send`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          from: walletAddress,
          to: to,
          amount: parseFloat(amount),
          currency: currency,
          cross_chain: crossChain,
          target_chain: targetChain
        })
      });

      const data = await response.json();
      if (data.success) {
        alert('Transaction sent successfully!\nTX ID: ' + data.data.transaction_id);
        loadBalance(walletAddress);
        loadTransactions(walletAddress);
      } else {
        alert('Transaction failed: ' + data.message);
      }
    } catch (error) {
      console.error('Error sending transaction:', error);
      alert('Error sending transaction');
    } finally {
      setIsLoading(false);
    }
  };

  // Dashboard View
  const DashboardView = () => (
    <div style={styles.dashboardContainer}>
      <h2>Dashboard</h2>
      {walletAddress ? (
        <>
          <div style={styles.addressCard}>
            <strong>Wallet Address:</strong>
            <div style={styles.address}>{walletAddress}</div>
          </div>
          
          <div style={styles.balancesContainer}>
            <h3>Multi-Currency Balances</h3>
            <div style={styles.balanceGrid}>
              {supportedCurrencies.map(currency => (
                <div key={currency} style={styles.balanceCard}>
                  <div style={styles.currencyIcon}>{currency}</div>
                  <div style={styles.balanceAmount}>
                    {balances[currency] ? balances[currency].toFixed(4) : '0.0000'}
                  </div>
                </div>
              ))}
            </div>
          </div>

          <button 
            style={styles.refreshButton}
            onClick={() => {
              loadBalance(walletAddress);
              loadTransactions(walletAddress);
            }}
          >
            Refresh
          </button>
        </>
      ) : (
        <div style={styles.noWallet}>
          <p>No wallet loaded. Please create or import a wallet.</p>
        </div>
      )}
    </div>
  );

  // Send View
  const SendView = () => {
    const [recipient, setRecipient] = useState('');
    const [amount, setAmount] = useState('');
    const [currency, setCurrency] = useState('BTN');
    const [crossChain, setCrossChain] = useState(false);
    const [targetChain, setTargetChain] = useState('');

    return (
      <div style={styles.sendContainer}>
        <h2>Send Transaction</h2>
        <form onSubmit={(e) => {
          e.preventDefault();
          sendTransaction(recipient, amount, currency, crossChain, targetChain);
          setRecipient('');
          setAmount('');
        }}>
          <div style={styles.formGroup}>
            <label>Recipient Address:</label>
            <input
              type="text"
              value={recipient}
              onChange={(e) => setRecipient(e.target.value)}
              style={styles.input}
              placeholder="BTN..."
              required
            />
          </div>

          <div style={styles.formGroup}>
            <label>Amount:</label>
            <input
              type="number"
              step="0.0001"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              style={styles.input}
              placeholder="0.00"
              required
            />
          </div>

          <div style={styles.formGroup}>
            <label>Currency:</label>
            <select 
              value={currency}
              onChange={(e) => setCurrency(e.target.value)}
              style={styles.select}
            >
              {supportedCurrencies.map(curr => (
                <option key={curr} value={curr}>{curr}</option>
              ))}
            </select>
          </div>

          <div style={styles.formGroup}>
            <label style={styles.checkboxLabel}>
              <input
                type="checkbox"
                checked={crossChain}
                onChange={(e) => setCrossChain(e.target.checked)}
              />
              Cross-Chain Transaction
            </label>
          </div>

          {crossChain && (
            <div style={styles.formGroup}>
              <label>Target Chain:</label>
              <select 
                value={targetChain}
                onChange={(e) => setTargetChain(e.target.value)}
                style={styles.select}
              >
                <option value="">Select chain...</option>
                {supportedCurrencies.filter(c => c !== currency).map(curr => (
                  <option key={curr} value={curr}>{curr}</option>
                ))}
              </select>
            </div>
          )}

          <button type="submit" style={styles.sendButton} disabled={isLoading}>
            {isLoading ? 'Sending...' : 'Send Transaction'}
          </button>
        </form>
      </div>
    );
  };

  // Transactions View
  const TransactionsView = () => (
    <div style={styles.transactionsContainer}>
      <h2>Transaction History</h2>
      {transactions.length > 0 ? (
        <div style={styles.transactionsList}>
          {transactions.map((tx, index) => (
            <div key={index} style={styles.transactionCard}>
              <div style={styles.txHeader}>
                <span style={styles.txId}>{tx.id}</span>
                <span style={styles.txCurrency}>{tx.currency}</span>
              </div>
              <div style={styles.txDetails}>
                <div>From: {tx.from.substring(0, 20)}...</div>
                <div>To: {tx.to.substring(0, 20)}...</div>
                <div>Amount: {tx.amount}</div>
                <div>Time: {new Date(tx.timestamp * 1000).toLocaleString()}</div>
                {tx.cross_chain && <div style={styles.crossChainBadge}>Cross-Chain</div>}
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p>No transactions yet.</p>
      )}
    </div>
  );

  // Settings View
  const SettingsView = () => (
    <div style={styles.settingsContainer}>
      <h2>Security Settings</h2>
      
      <div style={styles.securityCard}>
        <h3>Two-Factor Authentication (2FA)</h3>
        <p>Enhance your wallet security with 2FA</p>
        <label style={styles.checkboxLabel}>
          <input
            type="checkbox"
            checked={enable2FA}
            onChange={(e) => setEnable2FA(e.target.checked)}
          />
          Enable 2FA
        </label>
      </div>

      <div style={styles.securityCard}>
        <h3>Biometric Authentication</h3>
        <p>Use biometric authentication for transactions</p>
        <label style={styles.checkboxLabel}>
          <input
            type="checkbox"
            checked={enableBiometric}
            onChange={(e) => setEnableBiometric(e.target.checked)}
          />
          Enable Biometric Auth
        </label>
      </div>

      <div style={styles.securityCard}>
        <h3>Wallet Information</h3>
        {walletAddress ? (
          <>
            <p><strong>Address:</strong> {walletAddress}</p>
            <p><strong>2FA Status:</strong> {enable2FA ? 'Enabled' : 'Disabled'}</p>
            <p><strong>Biometric:</strong> {enableBiometric ? 'Enabled' : 'Disabled'}</p>
          </>
        ) : (
          <p>No wallet loaded</p>
        )}
      </div>
    </div>
  );

  // Create Wallet View
  const CreateWalletView = () => (
    <div style={styles.createWalletContainer}>
      <h2>Create New Wallet</h2>
      <p>Create a new universal wallet with multi-currency support</p>
      
      <div style={styles.securityOptions}>
        <label style={styles.checkboxLabel}>
          <input
            type="checkbox"
            checked={enable2FA}
            onChange={(e) => setEnable2FA(e.target.checked)}
          />
          Enable Two-Factor Authentication
        </label>

        <label style={styles.checkboxLabel}>
          <input
            type="checkbox"
            checked={enableBiometric}
            onChange={(e) => setEnableBiometric(e.target.checked)}
          />
          Enable Biometric Authentication
        </label>
      </div>

      <button 
        style={styles.createButton}
        onClick={createWallet}
        disabled={isLoading}
      >
        {isLoading ? 'Creating...' : 'Create Wallet'}
      </button>

      <div style={styles.infoBox}>
        <h4>Supported Features:</h4>
        <ul>
          <li>Multi-currency support (BTN, BTC, ETH, USDT, BNB)</li>
          <li>Cross-chain transactions</li>
          <li>Two-factor authentication</li>
          <li>Biometric security</li>
          <li>Encrypted key storage</li>
        </ul>
      </div>
    </div>
  );

  // Main render
  return (
    <div style={styles.container}>
      <header style={styles.header}>
        <h1>🔐 Bituncoin Universal Wallet</h1>
        <p>Multi-Currency • Cross-Chain • Secure</p>
      </header>

      <nav style={styles.nav}>
        <button 
          style={activeTab === 'create' ? styles.navButtonActive : styles.navButton}
          onClick={() => setActiveTab('create')}
        >
          Create Wallet
        </button>
        <button 
          style={activeTab === 'dashboard' ? styles.navButtonActive : styles.navButton}
          onClick={() => setActiveTab('dashboard')}
        >
          Dashboard
        </button>
        <button 
          style={activeTab === 'send' ? styles.navButtonActive : styles.navButton}
          onClick={() => setActiveTab('send')}
        >
          Send
        </button>
        <button 
          style={activeTab === 'transactions' ? styles.navButtonActive : styles.navButton}
          onClick={() => setActiveTab('transactions')}
        >
          Transactions
        </button>
        <button 
          style={activeTab === 'settings' ? styles.navButtonActive : styles.navButton}
          onClick={() => setActiveTab('settings')}
        >
          Settings
        </button>
      </nav>

      <main style={styles.main}>
        {activeTab === 'create' && <CreateWalletView />}
        {activeTab === 'dashboard' && <DashboardView />}
        {activeTab === 'send' && <SendView />}
        {activeTab === 'transactions' && <TransactionsView />}
        {activeTab === 'settings' && <SettingsView />}
      </main>

      <footer style={styles.footer}>
        <p>Bituncoin Universal Wallet v1.0 | Secure • Fast • Reliable</p>
      </footer>
    </div>
  );
};

// Styles
const styles = {
  container: {
    fontFamily: 'Arial, sans-serif',
    maxWidth: '1200px',
    margin: '0 auto',
    padding: '20px',
    backgroundColor: '#f5f5f5',
  },
  header: {
    background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    color: 'white',
    padding: '30px',
    borderRadius: '10px',
    textAlign: 'center',
    marginBottom: '20px',
  },
  nav: {
    display: 'flex',
    gap: '10px',
    marginBottom: '20px',
    flexWrap: 'wrap',
  },
  navButton: {
    padding: '12px 24px',
    border: 'none',
    borderRadius: '5px',
    cursor: 'pointer',
    backgroundColor: '#fff',
    fontWeight: 'bold',
    transition: 'all 0.3s',
  },
  navButtonActive: {
    padding: '12px 24px',
    border: 'none',
    borderRadius: '5px',
    cursor: 'pointer',
    backgroundColor: '#667eea',
    color: 'white',
    fontWeight: 'bold',
  },
  main: {
    backgroundColor: 'white',
    borderRadius: '10px',
    padding: '30px',
    minHeight: '400px',
  },
  footer: {
    textAlign: 'center',
    padding: '20px',
    color: '#666',
    marginTop: '20px',
  },
  dashboardContainer: {
    padding: '20px',
  },
  addressCard: {
    backgroundColor: '#f9f9f9',
    padding: '20px',
    borderRadius: '8px',
    marginBottom: '20px',
  },
  address: {
    fontFamily: 'monospace',
    fontSize: '14px',
    marginTop: '10px',
    wordBreak: 'break-all',
  },
  balancesContainer: {
    marginTop: '30px',
  },
  balanceGrid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))',
    gap: '15px',
    marginTop: '15px',
  },
  balanceCard: {
    backgroundColor: '#667eea',
    color: 'white',
    padding: '20px',
    borderRadius: '8px',
    textAlign: 'center',
  },
  currencyIcon: {
    fontSize: '24px',
    fontWeight: 'bold',
    marginBottom: '10px',
  },
  balanceAmount: {
    fontSize: '20px',
    fontWeight: 'bold',
  },
  refreshButton: {
    marginTop: '20px',
    padding: '10px 20px',
    backgroundColor: '#667eea',
    color: 'white',
    border: 'none',
    borderRadius: '5px',
    cursor: 'pointer',
    fontWeight: 'bold',
  },
  noWallet: {
    textAlign: 'center',
    padding: '40px',
    color: '#666',
  },
  sendContainer: {
    maxWidth: '600px',
    margin: '0 auto',
  },
  formGroup: {
    marginBottom: '20px',
  },
  input: {
    width: '100%',
    padding: '12px',
    border: '1px solid #ddd',
    borderRadius: '5px',
    fontSize: '16px',
    boxSizing: 'border-box',
  },
  select: {
    width: '100%',
    padding: '12px',
    border: '1px solid #ddd',
    borderRadius: '5px',
    fontSize: '16px',
    boxSizing: 'border-box',
  },
  checkboxLabel: {
    display: 'flex',
    alignItems: 'center',
    gap: '10px',
    cursor: 'pointer',
  },
  sendButton: {
    width: '100%',
    padding: '15px',
    backgroundColor: '#667eea',
    color: 'white',
    border: 'none',
    borderRadius: '5px',
    fontSize: '16px',
    fontWeight: 'bold',
    cursor: 'pointer',
  },
  transactionsContainer: {
    padding: '20px',
  },
  transactionsList: {
    display: 'flex',
    flexDirection: 'column',
    gap: '15px',
  },
  transactionCard: {
    border: '1px solid #ddd',
    borderRadius: '8px',
    padding: '15px',
    backgroundColor: '#f9f9f9',
  },
  txHeader: {
    display: 'flex',
    justifyContent: 'space-between',
    marginBottom: '10px',
  },
  txId: {
    fontFamily: 'monospace',
    fontSize: '12px',
    color: '#666',
  },
  txCurrency: {
    backgroundColor: '#667eea',
    color: 'white',
    padding: '2px 8px',
    borderRadius: '3px',
    fontSize: '12px',
    fontWeight: 'bold',
  },
  txDetails: {
    fontSize: '14px',
    lineHeight: '1.6',
  },
  crossChainBadge: {
    display: 'inline-block',
    backgroundColor: '#ff6b6b',
    color: 'white',
    padding: '2px 8px',
    borderRadius: '3px',
    fontSize: '12px',
    marginTop: '5px',
  },
  settingsContainer: {
    padding: '20px',
  },
  securityCard: {
    border: '1px solid #ddd',
    borderRadius: '8px',
    padding: '20px',
    marginBottom: '20px',
    backgroundColor: '#f9f9f9',
  },
  createWalletContainer: {
    maxWidth: '600px',
    margin: '0 auto',
    textAlign: 'center',
  },
  securityOptions: {
    textAlign: 'left',
    marginBottom: '30px',
    display: 'flex',
    flexDirection: 'column',
    gap: '15px',
  },
  createButton: {
    width: '100%',
    padding: '15px',
    backgroundColor: '#667eea',
    color: 'white',
    border: 'none',
    borderRadius: '5px',
    fontSize: '18px',
    fontWeight: 'bold',
    cursor: 'pointer',
    marginBottom: '30px',
  },
  infoBox: {
    textAlign: 'left',
    backgroundColor: '#f0f7ff',
    padding: '20px',
    borderRadius: '8px',
    border: '1px solid #b3d9ff',
  },
};

export default Wallet;
