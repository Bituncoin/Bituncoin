# 🧾 BTNG Merchant API Guide

## Overview

This guide helps merchant systems (POS, web checkout, mobile cashier apps) consume BTNG sovereign gold pricing and health signals for production use.

## Base Endpoints

- Gold API Base: `http://74.118.126.72:64799/api/btng/gold`
- Health Endpoint: `http://74.118.126.72:64799/api/health`

## Core Merchant Flows

### 1) Read Current Price

`GET /api/btng/gold/price/latest`

Use this before each checkout quote.

Example response fields:
- `base_price_gram`
- `currencies[]` (e.g., GHS, NGN)
- `timestamp`
- `fx_rates`

### 2) Validate Service Health

`GET /api/btng/gold/price/status`

Use this for pre-shift checks and uptime monitoring.

### 3) Fetch Audit History

`GET /api/btng/gold/price/history?limit=100&startTime=&endTime=`

Use this for reconciliation, accounting, and dispute resolution.

## Pricing Rules for POS

1. Pull latest price at cart open.
2. Store quote timestamp with order.
3. Mark quote stale if older than 5 minutes.
4. Refresh quote at payment confirmation.

## Recommended Merchant Data Model

```json
{
  "merchantOrderId": "POS-2026-000123",
  "quotedAt": 1771457774,
  "currency": "GHS",
  "goldPricePerGram": 1020.5,
  "grams": 0.25,
  "subtotal": 255.125,
  "source": "BTNG Sovereign Gold API",
  "priceTimestamp": 1739999999
}
```

## Error Handling

- If `/price/latest` fails: retry with exponential backoff.
- If price age > 5 minutes: show `Stale Price` warning and block checkout.
- If `/price/status` indicates degraded mode: switch merchant panel to warning state.

## Security Practices

- Use server-to-server API calls where possible.
- Do not trust client-side computed totals without server verification.
- Log quote inputs (`timestamp`, `currency`, `price_gram`) for audit trails.

## Ghana Pilot Defaults

- Default display currency: `GHS`
- Stale threshold: 5 minutes
- Volatility review trigger: >2% in one hour

## Quick Verification Commands

```bash
curl http://74.118.126.72:64799/api/btng/gold/price/latest
curl http://74.118.126.72:64799/api/btng/gold/price/status
curl "http://74.118.126.72:64799/api/btng/gold/price/history?limit=10"
```

## Go-Live Checklist

- [ ] Merchant POS fetches latest price successfully
- [ ] Stale-price warning appears when data is old
- [ ] Reconciliation report includes price timestamp
- [ ] Health endpoint integrated into merchant monitoring
- [ ] Ghana pilot currency defaults verified
