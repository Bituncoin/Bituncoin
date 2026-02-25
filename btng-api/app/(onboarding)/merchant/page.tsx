import Link from 'next/link'
import styles from './page.module.css'

export default function MerchantOnboardingPage() {
  return (
    <div className={styles.container}>
      <section className={styles.hero}>
        <div className="container">
          <h1 className={styles.title}>Become a BTNG Merchant</h1>
          <p className={styles.subtitle}>
            Accept sovereign payments and join the Trust Union network
          </p>
        </div>
      </section>

      <section className={styles.content}>
        <div className="container">
          <div className={styles.grid}>
            <div className={styles.info}>
              <h2>Why Join as a Merchant?</h2>
              <p>
                BTNG merchants gain access to a growing network of Gold Card holders
                across 8 countries with instant, low-cost payment processing.
              </p>

              <div className={styles.benefits}>
                <h3>Merchant Benefits</h3>
                <div className={styles.benefitItem}>
                  <span className={styles.icon}>💰</span>
                  <div>
                    <h4>Low Transaction Fees</h4>
                    <p>Competitive rates via Trust Union protocol</p>
                  </div>
                </div>
                <div className={styles.benefitItem}>
                  <span className={styles.icon}>⚡</span>
                  <div>
                    <h4>Instant Settlement</h4>
                    <p>Receive funds in real-time with QR wallet</p>
                  </div>
                </div>
                <div className={styles.benefitItem}>
                  <span className={styles.icon}>🌍</span>
                  <div>
                    <h4>Cross-Border Payments</h4>
                    <p>Accept payments from customers in any BTNG country</p>
                  </div>
                </div>
                <div className={styles.benefitItem}>
                  <span className={styles.icon}>🛡️</span>
                  <div>
                    <h4>Fraud Protection</h4>
                    <p>Zero-knowledge proof verification reduces fraud</p>
                  </div>
                </div>
                <div className={styles.benefitItem}>
                  <span className={styles.icon}>📊</span>
                  <div>
                    <h4>Business Dashboard</h4>
                    <p>Track sales, transactions, and analytics</p>
                  </div>
                </div>
                <div className={styles.benefitItem}>
                  <span className={styles.icon}>🤝</span>
                  <div>
                    <h4>Trust Score Boost</h4>
                    <p>Build reputation in the Trust Union network</p>
                  </div>
                </div>
              </div>

              <div className={styles.requirements}>
                <h3>Requirements</h3>
                <ul>
                  <li>Valid business registration in one of our 8 countries</li>
                  <li>Physical or online retail presence</li>
                  <li>Mobile money account or bank account</li>
                  <li>Business identity verification</li>
                  <li>Compliance with local payment regulations</li>
                </ul>
              </div>
            </div>

            <div className={styles.formSection}>
              <div className={styles.form}>
                <h3>Merchant Application</h3>
                <form className={styles.merchantForm}>
                  <div className={styles.formGroup}>
                    <label>Business Name</label>
                    <input type="text" placeholder="Enter your business name" />
                  </div>
                  <div className={styles.formGroup}>
                    <label>Business Type</label>
                    <select>
                      <option value="">Select business type</option>
                      <option value="retail">Retail Store</option>
                      <option value="restaurant">Restaurant/Café</option>
                      <option value="services">Services</option>
                      <option value="hospitality">Hospitality</option>
                      <option value="ecommerce">E-commerce</option>
                      <option value="other">Other</option>
                    </select>
                  </div>
                  <div className={styles.formGroup}>
                    <label>Country of Operation</label>
                    <select>
                      <option value="">Select country</option>
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
                    <label>Business Registration Number</label>
                    <input type="text" placeholder="Registration number" />
                  </div>
                  <div className={styles.formGroup}>
                    <label>Contact Person Name</label>
                    <input type="text" placeholder="Full name" />
                  </div>
                  <div className={styles.formGroup}>
                    <label>Business Phone</label>
                    <input type="tel" placeholder="+XXX XXX XXXX" />
                  </div>
                  <div className={styles.formGroup}>
                    <label>Business Email</label>
                    <input type="email" placeholder="business@example.com" />
                  </div>
                  <div className={styles.formGroup}>
                    <label>Physical Address</label>
                    <textarea rows={3} placeholder="Enter your business address"></textarea>
                  </div>
                  <div className={styles.formGroup}>
                    <label>Expected Monthly Transaction Volume</label>
                    <select>
                      <option value="">Select range</option>
                      <option value="low">Under $1,000</option>
                      <option value="medium">$1,000 - $10,000</option>
                      <option value="high">$10,000 - $100,000</option>
                      <option value="enterprise">Over $100,000</option>
                    </select>
                  </div>
                  <div className={styles.formGroup}>
                    <label>
                      <input type="checkbox" />
                      <span>I confirm that all information is accurate and agree to merchant terms</span>
                    </label>
                  </div>
                  <button type="submit" className="btn-primary" style={{ width: '100%' }}>
                    Submit Merchant Application
                  </button>
                </form>
              </div>

              <div className={styles.support}>
                <h4>Need Help?</h4>
                <p>Contact our merchant support team:</p>
                <div className={styles.contactItem}>
                  <span>📧</span>
                  <span>merchants@btng.global</span>
                </div>
                <div className={styles.contactItem}>
                  <span>💬</span>
                  <span>Live Chat Support</span>
                </div>

                <div className={styles.packageSection}>
                  <h4>Merchant App Deployment Package</h4>
                  <p>Download the latest install package and sovereign bootstrap config for your first 10 vendors.</p>
                  <div className={styles.packageActions}>
                    <a href="/api/merchant/deployment-package" className="btn-primary">
                      Download Package JSON
                    </a>
                    <Link href="/watchtower" className="btn-secondary">
                      Open Watchtower Map
                    </Link>
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
