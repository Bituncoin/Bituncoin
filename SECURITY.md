# BTNG Security Reporting Policy

## Purpose
This document outlines the process for responsibly reporting security vulnerabilities in the BTNG Sovereign Platform. We take security seriously and appreciate your efforts to protect our users and infrastructure.

## Scope
This policy applies to all components of the BTNG ecosystem:
- Sovereign node software
- Wallet applications
- Smart contracts
- API endpoints
- Governance modules
- Documentation and auxiliary tools

## Reporting Guidelines
Submit security reports to **security@bituncoin.org**.

## Vulnerability Report Template

**Subject**: Security Vulnerability Report - [Brief Title]

**Report Date**: [YYYY-MM-DD]

### Contact Information
**Reporter Name/Handle**: [Optional]

**Preferred Contact Method**: [e.g., Email, Telegram]

**PGP Key Fingerprint**: [Optional, for encrypted follow-up]

### Vulnerability Overview
**Component**: [e.g., Core Client (v1.x.x), Smart Contract 0x..., Wallet App]

**Type of Issue**: [e.g., Consensus Bypass, Smart Contract Reentrancy, Authentication Bypass]

**Severity**: [Critical/High/Medium/Low]

### Affected Products & Versions
**Package/Product Name**: [e.g., bituncoin-node]

**Affected Versions**: < 1.2.3 [List all confirmed versions]

**Patched Versions**: 1.2.3 [If known]

### Technical Details
**Description**: [Detailed explanation of vulnerability, conditions, security impact]

**Steps to Reproduce**: [Step-by-step instructions]

**Proof of Concept (PoC)**: [Code snippets, screenshots, scripts]

**Suggested Fix**: [Optional recommendations]

### Impact Assessment
**Attack Scenario**: [What an attacker could achieve]

**Affected Parties**: [e.g., All node operators, specific dApp users]

### Disclosure Preferences
**Public Disclosure**: ☐ Yes ☐ No (Anonymous)

**Embargo Requested**: [Duration, e.g., 90 days]

**CVSS Score**: [CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:H]

## What We Do After You Report
1. **Immediate Acknowledgment** (24 hours)
2. **Impact Assessment** (72 hours)
3. **Development of Fix** (Priority based on severity)
4. **Coordinated Disclosure** (Respecting your preferences)
5. **Credit & Thanks** (Per your preference)

## Eligible Severity Levels
- **Critical**: Remote code execution, consensus breaks, fund theft
- **High**: Authentication bypass, private key exposure, DoS impacting consensus
- **Medium**: Information disclosure, rate limiting bypass
- **Low**: Minor UI issues, documentation errors

## Rewards
Security researchers reporting valid vulnerabilities may receive:
- Bounty payments (BTC/BTNG)
- Recognition in release notes
- Swag and ceremonial acknowledgments

## Safe Harbor
Any security reports made in good faith are exempt from legal action.

## Questions?
Contact security@bituncoin.org for clarification.

---

*Updated: March 2026*
*Ceremonial License v1.0*
*Swiss Jurisdiction + ICC Arbitration*
