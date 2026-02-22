import styles from './page.module.css'

const countries = [
  { name: 'Kenya', flag: '🇰🇪', status: 'Active', mobileMoney: 'M-Pesa' },
  { name: 'Uganda', flag: '🇺🇬', status: 'Active', mobileMoney: 'MTN Mobile Money' },
  { name: 'Ghana', flag: '🇬🇭', status: 'Planned', mobileMoney: 'Vodafone Cash' },
  { name: 'Nigeria', flag: '🇳🇬', status: 'Planned', mobileMoney: 'Multiple providers' },
  { name: 'Tanzania', flag: '🇹🇿', status: 'Planned', mobileMoney: 'M-Pesa, Tigo Pesa' },
  { name: 'Rwanda', flag: '🇷🇼', status: 'Planned', mobileMoney: 'MTN Mobile Money' }
]

export default function CountryOnboardingPage() {
  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className="container">
          <h1 className={styles.title}>Country Onboarding</h1>
          <p className={styles.subtitle}>
            Expanding the BTNG Trust Union across sovereign nations
          </p>
        </div>
      </section>

      <section className={styles.content}>
        <div className="container">
          <div className={styles.intro}>
            <h2>Join the Trust Union</h2>
            <p>
              Countries joining BTNG gain access to:
            </p>
            <ul className={styles.benefits}>
              <li>✅ Sovereign identity infrastructure</li>
              <li>✅ Trust Union node participation</li>
              <li>✅ Mobile money integration</li>
              <li>✅ Cross-border value transfer</li>
              <li>✅ Debt-release protocol access</li>
              <li>✅ Institutional trust credentials</li>
            </ul>
          </div>

          <div className={styles.countries}>
            <h2>Country Status</h2>
            <div className={styles.countryGrid}>
              {countries.map((country) => (
                <div key={country.name} className={styles.countryCard}>
                  <div className={styles.countryHeader}>
                    <span className={styles.flag}>{country.flag}</span>
                    <h3>{country.name}</h3>
                  </div>
                  <div className={styles.countryInfo}>
                    <div className={styles.infoRow}>
                      <span className={styles.label}>Status:</span>
                      <span className={`${styles.status} ${styles[country.status.toLowerCase()]}`}>
                        {country.status}
                      </span>
                    </div>
                    <div className={styles.infoRow}>
                      <span className={styles.label}>Mobile Money:</span>
                      <span>{country.mobileMoney}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className={styles.process}>
            <h2>Onboarding Process</h2>
            <div className={styles.steps}>
              <div className={styles.step}>
                <div className={styles.stepNumber}>1</div>
                <h4>Government Agreement</h4>
                <p>Establish sovereign partnership with BTNG Trust Union</p>
              </div>
              <div className={styles.step}>
                <div className={styles.stepNumber}>2</div>
                <h4>Node Deployment</h4>
                <p>Deploy Trust Union node infrastructure</p>
              </div>
              <div className={styles.step}>
                <div className={styles.stepNumber}>3</div>
                <h4>Mobile Money Integration</h4>
                <p>Connect local mobile money providers</p>
              </div>
              <div className={styles.step}>
                <div className={styles.stepNumber}>4</div>
                <h4>Gold Card Rollout</h4>
                <p>Begin citizen onboarding with Gold Cards</p>
              </div>
            </div>
          </div>

          <div className={styles.cta}>
            <button className="btn-primary">Register Your Country</button>
          </div>
        </div>
      </section>
    </div>
  )
}
