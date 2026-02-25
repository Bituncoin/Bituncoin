# BTNG Freedom Mainnet Runbook

## 1) Start Watchtower Suite

```powershell
npm run start:watchtower-suite
```

Use the printed URL (for example `http://localhost:3001` or fallback port).

## 2) Verify Crypto + Readiness (Before Any Real Money)

```powershell
npm run verify:watchtower-crypto
npm run verify:mainnet-readiness -- --watchtowerUrl http://localhost:3001 --minOnlineNodes 2
```

Pass criteria:
- signed snapshot present
- anchored snapshot chain present
- `onlineNodes >= 2`

## 3) Generate Second-Node Join Package

```powershell
npm run generate:second-node-package -- --country Nigeria --nodeId nigeria-peer-0 --peerPort 8051 --caPort 8054 --watchtowerUrl http://localhost:3001
```

Output file:
- `artifacts/join-packages/nigeria-join-package.json`

## 4) Register the Second Node (Signed)

Use your ES256 signing flow to submit a signed registration payload to:
- `POST /api/watchtower/nodes`

Then re-run:

```powershell
npm run verify:mainnet-readiness -- --watchtowerUrl http://localhost:3001 --minOnlineNodes 2
```

## 5) Penny Mint (Dry-Run First, Live Only with Gates)

Dry-run preview:

```powershell
npm run momo:live -- --amount 0.10 --currency GHS --msisdn 233XXXXXXXXX
```

Live execution requires all gates:
- `BTNG_ENABLE_REAL_MONEY=true`
- valid `BTNG_MOMO_API_BASE_URL`
- valid `BTNG_MOMO_API_TOKEN`
- explicit `--executeLive`
- exact `--confirm "I UNDERSTAND REAL MONEY WILL MOVE"`

Example:

```powershell
$env:BTNG_ENABLE_REAL_MONEY='true'
$env:BTNG_MOMO_API_BASE_URL='https://<your-momo-endpoint>'
$env:BTNG_MOMO_API_TOKEN='<secure-token>'
npm run momo:live -- --amount 0.10 --currency GHS --msisdn 233XXXXXXXXX --executeLive --confirm "I UNDERSTAND REAL MONEY WILL MOVE"
```

## 6) Post-Mint Snapshot Capture

```powershell
Invoke-RestMethod -Uri 'http://localhost:3001/api/watchtower/nodes' | ConvertTo-Json -Depth 8
```

Archive:
- `snapshot.snapshotSequence`
- `snapshot.previousSnapshotHash`
- `snapshot.snapshotHash`
- `snapshotSigned`
