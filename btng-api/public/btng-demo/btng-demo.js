// BTNG Demo Sandbox - Complete Interactive Logic
(function() {
  'use strict';

  // Wait for DOM to be ready
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initBTNGDemo);
  } else {
    initBTNGDemo();
  }

  function initBTNGDemo() {
    // -------------------- STATE --------------------
    const S = {
      price: 4.27,
      priceChange: 8.34,
      supply: 1240000,
      blockHeight: 12457,
      totalTx: 84293,
      wallet: {
        address: '',
        balance: 125.487,
        sent: 47.3,
        received: 172.787
      },
      mining: {
        active: false,
        progress: 0,
        blocksFound: 0,
        earned: 0,
        difficulty: 4
      },
      priceHistory: [],
      blocks: [],
      transactions: [],
      ath: 7.84,
      atl: 0.32
    };

    // -------------------- API FUNCTIONS --------------------
    async function apiCall(endpoint, options = {}) {
      try {
        const response = await fetch(`/api/btng${endpoint}`, {
          headers: {
            'Content-Type': 'application/json',
          },
          ...options,
        });
        const data = await response.json();
        return data;
      } catch (error) {
        console.error('API call failed:', endpoint, error);
        return null;
      }
    }

    async function getWalletBalance(address) {
      const data = await apiCall(`/wallet/balance/${address}`);
      return data?.balance || 125.487;
    }

    async function getWalletTransactions(address) {
      const data = await apiCall(`/wallet/transactions/${address}`);
      return data?.transactions || [];
    }

    async function sendTransactionAPI(txData) {
      const data = await apiCall('/wallet/send', {
        method: 'POST',
        body: JSON.stringify(txData),
      });
      return data;
    }

    async function getBlock(height) {
      const data = await apiCall(`/explorer/block/${height}`);
      return data?.block || null;
    }

    async function getTransaction(hash) {
      const data = await apiCall(`/explorer/tx/${hash}`);
      return data?.transaction || null;
    }

    async function getAddressInfo(address) {
      const data = await apiCall(`/explorer/address/${address}`);
      return data?.address || null;
    }

    async function getMiningInfo() {
      const data = await apiCall('/mining/info');
      return data || { difficulty: 4, networkHashrate: '24.5 TH/s' };
    }

    async function getPriceData() {
      const data = await apiCall('/oracle/price');
      return data || { price: 4.27, change24h: 8.34 };
    }

    async function getMarketCapData() {
      const data = await apiCall('/oracle/marketcap');
      return data || { marketcap: 5108400, supply: 1240000 };
    }
    function r(a, b) { return Math.random() * (b - a) + a; }
    function ri(a, b) { return Math.floor(r(a, b)); }
    function sh(n = 10) {
      const c = '0123456789abcdef';
      return '0x' + Array.from({ length: n }, () => c[ri(0, 16)]).join('');
    }
    function ga() {
      const c = '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ';
      return 'BTNG1' + Array.from({ length: 32 }, () => c[ri(0, c.length)]).join('');
    }
    function fm(n, d = 2) { return n.toLocaleString('en-US', { minimumFractionDigits: d, maximumFractionDigits: d }); }
    function fu(n) { return '$' + fm(n, 2); }
    function fb(n) {
      if (n >= 1e9) return '$' + fm(n / 1e9, 2) + 'B';
      if (n >= 1e6) return '$' + fm(n / 1e6, 2) + 'M';
      if (n >= 1e3) return '$' + fm(n / 1e3, 2) + 'K';
      return '$' + fm(n, 2);
    }
    function ta(s) {
      if (s < 60) return s + 's ago';
      if (s < 3600) return Math.floor(s / 60) + 'm ago';
      return Math.floor(s / 3600) + 'h ago';
    }
    function gt() {
      const n = new Date();
      return [n.getHours(), n.getMinutes(), n.getSeconds()].map(x => String(x).padStart(2, '0')).join(':');
    }
    function showToast(msg) {
      const t = document.getElementById('toast');
      if (!t) return;
      const msgEl = document.getElementById('toast-msg');
      if (msgEl) msgEl.textContent = msg;
      t.classList.add('show');
      clearTimeout(t._t);
      t._t = setTimeout(() => t.classList.remove('show'), 3500);
    }

    function copyAddress() {
      const addr = S.wallet.address;
      if (navigator.clipboard) {
        navigator.clipboard.writeText(addr).then(() => showToast('✅ Address copied!')).catch(() => fallbackCopy());
      } else {
        fallbackCopy();
      }
      function fallbackCopy() {
        const textarea = document.createElement('textarea');
        textarea.value = addr;
        document.body.appendChild(textarea);
        textarea.select();
        try {
          document.execCommand('copy');
          showToast('✅ Address copied!');
        } catch (err) {
          showToast('❌ Could not copy');
        }
        document.body.removeChild(textarea);
      }
    }
    window.copyAddress = copyAddress;

    function trimArrays() {
      if (S.blocks.length > 50) S.blocks = S.blocks.slice(-50);
      if (S.transactions.length > 100) S.transactions = S.transactions.slice(-100);
    }

    // -------------------- QR CODE FUNCTIONS --------------------
    function generateWalletQR() {
      if (typeof QRCode === 'undefined') {
        setTimeout(generateWalletQR, 500);
        return;
      }
      const canvas = document.getElementById('wallet-qr');
      if (!canvas) return;
      QRCode.toCanvas(canvas, S.wallet.address, { width: 160, margin: 1 }, function(error) {
        if (error) console.error('QR error:', error);
      });
    }

    function generateAPKQR() {
      if (typeof QRCode === 'undefined') {
        setTimeout(generateAPKQR, 500);
        return;
      }
      const canvas = document.getElementById('apk-qr');
      if (!canvas) return;
      const apkUrl = 'https://example.com/btng-wallet.apk';
      QRCode.toCanvas(canvas, apkUrl, { width: 160, margin: 1 }, function(error) {
        if (error) console.error('APK QR error:', error);
      });
    }

    window.toggleWalletQR = function() {
      const canvas = document.getElementById('wallet-qr');
      const btns = document.querySelectorAll('.show-qr-btn');
      if (!canvas) return;
      if (canvas.style.display === 'none' || canvas.style.display === '') {
        canvas.style.display = 'block';
        btns.forEach(btn => btn.textContent = '📷 Hide QR');
      } else {
        canvas.style.display = 'none';
        btns.forEach(btn => btn.textContent = '📷 Show QR');
      }
    };

    // -------------------- INIT --------------------
    async function initData() {
      S.wallet.address = ga();

      // Try to load real price data
      const priceData = await getPriceData();
      S.price = priceData.price || 4.27;
      S.priceChange = priceData.change24h || 8.34;

      // Try to load market cap data
      const marketData = await getMarketCapData();
      S.supply = marketData.supply || 1240000;

      // Try to load mining info
      const miningData = await getMiningInfo();
      S.mining.difficulty = miningData.difficulty || 4;

      // Generate price history (keep mock for now, could be enhanced later)
      let p = S.price * 0.8;
      for (let i = 48; i >= 0; i--) {
        p *= (1 + r(-0.03, 0.04));
        p = Math.max(0.5, p);
        S.priceHistory.push(parseFloat(p.toFixed(4)));
      }
      S.price = S.priceHistory[S.priceHistory.length - 1];
      S.ath = Math.max(...S.priceHistory) * 1.4;
      S.atl = Math.min(...S.priceHistory) * 0.5;
      S.priceChange = ((S.price - S.priceHistory[0]) / S.priceHistory[0]) * 100;

      // Load real wallet data if available
      const balance = await getWalletBalance(S.wallet.address);
      S.wallet.balance = balance;

      const transactions = await getWalletTransactions(S.wallet.address);
      if (transactions.length > 0) {
        S.transactions = transactions.map(tx => ({
          hash: tx.hash,
          from: tx.from,
          to: tx.to,
          amount: tx.amount,
          status: tx.status || 'confirmed',
          time: tx.timestamp || ri(5, 600)
        }));
      } else {
        // Fallback to mock transactions
        for (let i = 0; i < 20; i++) {
          S.transactions.unshift({
            hash: sh(12),
            from: ga().slice(0, 20) + '...',
            to: ga().slice(0, 20) + '...',
            amount: parseFloat(r(0.01, 500).toFixed(4)),
            status: ['confirmed', 'confirmed', 'confirmed', 'pending', 'failed'][ri(0, 5)],
            time: ri(5, 600)
          });
        }
      }

      // Mock blocks (could be enhanced to load real blocks)
      for (let i = 0; i < 10; i++) {
        S.blocks.unshift({
          number: S.blockHeight - i,
          hash: sh(16),
          txCount: ri(3, 40),
          time: ri(10, 300 * i + 30),
          size: ri(800, 5000)
        });
      }

      trimArrays();
      setTimeout(() => {
        generateWalletQR();
        generateAPKQR();
      }, 500);
    }

    // -------------------- CHART --------------------
    function drawChart(id, data) {
      const canvas = document.getElementById(id);
      if (!canvas) return;
      const parent = canvas.parentElement;
      if (!parent) return;
      const rect = parent.getBoundingClientRect();
      if (rect.width === 0 || rect.height === 0) return;
      canvas.width = rect.width;
      canvas.height = rect.height;
      const ctx = canvas.getContext('2d');
      const W = canvas.width, H = canvas.height;
      const pad = { t: 10, r: 10, b: 30, l: 50 };
      const cw = W - pad.l - pad.r;
      const ch = H - pad.t - pad.b;
      ctx.clearRect(0, 0, W, H);

      const mn = Math.min(...data) * 0.97;
      const mx = Math.max(...data) * 1.03;
      const tx = i => pad.l + (i / (data.length - 1)) * cw;
      const ty = v => pad.t + ch - ((v - mn) / (mx - mn)) * ch;

      ctx.strokeStyle = 'rgba(255,255,255,0.04)';
      ctx.lineWidth = 1;
      for (let i = 0; i <= 4; i++) {
        const y = pad.t + (i / 4) * ch;
        ctx.beginPath();
        ctx.moveTo(pad.l, y);
        ctx.lineTo(pad.l + cw, y);
        ctx.stroke();
      }

      ctx.fillStyle = '#555';
      ctx.font = '10px JetBrains Mono';
      ctx.textAlign = 'right';
      for (let i = 0; i <= 3; i++) {
        const v = mn + ((3 - i) / 3) * (mx - mn);
        ctx.fillText('$' + v.toFixed(2), pad.l - 5, pad.t + (i / 3) * ch + 4);
      }

      const g = ctx.createLinearGradient(0, pad.t, 0, pad.t + ch);
      g.addColorStop(0, 'rgba(245,197,24,0.3)');
      g.addColorStop(1, 'rgba(245,197,24,0)');
      ctx.beginPath();
      data.forEach((v, i) => {
        const x = tx(i), y = ty(v);
        i === 0 ? ctx.moveTo(x, y) : ctx.lineTo(x, y);
      });
      ctx.lineTo(tx(data.length - 1), pad.t + ch);
      ctx.lineTo(tx(0), pad.t + ch);
      ctx.closePath();
      ctx.fillStyle = g;
      ctx.fill();

      ctx.beginPath();
      data.forEach((v, i) => {
        const x = tx(i), y = ty(v);
        i === 0 ? ctx.moveTo(x, y) : ctx.lineTo(x, y);
      });
      ctx.strokeStyle = '#F5C518';
      ctx.lineWidth = 2.5;
      ctx.lineJoin = 'round';
      ctx.lineCap = 'round';
      ctx.shadowColor = 'rgba(245,197,24,0.5)';
      ctx.shadowBlur = 8;
      ctx.stroke();
      ctx.shadowBlur = 0;

      const lx = tx(data.length - 1), ly = ty(data[data.length - 1]);
      ctx.beginPath();
      ctx.arc(lx, ly, 5, 0, Math.PI * 2);
      ctx.fillStyle = '#F5C518';
      ctx.fill();
      ctx.strokeStyle = '#0A0A0A';
      ctx.lineWidth = 2;
      ctx.stroke();

      ctx.fillStyle = '#555';
      ctx.font = '9px Inter';
      ctx.textAlign = 'center';
      ['24h ago', '18h', '12h', '6h', 'Now'].forEach((l, i) => ctx.fillText(l, pad.l + (i / 4) * cw, H - 6));
    }

    // -------------------- RENDER FUNCTIONS --------------------
    function renderDashboard() {
      const el = (id) => document.getElementById(id);
      if (el('d-price')) el('d-price').textContent = fu(S.price);
      const ch = S.priceChange;
      if (el('d-pch')) {
        el('d-pch').textContent = (ch >= 0 ? '+' : '') + ch.toFixed(2) + '% today';
        el('d-pch').style.color = ch >= 0 ? 'var(--green)' : 'var(--red)';
      }

      // Try to get real market cap data
      getMarketCapData().then(marketData => {
        const mcap = marketData?.marketcap || (S.price * S.supply);
        if (el('d-mcap')) el('d-mcap').textContent = fb(mcap);
      }).catch(() => {
        if (el('d-mcap')) el('d-mcap').textContent = fb(S.price * S.supply);
      });

      if (el('d-blocks')) el('d-blocks').textContent = S.blockHeight.toLocaleString();
      if (el('d-txcount')) el('d-txcount').textContent = S.totalTx.toLocaleString();
      if (el('h-price')) el('h-price').textContent = fu(S.price);
      if (el('h-block')) el('h-block').textContent = '#' + S.blockHeight.toLocaleString();

      const tbody = el('recent-tx-body');
      if (tbody) {
        tbody.innerHTML = '';
        S.transactions.slice(0, 6).forEach(tx => {
          const tr = document.createElement('tr');
          tr.innerHTML = `<td class="tx-hash">${tx.hash.slice(0, 14)}...</td><td>${fm(tx.amount, 4)}</td><td><span class="badge badge-${tx.status}">${tx.status}</span></td>`;
          tbody.appendChild(tr);
        });
      }
    }

    async function renderWallet() {
      const el = (id) => document.getElementById(id);
      const w = S.wallet;

      // Try to get real balance
      const realBalance = await getWalletBalance(w.address);
      w.balance = realBalance;

      if (el('w-balance')) el('w-balance').textContent = fm(w.balance, 4);
      if (el('w-usd')) el('w-usd').textContent = '≈ ' + fu(w.balance * S.price) + ' USD';
      if (el('w-address')) el('w-address').textContent = w.address.slice(0, 30) + '...';
      if (el('w-sent')) el('w-sent').textContent = fm(w.sent, 3) + ' BTNG';
      if (el('w-recv')) el('w-recv').textContent = fm(w.received, 3) + ' BTNG';

      generateWalletQR();

      // Try to get real transactions
      const realTransactions = await getWalletTransactions(w.address);
      if (realTransactions && realTransactions.length > 0) {
        // Use real transactions
        const tbody = el('wallet-tx-body');
        if (tbody) {
          tbody.innerHTML = '';
          realTransactions.slice(0, 12).forEach(tx => {
            const isSent = tx.from === w.address;
            const type = isSent ? '↑ Sent' : '↓ Received';
            const cls = isSent ? 'amount-negative' : 'amount-positive';
            const sign = isSent ? '-' : '+';
            const tr = document.createElement('tr');
            tr.innerHTML = `<td class="tx-hash">${tx.hash.slice(0, 12)}...</td><td style="color:var(--text2)">${type}</td><td class="${cls}">${sign}${fm(tx.amount, 4)}</td><td><span class="badge badge-${tx.status || 'confirmed'}">${tx.status || 'confirmed'}</span></td>`;
            tbody.appendChild(tr);
          });
        }
      } else {
        // Fallback to mock transactions
        const tbody = el('wallet-tx-body');
        if (tbody) {
          tbody.innerHTML = '';
          S.transactions.slice(0, 12).forEach(tx => {
            const isSent = tx.from.includes(S.wallet.address.slice(0, 20));
            const type = isSent ? '↑ Sent' : '↓ Received';
            const cls = isSent ? 'amount-negative' : 'amount-positive';
            const sign = isSent ? '-' : '+';
            const tr = document.createElement('tr');
            tr.innerHTML = `<td class="tx-hash">${tx.hash.slice(0, 12)}...</td><td style="color:var(--text2)">${type}</td><td class="${cls}">${sign}${fm(tx.amount, 4)}</td><td><span class="badge badge-${tx.status}">${tx.status}</span></td>`;
            tbody.appendChild(tr);
          });
        }
      }
    }

    function renderExplorer() {
      const el = (id) => document.getElementById(id);
      const bl = el('block-list');
      if (bl) {
        bl.innerHTML = '';
        [...S.blocks].reverse().slice(0, 8).forEach(b => {
          const d = document.createElement('div');
          d.className = 'block-card';
          d.innerHTML = `<div class="block-num">#${b.number.toLocaleString()}</div><div class="block-info"><div class="block-hash">${b.hash}</div><div class="block-meta">${ta(b.time)} &bull; ${b.size.toLocaleString()} bytes</div></div><div class="block-txcount">${b.txCount} txs</div>`;
          bl.appendChild(d);
        });
      }

      const tbody = el('exp-tx-body');
      if (tbody) {
        tbody.innerHTML = '';
        S.transactions.slice(0, 10).forEach(tx => {
          const tr = document.createElement('tr');
          tr.innerHTML = `<td class="tx-hash">${tx.hash.slice(0, 12)}...</td><td class="amount-positive">${fm(tx.amount, 4)}</td><td><span class="badge badge-${tx.status}">${tx.status}</span></td><td style="color:var(--text3);font-size:11px">${ta(tx.time)}</td>`;
          tbody.appendChild(tr);
        });
      }
    }

    function renderMarket() {
      const el = (id) => document.getElementById(id);
      const p = S.price, ch = S.priceChange;
      const lo = Math.min(...S.priceHistory);
      const hi = Math.max(...S.priceHistory);
      S.ath = Math.max(S.ath, p);
      S.atl = Math.min(S.atl, p);

      if (el('m-price')) el('m-price').textContent = fu(p);
      if (el('m-change')) {
        el('m-change').textContent = (ch >= 0 ? '+' : '') + ch.toFixed(2) + '%';
        el('m-change').className = ch >= 0 ? 'price-change-pos' : 'price-change-neg';
      }
      if (el('m-range')) el('m-range').textContent = fu(lo) + ' – ' + fu(hi);
      if (el('m-vol')) el('m-vol').textContent = fb(p * S.supply * 0.04);

      // Try to get real market cap data
      getMarketCapData().then(marketData => {
        const mcap = marketData?.marketcap || (p * S.supply);
        const supply = marketData?.supply || S.supply;
        if (el('m-mcap')) el('m-mcap').textContent = fb(mcap);
        if (el('m-supply')) el('m-supply').textContent = supply.toLocaleString() + ' BTNG';
      }).catch(() => {
        if (el('m-mcap')) el('m-mcap').textContent = fb(p * S.supply);
        if (el('m-supply')) el('m-supply').textContent = S.supply.toLocaleString() + ' BTNG';
      });

      if (el('m-ath')) el('m-ath').textContent = fu(S.ath);
      if (el('m-atl')) el('m-atl').textContent = fu(S.atl);
    }

    function updateNetworkHash() {
      // Try to get real mining data
      getMiningInfo().then(miningData => {
        const base = miningData?.networkHashrate ? parseFloat(miningData.networkHashrate) : 24.5;
        const factor = S.mining.active ? 1 + Math.sin(Date.now() / 5000) * 0.1 : 1;
        const val = (base * factor).toFixed(1);
        const el = (id) => document.getElementById(id);
        if (el('h-net-hash')) el('h-net-hash').textContent = val + ' TH/s';
        if (el('net-hash')) el('net-hash').textContent = val + ' TH/s';
      }).catch(() => {
        // Fallback
        const base = 24.5;
        const factor = S.mining.active ? 1 + Math.sin(Date.now() / 5000) * 0.1 : 1;
        const val = (base * factor).toFixed(1);
        const el = (id) => document.getElementById(id);
        if (el('h-net-hash')) el('h-net-hash').textContent = val + ' TH/s';
        if (el('net-hash')) el('net-hash').textContent = val + ' TH/s';
      });
    }

    // -------------------- SECTION SWITCH --------------------
    window.showSection = async function(name, btn) {
      document.querySelectorAll('.section').forEach(s => s.classList.remove('active'));
      document.querySelectorAll('.nav-btn').forEach(b => b.classList.remove('active'));
      const sec = document.getElementById('sec-' + name);
      if (sec) sec.classList.add('active');
      if (btn) btn.classList.add('active');
      
      if (name === 'dashboard') {
        renderDashboard();
        setTimeout(() => drawChart('priceChart', S.priceHistory), 50);
      }
      if (name === 'wallet') await renderWallet();
      if (name === 'explorer') renderExplorer();
      if (name === 'market') {
        renderMarket();
        setTimeout(() => drawChart('marketChart', S.priceHistory), 50);
      }
    };

    // -------------------- SEND TRANSACTION --------------------
    window.sendTransaction = async function() {
      const el = (id) => document.getElementById(id);
      const to = el('send-to')?.value.trim();
      const amt = parseFloat(el('send-amount')?.value);
      const fee = parseFloat(el('send-fee')?.value);
      
      if (!to) return showToast('⚠️ Enter a recipient address');
      if (!/^BTNG1[A-Z0-9]{32,}$/.test(to)) return showToast('⚠️ Address must start with BTNG1 and be at least 33 chars');
      if (!amt || amt <= 0) return showToast('⚠️ Enter a valid amount');
      if (amt + fee > S.wallet.balance) return showToast('❌ Insufficient balance');

      S.wallet.balance = parseFloat((S.wallet.balance - amt - fee).toFixed(6));
      S.wallet.sent = parseFloat((S.wallet.sent + amt + fee).toFixed(6));
      S.totalTx++;

      const tx = {
        hash: sh(12),
        from: S.wallet.address,
        to: to,
        amount: amt,
        fee: fee,
        status: 'pending',
        time: 5
      };
      S.transactions.unshift(tx);
      trimArrays();

      setTimeout(() => {
        const found = S.transactions.find(t => t.hash === tx.hash);
        if (found) found.status = 'confirmed';
        const walletSec = document.getElementById('sec-wallet');
        if (walletSec?.classList.contains('active')) renderWallet();
      }, 3000);

      if (el('send-to')) el('send-to').value = '';
      if (el('send-amount')) el('send-amount').value = '';
      renderWallet();
      showToast('🚀 Sent ' + fm(amt, 4) + ' BTNG! Confirming...');
    };

    // -------------------- MINING --------------------
    let mInt = null, lInt = null;
    const BLOCK_TARGET = 120;

    window.toggleMining = function() {
      const el = (id) => document.getElementById(id);
      const m = S.mining;
      if (!m.active) {
        m.active = true;
        m.progress = m.progress || 0;
        if (el('mine-btn')) {
          el('mine-btn').textContent = '🛑 Stop Mining';
          el('mine-btn').classList.add('mining');
        }
        if (el('mine-icon')) el('mine-icon').textContent = '⚙️';

        mInt = setInterval(() => {
          m.progress += ri(1, 6);
          const pct = Math.min(100, Math.round((m.progress / BLOCK_TARGET) * 100));
          if (el('mine-bar')) el('mine-bar').style.width = pct + '%';
          if (el('mine-pct')) el('mine-pct').textContent = pct + '%';
          if (el('hashrate-disp')) el('hashrate-disp').textContent = r(42, 58).toFixed(1) + ' MH/s';
          updateNetworkHash();

          if (m.progress >= BLOCK_TARGET) {
            m.progress = 0;
            m.blocksFound++;
            m.earned = parseFloat((m.earned + 50).toFixed(2));
            m.difficulty = Math.min(8, Math.ceil(m.blocksFound / 3) + 1);
            S.blockHeight++;
            S.wallet.balance = parseFloat((S.wallet.balance + 50).toFixed(6));
            S.wallet.received = parseFloat((S.wallet.received + 50).toFixed(6));

            const nb = {
              number: S.blockHeight,
              hash: sh(16),
              txCount: ri(5, 30),
              time: 5,
              size: ri(1000, 4000)
            };
            S.blocks.push(nb);
            trimArrays();

            if (el('h-block')) el('h-block').textContent = '#' + S.blockHeight.toLocaleString();
            addLog(`<span class="log-green">✅ BLOCK #${S.blockHeight} FOUND! +50 BTNG</span>`);
            addLog(`<span class="log-gold">   Hash: ${nb.hash}</span>`);
            showToast('🏆 Block #' + S.blockHeight + ' mined! +50 BTNG');

            const walletSec = document.getElementById('sec-wallet');
            if (walletSec?.classList.contains('active')) renderWallet();
          }

          if (el('mine-blocks')) el('mine-blocks').textContent = m.blocksFound;
          if (el('mine-earned')) el('mine-earned').textContent = m.earned + ' BTNG';
          if (el('mine-diff')) el('mine-diff').textContent = m.difficulty;
          if (el('halving-in')) el('halving-in').textContent = (21000 - S.blockHeight).toLocaleString();
        }, 200);

        lInt = setInterval(() => addLog(`<span class="log-dim">[${gt()}]</span> nonce:${ri(1000000, 9999999)} → ${sh(10)}`), 500);
        addLog(`<span class="log-gold">[${gt()}] Mining started! Difficulty: ${m.difficulty}</span>`);
      } else {
        m.active = false;
        clearInterval(mInt);
        clearInterval(lInt);
        if (el('mine-btn')) {
          el('mine-btn').textContent = '⛏️ Start Mining';
          el('mine-btn').classList.remove('mining');
        }
        if (el('mine-icon')) el('mine-icon').textContent = '⛏️';
        if (el('hashrate-disp')) el('hashrate-disp').textContent = '0 MH/s';
        if (el('mine-bar')) el('mine-bar').style.width = '0%';
        if (el('mine-pct')) el('mine-pct').textContent = '0%';
        updateNetworkHash();
        addLog(`<span class="log-red">[${gt()}] Mining stopped.</span>`);
      }
    };

    function addLog(html) {
      const box = document.getElementById('mine-log');
      if (!box) return;
      const d = document.createElement('div');
      d.className = 'log-line';
      d.innerHTML = html;
      box.appendChild(d);
      while (box.children.length > 80) box.removeChild(box.firstChild);
      box.scrollTop = box.scrollHeight;
    }

    // -------------------- SEARCH --------------------
    window.doSearch = async function() {
      const el = document.getElementById('exp-search');
      const q = el?.value.trim();
      if (!q) return;

      showToast('🔍 Searching...');

      try {
        if (/^\d+$/.test(q)) {
          const num = parseInt(q);
          const block = await getBlock(num);
          if (block) {
            showToast(`📦 Block #${num} — ${block.txCount || 'N/A'} txs, hash ${block.hash?.slice(0, 10) || 'N/A'}...`);
          } else {
            showToast(`❌ Block #${num} not found`);
          }
        }
        else if (q.startsWith('0x')) {
          const tx = await getTransaction(q);
          if (tx) {
            showToast(`🔍 TX ${q.slice(0, 14)}... — ${tx.status || 'confirmed'}, amount ${fm(tx.amount || 0, 4)} BTNG`);
          } else {
            showToast(`❌ TX not found`);
          }
        }
        else if (q.toUpperCase().startsWith('BTNG1')) {
          const addressInfo = await getAddressInfo(q);
          if (addressInfo) {
            showToast(`💼 Address ${q.slice(0, 12)}... — balance ${fm(addressInfo.balance || 0, 4)} BTNG`);
          } else {
            showToast(`❌ Address not found`);
          }
        }
        else {
          showToast(`🔍 Searching "${q}"...`);
        }
      } catch (error) {
        console.error('Search error:', error);
        showToast('❌ Search failed');
      }
    };
    
    const searchInput = document.getElementById('exp-search');
    if (searchInput) {
      searchInput.addEventListener('keydown', e => { if (e.key === 'Enter') window.doSearch(); });
    }

    // -------------------- LIVE PRICE TICK --------------------
    setInterval(async () => {
      // Try to get real price data
      const priceData = await getPriceData();
      if (priceData && priceData.price) {
        S.price = priceData.price;
        S.priceChange = priceData.change24h || S.priceChange;
      } else {
        // Fallback to mock price movement
        S.price = Math.max(0.01, parseFloat((S.price * (1 + r(-0.02, 0.025))).toFixed(4)));
      }

      S.priceHistory.push(S.price);
      if (S.priceHistory.length > 50) S.priceHistory.shift();
      S.priceChange = ((S.price - S.priceHistory[0]) / S.priceHistory[0]) * 100;

      const el = (id) => document.getElementById(id);
      if (el('h-price')) el('h-price').textContent = fu(S.price);
      
      const activeId = document.querySelector('.section.active')?.id;
      if (activeId === 'sec-dashboard') {
        renderDashboard();
        drawChart('priceChart', S.priceHistory);
      }
      if (activeId === 'sec-market') {
        renderMarket();
        drawChart('marketChart', S.priceHistory);
      }
      if (activeId === 'sec-wallet') {
        if (el('w-usd')) el('w-usd').textContent = '≈ ' + fu(S.wallet.balance * S.price) + ' USD';
      }
      updateNetworkHash();
    }, 4000);

    // -------------------- RESIZE HANDLER --------------------
    window.addEventListener('resize', () => {
      const activeId = document.querySelector('.section.active')?.id;
      if (activeId === 'sec-dashboard') drawChart('priceChart', S.priceHistory);
      if (activeId === 'sec-market') drawChart('marketChart', S.priceHistory);
    });

    // -------------------- START --------------------
    (async () => {
      await initData();
      renderDashboard();
      setTimeout(() => drawChart('priceChart', S.priceHistory), 100);
      updateNetworkHash();
    })();
  }
})();
