import QRWallet from '@/components/QRWallet'
import styles from './page.module.css'

export default function WalletPage() {
  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className="container">
          <h1 className={styles.title}>BTNG QR Wallet</h1>
          <p className={styles.subtitle}>
            Sovereign value transfer with zero-knowledge proof
          </p>
        </div>
      </section>

      <section className={styles.content}>
        <div className="container">
          <div className={styles.walletSection}>
            <div className={styles.info}>
              <h2>Your Sovereign Wallet</h2>
              <p>
                The QR Wallet enables instant, secure value transfer using your 
                Gold Card identity credentials.
              </p>

              <div className={styles.features}>
                <div className={styles.feature}>
                  <span className={styles.icon}>🔐</span>
                  <div>
                    <h4>Zero-Knowledge Proof</h4>
                    <p>Transactions validated without exposing identity</p>
                  </div>
                </div>
                <div className={styles.feature}>
                  <span className={styles.icon}>⚡</span>
                  <div>
                    <h4>Instant Transfer</h4>
                    <p>QR-based value transfer in seconds</p>
                  </div>
                </div>
                <div className={styles.feature}>
                  <span className={styles.icon}>🌍</span>
                  <div>
                    <h4>Global Reach</h4>
                    <p>Cross-border transactions via Trust Union</p>
                  </div>
                </div>
                <div className={styles.feature}>
                  <span className={styles.icon}>🛡️</span>
                  <div>
                    <h4>Sovereign Security</h4>
                    <p>Institutional-grade cryptographic protection</p>
                  </div>
                </div>
              </div>
            </div>

            <div className={styles.walletDisplay}>
              <QRWallet 
                walletId="BTNG-WALLET-DEMO"
                balance="0.00"
                currency="BTNG"
              />
            </div>
          </div>

          <div className={styles.actions}>
            <button className="btn-primary">Generate Payment QR</button>
            <button className="btn-secondary">Scan to Send</button>
            <button className="btn-secondary">Transaction History</button>
          </div>
        </div>
      </section>
    </div>
  )
}
