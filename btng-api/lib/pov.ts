import { ValueProfile, ProofOfValue } from '@/types/trust-union'

/**
 * Proof-of-Value Management Module
 * Handles POV credential creation, verification, and value calculations
 */

export type ProofOfValueType = 'work' | 'trade' | 'trust' | 'contribution'

export interface POVCredential extends ProofOfValue {
  id: string
  timestamp: string
}

/**
 * Create a new value profile for a trust profile
 */
export function createValueProfile(trustProfileId: string): ValueProfile {
  return {
    profileId: `vp_${Date.now()}`,
    trustProfileId,
    totalValue: 0,
    currency: 'BTNG',
    proofOfValueCredentials: [],
    lastUpdated: new Date().toISOString()
  }
}

/**
 * Add a proof-of-value credential to a profile
 */
export function addProofOfValue(
  profile: ValueProfile,
  proof: Omit<ProofOfValue, 'id' | 'timestamp'>
): ValueProfile {
  const newProof: ProofOfValue = {
    ...proof,
    id: `pov_${Date.now()}`,
    timestamp: new Date().toISOString()
  }

  return {
    ...profile,
    proofOfValueCredentials: [...profile.proofOfValueCredentials, newProof],
    totalValue: profile.totalValue + proof.amount,
    lastUpdated: new Date().toISOString()
  }
}

/**
 * Calculate total verified value from all proofs
 */
export function calculateTotalValue(profile: ValueProfile): number {
  return profile.proofOfValueCredentials.reduce(
    (total, proof) => total + proof.amount,
    0
  )
}

/**
 * Get value breakdown by proof type
 */
export function getValueBreakdown(profile: ValueProfile) {
  const breakdown = {
    work: 0,
    trade: 0,
    trust: 0,
    contribution: 0
  }

  profile.proofOfValueCredentials.forEach(proof => {
    breakdown[proof.type] += proof.amount
  })

  return breakdown
}

/**
 * Verify proof-of-value credential authenticity
 */
export async function verifyProof(proof: ProofOfValue): Promise<boolean> {
  // Placeholder for cryptographic verification
  // Future: Implement zero-knowledge proof verification
  console.log(`Verifying proof: ${proof.id}`)
  return true
}

/**
 * Generate a unique value profile ID
 */
export function generateValueProfileId(): string {
  return `vp_${Date.now()}_${Math.random().toString(36).substring(2, 9)}`
}

/**
 * Validate proof-of-value type
 */
export function isValidPOVType(type: string): type is ProofOfValueType {
  return ['work', 'trade', 'trust', 'contribution'].includes(type)
}

/**
 * Get recent proofs from a value profile
 */
export function getRecentProofs(profile: ValueProfile, limit: number = 10): ProofOfValue[] {
  return [...profile.proofOfValueCredentials]
    .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
    .slice(0, limit)
}
