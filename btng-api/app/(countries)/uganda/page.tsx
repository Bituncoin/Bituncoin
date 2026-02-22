import { countryData } from '@/lib/countries'
import styles from '../country.module.css'
import Link from 'next/link'

export default function UgandaPage() {
  const country = countryData.uganda

  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className="container">
          <div className={styles.heroContent}>
            <span className={styles.flag}>{country.flag}</span>
            <h1 className={styles.title}>BTNG in {country.name}</h1>
            <p className={styles.subtitle}>{country.description}</p>
            <div className={styles.status}>
              <span className={`${styles.statusBadge} ${styles[country.status]}`}>
                {country.status === 'active' ? '✅ Fully Operational' : 
                 country.status === 'launching' ? '🚀 Launching Soon' : 
                 '📋 Planned Deployment'}
              </span>
            </div>
          </div>
        </div>
      </section>

      <section className={styles.content}>
        <div className="container">
          <div className={styles.statsGrid}>
            <div className={styles.statCard}>
              <div className={styles.statValue}>{country.activeUsers.toLocaleString()}</div>
              <div className={styles.statLabel}>Active Gold Card Holders</div>
            </div>
            <div className={styles.statCard}>
              <div className={styles.statValue}>{country.trustNodes}</div>
              <div className={styles.statLabel}>Trust Union Nodes</div>
            </div>
            <div className={styles.statCard}>
              <div className={styles.statValue}>{country.mobileMoney.length}</div>
              <div className={styles.statLabel}>Mobile Money Providers</div>
            </div>
            <div className={styles.statCard}>
              <div className={styles.statValue}>{country.population}</div>
              <div className={styles.statLabel}>Population</div>
            </div>
          </div>

          <div className={styles.mainGrid}>
            <div className={styles.mainContent}>
              <div className={styles.section}>
                <h2>Country Overview</h2>
                <div className={styles.overviewGrid}>
                  <div className={styles.overviewItem}>
                    <div className={styles.label}>Capital</div>
                    <div className={styles.value}>{country.capital}</div>
                  </div>
                  <div className={styles.overviewItem}>
                    <div className={styles.label}>Currency</div>
                    <div className={styles.value}>{country.currency} ({country.currencyCode})</div>
                  </div>
                  <div className={styles.overviewItem}>
                    <div className={styles.label}>Population</div>
                    <div className={styles.value}>{country.population}</div>
                  </div>
                  <div className={styles.overviewItem}>
                    <div className={styles.label}>BTNG Status</div>
                    <div className={styles.value}>
                      <span className={country.status === 'active' ? styles.activeText : styles.launchingText}>
                        {country.status.charAt(0).toUpperCase() + country.status.slice(1)}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <div className={styles.section}>
                <h2>Highlights & Achievements</h2>
                <div className={styles.highlightsList}>
                  {country.highlights.map((highlight, index) => (
                    <div key={index} className={styles.highlightItem}>
                      <span className={styles.checkmark}>✓</span>
                      <span>{highlight}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className={styles.section}>
                <h2>Mobile Money Integration</h2>
                <div className={styles.providerList}>
                  {country.mobileMoney.map((provider, index) => (
                    <div key={index} className={styles.providerCard}>
                      <span className={styles.providerIcon}>📱</span>
                      <div>
                        <div className={styles.providerName}>{provider}</div>
                        <div className={styles.providerStatus}>
                          {country.status === 'active' ? 'Integrated' : 'Integration Pending'}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className={styles.sidebar}>
              <div className={styles.sidebarCard}>
                <h3>Get Started in {country.name}</h3>
                <p>Join thousands of Ugandan citizens using BTNG Gold Card.</p>
                <Link href="/(onboarding)/user" className="btn-primary" style={{width: '100%', display: 'block', textAlign: 'center'}}>
                  Apply for Gold Card
                </Link>
              </div>

              <div className={styles.sidebarCard}>
                <h3>For Merchants</h3>
                <p>Accept BTNG payments and join the Trust Union network.</p>
                <Link href="/(onboarding)/merchant" className="btn-secondary" style={{width: '100%', display: 'block', textAlign: 'center'}}>
                  Merchant Registration
                </Link>
              </div>

              <div className={styles.sidebarCard}>
                <h3>Need Help?</h3>
                <div className={styles.contactList}>
                  <div className={styles.contactItem}>
                    <span>📧</span>
                    <span>support-ug@btng.global</span>
                  </div>
                  <div className={styles.contactItem}>
                    <span>📞</span>
                    <span>+256 XX XXX XXXX</span>
                  </div>
                  <div className={styles.contactItem}>
                    <span>🏢</span>
                    <span>{country.capital} Office</span>
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
