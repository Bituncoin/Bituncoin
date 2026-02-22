// Trust Union Protocol Types
export interface TrustProfile {
  id: string
  holderName: string
  goldCardNumber: string
  issueDate: string
  trustScore: number
  verificationLevel: 'basic' | 'enhanced' | 'sovereign'
  activeCountries: string[]
}

export interface ValueProfile {
  profileId: string
  trustProfileId: string
  totalValue: number
  currency: string
  proofOfValueCredentials: ProofOfValue[]
  lastUpdated: string
}

export interface ProofOfValue {
  id: string
  type: 'work' | 'trade' | 'trust' | 'contribution'
  amount: number
  verifiedBy: string
  timestamp: string
  metadata: Record<string, any>
}

export interface TrustUnionNode {
  nodeId: string
  country: string
  status: 'active' | 'pending' | 'suspended'
  trustCapacity: number
  connectedProfiles: number
}

export interface WalletTransaction {
  txId: string
  from: string
  to: string
  amount: number
  currency: string
  timestamp: string
  status: 'pending' | 'confirmed' | 'failed'
  proofHash: string
}
