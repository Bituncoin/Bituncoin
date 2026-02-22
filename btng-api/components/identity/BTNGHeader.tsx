import Link from 'next/link'
import styles from './BTNGHeader.module.css'

export default function BTNGHeader() {
  return (
    <header className={styles.header}>
      <div className={styles.container}>
        <Link href="/" className={styles.logo}>
          <span className={styles.logoText}>BTNG</span>
          <span className={styles.logoTagline}>Building Trust</span>
        </Link>
        
        <nav className={styles.nav}>
          <div className={styles.dropdown}>
            <span className={styles.navLink}>Identity</span>
            <div className={styles.dropdownContent}>
              <Link href="/(identity)/card">Gold Card</Link>
              <Link href="/(identity)/wallet">QR Wallet</Link>
              <Link href="/(identity)/profile">My Profile</Link>
            </div>
          </div>
          
          <div className={styles.dropdown}>
            <span className={styles.navLink}>Countries</span>
            <div className={styles.dropdownContent}>
              <Link href="/(countries)/ghana">Ghana 🇬🇭</Link>
              <Link href="/(countries)/kenya">Kenya 🇰🇪</Link>
              <Link href="/(countries)/nigeria">Nigeria 🇳🇬</Link>
              <Link href="/(countries)/togo">Togo 🇹🇬</Link>
              <Link href="/(countries)/uganda">Uganda 🇺🇬</Link>
              <Link href="/(countries)/ivory-coast">Ivory Coast 🇨🇮</Link>
              <Link href="/(countries)/burkina-faso">Burkina Faso 🇧🇫</Link>
              <Link href="/(countries)/south-africa">South Africa 🇿🇦</Link>
            </div>
          </div>
          
          <Link href="/trust-union" className={styles.navLink}>
            Trust Union
          </Link>
          <Link href="/(onboarding)/user" className={`${styles.navLink} ${styles.goldCard}`}>
            Get Gold Card
          </Link>
        </nav>
      </div>
    </header>
  )
}
