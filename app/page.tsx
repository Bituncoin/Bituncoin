import Link from 'next/link'
import styles from './page.module.css'
const countryData = {
  nigeria: {
    name: 'Nigeria',
    flag: '🇳🇬',
    status: 'active',
    activeUsers: 820000,
    population: '223M',
  },
  ghana: {
    name: 'Ghana',
    flag: '🇬🇭',
    status: 'active',
    activeUsers: 310000,
    population: '34M',
  },
  kenya: {
    name: 'Kenya',
    flag: '🇰🇪',
    status: 'active',
    activeUsers: 275000,
    population: '55M',
  },
  rwanda: {
    name: 'Rwanda',
    flag: '🇷🇼',
    status: 'active',
    activeUsers: 98000,
    population: '14M',
  },
  uganda: {
    name: 'Uganda',
    flag: '🇺🇬',
    status: 'launching',
    activeUsers: 0,
    population: '48M',
  },
  tanzania: {
    name: 'Tanzania',
    flag: '🇹🇿',
    status: 'launching',
    activeUsers: 0,
    population: '67M',
  },
  ethiopia: {
    name: 'Ethiopia',
    flag: '🇪🇹',
    status: 'launching',
    activeUsers: 0,
    population: '126M',
  },
  southAfrica: {
    name: 'South Africa',
    flag: '🇿🇦',
    status: 'launching',
    activeUsers: 0,
    population: '60M',
  },
} as const

export default function Home() {
  const entries = Object.entries(countryData) as [string, (typeof countryData)[keyof typeof countryData]][]
  const activeCountries = entries.filter(([, country]) => country.status === 'active')
  const launchingCountries = entries.filter(([, country]) => country.status === 'launching')

  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <div className={styles.heroEmblem}>
            <span className="emblem-hero" aria-hidden="true" style={{ fontSize: 120 }}>
              🛡️
            </span>
          </div>
          <h1 className={styles.title}>
            Building Trust<span className={styles.gold}>.</span>
            <br />
            Nurturing Growth<span className={styles.gold}>.</span>
          </h1>
          <p className={styles.subtitle}>
            A sovereign identity and value platform for institutional-grade trust operations
          </p>
          <div className={styles.ctaGroup}>
            <Link href="/gold-coin" className="btn-primary">
              BTNG 54 Africa Gold Coin
            </Link>
            <Link href="/(onboarding)/user" className="btn-secondary">
              Get Your Gold Card
            </Link>
            <Link href="/(identity)/wallet" className="btn-secondary">
              Access QR Wallet
            </Link>
          </div>
        </div>
      </section>

      <section className={styles.goldCoinSection}>
        <div className={styles.heroContent}>
          <h2 className={styles.sectionTitle}>BTNG 54 Africa Gold Coin</h2>
          <p className={styles.subtitle}>
            Institutional-grade gold-backed value for trusted, cross-border operations.
          </p>
          <Link href="/gold-coin" className="btn-primary">
            Learn More
          </Link>
        </div>
      </section>

      <section className={styles.features}>
        <div className="container">
          <h2 className={styles.sectionTitle}>Sovereign Architecture</h2>
          <div className={styles.featureGrid}>
            <Link href="/(identity)/card" className={styles.feature}>
              <div className={styles.featureIcon}>🛡️</div>
              <h3>Universal Identity</h3>
              <p>Gold Card credentials for sovereign trust operations</p>
            </Link>
            <Link href="/(identity)/wallet" className={styles.feature}>
              <div className={styles.featureIcon}>📱</div>
              <h3>QR Wallet System</h3>
              <p>Instant value transfer with zero-knowledge proof</p>
            </Link>
            <Link href="/trust-union" className={styles.feature}>
              <div className={styles.featureIcon}>🤝</div>
              <h3>Trust Union Protocol</h3>
              <p>Institutional-grade trust architecture</p>
            </Link>
            <Link href="/(onboarding)/merchant" className={styles.feature}>
              <div className={styles.featureIcon}>🏪</div>
              <h3>Merchant Network</h3>
              <p>Accept BTNG payments across 8 countries</p>
            </Link>
          </div>
        </div>
      </section>

      <section className={styles.countries}>
        <div className="container">
          <h2 className={styles.sectionTitle}>Active Countries</h2>
          <div className={styles.countryGrid}>
            {activeCountries.map(([key, country]) => (
              <Link key={key} href={`/(countries)/${key}`} className={styles.countryCard}>
                <span className={styles.countryFlag}>{country.flag}</span>
                <h3>{country.name}</h3>
                <p>{country.activeUsers.toLocaleString()} Gold Card holders</p>
                <span className={styles.activeStatus}>✅ Fully Operational</span>
              </Link>
            ))}
          </div>
          
          {launchingCountries.length > 0 && (
            <>
              <h3 className={styles.subsectionTitle}>Launching Soon</h3>
              <div className={styles.countryGrid}>
                {launchingCountries.map(([key, country]) => (
                  <Link key={key} href={`/(countries)/${key}`} className={styles.countryCard}>
                    <span className={styles.countryFlag}>{country.flag}</span>
                    <h3>{country.name}</h3>
                    <p>{country.population} population</p>
                    <span className={styles.launchingStatus}>🚀 Launching Soon</span>
                  </Link>
                ))}
              </div>
            </>
          )}
        </div>
      </section>

      <section className={styles.callout}>
        <div className="gold-card">
          <h2>Ready to join the trust revolution?</h2>
          <p style={{ marginTop: 'var(--space-md)', marginBottom: 'var(--space-lg)' }}>
            BTNG combines sovereign identity, institutional trust, and proof-of-value workflows 
            into a unified platform spanning 8 African nations.
          </p>
          <div style={{ display: 'flex', gap: 'var(--space-md)', justifyContent: 'center', flexWrap: 'wrap' }}>
            <Link href="/(onboarding)/user" className="btn-primary">
              Get Gold Card
            </Link>
            <Link href="/(onboarding)/merchant" className="btn-secondary">
              Become a Merchant
            </Link>
          </div>
        </div>
      </section>
    </div>
  )
}
