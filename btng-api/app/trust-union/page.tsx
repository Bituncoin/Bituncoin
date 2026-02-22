import styles from './page.module.css'

export default function TrustUnionPage() {
  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className="container">
          <h1 className={styles.title}>Trust Union Protocol</h1>
          <p className={styles.subtitle}>
            Institutional-grade trust architecture for sovereign operations
          </p>
        </div>
      </section>

      <section className={styles.content}>
        <div className="container">
          <div className={styles.overview}>
            <h2>Protocol Overview</h2>
            <p>
              The Trust Union Protocol is a distributed trust architecture that enables 
              sovereign nations, institutions, and individuals to participate in a unified 
              trust network with zero-knowledge proof verification.
            </p>
          </div>

          <div className={styles.features}>
            <h2>Core Features</h2>
            <div className={styles.featureGrid}>
              <div className={styles.feature}>
                <h3>🌐 Distributed Trust Network</h3>
                <p>Decentralized nodes across sovereign nations</p>
              </div>
              <div className={styles.feature}>
                <h3>🔐 Zero-Knowledge Proofs</h3>
                <p>Privacy-preserving transaction validation</p>
              </div>
              <div className={styles.feature}>
                <h3>⚖️ Institutional Grade</h3>
                <p>Enterprise-level security and compliance</p>
              </div>
              <div className={styles.feature}>
                <h3>🤝 Cross-Border Operations</h3>
                <p>Seamless trust transfer across jurisdictions</p>
              </div>
            </div>
          </div>

          <div className={styles.architecture}>
            <h2>Architecture</h2>
            <div className={styles.layers}>
              <div className={styles.layer}>
                <h4>Identity Layer</h4>
                <p>Gold Card credentials with universal recognition</p>
              </div>
              <div className={styles.layer}>
                <h4>Trust Layer</h4>
                <p>Proof-of-value calculation and trust scoring</p>
              </div>
              <div className={styles.layer}>
                <h4>Transaction Layer</h4>
                <p>QR wallet with zero-knowledge validation</p>
              </div>
              <div className={styles.layer}>
                <h4>Protocol Layer</h4>
                <p>Distributed consensus and node management</p>
              </div>
            </div>
          </div>

          <div className={styles.callout}>
            <h3>Join the Trust Union</h3>
            <p>
              Countries, institutions, and organizations can become Trust Union nodes 
              to participate in the sovereign trust network.
            </p>
            <button className="btn-primary">Request Node Registration</button>
          </div>
        </div>
      </section>
    </div>
  )
}
