import BTNG54AfricaGoldCoin from '@/components/BTNG54AfricaGoldCoin'
import styles from './page.module.css'

export default function GoldCoinPage() {
  return (
    <div className={styles.container}>
      <BTNG54AfricaGoldCoin />
    </div>
  )
}

export const metadata = {
  title: 'BTNG 54 Africa Gold Coin | Sovereign Digital Gold Standard',
  description: 'The unified sovereign digital gold standard for all 54 African nations. One BTNG = One gram of pure African gold, backed by the collective wealth of Africa.',
  keywords: 'BTNG, Africa, Gold Coin, Sovereign, Digital Currency, African Union, Gold Standard',
  openGraph: {
    title: 'BTNG 54 Africa Gold Coin',
    description: 'The unified sovereign digital gold standard for all 54 African nations',
    type: 'website',
    images: [
      {
        url: '/gold-coin-og.png',
        width: 1200,
        height: 630,
        alt: 'BTNG 54 Africa Gold Coin'
      }
    ]
  }
}