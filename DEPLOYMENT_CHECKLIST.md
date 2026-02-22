# 🚀 BTNG Production Deployment Checklist

**Version:** 1.0.0
**Last Updated:** February 22, 2026
**Platform:** BTNG Sovereign Gold Standard

---

## 📋 Pre-Deployment Verification

### ✅ Environment Setup
- [ ] `.env.production` file created with production secrets
- [ ] `BTNG_RPC_URL` configured for mainnet/testnet
- [ ] `BTNG_TOKEN_ADDRESS` set to deployed contract address
- [ ] `JWT_SECRET` and `API_SECRET` are strong, random strings
- [ ] Database connection strings configured (MongoDB/PostgreSQL)

### ✅ Security Validation
- [ ] Run `npm run verify:all` - all checks pass
- [ ] JWT tokens expire appropriately (24h for access, 7d for refresh)
- [ ] HTTPS certificates configured (Let's Encrypt recommended)
- [ ] CORS policy restricts to allowed domains only
- [ ] Rate limiting implemented on API endpoints

### ✅ Smart Contract Verification
- [ ] Contracts deployed to Sepolia/Mainnet
- [ ] Source code verified on Etherscan
- [ ] Contract addresses documented
- [ ] Multi-sig wallet configured for upgrades
- [ ] Emergency pause functionality tested

---

## 🏗️ Build & Bundle

### Frontend (Next.js)
```bash
# Build static frontend
npm run build

# Output: .next/ directory with optimized bundles
# - Removes dev code and console.logs
# - Minifies JavaScript and CSS
# - Optimizes images and fonts
```

### Backend (Genesis)
```bash
# Build Genesis application
npm run build:backend

# Output: genesis-app/build/ with JAR files
# - Compiled Kotlin code
# - Optimized for production
# - Includes all dependencies
```

### Docker (Optional)
```dockerfile
# Multi-stage Dockerfile for production
FROM node:18-alpine AS frontend
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
RUN npm run build

FROM openjdk:11-jre-slim AS backend
# Genesis deployment

FROM nginx:alpine AS production
COPY --from=frontend /app/.next /usr/share/nginx/html
COPY --from=backend /app/genesis-app/build /app/genesis
# Nginx config for API proxying
```

---

## 🚀 Deployment Steps

### 1. Server Preparation
```bash
# Update system packages
sudo apt update && sudo apt upgrade

# Install Node.js 18+
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Install PM2 for process management
sudo npm install -g pm2

# Install Nginx for reverse proxy
sudo apt install nginx
```

### 2. Application Deployment
```bash
# Clone repository
git clone https://github.com/your-org/btng-sovereign-platform.git
cd btng-sovereign-platform

# Install dependencies
npm ci --only=production

# Copy environment file
cp .env.production .env

# Build application
npm run build
npm run build:backend
```

### 3. Process Management
```bash
# Start with PM2
pm2 start ecosystem.config.js

# Or start manually
pm2 start npm --name "btng-frontend" -- run start
pm2 start npm --name "btng-genesis" -- run start:backend

# Save PM2 configuration
pm2 save
pm2 startup
```

### 4. Reverse Proxy (Nginx)
```nginx
# /etc/nginx/sites-available/btng
server {
    listen 80;
    server_name your-domain.com;

    # Frontend (Next.js)
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # API routes
    location /api/ {
        proxy_pass http://localhost:3003;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
    }

    # Genesis backend
    location /genesis/ {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
    }
}
```

### 5. SSL Certificate (Let's Encrypt)
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal (runs twice daily)
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

---

## 📊 Health Checks & Monitoring

### Application Health
```bash
# Health endpoint
curl https://your-domain.com/api/btng/health

# Gold price status
curl https://your-domain.com/api/btng/gold/price/status

# Genesis bridge status
curl https://your-domain.com/genesis/health
```

### Process Monitoring
```bash
# PM2 status
pm2 status

# Logs
pm2 logs btng-frontend
pm2 logs btng-genesis

# Restart if needed
pm2 restart btng-frontend
```

### System Monitoring
```bash
# Disk usage
df -h

# Memory usage
free -h

# CPU usage
top -n 1

# Network connections
netstat -tlnp | grep :3000
```

---

## 🔧 Post-Deployment Validation

### Functional Tests
```bash
# Run integration tests against production
npm run test:integration

# Test gold price API
curl -H "Authorization: Bearer YOUR_JWT" \
  https://your-domain.com/api/btng/gold/prices
```

### Performance Tests
```bash
# Load testing with Artillery
npm install -g artillery
artillery quick --count 50 --num 10 https://your-domain.com/api/btng/health
```

### Security Audit
```bash
# SSL certificate check
openssl s_client -connect your-domain.com:443 -servername your-domain.com

# Headers security check
curl -I https://your-domain.com
# Should include: Strict-Transport-Security, X-Frame-Options, etc.
```

---

## 🚨 Emergency Procedures

### Rollback Plan
```bash
# Quick rollback to previous version
pm2 stop all
git checkout previous-tag
npm ci
npm run build
pm2 start ecosystem.config.js
```

### Incident Response
1. **Check logs**: `pm2 logs --lines 100`
2. **Monitor resources**: `htop` or `top`
3. **Database connectivity**: Test MongoDB/PostgreSQL connections
4. **External dependencies**: Check RPC endpoints and API keys
5. **User impact**: Monitor error rates and response times

### Contact Information
- **DevOps Team**: devops@btng.org
- **Security Team**: security@btng.org
- **Emergency Hotline**: +254-XXX-XXXX

---

## 📈 Scaling & Optimization

### Horizontal Scaling
```bash
# Load balancer configuration
upstream btng_backend {
    server 127.0.0.1:3003;
    server 127.0.0.1:3004;
    server 127.0.0.1:3005;
}

# Database connection pooling
# Redis for session storage
# CDN for static assets
```

### Performance Optimization
- [ ] Enable gzip compression in Nginx
- [ ] Set up Redis caching for API responses
- [ ] Configure database indexes
- [ ] Implement API rate limiting
- [ ] Set up monitoring dashboards (Grafana + Prometheus)

---

## 🔒 Security Hardening

### Server Security
```bash
# Firewall configuration
sudo ufw enable
sudo ufw allow ssh
sudo ufw allow 'Nginx Full'

# SSH hardening
sudo sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sudo systemctl restart sshd

# Automatic updates
sudo apt install unattended-upgrades
```

### Application Security
- [ ] Regular dependency updates (`npm audit fix`)
- [ ] Log monitoring and alerting
- [ ] Backup strategy (database + files)
- [ ] Intrusion detection (fail2ban)
- [ ] Regular security audits

---

## 📚 Additional Resources

| Resource | Purpose | Link |
|----------|---------|------|
| **PM2 Documentation** | Process management | https://pm2.keymetrics.io/ |
| **Nginx Config Guide** | Reverse proxy setup | https://nginx.org/en/docs/ |
| **Let's Encrypt** | SSL certificates | https://certbot.eff.org/ |
| **Node.js Production** | Best practices | https://nodejs.org/en/docs/guides/ |
| **Docker Best Practices** | Container deployment | https://docs.docker.com/develop/ |

---

**🇰🇪 BTNG Sovereign Gold Standard - Production Ready 🇰🇪**

*This checklist ensures your BTNG platform maintains the highest standards of security, reliability, and sovereignty in production.*