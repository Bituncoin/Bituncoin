import Link from 'next/link'
import styles from './SovereignFooter.module.css'

export default function SovereignFooter() {
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
            <h4 className={styles.heading}>Identity</h4>
            <nav className={styles.linkList}>
              <Link href="/(identity)/wallet">QR Wallet</Link>
              <Link href="/(onboarding)/user">Get Gold Card</Link>
              <Link href="/(identity)/profile">My Profile</Link>
              <Link href="/(identity)/card">Gold Card</Link>
            </nav>
          </div>

          <div className={styles.section}>
            <h4 className={styles.heading}>Countries</h4>
            <nav className={styles.linkList}>
              <Link href="/(countries)/ghana">Ghana</Link>
              <Link href="/(countries)/kenya">Kenya</Link>
              <Link href="/(countries)/nigeria">Nigeria</Link>
              <Link href="/(countries)/togo">Togo</Link>
            </nav>
          </div>

          <div className={styles.section}>
            <h4 className={styles.heading}>Platform</h4>
            <nav className={styles.linkList}>
              <Link href="/trust-union">Trust Union</Link>
              <Link href="/(onboarding)/merchant">Merchant Integration</Link>
              <Link href="/mobile-money">Mobile Money</Link>
              <Link href="/api/health">Health Status</Link>
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
