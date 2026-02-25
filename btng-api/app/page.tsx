import Link from 'next/link'
import SovereignEmblem from '@/components/SovereignEmblem'
import BTNG54AfricaGoldCoin from '@/components/BTNG54AfricaGoldCoin'
import styles from './page.module.css'
import { countryData } from '@/lib/countries'

export default function Home() {
  const activeCountries = Object.entries(countryData).filter(
    ([_, country]) => country.status === 'active'
  )
  const launchingCountries = Object.entries(countryData).filter(
    ([_, country]) => country.status === 'launching'
  )

  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <div className={styles.heroEmblem}>
            <SovereignEmblem size={120} className="emblem-hero" />
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
            <Link href="/watchtower" className="btn-secondary">
              Open Sovereign Watchtower
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
        <BTNG54AfricaGoldCoin />
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
            <Link href="/watchtower" className={styles.feature}>
              <div className={styles.featureIcon}>🛰️</div>
              <h3>Sovereign Watchtower</h3>
              <p>Monitor consensus depth and node heartbeat in real-time</p>
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
