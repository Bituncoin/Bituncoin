#!/usr/bin/env bash
set -euo pipefail

# BTNG Sovereign Gatekeeper Policy Authentication Script
# Tests the policy engine for sovereign access control

echo "🛡 BTNG Sovereign Gatekeeper Policy Authentication"
echo "================================================"

# Check if Gatekeeper policies exist
POLICY_DIR="./k8s"
if [ ! -d "$POLICY_DIR" ]; then
    echo "❌ Policy directory not found: $POLICY_DIR"
    exit 1
fi

echo "✅ Policy directory found"

# Check for required policy files
REQUIRED_POLICIES=(
    "gatekeeper-vuln-vex-policy.yaml"
    "gatekeeper-signature-policy.yaml"
    "gatekeeper-provenance-policy.yaml"
    "gatekeeper-runtime-hardening-policy.yaml"
    "gatekeeper-baseimage-allowlist-policy.yaml"
    "gatekeeper-constraints.yaml"
)

for policy in "${REQUIRED_POLICIES[@]}"; do
    if [ ! -f "$POLICY_DIR/$policy" ]; then
        echo "❌ Required policy file missing: $policy"
        exit 1
    fi
done

echo "✅ All required policy files present"

# Create a mock policy engine test
echo ""
echo "🔧 Testing policy engine logic..."

node - <<EOF
// Mock policy engine test for BTNG sovereign access control
const fs = require('fs');

function testPolicy(request) {
  console.log('Testing policy for request:', JSON.stringify(request, null, 2));

  // Mock policy rules (simulating Gatekeeper logic)
  const policies = {
    // Vulnerability policy
    vulnCheck: (req) => {
      if (req.annotations?.['trivy.security/btng-vuln-report'] === 'CRITICAL') {
        return { authorized: false, reason: 'CRITICAL vulnerabilities detected' };
      }
      return { authorized: true };
    },

    // Signature policy
    signatureCheck: (req) => {
      if (!req.annotations?.['cosign.sigstore.dev/signature']) {
        return { authorized: false, reason: 'Missing cosign signature' };
      }
      return { authorized: true };
    },

    // Runtime hardening policy
    runtimeCheck: (req) => {
      if (req.spec?.securityContext?.runAsNonRoot !== true) {
        return { authorized: false, reason: 'Must run as non-root user' };
      }
      if (!req.spec?.securityContext?.readOnlyRootFilesystem) {
        return { authorized: false, reason: 'Must have read-only root filesystem' };
      }
      return { authorized: true };
    },

    // Base image policy
    baseImageCheck: (req) => {
      const allowedImages = ['debian:13', 'gcr.io/distroless', 'btng.sovereign/base'];
      const image = req.spec?.containers?.[0]?.image || '';

      const isAllowed = allowedImages.some(allowed => image.startsWith(allowed));
      if (!isAllowed) {
        return { authorized: false, reason: \`Image '\${image}' not in approved allowlist\` };
      }
      return { authorized: true };
    }
  };

  // Test all policies
  const results = [];
  for (const [policyName, policyFn] of Object.entries(policies)) {
    const result = policyFn(request);
    results.push({ policy: policyName, ...result });
  }

  // Overall authorization
  const allAuthorized = results.every(r => r.authorized);
  const failedPolicies = results.filter(r => !r.authorized);

  return {
    authorized: allAuthorized,
    results: results,
    failedPolicies: failedPolicies,
    request: request
  };
}

// Test cases
const testCases = [
  {
    name: 'Valid sovereign deployment',
    request: {
      metadata: {
        annotations: {
          'trivy.security/btng-vuln-report': 'CLEAN',
          'cosign.sigstore.dev/signature': 'MEUCIQ...',
          'btng.sovereign/sbom-digest': 'sha256:...',
          'btng.sovereign/provenance-commit': 'abc123'
        }
      },
      spec: {
        securityContext: {
          runAsNonRoot: true,
          readOnlyRootFilesystem: true
        },
        containers: [{
          image: 'btng.sovereign/base:v1.0.0',
          securityContext: {
            allowPrivilegeEscalation: false
          }
        }]
      }
    }
  },
  {
    name: 'Invalid: critical vulnerability',
    request: {
      metadata: {
        annotations: {
          'trivy.security/btng-vuln-report': 'CRITICAL'
        }
      },
      spec: {
        securityContext: {
          runAsNonRoot: true,
          readOnlyRootFilesystem: true
        },
        containers: [{
          image: 'btng.sovereign/base:v1.0.0'
        }]
      }
    }
  },
  {
    name: 'Invalid: unapproved base image',
    request: {
      metadata: {
        annotations: {
          'trivy.security/btng-vuln-report': 'CLEAN',
          'cosign.sigstore.dev/signature': 'MEUCIQ...'
        }
      },
      spec: {
        securityContext: {
          runAsNonRoot: true,
          readOnlyRootFilesystem: true
        },
        containers: [{
          image: 'ubuntu:latest'
        }]
      }
    }
  }
];

// Run tests
let passedTests = 0;
let totalTests = testCases.length;

testCases.forEach((testCase, index) => {
  console.log(\`\\n🧪 Test \${index + 1}: \${testCase.name}\`);
  const result = testPolicy(testCase.request);

  if (result.authorized) {
    console.log('✅ PASSED - Request authorized');
    passedTests++;
  } else {
    console.log('❌ FAILED - Request denied');
    console.log('Failed policies:', result.failedPolicies.map(p => p.policy + ': ' + p.reason));
  }
});

console.log(\`\\n📊 Test Results: \${passedTests}/\${totalTests} tests passed\`);

if (passedTests === totalTests) {
  console.log('🎉 All policy tests PASSED!');
  process.exit(0);
} else {
  console.log('❌ Some policy tests FAILED!');
  process.exit(1);
}
EOF

if [ $? -eq 0 ]; then
    echo ""
    echo "🎉 BTNG Gatekeeper policy authentication PASSED!"
    echo ""
    echo "📋 Next steps:"
    echo "1. Deploy policies: ./scripts/deploy-gatekeeper-policies.sh"
    echo "2. Test with real Kubernetes cluster"
else
    echo "❌ BTNG Gatekeeper policy authentication FAILED!"
    exit 1
fi