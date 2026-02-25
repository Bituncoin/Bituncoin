# Bituncoin Universal Wallet - Multi-Platform Deployment Guide

## Overview

The Bituncoin Universal Wallet supports deployment across 7 platforms with automated CI/CD pipelines. This guide covers building, testing, and deploying the wallet on all supported platforms.

## Prerequisites

### Development Environment
- **Node.js**: 18.x or 20.x (LTS recommended)
- **Go**: 1.21 or later
- **Git**: 2.30 or later
- **Docker**: 20.10 or later (for containerized builds)

### Platform-Specific Requirements

#### Windows
- Windows 10/11
- Visual Studio Build Tools 2022
- Windows SDK 10.0.19041.0 or later

#### macOS
- macOS 11.0 or later
- Xcode 13.0 or later
- Command Line Tools for Xcode

#### Linux
- Ubuntu 18.04 or later / CentOS 7 or later
- GCC 7.0 or later
- GTK development libraries

#### Android
- Android Studio 2021.3.1 or later
- Android SDK API level 21 or later
- JDK 11 or later

#### iOS
- macOS with Xcode 13.0 or later
- iOS Simulator
- Apple Developer Account (for App Store deployment)

### Build Tools
```bash
# Install global dependencies
npm install -g electron-builder @electron-forge/cli react-native-cli

# Install Go dependencies
go install github.com/electron-userland/electron-builder@latest
```

## CI/CD Pipeline

### GitHub Actions Setup

The wallet uses GitHub Actions for automated multi-platform builds:

```yaml
# .github/workflows/build.yml
name: Multi-Platform Build

on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:

jobs:
  build-web:
    runs-on: ubuntu-latest
    # ... build configuration

  build-windows:
    runs-on: windows-latest
    # ... build configuration

  # ... other platform builds
```

### Build Triggers
- **Version Tags**: `git tag v1.0.0 && git push --tags`
- **Manual Dispatch**: GitHub UI or API call
- **Scheduled Builds**: Weekly on Sundays for testing

## Web Application

### Development Build
```bash
cd Bituncoin/wallet
npm install
npm start
# Opens http://localhost:3000
```

### Production Build
```bash
cd Bituncoin/wallet
npm run build
# Output: build/ directory
```

### Deployment Options

#### Static Hosting
```bash
# Deploy to Netlify
npm install -g netlify-cli
netlify deploy --prod --dir=build

# Deploy to Vercel
npm install -g vercel
vercel --prod

# Deploy to AWS S3 + CloudFront
aws s3 sync build/ s3://your-bucket-name
aws cloudfront create-invalidation --distribution-id YOUR_DISTRIBUTION_ID --paths "/*"
```

#### Docker Container
```dockerfile
FROM nginx:alpine
COPY build/ /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

```bash
docker build -t bituncoin-wallet-web .
docker run -p 8080:80 bituncoin-wallet-web
```

## Desktop Applications (Electron)

### Development
```bash
cd Bituncoin/wallet
npm install
npm run electron-dev
```

### Production Builds

#### Windows
```bash
npm run build:win
# Output: dist/Bituncoin Wallet 1.0.0.exe
# Output: dist/Bituncoin Wallet 1.0.0.exe.blockmap
```

Windows build artifacts:
- `Bituncoin Wallet 1.0.0.exe` - Main installer
- `Bituncoin Wallet 1.0.0.exe.blockmap` - Update manifest
- `win-unpacked/` - Portable version

#### macOS Intel
```bash
npm run build:mac-intel
# Output: dist/Bituncoin Wallet-1.0.0.dmg
# Output: dist/Bituncoin Wallet-1.0.0-mac.zip
```

macOS build artifacts:
- `Bituncoin Wallet-1.0.0.dmg` - Disk image installer
- `Bituncoin Wallet-1.0.0-mac.zip` - Archive for manual installation

#### macOS ARM (Apple Silicon)
```bash
npm run build:mac-arm
# Output: dist/Bituncoin Wallet-1.0.0-arm64.dmg
```

#### Linux
```bash
npm run build:linux
# Output: dist/Bituncoin Wallet-1.0.0.AppImage
# Output: dist/bituncoin-wallet_1.0.0_amd64.deb
# Output: dist/bituncoin-wallet-1.0.0.x86_64.rpm
```

Linux build artifacts:
- `Bituncoin Wallet-1.0.0.AppImage` - Universal Linux app
- `bituncoin-wallet_1.0.0_amd64.deb` - Debian package
- `bituncoin-wallet-1.0.0.x86_64.rpm` - RPM package

### Code Signing

#### Windows Code Signing
```bash
# Install certificate
certutil -addstore -f "TRUSTEDPUBLISHER" your-cert.p12

# Build with signing
electron-builder --win --cert your-cert.p12 --certPassword your-password
```

#### macOS Code Signing
```bash
# Install certificate
security import your-cert.p12 -k ~/Library/Keychains/login.keychain

# Build with signing
electron-builder --mac --cert your-cert.p12 --certPassword your-password
```

### Auto-Updates

The wallet includes auto-update functionality:

```javascript
const { autoUpdater } = require('electron-updater');

autoUpdater.checkForUpdatesAndNotify();

// Check manually
autoUpdater.checkForUpdates();

// Handle updates
autoUpdater.on('update-available', () => {
  // Show update dialog
});

autoUpdater.on('update-downloaded', () => {
  autoUpdater.quitAndInstall();
});
```

## Mobile Applications (React Native)

### Development Setup

#### Android
```bash
# Install Android SDK
# Set ANDROID_HOME environment variable

# Start Metro bundler
npx react-native start

# Run on Android emulator
npx react-native run-android
```

#### iOS
```bash
# Install CocoaPods
cd ios && pod install

# Start Metro bundler
npx react-native start

# Run on iOS simulator
npx react-native run-ios
```

### Production Builds

#### Android APK
```bash
# Debug build
npx react-native run-android --variant=debug

# Release build
cd android
./gradlew assembleRelease
# Output: android/app/build/outputs/apk/release/app-release.apk
```

#### Android AAB (Google Play)
```bash
cd android
./gradlew bundleRelease
# Output: android/app/build/outputs/bundle/release/app-release.aab
```

#### iOS Archive
```bash
# Archive for TestFlight/App Store
cd ios
xcodebuild -workspace BituncoinWallet.xcworkspace -scheme BituncoinWallet -configuration Release -archivePath build/BituncoinWallet.xcarchive archive

# Export IPA
xcodebuild -exportArchive -archivePath build/BituncoinWallet.xcarchive -exportPath build -exportOptionsPlist exportOptions.plist
```

### App Store Deployment

#### Google Play Store
1. **Prepare Release**:
   ```bash
   # Generate signed AAB
   cd android
   ./gradlew bundleRelease
   ```

2. **Upload to Play Console**:
   - Go to Google Play Console
   - Create new release
   - Upload AAB file
   - Fill release notes
   - Submit for review

#### Apple App Store
1. **Prepare Archive**:
   ```bash
   # Create archive
   xcodebuild -workspace BituncoinWallet.xcworkspace -scheme BituncoinWallet -configuration Release -archivePath build/BituncoinWallet.xcarchive archive
   ```

2. **Upload to App Store Connect**:
   ```bash
   # Using Xcode
   # Xcode -> Product -> Archive
   # Distribute App -> App Store Connect

   # Using Transporter
   xcrun altool --upload-app --type ios --file "BituncoinWallet.ipa" --username "your-apple-id" --password "app-specific-password"
   ```

3. **Submit for Review**:
   - Go to App Store Connect
   - Select build
   - Fill app information
   - Submit for review

## Testing

### Automated Testing
```bash
# Run all tests
npm run test

# Run with coverage
npm run test:coverage

# Run E2E tests
npm run test:e2e
```

### Platform-Specific Testing

#### Desktop Testing
```bash
# Test Electron app
npm run test:electron

# Test auto-updates
npm run test:updates
```

#### Mobile Testing
```bash
# Android instrumentation tests
cd android && ./gradlew test

# iOS unit tests
xcodebuild test -workspace BituncoinWallet.xcworkspace -scheme BituncoinWalletTests

# E2E tests
npx detox test
```

### Performance Testing
```bash
# Lighthouse for web
npx lighthouse http://localhost:3000

# Electron performance
npm run test:performance

# Mobile performance
npx react-native-performance-monitor
```

## Distribution

### Desktop Distribution

#### Windows
- **Microsoft Store**: Submit MSIX package
- **Direct Download**: Host on GitHub Releases
- **Enterprise**: Use SCCM or Group Policy

#### macOS
- **Mac App Store**: Submit through Xcode
- **Direct Download**: Host DMG on website
- **Enterprise**: Use MDM solutions

#### Linux
- **Snap Store**: `snapcraft upload bituncoin-wallet.snap`
- **Flathub**: Submit Flatpak
- **Direct Download**: Host AppImage

### Mobile Distribution

#### Android
- **Google Play Store**: Main distribution
- **Alternative Stores**: Amazon Appstore, Huawei AppGallery
- **Enterprise**: APK sideloading

#### iOS
- **App Store**: Primary distribution
- **TestFlight**: Beta testing
- **Enterprise**: In-house distribution

## Monitoring and Analytics

### Crash Reporting
```javascript
// Sentry integration
import * as Sentry from '@sentry/electron';

Sentry.init({
  dsn: 'your-dsn',
  environment: process.env.NODE_ENV,
});
```

### Usage Analytics
```javascript
// Mixpanel integration
import mixpanel from 'mixpanel-browser';

mixpanel.init('your-token');
mixpanel.track('app_launched');
```

### Performance Monitoring
```javascript
// Application Insights
import { ApplicationInsights } from '@microsoft/applicationinsights-web';

const appInsights = new ApplicationInsights({
  config: {
    instrumentationKey: 'your-key'
  }
});
appInsights.loadAppInsights();
```

## Security Considerations

### Code Signing
- All desktop builds are code-signed
- Mobile apps use proper certificates
- Web app served over HTTPS

### Binary Integrity
- SHA256 checksums provided for all downloads
- Automatic integrity verification
- Secure update channels

### Runtime Security
- CSP (Content Security Policy) headers
- Sandboxed execution environments
- Regular security audits

## Troubleshooting

### Common Build Issues

#### Electron Build Failures
```bash
# Clear cache
npx electron-builder install-app-deps --force

# Rebuild native modules
npm rebuild

# Check Node.js version compatibility
node --version
npm --version
```

#### React Native Issues
```bash
# Clear Metro cache
npx react-native start --reset-cache

# Clean Android build
cd android && ./gradlew clean

# Reset iOS dependencies
cd ios && rm -rf Pods && pod install
```

#### Code Signing Problems
```bash
# Verify certificate
security find-identity -v

# Check provisioning profile
xcodebuild -showBuildSettings
```

### Performance Optimization

#### Bundle Size
```bash
# Analyze bundle
npx webpack-bundle-analyzer build/static/js/*.js

# Optimize imports
# Use dynamic imports for large components
const Component = lazy(() => import('./Component'));
```

#### Startup Time
```javascript
// Preload critical resources
<link rel="preload" href="critical.css" as="style">
<link rel="modulepreload" href="critical.js">
```

## Version Information

- **Wallet Version**: 1.0.0
- **Electron Version**: 25.0.0
- **React Native Version**: 0.72.0
- **Node.js Requirement**: 18.x+
- **Go Requirement**: 1.21+

## Support

- **Documentation**: https://docs.bituncoin.com/deployment
- **Community**: https://forum.bituncoin.com/c/deployment
- **Issues**: https://github.com/bituncoin/wallet/issues
- **Security**: security@bituncoin.com

For additional support, please contact the Bituncoin development team.
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
