# BTNG Demo Sandbox Integration

## 📋 Overview

The BTNG Demo Sandbox is now fully integrated into your sovereign Next.js 14 platform at **`/btng-demo`**.

## 🏗️ Architecture

### File Structure
```
app/
  (sandbox)/
    btng-demo/
      page.tsx              # Main demo page (server component)
      BTNGDemoClient.tsx    # Client-side event handler wrapper
      btng-demo.css         # Demo-specific styles

public/
  btng-demo/
    btng-demo.js            # Complete demo logic (wallet, mining, explorer, market)

components/
  SovereignHeader.tsx       # Updated with 🎮 Demo link
```

### Integration Points

1. **Sovereign Frame** — The demo renders inside your global identity layout
2. **Navigation** — Added "🎮 Demo" link in the main header between "Trust Union" and "Get Gold Card"
3. **Emblem & Chime** — Verification events and sovereign emblem remain active throughout the demo
4. **Isolated State** — Demo state is self-contained and doesn't interfere with production modules

## ✅ Features Implemented

### Core Functionality
- ✅ **Wallet** — Send/receive BTNG, QR address display, transaction history
- ✅ **Mining** — Simulated proof-of-work with difficulty adjustment
- ✅ **Explorer** — Search blocks, transactions, and addresses
- ✅ **Market** — Live price charts, market cap, volume tracking
- ✅ **QR Codes** — Wallet address QR + APK download QR

### Improvements
- ✅ Real transaction direction logic (sent/received based on address)
- ✅ Live wallet updates during mining
- ✅ Dynamic ATH/ATL recalculation
- ✅ Clipboard fallback for older browsers
- ✅ Resize-aware charts
- ✅ Memory-safe arrays (blocks capped at 50, transactions at 100)
- ✅ Fee included in "Total Sent"
- ✅ Address validation (BTNG1 + 32+ chars)
- ✅ Real search functionality
- ✅ Mining progress persistence
- ✅ Dynamic network hash rate

## 🚀 Usage

### Accessing the Demo
Navigate to: **`http://localhost:3000/btng-demo`** (or your deployment URL)

### Navigation
- **Dashboard** — Overview of price, market cap, blocks, and transactions
- **Wallet** — Manage BTNG, send transactions, view QR code
- **Mining** — Start/stop mining simulation, view logs
- **Explorer** — Search and explore blocks and transactions
- **Market** — Price charts, market stats, download APK

## 🔧 Customization

### APK Download Link
Update the APK download URL in two places:

1. **JavaScript:** `public/btng-demo/btng-demo.js`
   ```js
   const apkUrl = 'https://example.com/btng-wallet.apk'; // Line ~133
   ```

2. **HTML:** `app/(sandbox)/btng-demo/page.tsx`
   ```tsx
   <a href="https://example.com/btng-wallet.apk" ... >
   ```

### Styling
All demo-specific styles are in `app/(sandbox)/btng-demo/btng-demo.css`

CSS variables:
```css
--gold: #F5C518;
--gold-light: #FFD700;
--gold-dark: #B8960C;
--bg: #0A0A0A;
--green: #00D68F;
--red: #FF4757;
```

## 🎯 Next Steps

### Evolution Path
This demo can evolve into production modules:

1. **Wallet Logic** → Migrate to `app/(identity)/wallet`
2. **Mining Simulation** → Educational tool for miners
3. **Explorer** → Full blockchain explorer
4. **Market Data** → Connect to real price feeds
5. **QR Identity** → Integrate with actual sovereign identity system

### Connecting to Real Data
To connect the demo to real blockchain data:

1. Replace mock functions in `btng-demo.js` with API calls
2. Connect to `app/api/wallet/*` endpoints
3. Integrate with `app/api/pov/route.ts` for proof-of-value
4. Link to actual identity system in `app/(identity)/`

## 🧪 Testing

### Local Development
```bash
npm run dev
# Navigate to http://localhost:3000/btng-demo
```

### Production Build
```bash
npm run build
npm start
```

## 📝 Notes

- The demo uses **external JavaScript** loaded via Next.js `<Script>` component for optimal performance
- **QR Code generation** uses the `qrcode` library from CDN
- All interactive elements are **progressively enhanced** — the page loads without JavaScript, then enhances
- The demo is **mobile-responsive** and works on all modern browsers

## 🔗 Related Modules

- `app/(identity)/wallet/page.tsx` — Production wallet
- `app/(identity)/profile/page.tsx` — User profiles
- `app/api/wallet/transaction/route.ts` — Transaction API
- `app/api/pov/route.ts` — Proof-of-value verification

---

**Status:** ✅ Fully integrated and ready for demonstration

**Access:** Public sandbox visible in main navigation

**Performance:** Optimized with lazy-loaded scripts and CSS isolation
