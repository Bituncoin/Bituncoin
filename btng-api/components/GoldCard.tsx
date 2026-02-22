import styles from './GoldCard.module.css'

interface GoldCardProps {
  holderName: string
  cardNumber: string
  issueDate: string
  isPreview?: boolean
}

export default function GoldCard({ 
  holderName, 
  cardNumber, 
  issueDate, 
  isPreview = false 
}: GoldCardProps) {
  return (
    <div className={styles.card}>
      <div className={styles.header}>
        <div className={styles.logo}>BTNG</div>
        {isPreview && <div className={styles.previewBadge}>PREVIEW</div>}
      </div>
      
      <div className={styles.chip}>
        <div className={styles.chipPattern}></div>
      </div>

      <div className={styles.cardNumber}>{cardNumber}</div>

      <div className={styles.footer}>
        <div className={styles.info}>
          <div className={styles.label}>Cardholder</div>
          <div className={styles.name}>{holderName}</div>
        </div>
        <div className={styles.info}>
          <div className={styles.label}>Issued</div>
          <div className={styles.date}>{issueDate}</div>
        </div>
      </div>

      <div className={styles.badge}>
        <span className={styles.badgeText}>🛡️ SOVEREIGN</span>
      </div>
    </div>
  )
}
