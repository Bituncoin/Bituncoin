import GoldCard from '@/components/GoldCard'
import styles from './page.module.css'

export default function OnboardingPage() {
  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className="container">
          <h1 className={styles.title}>Get Your BTNG Gold Card</h1>
          <p className={styles.subtitle}>
            Universal identity credentials for sovereign trust operations
          </p>
        </div>
      </section>

      <section className={styles.content}>
        <div className="container">
          <div className={styles.grid}>
            <div className={styles.info}>
              <h2>What is the BTNG Gold Card?</h2>
              <p>
                The Gold Card is your sovereign identity credential—a universal trust 
                passport that enables:
              </p>
              <ul className={styles.benefits}>
                <li>✅ Institutional-grade identity verification</li>
                <li>✅ Access to QR wallet and value transfer</li>
                <li>✅ Proof-of-value credentials</li>
                <li>✅ Trust Union protocol participation</li>
                <li>✅ Cross-border trust operations</li>
              </ul>

              <div className={styles.steps}>
                <h3>Onboarding Process</h3>
                <ol>
                  <li>Complete identity verification</li>
                  <li>Generate sovereign credentials</li>
                  <li>Receive Gold Card with QR code</li>
                  <li>Activate wallet and value profile</li>
                </ol>
              </div>
            </div>

            <div className={styles.cardPreview}>
              <GoldCard 
                holderName="Your Name"
                cardNumber="BTNG-XXXX-XXXX"
                issueDate="Launch Phase"
                isPreview={true}
              />
              <p className={styles.note}>
                Preview of your sovereign Gold Card identity
              </p>
            </div>
          </div>

          <div className={styles.cta}>
            <button className="btn-primary" style={{ fontSize: '1.1rem', padding: 'var(--space-lg) var(--space-2xl)' }}>
              Begin Onboarding Process
            </button>
            <p className={styles.subtext}>
              Secure • Sovereign • Institutional-Grade
            </p>
          </div>
        </div>
      </section>
    </div>
  )
}
