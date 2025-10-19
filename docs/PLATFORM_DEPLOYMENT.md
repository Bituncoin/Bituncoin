# Multi-Platform Deployment Guide

## Overview
This guide provides instructions for building and deploying the Bituncoin wallet across all supported platforms.

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Web Application](#web-application)
3. [Desktop Applications](#desktop-applications)
4. [Mobile Applications](#mobile-applications)
5. [Backend Services](#backend-services)
6. [CI/CD Automation](#cicd-automation)

## Prerequisites

### Development Tools
- **Go**: 1.18 or later
- **Node.js**: 16 or later
- **npm** or **yarn**: Latest version
- **Git**: For version control

### Platform-Specific Requirements

#### macOS
- Xcode Command Line Tools
- CocoaPods (for iOS development)

#### Windows
- Visual Studio Build Tools
- Windows SDK

#### Linux
- Build essentials (`build-essential` on Debian/Ubuntu)
- Development libraries

## Web Application

### Development

```bash
cd wallet
npm install
npm start
```

The wallet will be available at `http://localhost:3000`

### Production Build

```bash
cd wallet
npm run build
```

This creates an optimized production build in `wallet/build/`

### Deployment

#### Option 1: Static Hosting (Netlify, Vercel, GitHub Pages)

```bash
# Build the application
cd wallet
npm run build

# Deploy to your hosting provider
# For Netlify:
netlify deploy --prod --dir=build

# For Vercel:
vercel --prod

# For GitHub Pages:
npm install -g gh-pages
gh-pages -d build
```

#### Option 2: Docker Container

Create `Dockerfile` in the wallet directory:

```dockerfile
FROM node:16-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

Build and run:

```bash
docker build -t bituncoin-wallet .
docker run -p 80:80 bituncoin-wallet
```

## Desktop Applications

### Electron Setup

First, install Electron and builder:

```bash
cd wallet
npm install --save-dev electron electron-builder
```

Create `wallet/electron.js`:

```javascript
const { app, BrowserWindow } = require('electron');
const path = require('path');

function createWindow() {
  const win = new BrowserWindow({
    width: 1200,
    height: 800,
    webPreferences: {
      nodeIntegration: false,
      contextIsolation: true,
      preload: path.join(__dirname, 'preload.js')
    }
  });

  // In production, load the built files
  if (app.isPackaged) {
    win.loadFile(path.join(__dirname, 'build/index.html'));
  } else {
    // In development, load from dev server
    win.loadURL('http://localhost:3000');
  }
}

app.whenReady().then(createWindow);

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit();
  }
});

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow();
  }
});
```

Update `wallet/package.json`:

```json
{
  "main": "electron.js",
  "scripts": {
    "electron-dev": "ELECTRON_START_URL=http://localhost:3000 electron .",
    "electron-build": "npm run build && electron-builder"
  },
  "build": {
    "appId": "com.bituncoin.wallet",
    "productName": "Bituncoin Wallet",
    "files": [
      "build/**/*",
      "electron.js"
    ],
    "directories": {
      "buildResources": "assets"
    },
    "mac": {
      "category": "public.app-category.finance",
      "target": ["dmg", "zip"]
    },
    "win": {
      "target": ["nsis", "portable"]
    },
    "linux": {
      "target": ["AppImage", "deb", "rpm"],
      "category": "Finance"
    }
  }
}
```

### Building for Windows

```bash
cd wallet
npm run build
npm run electron-build -- --win
```

Output: `wallet/dist/*.exe`

### Building for macOS

```bash
cd wallet
npm run build
npm run electron-build -- --mac
```

Output: `wallet/dist/*.dmg`

### Building for Linux

```bash
cd wallet
npm run build
npm run electron-build -- --linux
```

Output: `wallet/dist/*.AppImage`, `*.deb`, `*.rpm`

## Mobile Applications

### React Native Setup

Convert the React web app to React Native for mobile platforms.

#### iOS Development

Prerequisites:
- macOS computer
- Xcode 12+
- CocoaPods

```bash
# Initialize React Native project
npx react-native init BituncoinWallet
cd BituncoinWallet

# Install dependencies
npm install

# Install iOS dependencies
cd ios
pod install
cd ..

# Run on iOS simulator
npx react-native run-ios
```

#### Android Development

Prerequisites:
- Android Studio
- Android SDK
- Java Development Kit (JDK)

```bash
# Run on Android emulator or device
npx react-native run-android
```

#### Building for Production

**iOS:**

```bash
# Open in Xcode
open ios/BituncoinWallet.xcworkspace

# In Xcode:
# 1. Select "Any iOS Device" as target
# 2. Product > Archive
# 3. Distribute App > App Store Connect
```

**Android:**

```bash
cd android
./gradlew assembleRelease

# Output: android/app/build/outputs/apk/release/app-release.apk
```

## Backend Services

### API Node

#### Development

```bash
cd /path/to/Bituncoin
go run api/btnnode.go
```

#### Production Build

Build for current platform:

```bash
go build -o btnnode api/btnnode.go
```

Build for specific platforms:

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o btnnode-linux-amd64 api/btnnode.go

# Windows
GOOS=windows GOARCH=amd64 go build -o btnnode-windows-amd64.exe api/btnnode.go

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o btnnode-darwin-amd64 api/btnnode.go

# macOS Apple Silicon
GOOS=darwin GOARCH=arm64 go build -o btnnode-darwin-arm64 api/btnnode.go

# Linux ARM (Raspberry Pi)
GOOS=linux GOARCH=arm64 go build -o btnnode-linux-arm64 api/btnnode.go
```

#### Docker Deployment

Create `Dockerfile`:

```dockerfile
FROM golang:1.18-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o btnnode api/btnnode.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/btnnode .
COPY config.yml .

EXPOSE 8080
CMD ["./btnnode"]
```

Build and run:

```bash
docker build -t bituncoin-node .
docker run -p 8080:8080 bituncoin-node
```

#### Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./config.yml:/root/config.yml
      - ./data:/root/data
    restart: unless-stopped

  wallet:
    build: ./wallet
    ports:
      - "80:80"
    depends_on:
      - api
    restart: unless-stopped
```

Run:

```bash
docker-compose up -d
```

## CI/CD Automation

### GitHub Actions

The repository includes automated CI/CD workflows:

#### Test Workflow (`.github/workflows/test.yml`)
- Runs on push to main/develop
- Tests Go backend
- Tests wallet frontend
- Runs linting

#### Build Workflow (`.github/workflows/build.yml`)
- Triggered on version tags (v*)
- Builds for all platforms
- Creates GitHub releases
- Uploads artifacts

### Manual Release

```bash
# Tag a new version
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0

# GitHub Actions will automatically:
# 1. Build for all platforms
# 2. Create a release
# 3. Upload binaries
```

### Platform-Specific Notes

#### Windows

Installers support:
- Silent installation: `/S` flag
- Custom install directory: `/D=C:\CustomPath`

#### macOS

DMG files include:
- Application bundle
- Drag-to-Applications folder
- Code signing (requires Apple Developer account)

#### Linux

AppImage:
- Portable, no installation required
- `chmod +x` and run
- Works on most distributions

DEB packages:
- For Debian/Ubuntu
- `sudo dpkg -i bituncoin-wallet.deb`

RPM packages:
- For Fedora/RedHat/CentOS
- `sudo rpm -i bituncoin-wallet.rpm`

## Environment Configuration

### Production Environment Variables

Create `.env` file:

```bash
# API Configuration
API_HOST=0.0.0.0
API_PORT=8080
API_NETWORK=mainnet

# Database
DB_PATH=/var/lib/bituncoin/data

# Security
SESSION_SECRET=your-secret-key-here
ENCRYPTION_KEY=your-encryption-key-here

# External Services
MARKET_DATA_API_KEY=your-api-key
```

### Configuration File

Edit `config.yml`:

```yaml
network:
  name: "bituncoin-mainnet"
  type: "mainnet"
  
api:
  host: "0.0.0.0"
  port: 8080
  cors_origins:
    - "https://wallet.bituncoin.com"
    
security:
  session_timeout: 86400  # 24 hours
  max_login_attempts: 5
  enable_2fa: true
  
blockchain:
  block_time: 10
  consensus: "pos"
```

## Monitoring and Maintenance

### Health Checks

```bash
# Check API health
curl http://localhost:8080/api/health

# Check node info
curl http://localhost:8080/api/info
```

### Logs

```bash
# Docker logs
docker logs bituncoin-node

# System logs (systemd)
journalctl -u bituncoin-node -f
```

### Updates

```bash
# Pull latest changes
git pull origin main

# Rebuild
go build -o btnnode api/btnnode.go

# Restart service
systemctl restart bituncoin-node
```

## Security Checklist

Before deploying to production:

- [ ] Enable HTTPS/SSL
- [ ] Configure firewall rules
- [ ] Set up regular backups
- [ ] Enable logging and monitoring
- [ ] Use strong session secrets
- [ ] Enable rate limiting
- [ ] Configure CORS properly
- [ ] Use environment variables for secrets
- [ ] Enable 2FA for admin accounts
- [ ] Set up DDoS protection
- [ ] Regular security audits
- [ ] Keep dependencies updated

## Support

For deployment assistance:
- GitHub Issues: https://github.com/Bituncoin/Bituncoin/issues
- Documentation: /docs directory
- Community: Discord channel

## Next Steps

After deployment:
1. Test all functionality
2. Set up monitoring
3. Configure backups
4. Plan for scaling
5. Document custom configurations
