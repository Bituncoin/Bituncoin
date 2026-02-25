# BTNG Sovereign Live Test Runbook

## 🎯 Objective

Execute a controlled end-to-end test of the complete BTNG sovereign flow using real money (0.10 GHS) to validate:

- MoMo payment processing
- Token minting on Fabric
- Watchtower snapshot generation
- Oracle price feed integration

## 📋 Prerequisites

- [ ] Production environment configured
- [ ] `BTNG_ENABLE_REAL_MONEY=true` set
- [ ] Canary test account with 0.10 GHS balance
- [ ] All services running and healthy
- [ ] Monitoring dashboards active

## 🚀 Execution Steps

### Step 1: Pre-Test Validation (5 minutes)

```bash
# Verify environment
echo $BTNG_ENABLE_REAL_MONEY  # Should be "true" ONLY in production

# IMPORTANT: Never set this in development environments
# This flag enables real money movement - use extreme caution

# Check service health
curl -s http://localhost:3000/api/health | jq .status  # Should be "operational"

# Confirm fabric network
npm run test-fabric

# Verify MoMo credentials
node scripts/momo_live.js --help  # Should show all options
```

### Security Gating for Real-Money Flag

**CRITICAL:** The `BTNG_ENABLE_REAL_MONEY` environment variable should only be set in production environments with:

- Multi-approval process
- Audit logging
- Restricted access
- Automatic rollback capability

**Recommended Implementation:**

```bash
# In production deployment script
if [ "$ENVIRONMENT" = "production" ] && [ "$APPROVED_BY_CEO" = "true" ]; then
  export BTNG_ENABLE_REAL_MONEY=true
  echo "Real-money mode enabled - $(date)" >> audit.log
else
  export BTNG_ENABLE_REAL_MONEY=false
  echo "Real-money mode blocked - $(date)" >> audit.log
fi
```

### Step 2: Execute Penny Transaction (2 minutes)

```bash
# Run the live transaction
node scripts/momo_live.js \
  --amount 0.10 \
  --currency GHS \
  --msisdn +233XXXXXXXXX \  # Replace with canary number
  --payerMessage "BTNG sovereign live test" \
  --executeLive \
  --confirm "I UNDERSTAND REAL MONEY WILL MOVE"
```

**Expected Output:**

```
MOMO_LIVE_EXECUTION=SUCCESS
Transaction ID: momo-tx-XXXXXXXX
Status: completed
```

### Step 3: Validate End-to-End Flow (10 minutes)

#### 3.1 Check Payment Processing

```bash
# Verify MoMo callback received
tail -f logs/momo-callback.log | grep "btng-live-"

# Check payment record in database
mongosh btng --eval "db.payments.find({externalId:/btng-live-/}).sort({timestamp:-1}).limit(1)"
```

#### 3.2 Verify Token Minting

```bash
# Check Fabric chaincode invocation
curl -X POST http://localhost:3000/api/btng/fabric/chaincode \
  -H "Content-Type: application/json" \
  -d '{
    "chaincode": "btng-gold-token",
    "function": "BalanceOf",
    "args": ["BTNGXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"]
  }'
```

#### 3.3 Confirm Watchtower Update

```bash
# Check latest snapshot
curl -s http://localhost:3000/api/watchtower/nodes | jq .snapshot

# Verify sequence incremented
curl -s http://localhost:3000/api/watchtower/nodes | jq .snapshot.sequence
```

#### 3.4 Validate Oracle Integration

```bash
# Check gold price update
curl -s http://localhost:3000/api/btng/gold/price/latest | jq .

# Verify price timestamp recent
curl -s http://localhost:3000/api/btng/gold/price/latest | jq .timestamp
```

## ✅ Success Criteria

### Payment Success

- [ ] MoMo transaction status: "completed"
- [ ] Payment record created in database
- [ ] Callback processed without errors

### Token Minting

- [ ] Fabric transaction successful
- [ ] Token balance increased by 0.10 GHS worth
- [ ] Transaction ID recorded

### Watchtower Integrity

- [ ] New snapshot generated
- [ ] Sequence number incremented
- [ ] Hash continuity maintained

### Oracle Updates

- [ ] Gold price feed current
- [ ] Price data validated
- [ ] Timestamp within last 5 minutes

## ❌ Failure Scenarios & Recovery

### Payment Failure

```bash
# Check MoMo error codes
tail -f logs/momo-errors.log

# Retry with different amount/currency
node scripts/momo_live.js --amount 0.05 --currency GHS ...
```

### Minting Failure

```bash
# Check Fabric logs
docker logs btng-fabric-peer

# Manual minting fallback
curl -X POST http://localhost:3000/api/btng/fabric/chaincode \
  -H "Content-Type: application/json" \
  -d '{
    "chaincode": "btng-gold-token",
    "function": "Mint",
    "args": ["10", "BTNGXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"]
  }'
```

### Watchtower Issues

```bash
# Force snapshot regeneration
curl -X POST http://localhost:3000/api/watchtower/nodes/refresh

# Verify cryptographic integrity
npm run verify:watchtower-crypto
```

## 📊 Test Results Documentation

### Transaction Details

- **Transaction ID:** ____________________
- **Amount:** ____________________
- **Timestamp:** ____________________
- **Status:** ____________________

### Validation Results

- [ ] Payment Processing: PASS / FAIL
- [ ] Token Minting: PASS / FAIL
- [ ] Watchtower Update: PASS / FAIL
- [ ] Oracle Integration: PASS / FAIL

### Notes

```
[Record any issues, observations, or deviations from expected behavior]
```

## 🎯 Go/No-Go Decision

**Test Status:** ⏳ PENDING

**Recommendation:** ____________________

**Approved By:** ____________________
**Date/Time:** ____________________

---

**Test Conductor:** [Your Name]
**Canary Account:** [Phone Number]
**Environment:** Production</content>
<parameter name="filePath">c:\BTNGAI_files\LIVE_TEST_RUNBOOK.md
