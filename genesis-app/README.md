# BTNG Genesis Bridge Application

🇰🇪 **Genesis Platform Integration for BTNG Sovereign Gold Standard**

This Genesis application serves as a bridge between your Node.js BTNG API and the Genesis trading platform, enabling real-time gold price data streaming to financial applications.

## 🏗 Architecture Overview

```
BTNG Node.js API (Port 3003)
        ↓ (HTTP + JWT Auth)
    BTNG Price Bridge Service
        ↓ (Database)
    Genesis GOLD_PRICE_CACHE Table
        ↓ (Data Server Streaming)
    Frontend Components (Real-time Updates)
```

## 📁 Project Structure

```
genesis-app/
├── server/
│   ├── src/main/genesis/
│   │   ├── BTNGPriceBridge.kts    # Bridge service with JWT auth
│   │   └── DataServer.kts         # Real-time data streaming
│   └── src/main/kotlin/
│       └── global/genesis/
│           └── GoldPriceCache.kt  # Database table definition
└── client/                        # Frontend components (future)
```

## 🚀 Setup Instructions

### 1. Initialize Genesis Project

Use the [Genesis Create](https://genesis.global/docs/develop/development-environment/launchpad/genesis-create/) web tool to generate your project:

1. **Select Template**: Choose "Financial Services" or "Trading Platform"
2. **Configure Database**: Select PostgreSQL or FoundationDB
3. **Add Dependencies**: Include HTTP Client and JSON processing
4. **Download & Extract**: Place files in this `genesis-app/` directory

### 2. Configure Dependencies

Add to your `server/build.gradle.kts`:

```kotlin
dependencies {
    implementation("global.genesis:genesis-http-client")
    implementation("org.jetbrains.kotlinx:kotlinx-serialization-json")
}
```

### 3. Database Setup

The `GoldPriceCache` table will be automatically created when you run the Genesis server for the first time.

### 4. Start the Services

**Terminal 1 - BTNG API:**
```bash
cd C:\BTNGAI_files
npm run dev
```

**Terminal 2 - Genesis Server:**
```bash
cd C:\BTNGAI_files\genesis-app
./gradlew server
```

## 🔧 Bridge Service Details

### BTNG Price Poller
- **Schedule**: Every 60 seconds (`"0 * * * * ?"`)
- **Authentication**: Automatic JWT login to BTNG API
- **Data Flow**:
  1. Authenticate with BTNG API (`/api/auth/login`)
  2. Fetch price data (`/api/btng/gold/price/status`)
  3. Store in Genesis database (`GOLD_PRICE_CACHE`)
  4. Log success/failure

### Health Check
- **Schedule**: Every 5 minutes
- **Purpose**: Monitor BTNG API connectivity
- **Logs**: Authentication status and connection health

## 📊 Data Server Queries

### Available Streams

1. **`ALL_GOLD_PRICES`**: All price updates as they occur
2. **`LATEST_GOLD_PRICE`**: Most recent price data
3. **`GOLD_PRICE_HISTORY_24H`**: Price history for last 24 hours

### Frontend Integration

Subscribe to real-time updates in your Genesis client:

```javascript
// Subscribe to latest gold price
genesis.subscribe('LATEST_GOLD_PRICE', (data) => {
    console.log('Gold price updated:', data.PRICE);
    updateGoldPriceDisplay(data);
});
```

## 🔍 Monitoring & Debugging

### Check Bridge Logs
```bash
# In Genesis terminal
mon
```

Look for:
- `BTNG Bridge: Successfully synced gold price`
- `BTNG Bridge Health: ✅ Authentication successful`

### Inspect Database
```bash
# Use Genesis Db Browser or query directly
SELECT * FROM GOLD_PRICE_CACHE ORDER BY LAST_UPDATE DESC LIMIT 10;
```

### Test API Connectivity
```bash
# Test BTNG API directly
curl -s http://localhost:3003/api/btng/gold/price/status
```

## 🛠 Troubleshooting

### Bridge Not Polling
1. **Check BTNG API**: Ensure Node.js server is running on port 3003
2. **Verify JWT Auth**: Test login endpoint manually
3. **Check Genesis Logs**: Look for authentication errors

### Data Not Streaming
1. **Data Server**: Ensure DataServer.kts is loaded
2. **Client Subscription**: Verify frontend subscription code
3. **Database**: Check if records are being inserted

### Authentication Failures
1. **Credentials**: Verify admin username/password match BTNG API
2. **JWT Secret**: Ensure BTNG API JWT_SECRET is consistent
3. **Network**: Check localhost connectivity between services

## 🔐 Security Considerations

- **JWT Tokens**: Bridge service handles authentication automatically
- **Network Security**: Services communicate over localhost
- **Data Validation**: Price data is validated before storage
- **Error Handling**: Failed polls are logged but don't crash the service

## 📈 Performance Notes

- **Polling Frequency**: 1 minute intervals balance real-time needs with API limits
- **Data Retention**: Consider implementing data cleanup for old records
- **Connection Pooling**: Genesis HTTP client handles connection management
- **Error Recovery**: Service continues polling even after temporary failures

## 🎯 Next Steps

1. **Initialize Genesis Project**: Use Genesis Create tool
2. **Test Bridge Service**: Verify data polling and storage
3. **Implement Frontend**: Subscribe to data streams
4. **Add Monitoring**: Set up alerts for bridge failures
5. **Production Deployment**: Configure for production environment

---

**🇰🇪 BTNG Sovereign Gold Standard - Genesis Integration Complete**