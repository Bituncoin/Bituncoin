# 🚀 BTNG Immediate Live Operations Guide

## 🎯 MISSION: Activate Sovereign Gold Operations

**Goal**: Get BTNG backend live and connected to frontend demo within minutes

### Step 1: Start Your BTNG Backend Server
```bash
# On your BTNG server machine (connected to MTN router):
# Ensure BTNG node or API server is running
# Bound to 0.0.0.0:64799 (not just localhost)
# Firewall allows inbound TCP on port 64799
```

### Step 2: Verify Backend Health
```bash
# Test from any machine:
curl http://74.118.126.72:64799/health

# Expected response:
{
  "status": "operational",
  "version": "1.0.0",
  "network": "BTNG Sovereign"
}
```

### Step 3: Start Frontend Demo
```bash
# In your development environment:
npm run dev

# Visit: http://localhost:3000/btng-demo
```

### Step 4: Test Live Integration
- **Wallet Operations**: Balance, send, transactions
- **Explorer Functions**: Blocks, transactions, addresses
- **Mining Data**: Hashrate, difficulty
- **Oracle Feeds**: Gold price, market cap
- **PoV Verification**: Signature validation

## 🔧 Backend Startup Checklist

### For BTNG Node:
- [ ] BTNG daemon running
- [ ] RPC server enabled on port 64799
- [ ] Wallet functionality active
- [ ] Explorer API enabled
- [ ] Mining operations active
- [ ] Oracle feeds connected

### For API Server:
- [ ] BTNG API server running
- [ ] Bound to 0.0.0.0:64799
- [ ] Health endpoint responding
- [ ] All API endpoints functional
- [ ] Database connected
- [ ] External access confirmed

## 📊 Expected Live Operations

Once backend is live, the demo will show:
- **Real BTNG balances** (not simulated)
- **Live transaction processing**
- **Actual mining statistics**
- **Current gold price feeds**
- **Sovereign verification**

## 🚨 Quick Troubleshooting

### Backend Not Starting:
```bash
# Check BTNG logs
tail -f ~/.btng/debug.log

# Verify port binding
netstat -tlnp | grep 64799
```

### Health Endpoint Failing:
```bash
# Test locally first
curl http://localhost:64799/health

# Check firewall
sudo ufw status
sudo ufw allow 64799
```

### Network Issues:
- Confirm MTN public IP: 74.118.126.72
- Check router port forwarding
- Verify VPN connection

## 🎉 Success Indicators

- ✅ Health endpoint returns valid JSON
- ✅ Frontend demo shows real data
- ✅ Transactions process successfully
- ✅ Gold price updates live
- ✅ Mining stats display current data

## 📞 Support

If backend startup issues persist:
1. Check BTNG configuration files
2. Verify network settings
3. Review firewall rules
4. Confirm MTN router configuration

---
*Immediate live operations activation guide*</content>
<parameter name="filePath">c:\BTNGAI_files\IMMEDIATE_LIVE_OPERATIONS.md