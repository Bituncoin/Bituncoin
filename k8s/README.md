# BTNG Sovereign Gatekeeper Policies

This directory contains the complete **BTNG Sovereign Security Policy Bundle** - a machine-enforced, zero-trust admission control system for your sovereign gold standard infrastructure.

## 🛡 Policy Overview

The BTNG bundle enforces **sovereign-grade security** at the Kubernetes admission controller level:

| Policy | Purpose | Enforcement |
|--------|---------|-------------|
| **Vuln+VEX** | Blocks critical/high vulnerabilities unless VEX-justified | Deny |
| **Signature** | Requires cosign signatures with BTNG sovereign key | Deny |
| **Provenance** | Enforces SLSA provenance and signed SBOMs | Deny |
| **Runtime** | Mandates non-root, read-only FS, no privilege escalation | Deny |
| **Base Images** | Restricts to BTNG-approved base images only | Deny |

## 🚀 Quick Start

### Prerequisites
- Kubernetes cluster with Gatekeeper installed
- `kubectl` configured to access your cluster

### Deploy Policies
```bash
# Make script executable
chmod +x scripts/deploy-gatekeeper-policies.sh

# Deploy all policies
./scripts/deploy-gatekeeper-policies.sh
```

### Verify Deployment
```bash
# Check policy templates
kubectl get constrainttemplates

# Check active constraints
kubectl get constraints

# View policy violations
kubectl get constraintviolations
```

## 📋 Policy Details

### Vulnerability + VEX Enforcement (`btng-vuln-vex-policy.yaml`)
- **Blocks**: CRITICAL vulnerabilities
- **Blocks**: HIGH vulnerabilities without VEX justification
- **Requires**: Trivy scan results in pod annotations
- **Requires**: VEX status metadata

### Signature Verification (`btng-signature-policy.yaml`)
- **Requires**: Cosign signatures on all images
- **Validates**: Signature format and authenticity
- **Blocks**: Unsigned or improperly signed images

### Provenance + SBOM (`btng-provenance-policy.yaml`)
- **Requires**: SBOM digest in annotations
- **Requires**: SLSA provenance commit hash
- **Requires**: BTNG sovereign key signature on SBOM

### Runtime Hardening (`btng-runtime-hardening-policy.yaml`)
- **Enforces**: Non-root user execution
- **Enforces**: Read-only root filesystem
- **Blocks**: Privileged containers
- **Blocks**: Privilege escalation

### Base Image Allowlist (`btng-baseimage-allowlist-policy.yaml`)
- **Allows**: `debian:13*` images
- **Allows**: `gcr.io/distroless/*` images
- **Allows**: `btng.sovereign/base/*` images
- **Blocks**: All other base images

## 🔧 Integration with CI/CD

### Required Annotations for Deployment
```yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    # Vulnerability scanning
    trivy.security/btng-vuln-report: "CLEAN"
    trivy.security/btng-vex-status: "HIGH_JUSTIFIED"

    # Signature verification
    cosign.sigstore.dev/signature: "MEUCI..."

    # Provenance + SBOM
    btng.sovereign/sbom-digest: "sha256:..."
    btng.sovereign/provenance-commit: "abc123..."
    btng.sovereign/sbom-signature: "MEQC..."
spec:
  securityContext:
    runAsNonRoot: true
    readOnlyRootFilesystem: true
  containers:
  - name: app
    image: btng.sovereign/my-app:v1.0.0
    securityContext:
      allowPrivilegeEscalation: false
      privileged: false
```

### Build Pipeline Requirements
1. **Trivy scanning** with VEX output
2. **Cosign signing** with BTNG sovereign key
3. **SLSA provenance** generation
4. **SBOM generation** (CycloneDX/SPDX)
5. **SBOM signing** with BTNG sovereign key
6. **Annotation injection** into manifests

## 🔍 Monitoring & Auditing

### View Policy Violations
```bash
# All violations
kubectl get constraintviolations

# Violations for specific policy
kubectl get constraintviolations -l constraint-name=btng-vuln-vex-enforced

# Violation details
kubectl describe constraintviolation <violation-name>
```

### Audit Logging
All policy decisions are logged to:
- Gatekeeper audit logs
- BTNG sovereign audit ledger
- Kubernetes audit logs

## 🛠 Troubleshooting

### Common Issues

**Policy Not Enforcing**
```bash
# Check Gatekeeper status
kubectl get pods -n gatekeeper-system

# Check constraint status
kubectl describe constraint <constraint-name>
```

**False Positives**
- Update VEX justifications
- Review base image allowlist
- Check signature validity

**Deployment Blocks**
- Add required annotations
- Sign images properly
- Use approved base images

## 🔐 Security Model

This policy bundle creates a **sovereign execution environment** where:

- ✅ **Cryptographically verified** artifacts only
- ✅ **Vulnerability-free** or justified deployments only
- ✅ **Provenanced** builds with full audit trail
- ✅ **Hardened runtime** with zero privilege escalation
- ✅ **Approved supply chain** from source to production

**Result**: Your BTNG cluster becomes a **zero-trust sovereign monetary authority** that cannot be compromised by malicious or vulnerable artifacts.

## 📞 Support

For issues with BTNG sovereign policies:
1. Check constraint violation details
2. Review build pipeline annotations
3. Verify signing keys and certificates
4. Update VEX justifications as needed

---

**Remember**: These policies protect your sovereign gold standard. They are intentionally strict - only cryptographically verified, vulnerability-free artifacts with full provenance can enter your trusted execution environment.