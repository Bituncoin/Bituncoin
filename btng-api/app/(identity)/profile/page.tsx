import styles from './page.module.css'

export default function ProfilePage() {
  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className="container">
          <h1 className={styles.title}>Your Identity Profile</h1>
          <p className={styles.subtitle}>
            Universal identity and value credentials
          </p>
        </div>
      </section>

      <section className={styles.content}>
        <div className="container">
          <div className={styles.profileGrid}>
            <div className={styles.sidebar}>
              <div className={styles.avatarSection}>
                <div className={styles.avatar}>
                  <span className={styles.avatarIcon}>👤</span>
                </div>
                <h2 className={styles.name}>Your Name</h2>
                <p className={styles.cardNumber}>BTNG-XXXX-XXXX</p>
                <div className={styles.trustBadge}>
                  <span className="trust-badge">Trust Score: 847</span>
                </div>
              </div>

              <div className={styles.quickActions}>
                <h3>Quick Actions</h3>
                <button className="btn-primary" style={{width: '100%', marginBottom: 'var(--space-sm)'}}>
                  View Gold Card
                </button>
                <button className="btn-secondary" style={{width: '100%', marginBottom: 'var(--space-sm)'}}>
                  Open Wallet
                </button>
                <button className="btn-secondary" style={{width: '100%'}}>
                  Add Proof of Value
                </button>
              </div>
            </div>

            <div className={styles.mainContent}>
              <div className={styles.section}>
                <h2>Identity Information</h2>
                <div className={styles.infoGrid}>
                  <div className={styles.infoItem}>
                    <div className={styles.label}>Full Name</div>
                    <div className={styles.value}>Your Name</div>
                  </div>
                  <div className={styles.infoItem}>
                    <div className={styles.label}>Card Number</div>
                    <div className={styles.value}>BTNG-XXXX-XXXX</div>
                  </div>
                  <div className={styles.infoItem}>
                    <div className={styles.label}>Issue Date</div>
                    <div className={styles.value}>January 2026</div>
                  </div>
                  <div className={styles.infoItem}>
                    <div className={styles.label}>Verification Level</div>
                    <div className={styles.value}>
                      <span className={styles.sovereign}>Sovereign</span>
                    </div>
                  </div>
                </div>
              </div>

              <div className={styles.section}>
                <h2>Trust & Value Metrics</h2>
                <div className={styles.metricsGrid}>
                  <div className={styles.metric}>
                    <div className={styles.metricValue}>847</div>
                    <div className={styles.metricLabel}>Trust Score</div>
                    <div className={styles.metricChange}>+23 this month</div>
                  </div>
                  <div className={styles.metric}>
                    <div className={styles.metricValue}>12,450</div>
                    <div className={styles.metricLabel}>Total Value (BTNG)</div>
                    <div className={styles.metricChange}>+1,200 this month</div>
                  </div>
                  <div className={styles.metric}>
                    <div className={styles.metricValue}>5</div>
                    <div className={styles.metricLabel}>Active Countries</div>
                    <div className={styles.metricChange}>Expanded recently</div>
                  </div>
                  <div className={styles.metric}>
                    <div className={styles.metricValue}>142</div>
                    <div className={styles.metricLabel}>Total Transactions</div>
                    <div className={styles.metricChange}>18 this week</div>
                  </div>
                </div>
              </div>

              <div className={styles.section}>
                <h2>Active Countries</h2>
                <div className={styles.countryList}>
                  {['Kenya 🇰🇪', 'Uganda 🇺🇬', 'Ghana 🇬🇭', 'Nigeria 🇳🇬', 'South Africa 🇿🇦'].map(country => (
                    <div key={country} className={styles.countryItem}>
                      <span>{country}</span>
                      <span className={styles.activeStatus}>Active</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className={styles.section}>
                <h2>Recent Proof of Value</h2>
                <div className={styles.povList}>
                  <div className={styles.povItem}>
                    <div className={styles.povIcon}>💼</div>
                    <div className={styles.povDetails}>
                      <div className={styles.povType}>Work Contribution</div>
                      <div className={styles.povAmount}>+500 BTNG</div>
                      <div className={styles.povDate}>2 days ago</div>
                    </div>
                  </div>
                  <div className={styles.povItem}>
                    <div className={styles.povIcon}>🤝</div>
                    <div className={styles.povDetails}>
                      <div className={styles.povType}>Trust Verification</div>
                      <div className={styles.povAmount}>+15 Trust Points</div>
                      <div className={styles.povDate}>5 days ago</div>
                    </div>
                  </div>
                  <div className={styles.povItem}>
                    <div className={styles.povIcon}>💱</div>
                    <div className={styles.povDetails}>
                      <div className={styles.povType}>Trade Activity</div>
                      <div className={styles.povAmount}>+300 BTNG</div>
                      <div className={styles.povDate}>1 week ago</div>
                    </div>
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
