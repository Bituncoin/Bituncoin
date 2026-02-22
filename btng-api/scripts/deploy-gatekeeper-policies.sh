#!/bin/bash

# BTNG Sovereign Gatekeeper Policy Deployment
# This script deploys the complete BTNG security policy bundle

set -e

echo "🛡 Deploying BTNG Sovereign Gatekeeper Policies..."

# Check if Gatekeeper is installed
if ! kubectl get crd constrainttemplates.templates.gatekeeper.sh >/dev/null 2>&1; then
    echo "❌ Gatekeeper not found. Install Gatekeeper first:"
    echo "kubectl apply -f https://raw.githubusercontent.com/open-policy-agent/gatekeeper/release-3.13/deploy/gatekeeper.yaml"
    exit 1
fi

# Deploy policy templates
echo "📋 Deploying policy templates..."
kubectl apply -f k8s/gatekeeper-vuln-vex-policy.yaml
kubectl apply -f k8s/gatekeeper-signature-policy.yaml
kubectl apply -f k8s/gatekeeper-provenance-policy.yaml
kubectl apply -f k8s/gatekeeper-runtime-hardening-policy.yaml
kubectl apply -f k8s/gatekeeper-baseimage-allowlist-policy.yaml

# Wait for templates to be ready
echo "⏳ Waiting for policy templates..."
kubectl wait --for=condition=established --timeout=60s crd/btngvulnvs.constraints.gatekeeper.sh
kubectl wait --for=condition=established --timeout=60s crd/btngsignaturerequireds.constraints.gatekeeper.sh
kubectl wait --for=condition=established --timeout=60s crd/btngprovenancerequireds.constraints.gatekeeper.sh
kubectl wait --for=condition=established --timeout=60s crd/btnruntimehardenings.constraints.gatekeeper.sh
kubectl wait --for=condition=established --timeout=60s crd/btngbaseimageallowlists.constraints.gatekeeper.sh

# Deploy constraints
echo "🔒 Deploying policy constraints..."
kubectl apply -f k8s/gatekeeper-constraints.yaml

# Verify deployment
echo "✅ Verifying deployment..."
kubectl get constrainttemplates
kubectl get constraints

echo ""
echo "🎯 BTNG Sovereign Security Policies Active!"
echo "Your cluster now enforces:"
echo "  • Vulnerability + VEX validation"
echo "  • Cryptographic signature verification"
echo "  • SBOM + provenance requirements"
echo "  • Runtime hardening (non-root, read-only FS)"
echo "  • Approved base image allowlist"
echo ""
echo "Only sovereign-grade artifacts can now enter your cluster."