# Emblem Integration Examples

## 1. Gold Card Verification Success

```tsx
'use client'

import { useState } from 'react'
import SovereignEmblem from '@/components/SovereignEmblem'
import { triggerEmblemPulse } from '@/lib/emblem-triggers'

export default function CardVerificationSuccess() {
  const [showEmblem, setShowEmblem] = useState(false)

  const handleVerification = async () => {
    // Trigger emblem globally
    triggerEmblemPulse('identity-verification', { 
      profileId: 'tp_123',
      cardNumber: 'BTNG-XYZ1-ABC2'
    })
    
    setShowEmblem(true)
  }

  return (
    <div className="verification-success">
      {showEmblem && (
        <div className="emblem-container">
          <SovereignEmblem size={80} />
        </div>
      )}
      <h2>Gold Card Verified</h2>
      <button onClick={handleVerification}>Complete Verification</button>
    </div>
  )
}
```

## 2. Wallet Transaction Confirmation

```tsx
'use client'

import { useEmblemTrigger } from '@/hooks/useEmblemTrigger'
import { triggerEmblemPulse } from '@/lib/emblem-triggers'
import SovereignEmblem from '@/components/SovereignEmblem'

export default function WalletTransaction() {
  const { trigger } = useEmblemTrigger()

  const sendPayment = async (amount: number) => {
    // Process transaction
    await processWalletPayment(amount)
    
    // Trigger emblem pulse
    triggerEmblemPulse('wallet-transaction', { 
      amount, 
      currency: 'BTNG' 
    })
  }

  return (
    <div className="wallet-interface">
      <SovereignEmblem size={60} trigger={trigger} className="emblem-card" />
      <h3>Send Payment</h3>
      <button onClick={() => sendPayment(100)}>Send 100 BTNG</button>
    </div>
  )
}
```

## 3. Trust Score Update Animation

```tsx
'use client'

import { useState, useEffect } from 'react'
import SovereignEmblem from '@/components/SovereignEmblem'
import { triggerEmblemPulse } from '@/lib/emblem-triggers'

export default function TrustScoreDisplay({ userId }: { userId: string }) {
  const [score, setScore] = useState(500)
  const [trigger, setTrigger] = useState(false)

  const updateTrustScore = async (newScore: number) => {
    setScore(newScore)
    setTrigger(true)
    
    // Trigger global emblem pulse
    triggerEmblemPulse('trust-score-update', { 
      oldScore: score, 
      newScore,
      userId 
    })
    
    setTimeout(() => setTrigger(false), 800)
  }

  return (
    <div className="trust-score-panel">
      <SovereignEmblem 
        size={50} 
        trigger={trigger}
        onTriggerComplete={() => console.log('Score update complete')}
      />
      <div className="trust-score">{score}</div>
      <button onClick={() => updateTrustScore(score + 50)}>
        Add Proof of Value
      </button>
    </div>
  )
}
```

## 4. Country Activation Announcement

```tsx
'use client'

import { useState } from 'react'
import SovereignEmblem from '@/components/SovereignEmblem'
import { triggerEmblemPulse } from '@/lib/emblem-triggers'

export default function CountryActivation() {
  const [activated, setActivated] = useState(false)

  const activateCountry = (country: string) => {
    setActivated(true)
    
    // Longest trigger duration for major events
    triggerEmblemPulse('country-activation', { 
      country,
      nodeCount: 15,
      timestamp: new Date().toISOString()
    })
  }

  return (
    <div className="country-activation">
      <div className="activation-emblem">
        <SovereignEmblem size={120} trigger={activated} className="emblem-hero" />
      </div>
      <h1>Nigeria Trust Union Activated</h1>
      <button onClick={() => activateCountry('nigeria')}>
        Activate Node
      </button>
    </div>
  )
}
```

## 5. Proof-of-Value Submission

```tsx
'use client'

import { useManualEmblemTrigger } from '@/hooks/useEmblemTrigger'
import SovereignEmblem from '@/components/SovereignEmblem'
import { triggerEmblemPulse } from '@/lib/emblem-triggers'

export default function POVSubmission() {
  const { trigger, activate } = useManualEmblemTrigger()

  const submitPOV = async (data: POVData) => {
    const response = await fetch('/api/pov', {
      method: 'POST',
      body: JSON.stringify(data)
    })
    
    if (response.ok) {
      // Activate local emblem
      activate()
      
      // Trigger global pulse
      triggerEmblemPulse('pov-created', { 
        type: data.type, 
        amount: data.amount 
      })
    }
  }

  return (
    <form onSubmit={(e) => {
      e.preventDefault()
      submitPOV({ type: 'work', amount: 500 })
    }}>
      <SovereignEmblem size={40} trigger={trigger} />
      <h3>Submit Proof of Value</h3>
      <button type="submit">Submit</button>
    </form>
  )
}
```

## 6. Merchant Approval Notification

```tsx
'use client'

import { useEffect, useState } from 'react'
import SovereignEmblem from '@/components/SovereignEmblem'
import { onEmblemTrigger } from '@/lib/emblem-triggers'

export default function MerchantDashboard() {
  const [lastEvent, setLastEvent] = useState<string>('')
  const [trigger, setTrigger] = useState(false)

  useEffect(() => {
    // Listen for merchant approval events
    const unsubscribe = onEmblemTrigger((context) => {
      if (context.event === 'merchant-approval') {
        setLastEvent(`Merchant ${context.metadata?.merchantId} approved`)
        setTrigger(true)
        setTimeout(() => setTrigger(false), 800)
      }
    })

    return unsubscribe
  }, [])

  return (
    <div className="merchant-dashboard">
      <header>
        <SovereignEmblem size={45} trigger={trigger} />
        <h2>Merchant Portal</h2>
      </header>
      {lastEvent && (
        <div className="notification">
          {lastEvent}
        </div>
      )}
    </div>
  )
}
```

## 7. Global Header Integration (Already Implemented)

```tsx
import SovereignEmblem from './SovereignEmblem'

export default function SovereignHeader() {
  return (
    <header>
      <div className="logo">
        <SovereignEmblem size={40} className="emblem-header" />
        <div className="logo-text">
          <span>BTNG</span>
          <span>Building Trust</span>
        </div>
      </div>
    </header>
  )
}
```

## Styling Recommendations

```css
/* Emblem container centering */
.emblem-container {
  display: flex;
  justify-content: center;
  align-items: center;
  margin: var(--space-xl) 0;
}

/* Success screen emblem */
.verification-success .emblem-container {
  padding: var(--space-2xl);
  background: radial-gradient(
    circle,
    rgba(212, 175, 55, 0.1) 0%,
    transparent 70%
  );
}

/* Inline emblem with text */
.trust-score-panel {
  display: flex;
  align-items: center;
  gap: var(--space-md);
}

/* Hero emblem with dramatic entrance */
.activation-emblem {
  margin: var(--space-2xl) auto;
  animation: emblem-entrance 0.8s ease-out;
}

@keyframes emblem-entrance {
  0% {
    opacity: 0;
    transform: scale(0.5);
  }
  100% {
    opacity: 1;
    transform: scale(1);
  }
}
```

## Best Practices

1. **Use triggers sparingly** - Only for meaningful state changes
2. **Match trigger to event importance** - Higher multiplier for major events
3. **Provide completion callbacks** - For chaining animations
4. **Consider global vs local** - Global events affect all emblems
5. **Test rhythm synchronization** - Ensure baseline returns smoothly
6. **Add screen reader announcements** - "Identity verified", "Transaction complete"

## Accessibility

```tsx
<div role="status" aria-live="polite">
  <SovereignEmblem 
    size={60} 
    trigger={isVerifying}
    onTriggerComplete={() => {
      // Announce to screen readers
      const announcement = document.createElement('div')
      announcement.setAttribute('role', 'status')
      announcement.textContent = 'Verification complete'
      document.body.appendChild(announcement)
      setTimeout(() => announcement.remove(), 1000)
    }}
  />
  {isVerifying && <span className="sr-only">Verifying identity...</span>}
</div>
```
