/**
 * Sovereign Emblem Trigger System
 * Manages emblem intensification events across the platform
 */

export type EmblemTriggerEvent = 
  | 'identity-verification'
  | 'wallet-transaction'
  | 'pov-created'
  | 'trust-score-update'
  | 'country-activation'
  | 'merchant-approval'

export interface EmblemTriggerContext {
  event: EmblemTriggerEvent
  timestamp: string
  metadata?: Record<string, any>
}

/**
 * Global emblem trigger listeners
 */
const triggerListeners = new Set<(context: EmblemTriggerContext) => void>()

/**
 * Register a listener for emblem trigger events
 */
export function onEmblemTrigger(callback: (context: EmblemTriggerContext) => void): () => void {
  triggerListeners.add(callback)
  return () => triggerListeners.delete(callback)
}

/**
 * Dispatch an emblem trigger event
 */
export function triggerEmblemPulse(event: EmblemTriggerEvent, metadata?: Record<string, any>): void {
  const context: EmblemTriggerContext = {
    event,
    timestamp: new Date().toISOString(),
    metadata
  }

  triggerListeners.forEach(listener => {
    try {
      listener(context)
    } catch (error) {
      console.error('Emblem trigger listener error:', error)
    }
  })
}

/**
 * Emblem trigger configurations
 * Maps events to their intensification characteristics
 */
export const TRIGGER_CONFIGS: Record<EmblemTriggerEvent, {
  duration: number
  multiplier: number
  description: string
}> = {
  'identity-verification': {
    duration: 600,
    multiplier: 1.6,
    description: 'Gold Card or profile verification'
  },
  'wallet-transaction': {
    duration: 600,
    multiplier: 1.6,
    description: 'QR Wallet transaction processed'
  },
  'pov-created': {
    duration: 600,
    multiplier: 1.6,
    description: 'New proof-of-value credential added'
  },
  'trust-score-update': {
    duration: 800,
    multiplier: 1.8,
    description: 'Trust score calculation completed'
  },
  'country-activation': {
    duration: 1000,
    multiplier: 2.0,
    description: 'New country node activated'
  },
  'merchant-approval': {
    duration: 800,
    multiplier: 1.8,
    description: 'Merchant application approved'
  }
}

/**
 * Check if emblem should be in trigger mode
 * Used for real-time UI updates
 */
export function shouldTriggerEmblem(event: EmblemTriggerEvent): boolean {
  // Future: Add conditional logic based on user permissions, system state, etc.
  return true
}

/**
 * Get trigger configuration for an event
 */
export function getTriggerConfig(event: EmblemTriggerEvent) {
  return TRIGGER_CONFIGS[event]
}
