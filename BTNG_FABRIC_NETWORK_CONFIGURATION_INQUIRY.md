# 🧩 BTNG Fabric Network Configuration Inquiry

## Purpose

This document is the canonical quick-reference for BTNG Fabric topology, API routes, and operational commands used during production verification and onboarding.

## Key Operational Docs

- [BTNG Gold Price System Runbook](BTNG_GOLD_PRICE_SYSTEM.md)
- [BTNG Merchant API Guide](BTNG_MERCHANT_API_GUIDE.md)
- [BTNG Gold System Completion Notes](BTNG_GOLD_SYSTEM_COMPLETE.md)

## PM2 Oracle Service Commands

- `npm run pm2:oracle:start` — Start gold oracle updater service
- `npm run pm2:oracle:status` — View current PM2 service state
- `npm run pm2:oracle:logs` — Tail updater logs
- `npm run pm2:oracle:save` — Persist PM2 process list across reboot

## Fabric Network Configuration Source

Primary configuration file:
- [btng-fabric-network.json](btng-fabric-network.json)

Contains:
- Organizations and MSP IDs
- Peer and orderer endpoints
- Channel membership
- Chaincode metadata
- Documentation identity hash/signature block

## API Route Clarification

If route lookup at `c:/BTNGAI_files/app/api/auth/login` returns no results, use the API workspace path:
- [btng-api/app/api/auth/login/route.ts](btng-api/app/api/auth/login/route.ts)

Fabric-related API routes:
- [btng-api/app/api/btng/fabric/network/route.ts](btng-api/app/api/btng/fabric/network/route.ts)
- [btng-api/app/api/btng/fabric/node/route.ts](btng-api/app/api/btng/fabric/node/route.ts)
- [btng-api/app/api/btng/fabric/chaincode/route.ts](btng-api/app/api/btng/fabric/chaincode/route.ts)

## All-in-One BTNG Network Verification

Run the consolidated verifier:

- `npm run verify:btng-network` (strict)
- `npm run verify:btng-network:strict`
- `npm run verify:btng-network:soft`
- `npm run verify:btng-network:ci`

Checks include:
- Local signed documentation identity validity
- Watchtower metadata and documentation identity publication
- Snapshot cryptographic verification
- Mainnet readiness checks per discovered watchtower endpoint

Artifact output:
- Verifier writes a JSON artifact on every run
- Default strict artifact: `cache/verify-btng-network-strict.json`
- Default soft artifact: `cache/verify-btng-network-soft.json`
- CI artifact command path: `cache/verify-btng-network-ci.json`
- Override path with `--output <path>` or `BTNG_NETWORK_VERIFY_OUTPUT`

## Current Known Constraint

The verifier can still fail while cryptographic/documentation checks pass if remote node endpoints are unreachable (offline/timeouts). In that case, fix network ingress/port exposure on remote validator nodes first.

## BTNG Network Triage (3-Step)

1. **Identity Layer**
	- Run `npm run docs:identity:refresh`
	- Confirm Watchtower metadata is verified: `GET /api/watchtower/meta`

2. **Local Runtime Layer**
	- Run `npm run verify:btng-network:soft`
	- Confirm all non-readiness checks pass (signatures, snapshot verification, documentation identity)

3. **Remote Reachability Layer**
	- Open/forward validator ports and validate `/<health>` responses from remote nodes
	- Re-run `npm run verify:btng-network:strict` to enforce production readiness
