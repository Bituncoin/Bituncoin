# Sovereign Emblem - Quick Reference

## Import & Basic Usage

```tsx
import SovereignEmblem from '@/components/SovereignEmblem'

// Static display
<SovereignEmblem size={40} />

// With trigger event
<SovereignEmblem size={60} trigger={isVerifying} />
```

## Metronomic Rhythm Characteristics

**Baseline:** 2.4s linear cycle (continuous sovereignty signal)  
**Trigger:** 1.5s intensified cycle (verification in progress)  
**Timing:** Linear—no easing, exact institutional cadence  
**Return:** Seamless baseline restoration after events

## Trigger Events

```typescript
import { triggerEmblemPulse } from '@/lib/emblem-triggers'

// Identity verification
triggerEmblemPulse('identity-verification', { profileId })

// Wallet transaction
triggerEmblemPulse('wallet-transaction', { amount })

// Trust score update
triggerEmblemPulse('trust-score-update', { newScore })

// POV credential
triggerEmblemPulse('pov-created', { povId })

// Country activation
triggerEmblemPulse('country-activation', { country })

// Merchant approval
triggerEmblemPulse('merchant-approval', { merchantId })
```

## Hook Usage

```tsx
import { useEmblemTrigger } from '@/hooks/useEmblemTrigger'

function Component() {
  const { trigger } = useEmblemTrigger()
  
  return <SovereignEmblem size={80} trigger={trigger} />
}
```

## Design Tokens

```css
--emblem-pulse-cycle: 2.4s;
--emblem-pulse-timing: linear;
--emblem-intensify-multiplier: 1.6;
```

## Size Classes

- `emblem-header` - 40px (navigation)
- `emblem-card` - 32px (inline)
- `emblem-hero` - 120px (landing)

See [EMBLEM_RHYTHM.md](EMBLEM_RHYTHM.md) for complete documentation.
