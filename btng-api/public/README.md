# BTNG Public Assets

This directory contains static assets for the BTNG platform.

## Directory Structure

### `/card/`
Gold Card design assets:
- Card backgrounds
- Chip images
- Security features overlays
- Preview templates

### `/icons/`
Platform iconography:
- Navigation icons
- Feature icons
- Status indicators
- Country flags (high resolution)
- Mobile money provider logos

### Root Files
- `emblem.svg` - Official BTNG platform emblem
- `favicon.ico` - Browser favicon
- `og-image.png` - Open Graph social sharing image

## Asset Guidelines

### Gold Card Assets
- Use official BTNG gold color: `#D4AF37`
- Maintain 400x250px aspect ratio for card designs
- Include security watermarks where appropriate

### Icons
- SVG format preferred for scalability
- 24x24px default size
- Use platform colors from design tokens
- Ensure accessibility with proper contrast ratios

### Branding
- Logo files should maintain sovereign black and gold color scheme
- Trust Union blue (`#1A4B8C`) for trust-related elements
- Always include proper attribution and licensing

## Usage in Code

```tsx
// Import from public directory
import Image from 'next/image'

<Image 
  src="/emblem.svg" 
  alt="BTNG Emblem" 
  width={200} 
  height={200} 
/>
```

## Network Distribution
Assets in this directory are served statically and cached at the CDN level for optimal performance across the Trust Union network.
