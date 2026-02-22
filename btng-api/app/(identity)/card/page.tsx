import styles from './page.module.css'
import GoldCard from '@/components/GoldCard'

export default function CardPage() {
  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className="container">
          <h1 className={styles.title}>BTNG Gold Card</h1>
          <p className={styles.subtitle}>
            Your sovereign identity credential for institutional trust operations
          </p>
        </div>
      </section>

      <section className={styles.content}>
        <div className="container">
          <div className={styles.grid}>
            <div className={styles.cardDisplay}>
              <GoldCard 
                holderName="Cardholder Name"
                cardNumber="BTNG-XXXX-XXXX"
                issueDate="2026"
                isPreview={true}
              />
              <div className={styles.cardActions}>
                <button className="btn-primary">View Digital Card</button>
                <button className="btn-secondary">Download as Image</button>
              </div>
            </div>

            <div className={styles.info}>
              <h2>Your Gold Card Identity</h2>
              <p>
                The BTNG Gold Card is your universal passport to trust operations 
                across the sovereign network. Each card contains:
              </p>

              <div className={styles.features}>
                <div className={styles.feature}>
                  <span className={styles.icon}>🔐</span>
                  <div>
                    <h4>Cryptographic Identity</h4>
                    <p>Unique identifier secured with institutional-grade encryption</p>
                  </div>
                </div>
                <div className={styles.feature}>
                  <span className={styles.icon}>📊</span>
                  <div>
                    <h4>Trust Score</h4>
                    <p>Dynamic reputation calculated from proof-of-value credentials</p>
                  </div>
                </div>
                <div className={styles.feature}>
                  <span className={styles.icon}>🌍</span>
                  <div>
                    <h4>Cross-Border Recognition</h4>
                    <p>Accepted across all Trust Union member countries</p>
                  </div>
                </div>
                <div className={styles.feature}>
                  <span className={styles.icon}>🛡️</span>
                  <div>
                    <h4>Sovereign Protection</h4>
                    <p>Privacy-preserving zero-knowledge proof verification</p>
                  </div>
                </div>
              </div>

              <div className={styles.stats}>
                <h3>Card Statistics</h3>
                <div className={styles.statGrid}>
                  <div className={styles.stat}>
                    <div className={styles.statValue}>847</div>
                    <div className={styles.statLabel}>Trust Score</div>
                  </div>
                  <div className={styles.stat}>
                    <div className={styles.statValue}>5</div>
                    <div className={styles.statLabel}>Active Countries</div>
                  </div>
                  <div className={styles.stat}>
                    <div className={styles.statValue}>142</div>
                    <div className={styles.statLabel}>Transactions</div>
                  </div>
                  <div className={styles.stat}>
                    <div className={styles.statValue}>Sovereign</div>
                    <div className={styles.statLabel}>Verification Level</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  )
}
