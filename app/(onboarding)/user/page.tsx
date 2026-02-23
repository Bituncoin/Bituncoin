import GoldCard from '../../../components/GoldCard'
import styles from './page.module.css'

export default function UserOnboardingPage() {
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
                  <li>Select your country and provide basic information</li>
                  <li>Complete identity verification</li>
                  <li>Generate sovereign credentials</li>
                  <li>Receive Gold Card with QR code</li>
                  <li>Activate wallet and value profile</li>
                </ol>
              </div>

              <div className={styles.form}>
                <h3>Start Your Application</h3>
                <form className={styles.onboardingForm}>
                  <div className={styles.formGroup}>
                    <label>Full Name</label>
                    <input type="text" placeholder="Enter your full name" />
                  </div>
                  <div className={styles.formGroup}>
                    <label>Country</label>
                    <select>
                      <option value="">Select your country</option>
                      <option value="ghana">Ghana 🇬🇭</option>
                      <option value="kenya">Kenya 🇰🇪</option>
                      <option value="nigeria">Nigeria 🇳🇬</option>
                      <option value="togo">Togo 🇹🇬</option>
                      <option value="uganda">Uganda 🇺🇬</option>
                      <option value="ivory-coast">Ivory Coast 🇨🇮</option>
                      <option value="burkina-faso">Burkina Faso 🇧🇫</option>
                      <option value="south-africa">South Africa 🇿🇦</option>
                    </select>
                  </div>
                  <div className={styles.formGroup}>
                    <label>Mobile Number</label>
                    <input type="tel" placeholder="+XXX XXX XXXX" />
                  </div>
                  <div className={styles.formGroup}>
                    <label>Email Address</label>
                    <input type="email" placeholder="your@email.com" />
                  </div>
                  <div className={styles.formGroup}>
                    <label>
                      <input type="checkbox" />
                      <span>I agree to the BTNG Trust Union terms and privacy policy</span>
                    </label>
                  </div>
                  <button type="submit" className="btn-primary" style={{width: '100%'}}>
                    Begin Onboarding Process
                  </button>
                </form>
              </div>
            </div>

            <div className={styles.cardPreview}>
              <GoldCard 
                holderName="Your Name"
                cardNumber="BTNG-XXXX-XXXX"
                issueDate="2026"
                isPreview={true}
              />
              <p className={styles.note}>
                Preview of your sovereign Gold Card identity
              </p>
              
              <div className={styles.trustFeatures}>
                <h4>Why Trust BTNG?</h4>
                <div className={styles.trustItem}>
                  <span className={styles.trustIcon}>🛡️</span>
                  <div>
                    <strong>Sovereign Architecture</strong>
                    <p>Your identity, your control</p>
                  </div>
                </div>
                <div className={styles.trustItem}>
                  <span className={styles.trustIcon}>🔐</span>
                  <div>
                    <strong>Zero-Knowledge Security</strong>
                    <p>Privacy-preserving verification</p>
                  </div>
                </div>
                <div className={styles.trustItem}>
                  <span className={styles.trustIcon}>🌍</span>
                  <div>
                    <strong>Global Recognition</strong>
                    <p>Accepted across 8 countries</p>
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
