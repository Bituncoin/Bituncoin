/**
 * Wallet Management Module
 * Handles QR wallet operations, transactions, and balance management
 */

export interface WalletData {
  walletId: string
  balance: string
  currency: string
  ownerId: string
  createdAt: string
}

export interface Transaction {
  id: string
  fromWallet: string
  toWallet: string
  amount: number
  currency: string
  timestamp: string
  status: 'pending' | 'completed' | 'failed'
  type: 'send' | 'receive'
}

/**
 * Generate a unique wallet ID
 */
export function generateWalletId(): string {
  const prefix = 'BTNG'
  const timestamp = Date.now().toString(36).toUpperCase()
  const random = Math.random().toString(36).substring(2, 8).toUpperCase()
  return `${prefix}-${timestamp}-${random}`
}

/**
 * Create a new wallet for a Gold Card holder
 */
export function createWallet(ownerId: string, currency: string = 'BTNG'): WalletData {
  return {
    walletId: generateWalletId(),
    balance: '0.00',
    currency,
    ownerId,
    createdAt: new Date().toISOString()
  }
}

/**
 * Validate wallet ID format
 */
export function validateWalletId(walletId: string): boolean {
  const pattern = /^BTNG-[A-Z0-9]+-[A-Z0-9]+$/
  return pattern.test(walletId)
}

/**
 * Format balance for display
 */
export function formatBalance(balance: number, currency: string = 'BTNG'): string {
  return `${balance.toFixed(2)} ${currency}`
}

/**
 * Create a transaction record
 */
export function createTransaction(
  fromWallet: string,
  toWallet: string,
  amount: number,
  currency: string = 'BTNG'
): Transaction {
  return {
    id: `txn_${Date.now()}`,
    fromWallet,
    toWallet,
    amount,
    currency,
    timestamp: new Date().toISOString(),
    status: 'pending',
    type: 'send'
  }
}

/**
 * Validate transaction amount
 */
export function validateTransactionAmount(
  amount: number,
  balance: number
): { valid: boolean; error?: string } {
  if (amount <= 0) {
    return { valid: false, error: 'Amount must be greater than zero' }
  }
  if (amount > balance) {
    return { valid: false, error: 'Insufficient balance' }
  }
  return { valid: true }
}

/**
 * Generate QR code data for wallet
 */
export function generateWalletQRData(walletId: string): string {
  return `btng:wallet:${walletId}`
}

/**
 * Parse QR code data to extract wallet ID
 */
export function parseWalletQRData(qrData: string): string | null {
  const match = qrData.match(/^btng:wallet:(.+)$/)
  return match ? match[1] : null
}

/**
 * Calculate transaction fee
 */
export function calculateTransactionFee(
  amount: number,
  crossBorder: boolean = false
): number {
  const baseRate = 0.001 // 0.1%
  const crossBorderRate = 0.005 // 0.5%
  const rate = crossBorder ? crossBorderRate : baseRate
  return Math.max(0.01, amount * rate) // Minimum fee 0.01
}
