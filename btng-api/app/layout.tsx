import '../styles/globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import SovereignHeader from '@/components/SovereignHeader'
import SovereignFooter from '@/components/SovereignFooter'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'BTNG | Building Trust. Nurturing Growth.',
  description: 'Sovereign identity and value platform for institutional-grade trust operations',
  keywords: ['BTNG', 'sovereign identity', 'trust union', 'value platform'],
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <SovereignHeader />
        <main style={{ minHeight: '80vh' }}>
          {children}
        </main>
        <SovereignFooter />
      </body>
    </html>
  )
}
