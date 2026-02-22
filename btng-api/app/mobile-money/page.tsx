import styles from './page.module.css'

export default function MobileMoneyPage() {
  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className="container">
          <h1 className={styles.title}>Mobile Money Integration</h1>
          <p className={styles.subtitle}>
            Connecting BTNG wallet with country-specific mobile money providers
          </p>
        </div>
      </section>

      <section className={styles.content}>
        <div className="container">
          <div className={styles.overview}>
            <h2>Universal Value Transfer</h2>
            <p>
              BTNG integrates with mobile money providers across Africa and beyond, 
              enabling seamless value transfer between the QR wallet and local 
              mobile money accounts.
            </p>
          </div>

          <div className={styles.providers}>
            <h2>Supported Providers</h2>
            <div className={styles.providerGrid}>
              <div className={styles.provider}>
                <div className={styles.providerHeader}>
                  <span className={styles.flag}>🇰🇪</span>
                  <h3>M-Pesa (Kenya)</h3>
                </div>
                <div className={styles.providerInfo}>
                  <p><strong>Status:</strong> <span className={styles.active}>Active</span></p>
                  <p><strong>Currency:</strong> KES</p>
                  <p><strong>Daily Limit:</strong> 150,000 KES</p>
                </div>
              </div>

              <div className={styles.provider}>
                <div className={styles.providerHeader}>
                  <span className={styles.flag}>🇺🇬</span>
                  <h3>MTN Mobile Money (Uganda)</h3>
                </div>
                <div className={styles.providerInfo}>
                  <p><strong>Status:</strong> <span className={styles.active}>Active</span></p>
                  <p><strong>Currency:</strong> UGX</p>
                  <p><strong>Daily Limit:</strong> 5,000,000 UGX</p>
                </div>
              </div>

              <div className={styles.provider}>
                <div className={styles.providerHeader}>
                  <span className={styles.flag}>🇬🇭</span>
                  <h3>Vodafone Cash (Ghana)</h3>
                </div>
                <div className={styles.providerInfo}>
                  <p><strong>Status:</strong> <span className={styles.planned}>Planned</span></p>
                  <p><strong>Currency:</strong> GHS</p>
                  <p><strong>Daily Limit:</strong> 10,000 GHS</p>
                </div>
              </div>

              <div className={styles.provider}>
                <div className={styles.providerHeader}>
                  <span className={styles.flag}>🇹🇿</span>
                  <h3>M-Pesa (Tanzania)</h3>
                </div>
                <div className={styles.providerInfo}>
                  <p><strong>Status:</strong> <span className={styles.planned}>Planned</span></p>
                  <p><strong>Currency:</strong> TZS</p>
                  <p><strong>Daily Limit:</strong> 3,000,000 TZS</p>
                </div>
              </div>
            </div>
          </div>

          <div className={styles.howItWorks}>
            <h2>How It Works</h2>
            <div className={styles.workflow}>
              <div className={styles.workflowStep}>
                <div className={styles.stepIcon}>📱</div>
                <h4>Link Your Account</h4>
                <p>Connect your mobile money number to your BTNG Gold Card</p>
              </div>
              <div className={styles.arrow}>→</div>
              <div className={styles.workflowStep}>
                <div className={styles.stepIcon}>💰</div>
                <h4>Transfer Value</h4>
                <p>Send money from mobile money to QR wallet or vice versa</p>
              </div>
              <div className={styles.arrow}>→</div>
              <div className={styles.workflowStep}>
                <div className={styles.stepIcon}>✅</div>
                <h4>Instant Confirmation</h4>
                <p>Receive instant confirmation via Trust Union Protocol</p>
              </div>
            </div>
          </div>

          <div className={styles.features}>
            <h2>Integration Features</h2>
            <div className={styles.featureList}>
              <div className={styles.feature}>
                <span className={styles.check}>✓</span>
                <div>
                  <h4>Bidirectional Transfer</h4>
                  <p>Send and receive between BTNG wallet and mobile money</p>
                </div>
              </div>
              <div className={styles.feature}>
                <span className={styles.check}>✓</span>
                <div>
                  <h4>Real-Time Exchange Rates</h4>
                  <p>Automatic currency conversion at market rates</p>
                </div>
              </div>
              <div className={styles.feature}>
                <span className={styles.check}>✓</span>
                <div>
                  <h4>Low Fees</h4>
                  <p>Competitive transaction fees via Trust Union</p>
                </div>
              </div>
              <div className={styles.feature}>
                <span className={styles.check}>✓</span>
                <div>
                  <h4>Provider Agnostic</h4>
                  <p>Works with any connected mobile money provider</p>
                </div>
              </div>
            </div>
          </div>

          <div className={styles.callout}>
            <h3>Ready to connect your mobile money?</h3>
            <p>Link your mobile money account to your BTNG QR Wallet</p>
            <button className="btn-primary">Link Mobile Money Account</button>
          </div>
        </div>
      </section>
    </div>
  )
}
