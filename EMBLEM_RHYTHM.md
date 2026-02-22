# Sovereign Emblem Rhythm System

**Implementation Date:** February 18, 2026  
**Status:** ✅ Operational

## Overview

The BTNG platform emblem uses a **metronomic pulse rhythm** to signal institutional reliability, operational state, and verification events. This creates a living seal that reinforces sovereign authority through precise, disciplined animation.

## Core Principles

### 1. Metronomic Cadence
- **Fixed cycle:** 2.4 seconds (controlled by `--emblem-pulse-cycle` token)
- **Linear timing:** No easing curves—exact, predictable intervals
- **Global synchronization:** All emblems across the platform pulse in unison
- **Institutional reliability:** Rhythmic consistency signals precision-driven governance

### 2. Dual Glow Modes

#### Continuous Mode (Baseline)
The emblem maintains a steady sovereign pulse at all times:
- **Glow range:** Base (30% opacity) → Peak (80% opacity)
- **Timing:** Fixed 2.4s cycle with linear progression
- **Visual effect:** Drop shadow pulses from 8px to 16px
- **Opacity shift:** 0.9 → 1.0 → 0.9
- **Meaning:** System is active, identity layer is sovereign

#### Trigger Mode (Intensification)
During verification events, the emblem intensifies:
- **Speed multiplier:** 1.6x faster (cycle becomes 1.5s)
- **Burst effect:** Radial glow expands outward
- **Duration:** 600ms overlay on top of baseline pulse
- **Return:** Seamless transition back to continuous mode
- **Meaning:** Active verification/transaction in progress

### 3. Architecture

```
Continuous Pulse (Always Active)
    │
    ├──> Baseline rhythm: 2.4s linear cycle
    │
    └──> Trigger Events
         │
         ├──> Intensification overlay (1.6x multiplier)
         ├──> Burst animation (600ms)
         └──> Return to baseline (no drift)
```

## Implementation

### CSS Tokens (styles/tokens.css)
```css
--emblem-pulse-cycle: 2.4s;              /* Fixed metronomic interval */
--emblem-pulse-timing: linear;            /* No easing—exact rhythm */
--emblem-glow-base: rgba(212, 175, 55, 0.3);
--emblem-glow-peak: rgba(212, 175, 55, 0.8);
--emblem-glow-trigger: rgba(212, 175, 55, 1.0);
--emblem-intensify-multiplier: 1.6;       /* Trigger speed increase */
```

### Animation Keyframes (styles/identity.css)
```css
@keyframes emblem-sovereign-pulse {
  0%   { filter: drop-shadow(0 0 8px var(--emblem-glow-base)); opacity: 0.9; }
  50%  { filter: drop-shadow(0 0 16px var(--emblem-glow-peak)); opacity: 1; }
  100% { filter: drop-shadow(0 0 8px var(--emblem-glow-base)); opacity: 0.9; }
}
```

### Component (components/SovereignEmblem.tsx)
React component that manages emblem state and trigger events:
- **Props:** `size`, `trigger`, `onTriggerComplete`
- **Auto-reset:** Returns to baseline after 600ms
- **Reactive:** Responds to global trigger events

### Trigger System (lib/emblem-triggers.ts)
Event-driven system for emblem intensification:
```typescript
triggerEmblemPulse('identity-verification', { profileId: 'xyz' })
triggerEmblemPulse('wallet-transaction', { amount: 500 })
triggerEmblemPulse('trust-score-update', { newScore: 850 })
```

## Trigger Events

| Event | Duration | Multiplier | Context |
|-------|----------|------------|---------|
| `identity-verification` | 600ms | 1.6x | Gold Card creation, profile verification |
| `wallet-transaction` | 600ms | 1.6x | QR Wallet send/receive completion |
| `pov-created` | 600ms | 1.6x | New proof-of-value credential added |
| `trust-score-update` | 800ms | 1.8x | Trust score recalculation |
| `country-activation` | 1000ms | 2.0x | New Trust Union node goes live |
| `merchant-approval` | 800ms | 1.8x | Merchant application approved |

## Usage Patterns

### Static Display (Header, Footer)
```tsx
import SovereignEmblem from '@/components/SovereignEmblem'

<SovereignEmblem size={40} className="emblem-header" />
```

### With Manual Trigger
```tsx
import SovereignEmblem from '@/components/SovereignEmblem'
import { useState } from 'react'

const [trigger, setTrigger] = useState(false)

const handleVerification = () => {
  setTrigger(true)
  // Trigger resets automatically after 600ms
}

<SovereignEmblem 
  size={60} 
  trigger={trigger}
  onTriggerComplete={() => console.log('Pulse complete')}
/>
```

### With Global Event System
```tsx
import SovereignEmblem from '@/components/SovereignEmblem'
import { useEmblemTrigger } from '@/hooks/useEmblemTrigger'
import { triggerEmblemPulse } from '@/lib/emblem-triggers'

const { trigger } = useEmblemTrigger()

const handleTransaction = async () => {
  await processTransaction()
  triggerEmblemPulse('wallet-transaction', { amount: 100 })
}

<SovereignEmblem size={80} trigger={trigger} />
```

## Size Variants

| Class | Size | Context |
|-------|------|---------|
| `emblem-header` | 40px | Navigation header |
| `emblem-card` | 32px | Card corners, inline icons |
| `emblem-hero` | 120px | Landing page hero section |
| Custom | Any | Pass `size` prop directly |

## Behavioral Characteristics

### Predictability
- Same cycle duration across all pages
- Linear timing ensures exact intervals
- No randomness or organic variation

### Authority Signaling
- Glow pattern suggests institutional governance
- Rhythmic precision mirrors value verification
- Visual stability reinforces trust

### Operational Clarity
- Continuous pulse = system active
- Intensified pulse = verification in progress
- Smooth return = operation complete

### Institutional Feel
- Engineered rhythm (not decorative)
- Heartbeat of sovereign system
- Living seal that never drifts

## Integration Points

### Current Locations
1. **Header** - [components/SovereignHeader.tsx](components/SovereignHeader.tsx)
2. **Landing Page Hero** - [app/page.tsx](app/page.tsx)

### Recommended Additions
- Gold Card verification success screen
- Wallet transaction confirmation
- Proof-of-value submission feedback
- Country activation announcements
- Merchant approval notifications
- Trust score milestone celebrations

## Technical Notes

### Performance
- Uses CSS `will-change: filter, opacity` for GPU acceleration
- Minimal repaints due to isolated stacking context
- No JavaScript animation loops (pure CSS)

### Accessibility
- Animation respects `prefers-reduced-motion`
- Trigger events include screen reader announcements
- Visual feedback paired with status text

### Browser Support
- Modern browsers with CSS animation support
- SVG rendering for emblem graphics
- Graceful degradation for older browsers

## Future Enhancements

### Phase 2 (Planned)
- Multi-emblem synchronization across tabs
- Adaptive rhythm based on system load
- Color shift for different trust levels
- Audio feedback for trigger events

### Phase 3 (Exploratory)
- Haptic feedback on mobile devices
- Emblem as Web3 identity signature
- Real-time network consensus visualization

## Design Philosophy

The emblem's metronomic rhythm is not decorative—it's a **functional identity signal**. Every pulse communicates:
- The system is operational
- Sovereign governance is active
- Trust is being continuously verified
- Institutional precision is maintained

This rhythm shapes how users perceive BTNG: not as a static platform, but as a **living sovereign system** that governs identity, value, and trust with disciplined, metronomic reliability.

---

**Next Steps:**
1. Add emblem to verification success screens
2. Integrate trigger events in API responses
3. Test synchronization across multiple tabs
4. Gather user feedback on rhythm timing
