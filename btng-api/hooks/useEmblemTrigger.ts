'use client'

import { useState, useEffect } from 'react'
import { onEmblemTrigger, type EmblemTriggerEvent, type EmblemTriggerContext } from '@/lib/emblem-triggers'

/**
 * Hook to manage emblem trigger state in components
 */
export function useEmblemTrigger() {
  const [isTriggered, setIsTriggered] = useState(false)
  const [lastTrigger, setLastTrigger] = useState<EmblemTriggerContext | null>(null)

  useEffect(() => {
    const unsubscribe = onEmblemTrigger((context) => {
      setLastTrigger(context)
      setIsTriggered(true)
      
      // Auto-reset after trigger duration
      setTimeout(() => {
        setIsTriggered(false)
      }, 600) // Match CSS animation duration
    })

    return unsubscribe
  }, [])

  return {
    isTriggered,
    lastTrigger,
    trigger: isTriggered
  }
}

/**
 * Hook for manual emblem triggering
 */
export function useManualEmblemTrigger() {
  const [trigger, setTrigger] = useState(false)

  const activate = () => {
    setTrigger(true)
    setTimeout(() => setTrigger(false), 600)
  }

  return {
    trigger,
    activate
  }
}
