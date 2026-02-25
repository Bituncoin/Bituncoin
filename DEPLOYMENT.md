# BTNG Platform - Deployment Guide

## 🚀 Quick Start

### Prerequisites
- Node.js 18+ installed
- npm or yarn package manager
- Git (optional)

### Installation

1. **Install Dependencies**
```bash
npm install
```

2. **Environment Setup**
```bash
# Copy the example environment file
copy .env.example .env

# Edit .env with your configuration (optional for development)
```

3. **Run Development Server**
```bash
npm run dev
```

Access the platform at: `http://localhost:3000`

### Build for Production

```bash
# Create optimized production build
npm run build

# Start production server
npm start
```

## 📦 Deployment Options

### Option 1: Vercel (Recommended for Quick Deploy)

1. **Install Vercel CLI**
```bash
npm install -g vercel
```

2. **Deploy**
```bash
vercel
```

Follow the prompts to complete deployment. Your BTNG platform will be live in minutes.

### Option 2: Self-Hosted (Docker)

Create `Dockerfile`:
```dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/public ./public

EXPOSE 3000
CMD ["npm", "start"]
```

Build and run:
```bash
docker build -t btng-platform .
docker run -p 3000:3000 btng-platform
```

### Option 3: Traditional VPS/Server

1. **SSH into your server**
2. **Install Node.js 18+**
3. **Clone/Upload project files**
4. **Install dependencies**: `npm install`
5. **Build**: `npm run build`
6. **Use PM2 for process management**:
```bash
npm install -g pm2
pm2 start npm --name "btng-platform" -- start
pm2 save
pm2 startup
```

## 🌐 Domain & SSL

### With Vercel
- Automatic SSL with Let's Encrypt
- Custom domain configuration via Vercel dashboard

### Self-Hosted with Nginx

Nginx configuration (`/etc/nginx/sites-available/btng`):
```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## 📊 Health Monitoring

### Built-in Health Check
```bash
npm run health
```

### API Health Endpoint
```
GET https://your-domain.com/api/health
```

Returns:
```json
{
  "status": "operational",
  "version": "0.1.0",
  "services": {
    "platform": { "status": "healthy" },
    "identity": { "status": "healthy" },
    "wallet": { "status": "healthy" },
    "trustUnion": { "status": "healthy" }
  }
}
```

### Monitoring Solutions
- **Uptime Monitoring**: UptimeRobot, Pingdom
- **Application Monitoring**: Sentry, LogRocket
- **Infrastructure**: DataDog, New Relic

## 🔒 Security Checklist

- [ ] Environment variables secured (never commit .env)
- [ ] HTTPS enabled with valid SSL certificate
- [ ] API rate limiting configured
- [ ] Security headers configured (already in next.config.js)
- [ ] Database credentials secured (when implemented)
- [ ] Regular dependency updates (`npm audit`)
- [ ] CORS policies configured for production

## 🔐 Environment Variables

### Required for Production
```bash
NODE_ENV=production
NEXT_PUBLIC_APP_URL=https://your-domain.com
NEXT_PUBLIC_API_URL=https://your-domain.com
```

### Optional (Future Enhancements)
```bash
DATABASE_URL=postgresql://...
REDIS_URL=redis://...
MOBILE_MONEY_API_KEY=...
```

## 📈 Performance Optimization

### Already Configured
- ✅ Next.js automatic code splitting
- ✅ Image optimization (when images added)
- ✅ Compression enabled
- ✅ React strict mode

### Additional Recommendations
- Enable CDN for static assets
- Configure caching headers
- Implement database connection pooling (when DB added)
- Use Redis for session management

## 🔄 Continuous Deployment

### GitHub Actions Example

Create `.github/workflows/deploy.yml`:
```yaml
name: Deploy BTNG Platform

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
      - run: npm ci
      - run: npm run build
      - run: npm test # when tests added
      - uses: amondnet/vercel-action@v20
        with:
          vercel-token: ${{ secrets.VERCEL_TOKEN }}
          vercel-org-id: ${{ secrets.ORG_ID }}
          vercel-project-id: ${{ secrets.PROJECT_ID }}
```

## 🐛 Troubleshooting

### Build Errors
```bash
# Clear Next.js cache
rm -rf .next

# Clear node modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Port Already in Use
```bash
# Kill process on port 3000 (Windows)
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Change port
$env:PORT=3001; npm run dev
```

### Module Not Found Errors
- Check `tsconfig.json` paths configuration
- Verify all imports use correct path aliases
- Rebuild project: `npm run build`

## 📞 Support & Resources

- **Documentation**: See README.md and ARCHITECTURE.md
- **Health Check**: `/api/health`
- **Version**: Check package.json or health endpoint

---

## 🌐 Network Anchors

- **Primary Backend Endpoint**: `http://74.118.126.72:64799`
- **Genesis Transaction Hash**: `0x1111111111111111111111111111111111111111111111111111111111111111`
- **Genesis Explorer URL**: `http://74.118.126.72:64799/explorer/tx/0x1111111111111111111111111111111111111111111111111111111111111111`
- **Genesis Block Height**: `12458`
- **Genesis Timestamp**: `1771457774`

---

**BTNG Platform** — Sovereign. Stable. Operational.
