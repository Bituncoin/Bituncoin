/**
 * Mobile Money Adapter Interface
 * Enables integration with country-specific mobile money providers
 */

export interface MobileMoneyProvider {
  providerId: string
  name: string
  country: string
  currency: string
  apiEndpoint: string
}

export interface MobileMoneyAdapter {
  initialize(): Promise<void>
  sendMoney(to: string, amount: number): Promise<string>
  receiveMoney(from: string, amount: number): Promise<string>
  checkBalance(): Promise<number>
  getTransactionStatus(txId: string): Promise<string>
}

/**
 * Generic Mobile Money Adapter
 * Base implementation for country-specific adapters
 */
export class GenericMobileMoneyAdapter implements MobileMoneyAdapter {
  private provider: MobileMoneyProvider
  private initialized: boolean = false

  constructor(provider: MobileMoneyProvider) {
    this.provider = provider
  }

  async initialize(): Promise<void> {
    console.log(`Initializing ${this.provider.name} adapter for ${this.provider.country}`)
    // Future: Implement provider-specific initialization
    this.initialized = true
  }

  async sendMoney(to: string, amount: number): Promise<string> {
    if (!this.initialized) {
      throw new Error('Adapter not initialized')
    }
    console.log(`Sending ${amount} ${this.provider.currency} to ${to}`)
    // Future: Implement provider API integration
    return `mma_tx_${Date.now()}`
  }

  async receiveMoney(from: string, amount: number): Promise<string> {
    if (!this.initialized) {
      throw new Error('Adapter not initialized')
    }
    console.log(`Receiving ${amount} ${this.provider.currency} from ${from}`)
    return `mma_rx_${Date.now()}`
  }

  async checkBalance(): Promise<number> {
    if (!this.initialized) {
      throw new Error('Adapter not initialized')
    }
    // Future: Query provider API for balance
    return 0
  }

  async getTransactionStatus(txId: string): Promise<string> {
    console.log(`Checking status for transaction: ${txId}`)
    // Future: Query provider API for transaction status
    return 'confirmed'
  }
}

/**
 * Country-specific adapter examples (placeholders for future expansion)
 */

// Kenya - M-Pesa
export const MPesaAdapter = (apiKey: string): MobileMoneyAdapter => {
  return new GenericMobileMoneyAdapter({
    providerId: 'mpesa_ke',
    name: 'M-Pesa',
    country: 'Kenya',
    currency: 'KES',
    apiEndpoint: 'https://api.safaricom.co.ke/mpesa'
  })
}

// Uganda - MTN Mobile Money
export const MTNMobileMoneyAdapter = (apiKey: string): MobileMoneyAdapter => {
  return new GenericMobileMoneyAdapter({
    providerId: 'mtn_ug',
    name: 'MTN Mobile Money',
    country: 'Uganda',
    currency: 'UGX',
    apiEndpoint: 'https://api.mtn.com/mobile-money'
  })
}

// Ghana - Vodafone Cash
export const VodafoneCashAdapter = (apiKey: string): MobileMoneyAdapter => {
  return new GenericMobileMoneyAdapter({
    providerId: 'vodafone_gh',
    name: 'Vodafone Cash',
    country: 'Ghana',
    currency: 'GHS',
    apiEndpoint: 'https://api.vodafone.com.gh/mobile-money'
  })
}
