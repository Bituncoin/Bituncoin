# BTNG Platform - Current Structure

Last Updated: February 18, 2026

## Complete Directory Tree

```
BTNGAI_files/
│
├── app/
│   ├── layout.tsx                          # Root layout with SovereignHeader/Footer
│   ├── page.tsx                            # Landing page
│   ├── page.module.css
│   │
│   ├── (identity)/                         # Identity route group
│   │   ├── card/
│   │   │   ├── page.tsx                    # Gold Card display
│   │   │   └── page.module.css
│   │   ├── wallet/
│   │   │   ├── page.tsx                    # QR Wallet interface
│   │   │   └── page.module.css
│   │   └── profile/
│   │       ├── page.tsx                    # Universal identity profile
│   │       └── page.module.css
│   │
│   ├── (countries)/                        # Countries route group
│   │   ├── ghana/
│   │   │   └── page.tsx                    # Ghana BTNG presence
│   │   ├── kenya/
│   │   │   └── page.tsx                    # Kenya BTNG presence
│   │   ├── nigeria/
│   │   │   └── page.tsx                    # Nigeria BTNG presence
│   │   ├── togo/
│   │   │   └── page.tsx                    # Togo BTNG presence
│   │   ├── uganda/
│   │   │   └── page.tsx                    # Uganda BTNG presence
│   │   ├── ivory-coast/
│   │   │   └── page.tsx                    # Ivory Coast BTNG presence
│   │   ├── burkina-faso/
│   │   │   └── page.tsx                    # Burkina Faso BTNG presence
│   │   ├── south-africa/
│   │   │   └── page.tsx                    # South Africa BTNG presence
│   │   └── country.module.css              # Shared country styles
│   │
│   ├── (onboarding)/                       # Onboarding route group
│   │   ├── user/
│   │   │   ├── page.tsx                    # User Gold Card application
│   │   │   └── page.module.css
│   │   └── merchant/
│   │       ├── page.tsx                    # Merchant registration
│   │       └── page.module.css
│   │
│   └── api/                                # API routes
│       ├── health/
│       │   └── route.ts                    # Health check endpoint
│       ├── identity/
│       │   └── create/
│       │       └── route.ts                # Create trust profile
│       ├── pov/
│       │   └── route.ts                    # Proof-of-value CRUD
│       └── wallet/
│           └── transaction/
│               └── route.ts                # Wallet transactions
│
├── components/                             # Root-level components
│   ├── SovereignHeader.tsx                 # Platform header (navigation)
│   ├── SovereignHeader.module.css
│   ├── SovereignFooter.tsx                 # Platform footer
│   ├── SovereignFooter.module.css
│   ├── GoldCard.tsx                        # Gold Card component
│   ├── GoldCard.module.css
│   ├── QRWallet.tsx                        # QR Wallet component
│   └── QRWallet.module.css
│
├── lib/                                    # Business logic modules
│   ├── identity.ts                         # Identity & trust scoring
│   ├── pov.ts                              # Proof-of-value management
│   ├── wallet.ts                           # Wallet operations
│   ├── countries.ts                        # Country data (8 nations)
│   └── validators.ts                       # Validation utilities
│
├── public/                                 # Static assets
│   ├── emblem.svg                          # BTNG platform emblem
│   ├── README.md                           # Asset guidelines
│   ├── card/                               # Gold Card assets
│   │   └── .gitkeep.md
│   └── icons/                              # Platform icons
│       └── .gitkeep.md
│
├── styles/                                 # Global styles
│   ├── globals.css                         # Base styles + imports
│   ├── tokens.css                          # Design system tokens
│   └── identity.css                        # Identity-specific styles
│
├── types/                                  # TypeScript types
│   └── trust-union.ts                      # Trust Union type definitions
│
├── scripts/                                # Utility scripts
│   └── health-check.js                     # Platform health check
│
├── .env.example                            # Environment template
├── .gitignore                              # Git ignore rules
├── next.config.js                          # Next.js configuration
├── package.json                            # Dependencies & scripts
├── tsconfig.json                           # TypeScript configuration
├── README.md                               # Main documentation
├── ARCHITECTURE.md                         # Architecture overview
├── DEPLOYMENT.md                           # Deployment guide
└── RESTRUCTURE.md                          # Restructure documentation
```

## Key Modules

### Components
| File | Purpose | Exports |
|------|---------|---------|
| `SovereignHeader.tsx` | Navigation with dropdown menus | `default SovereignHeader` |
| `SovereignFooter.tsx` | Footer with links and branding | `default SovereignFooter` |
| `GoldCard.tsx` | Gold Card display | `default GoldCard` |
| `QRWallet.tsx` | QR code wallet interface | `default QRWallet` |

### Library Modules
| File | Purpose | Key Exports |
|------|---------|-------------|
| `identity.ts` | Gold Card & trust operations | `generateGoldCardNumber`, `calculateTrustScore`, `verifyProfile`, `createGoldCardCredentials` |
| `pov.ts` | Proof-of-value credentials | `createValueProfile`, `addProofOfValue`, `calculateTotalValue`, `verifyProof` |
| `wallet.ts` | Wallet & transactions | `generateWalletId`, `createWallet`, `createTransaction`, `calculateTransactionFee` |
| `countries.ts` | Country data & queries | `countryData`, `getCountriesByStatus`, `getTotalActiveUsers`, `getTotalTrustNodes` |
| `validators.ts` | Input validation | `validateEmail`, `validatePhone`, `validateGoldCard`, `validateWallet`, `validatePOVType` |

### Routes Overview
| Route Pattern | Component | Purpose |
|---------------|-----------|---------|
| `/` | `app/page.tsx` | Landing page |
| `/(identity)/card` | `app/(identity)/card/page.tsx` | Gold Card display |
| `/(identity)/wallet` | `app/(identity)/wallet/page.tsx` | QR Wallet |
| `/(identity)/profile` | `app/(identity)/profile/page.tsx` | User profile |
| `/(countries)/{country}` | `app/(countries)/{country}/page.tsx` | Country pages (8 total) |
| `/(onboarding)/user` | `app/(onboarding)/user/page.tsx` | User registration |
| `/(onboarding)/merchant` | `app/(onboarding)/merchant/page.tsx` | Merchant application |

### API Endpoints
| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/health` | GET | Health check |
| `/api/identity/create` | POST | Create trust profile & Gold Card |
| `/api/pov` | GET, POST, PUT | Proof-of-value CRUD operations |
| `/api/wallet/transaction` | POST | Process wallet transactions |

## Technology Stack

- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript 5.3
- **UI Library:** React 18
- **Styling:** CSS Modules + Design Tokens
- **QR Codes:** qrcode.react
- **Runtime:** Node.js ≥18.0.0

## Design System

### Color Tokens (from `styles/tokens.css`)
- **Gold:** `#D4AF37` (primary brand)
- **Trust Blue:** `#1A4B8C` (institutional)
- **Sovereign Black:** `#0A0A0A`
- **Sovereign White:** `#FEFEFE`

### Spacing Scale
- `--space-xs`: 0.25rem (4px)
- `--space-sm`: 0.5rem (8px)
- `--space-md`: 1rem (16px)
- `--space-lg`: 1.5rem (24px)
- `--space-xl`: 2rem (32px)
- `--space-2xl`: 3rem (48px)

### Typography
- **Display/Body:** Inter, system fallbacks  
- **Monospace:** Menlo, Monaco, Courier New

## Country Coverage

### Active (5 countries)
1. 🇬🇭 Ghana - 145,000 users
2. 🇰🇪 Kenya - 320,000 users
3. 🇹🇬 Togo - 42,000 users
4. 🇺🇬 Uganda - 125,000 users

### Launching (3 countries)
5. 🇳🇬 Nigeria - 89,000 users
6. 🇨🇮 Ivory Coast - 34,000 users
7. 🇿🇦 South Africa - 67,000 users

### Planned (1 country)
8. 🇧🇫 Burkina Faso - 8,500 users

**Total:** 830,500+ Gold Card holders across 8 African nations

## Development Commands

```bash
# Install dependencies
npm install

# Development server (http://localhost:3000)
npm run dev

# Production build
npm run build

# Production server
npm start

# Type checking
npx tsc --noEmit

# Linting
npm run lint

# Health check
npm run health
```

## Environment Variables

Create `.env.local` based on `.env.example`:

```env
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:3000/api

# Trust Union Network
TRUST_UNION_NODE_ID=
TRUST_UNION_NETWORK=production

# Database (future)
DATABASE_URL=

# Authentication (future)
AUTH_SECRET=
```

## Import Patterns

### Component Imports
```typescript
import SovereignHeader from '@/components/SovereignHeader'
import GoldCard from '@/components/GoldCard'
import QRWallet from '@/components/QRWallet'
```

### Library Imports
```typescript
import { generateGoldCardNumber, calculateTrustScore } from '@/lib/identity'
import { createValueProfile, addProofOfValue } from '@/lib/pov'
import { countryData, getCountriesByStatus } from '@/lib/countries'
import { validateEmail, validatePhone } from '@/lib/validators'
```

### Type Imports
```typescript
import { TrustProfile, ProofOfValue, ValueProfile } from '@/types/trust-union'
```

## Status

✅ **Architecture:** Fully restructured  
✅ **Components:** 4 root-level components created  
✅ **Libraries:** 5 consolidated modules  
✅ **Routes:** 11 pages + 4 API endpoints  
✅ **Countries:** 8 nation-specific pages  
✅ **Styles:** Modular with tokens  
✅ **Assets:** Organized structure  

**Ready for:** Development, testing, deployment

---

*For detailed migration notes, see [RESTRUCTURE.md](./RESTRUCTURE.md)*  
*For architecture overview, see [ARCHITECTURE.md](./ARCHITECTURE.md)*  
*For deployment instructions, see [DEPLOYMENT.md](./DEPLOYMENT.md)*
