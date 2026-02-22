'use client'

import { useState } from 'react'
import QRCode from 'qrcode.react'
import styles from './QRWallet.module.css'

interface QRWalletProps {
  walletId: string
  balance: string
  currency: string
}

export default function QRWallet({ walletId, balance, currency }: QRWalletProps) {
  const [showQR, setShowQR] = useState(true)

  return (
    <div className={styles.wallet}>
      <div className={styles.header}>
        <h3 className={styles.title}>Your Wallet</h3>
        <span className={styles.id}>{walletId}</span>
      </div>

      <div className={styles.balanceSection}>
        <div className={styles.label}>Available Balance</div>
        <div className={styles.balance}>
          {balance} <span className={styles.currency}>{currency}</span>
        </div>
      </div>

      {showQR && (
        <div className={styles.qrSection}>
          <div className={styles.qrContainer}>
            <QRCode 
              value={`btng:wallet:${walletId}`}
              size={200}
              level="H"
              includeMargin={true}
            />
          </div>
          <p className={styles.qrLabel}>Scan to send to this wallet</p>
        </div>
      )}

      <div className={styles.actions}>
        <button 
          className={styles.toggleButton}
          onClick={() => setShowQR(!showQR)}
        >
          {showQR ? 'Hide QR Code' : 'Show QR Code'}
        </button>
      </div>

      <div className={styles.status}>
        <span className={styles.indicator}>●</span> Active & Secure
      </div>
    </div>
  )
}
