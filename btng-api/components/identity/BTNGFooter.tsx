import Link from 'next/link'
import styles from './BTNGFooter.module.css'

export default function BTNGFooter() {
  return (
    <footer className={styles.footer}>
      <div className={styles.container}>
        <div className={styles.grid}>
          <div className={styles.section}>
            <h3 className={styles.sectionTitle}>BTNG Platform</h3>
            <p className={styles.description}>
              Building Trust. Nurturing Growth.
              <br />
              Sovereign identity for institutional-grade operations.
            </p>
            <div className={styles.badge}>
              <span className="trust-badge">🛡️ Sovereign Architecture</span>
            </div>
          </div>

          <div className={styles.section}>
            <h4 className={styles.heading}>Platform</h4>
            <nav className={styles.linkList}>
              <Link href="/wallet">QR Wallet</Link>
              <Link href="/onboarding">Get Gold Card</Link>
              <Link href="/trust-union">Trust Union</Link>
              <Link href="/api/health">Health Status</Link>
            </nav>
          </div>

          <div className={styles.section}>
            <h4 className={styles.heading}>Expansion</h4>
            <nav className={styles.linkList}>
              <Link href="/onboarding/country">Country Onboarding</Link>
              <Link href="/onboarding/merchant">Merchant Integration</Link>
              <Link href="/mobile-money">Mobile Money</Link>
              <Link href="/proof-of-value">Proof of Value</Link>
            </nav>
          </div>

          <div className={styles.section}>
            <h4 className={styles.heading}>Trust Operations</h4>
            <nav className={styles.linkList}>
              <Link href="/identity">Universal Identity</Link>
              <Link href="/value-profile">Value Profile</Link>
              <Link href="/debt-release">Debt Release</Link>
              <Link href="/trust-union/protocol">Protocol Docs</Link>
            </nav>
          </div>
        </div>

        <div className={styles.bottom}>
          <p className={styles.copyright}>
            © {new Date().getFullYear()} BTNG Platform. Sovereign Architecture.
          </p>
          <p className={styles.version}>
            v0.1.0 | Launch Phase
          </p>
        </div>
      </div>
    </footer>
  )
}
