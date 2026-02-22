import type { Metadata } from 'next';
import Script from 'next/script';
import BTNGDemoClient from './BTNGDemoClient';
import './btng-demo.css';

export const metadata: Metadata = {
  title: 'BTNG Demo Sandbox — Bituncoin Gold Blockchain',
  description: 'Interactive demonstration of BTNG wallet, mining, explorer, market and QR identity features.',
};

export default function BTNGDemoPage() {
  return (
    <>
      <BTNGDemoClient />
      
      {/* Load QR Code Library */}
      <Script
        src="https://cdn.jsdelivr.net/npm/qrcode@1.5.1/build/qrcode.min.js"
        strategy="beforeInteractive"
      />
      
      {/* Demo Logic Script */}
      <Script
        src="/btng-demo/btng-demo.js"
        strategy="lazyOnload"
      />
      
      <div className="btng-demo-wrapper">
      {/* Demo Toast */}
      <div id="toast"><span id="toast-msg"></span></div>

      {/* Demo Header (embedded, sovereign header remains above) */}
      <header className="demo-header">
        <div className="logo">
          <div className="logo-icon">₿G</div>
          <div className="logo-text">
            <span>Bituncoin Gold</span>
            <span>BTNG Blockchain Sandbox</span>
          </div>
        </div>
        <div className="header-stats">
          <div className="hstat">
            <span className="hstat-label">Price</span>
            <span className="hstat-val" id="h-price">$0.00</span>
          </div>
          <div className="hstat">
            <span className="hstat-label">Block</span>
            <span className="hstat-val" id="h-block">#0</span>
          </div>
          <div className="hstat">
            <span className="hstat-label">Network</span>
            <span className="hstat-val" id="h-net-hash">24.5 TH/s</span>
          </div>
          <div className="live-badge">
            <div className="live-dot"></div>LIVE
          </div>
        </div>
      </header>

      {/* Demo Navigation */}
      <nav className="demo-nav">
        <button className="nav-btn active" data-section="dashboard">
          📊 Dashboard
        </button>
        <button className="nav-btn" data-section="wallet">
          💰 Wallet
        </button>
        <button className="nav-btn" data-section="mining">
          ⛏️ Mining
        </button>
        <button className="nav-btn" data-section="explorer">
          🔍 Explorer
        </button>
        <button className="nav-btn" data-section="market">
          📈 Market
        </button>
      </nav>

      <main className="demo-main">
        {/* DASHBOARD */}
        <div className="section active" id="sec-dashboard">
          <div className="grid-4">
            <div className="card">
              <div className="card-title">BTNG Price</div>
              <div className="card-value" id="d-price">$0.00</div>
              <div className="card-sub" id="d-pch">+0.00% today</div>
            </div>
            <div className="card">
              <div className="card-title">Market Cap</div>
              <div className="card-value" id="d-mcap">$0</div>
              <div className="card-sub">Circulating supply</div>
            </div>
            <div className="card">
              <div className="card-title">Blocks Mined</div>
              <div className="card-value" id="d-blocks">0</div>
              <div className="card-sub">~120s avg time</div>
            </div>
            <div className="card">
              <div className="card-title">Total Transactions</div>
              <div className="card-value" id="d-txcount">0</div>
              <div className="card-sub">All time</div>
            </div>
          </div>
          <div className="grid-2">
            <div className="card">
              <div className="section-title">Price Chart (24h)</div>
              <div className="chart-wrapper">
                <canvas id="priceChart"></canvas>
              </div>
            </div>
            <div className="card">
              <div className="section-title">Recent Transactions</div>
              <table className="tx-table">
                <thead>
                  <tr>
                    <th>Hash</th>
                    <th>Amount</th>
                    <th>Status</th>
                  </tr>
                </thead>
                <tbody id="recent-tx-body"></tbody>
              </table>
            </div>
          </div>
        </div>

        {/* WALLET (with QR) */}
        <div className="section" id="sec-wallet">
          <div className="grid-2">
            <div className="card">
              <div className="section-title">My Wallet</div>
              <div className="balance-display">
                <div className="balance-btng" id="w-balance">0.000</div>
                <div className="balance-usd" id="w-usd">≈ $0.00 USD</div>
              </div>
              <div className="card-title">Wallet Address</div>
              <div className="wallet-address-box">
                <span id="w-address">Loading...</span>
                <button className="copy-btn" id="copy-address-btn">Copy</button>
              </div>

              <div className="qr-container">
                <canvas id="wallet-qr" className="qr-canvas" width="160" height="160"></canvas>
                <button className="copy-btn show-qr-btn" id="toggle-qr-btn">
                  📷 Show QR
                </button>
              </div>

              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '10px', marginTop: '8px' }}>
                <div className="mcard">
                  <div className="mcard-label">Total Sent</div>
                  <div className="mcard-val" id="w-sent">0 BTNG</div>
                </div>
                <div className="mcard">
                  <div className="mcard-label">Total Received</div>
                  <div className="mcard-val" id="w-recv">0 BTNG</div>
                </div>
              </div>
            </div>
            <div className="card">
              <div className="section-title">Send BTNG</div>
              <div className="input-group">
                <label>Recipient Address</label>
                <input type="text" id="send-to" placeholder="BTNG1..." />
              </div>
              <div className="input-group">
                <label>Amount (BTNG)</label>
                <input type="number" id="send-amount" placeholder="0.00" step="0.01" min="0" />
              </div>
              <div className="input-group">
                <label>Network Fee</label>
                <select id="send-fee">
                  <option value="0.001">Slow — 0.001 BTNG (~5 min)</option>
                  <option value="0.005">Standard — 0.005 BTNG (~2 min)</option>
                  <option value="0.01">Fast — 0.01 BTNG (~30 sec)</option>
                </select>
              </div>
              <button className="btn-primary" id="send-transaction-btn">
                🚀 Send Transaction
              </button>
              <div className="section-title" style={{ marginTop: '24px' }}>History</div>
              <div style={{ overflow: 'auto', maxHeight: '220px' }}>
                <table className="tx-table">
                  <thead>
                    <tr>
                      <th>Hash</th>
                      <th>Type</th>
                      <th>Amount</th>
                      <th>Status</th>
                    </tr>
                  </thead>
                  <tbody id="wallet-tx-body"></tbody>
                </table>
              </div>
            </div>
          </div>
        </div>

        {/* MINING */}
        <div className="section" id="sec-mining">
          <div className="grid-2">
            <div className="card">
              <div className="section-title">Mining Control</div>
              <div className="mining-display">
                <span className="mining-icon" id="mine-icon">⛏️</span>
                <div className="hash-rate-display" id="hashrate-disp">0 MH/s</div>
                <div style={{ fontSize: '12px', color: 'var(--text2)', marginTop: '4px' }}>Current Hash Rate</div>
              </div>
              <div className="progress-bar-wrapper">
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '8px' }}>
                  <span style={{ fontSize: '12px', color: 'var(--text2)' }}>Block Progress</span>
                  <span style={{ fontSize: '12px', color: 'var(--gold)' }} id="mine-pct">0%</span>
                </div>
                <div className="progress-bar-track">
                  <div className="progress-bar-fill" id="mine-bar"></div>
                </div>
              </div>
              <button className="mine-btn" id="mine-btn">
                ⛏️ Start Mining
              </button>
              <div className="mining-stats">
                <div className="mstat">
                  <div className="mstat-label">Blocks Found</div>
                  <div className="mstat-val" id="mine-blocks">0</div>
                </div>
                <div className="mstat">
                  <div className="mstat-label">Total Earned</div>
                  <div className="mstat-val" id="mine-earned">0 BTNG</div>
                </div>
                <div className="mstat">
                  <div className="mstat-label">Difficulty</div>
                  <div className="mstat-val" id="mine-diff">1</div>
                </div>
              </div>
            </div>
            <div className="card">
              <div className="section-title">Mining Log</div>
              <div className="log-box" id="mine-log">
                <div className="log-line">
                  <span className="log-dim">[00:00:00]</span> Mining engine ready. Press Start to begin.
                </div>
              </div>
              <div className="section-title" style={{ marginTop: '20px' }}>Network Info</div>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '12px' }}>
                <div className="mcard">
                  <div className="mcard-label">Network Hash Rate</div>
                  <div className="mcard-val" id="net-hash">24.5 TH/s</div>
                </div>
                <div className="mcard">
                  <div className="mcard-label">Block Reward</div>
                  <div className="mcard-val">50 BTNG</div>
                </div>
                <div className="mcard">
                  <div className="mcard-label">Next Halving</div>
                  <div className="mcard-val" id="halving-in">21,000</div>
                </div>
                <div className="mcard">
                  <div className="mcard-label">Avg Block Time</div>
                  <div className="mcard-val">~120s</div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* EXPLORER */}
        <div className="section" id="sec-explorer">
          <div className="search-bar">
            <input type="text" id="exp-search" placeholder="Search block number, TX hash, or address..." />
            <button id="search-btn">🔍 Search</button>
          </div>
          <div className="grid-2">
            <div className="card">
              <div className="section-title">Latest Blocks</div>
              <div id="block-list"></div>
            </div>
            <div className="card">
              <div className="section-title">Latest Transactions</div>
              <table className="tx-table">
                <thead>
                  <tr>
                    <th>Hash</th>
                    <th>Amount</th>
                    <th>Status</th>
                    <th>Age</th>
                  </tr>
                </thead>
                <tbody id="exp-tx-body"></tbody>
              </table>
            </div>
          </div>
        </div>

        {/* MARKET (with APK QR) */}
        <div className="section" id="sec-market">
          <div className="card" style={{ marginBottom: '20px' }}>
            <div className="price-hero">
              <div>
                <div className="card-title">BTNG / USD</div>
                <div className="price-main" id="m-price">$0.00</div>
                <div className="price-change-pos" id="m-change">+0.00%</div>
              </div>
              <div style={{ textAlign: 'right' }}>
                <div style={{ fontSize: '12px', color: 'var(--text3)', marginBottom: '6px' }}>24H Range</div>
                <div style={{ fontFamily: 'JetBrains Mono, monospace', fontSize: '14px', color: 'var(--text2)' }} id="m-range">
                  $0 – $0
                </div>
                <div style={{ fontSize: '12px', color: 'var(--text3)', marginTop: '8px', marginBottom: '6px' }}>
                  Volume (24H)
                </div>
                <div style={{ fontFamily: 'JetBrains Mono, monospace', fontSize: '14px', color: 'var(--text2)' }} id="m-vol">
                  $0
                </div>
              </div>
            </div>
            <div className="chart-wrapper" style={{ height: '220px' }}>
              <canvas id="marketChart"></canvas>
            </div>
          </div>
          <div className="market-grid">
            <div className="mcard">
              <div className="mcard-label">Market Cap</div>
              <div className="mcard-val" id="m-mcap">$0</div>
            </div>
            <div className="mcard">
              <div className="mcard-label">Circulating Supply</div>
              <div className="mcard-val" id="m-supply">0 BTNG</div>
            </div>
            <div className="mcard">
              <div className="mcard-label">Max Supply</div>
              <div className="mcard-val">21,000,000 BTNG</div>
            </div>
            <div className="mcard">
              <div className="mcard-label">All-Time High</div>
              <div className="mcard-val" id="m-ath">$0.00</div>
            </div>
            <div className="mcard">
              <div className="mcard-label">All-Time Low</div>
              <div className="mcard-val" id="m-atl">$0.00</div>
            </div>
            <div className="mcard">
              <div className="mcard-label">Global Rank</div>
              <div className="mcard-val">#247</div>
            </div>
          </div>
          <div className="card" style={{ marginTop: '20px' }}>
            <div className="section-title">Download Mobile Banking App</div>
            <div style={{ textAlign: 'center' }}>
              <canvas id="apk-qr" className="qr-canvas" width="160" height="160" style={{ display: 'block', margin: '0 auto' }}></canvas>
              <p style={{ color: 'var(--text2)', fontSize: '12px', marginTop: '8px' }}>
                Scan to download BTNG Wallet APK
              </p>
              <a href="https://example.com/btng-wallet.apk" target="_blank" rel="noopener noreferrer" className="copy-btn" style={{display: 'inline-block', textDecoration: 'none'}}>
                📲 Download APK
              </a>
            </div>
          </div>
        </div>
      </main>
    </div>
    </>
  );
}
