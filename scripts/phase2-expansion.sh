#!/bin/bash

# BTNG Sovereign Node Joiner Script
# Automated deployment for new nation nodes joining the mesh

set -e

# Configuration
PRIMARY_NODE="http://74.118.126.72:64799"
JOIN_LOG="node-join.log"
DEPLOY_TIMEOUT=600  # 10 minutes

# Target nations for Phase 2 expansion
declare -a PHASE2_NATIONS=(
    "Ethiopia:197.156.108.1"
    "Morocco:105.73.32.1"
    "Senegal:41.208.104.1"
    "Tanzania:154.118.230.1"
    "Zambia:41.63.0.1"
)

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$JOIN_LOG"
}

error_exit() {
    echo "ERROR: $1" >&2
    log "JOIN FAILED: $1"
    exit 1
}

success() {
    echo "SUCCESS: $1"
    log "JOIN SUCCESS: $1"
}

# Validate primary node connectivity
validate_primary() {
    log "Validating primary node connectivity..."

    if ! curl -s --max-time 30 "$PRIMARY_NODE/api/health" > /dev/null; then
        error_exit "Cannot connect to primary node $PRIMARY_NODE"
    fi

    local health_response
    health_response=$(curl -s "$PRIMARY_NODE/api/health")
    local status
    status=$(echo "$health_response" | jq -r '.status // empty' 2>/dev/null)

    if [ "$status" != "operational" ]; then
        error_exit "Primary node not operational (status: $status)"
    fi

    success "Primary node validated"
}

# Deploy Sovereign OS to new node
deploy_sovereign_os() {
    local nation_info="$1"
    local nation_name
    local node_ip

    nation_name=$(echo "$nation_info" | cut -d':' -f1)
    node_ip=$(echo "$nation_info" | cut -d':' -f2)
    local node_url="http://$node_ip:64799"

    log "Deploying Sovereign OS to $nation_name ($node_url)..."

    # Check if node is already online
    if curl -s --max-time 10 "$node_url/api/health" > /dev/null 2>&1; then
        log "Node $nation_name already online, updating to Sovereign OS..."
    else
        log "Node $nation_name not yet online, preparing for remote deployment..."
        # In real scenario, this would trigger actual server provisioning
        log "Simulating node provisioning for $nation_name..."
        sleep 5
    fi

    # Deploy Private Banker Agent
    log "Deploying Private Banker Agent..."
    local banker_payload='{
        "version": "1.0.0",
        "agents": ["liquidity-optimizer", "auditor", "predictive-alerting"],
        "mpcEnabled": true,
        "primaryNode": "'"$PRIMARY_NODE"'"
    }'

    if ! curl -s -X POST "$node_url/api/os/deploy-banker" \
        -H "Content-Type: application/json" \
        -d "$banker_payload" \
        --max-time 60 > /dev/null; then
        log "Banker deployment to $nation_name failed, will retry..."
        return 1
    fi

    # Enable Self-Custody Features
    log "Enabling Self-Custody features..."
    local custody_payload='{
        "mpcSharding": true,
        "offlineSigning": true,
        "hardwareIsolation": true,
        "ghostKeySystem": true
    }'

    if ! curl -s -X POST "$node_url/api/os/enable-custody" \
        -H "Content-Type: application/json" \
        -d "$custody_payload" \
        --max-time 60 > /dev/null; then
        log "Custody enablement for $nation_name failed, will retry..."
        return 1
    fi

    # Register with Watchtower
    log "Registering with Watchtower..."
    local watchtower_payload='{
        "nodeId": "'"$nation_name"'",
        "endpoint": "'"$node_url"'",
        "country": "'"$nation_name"'",
        "role": "sovereign-validator",
        "weight": 1.0
    }'

    if ! curl -s -X POST "$PRIMARY_NODE/api/watchtower/nodes/register" \
        -H "Content-Type: application/json" \
        -d "$watchtower_payload" \
        --max-time 60 > /dev/null; then
        log "Watchtower registration for $nation_name failed"
        return 1
    fi

    # Update Network Services
    log "Updating network services..."
    local services_payload='{
        "unifiedBanking": true,
        "microKernel": true,
        "continentalSync": true,
        "phase2Enabled": true
    }'

    if ! curl -s -X POST "$node_url/api/os/update-services" \
        -H "Content-Type: application/json" \
        -d "$services_payload" \
        --max-time 60 > /dev/null; then
        log "Services update for $nation_name failed"
        return 1
    fi

    success "Sovereign OS deployed to $nation_name"
    return 0
}

# Main deployment process
log "Starting BTNG Phase 2 Expansion - Deploying to 5 nations"

validate_primary

total_nations=${#PHASE2_NATIONS[@]}
successful_deploys=0
failed_deploys=0

for nation_info in "${PHASE2_NATIONS[@]}"; do
    nation_name=$(echo "$nation_info" | cut -d':' -f1)

    log "Processing $nation_name..."

    # Attempt deployment with retry
    max_retries=3
    retry_count=0
    deployed=false

    while [ $retry_count -lt $max_retries ] && [ "$deployed" = false ]; do
        if deploy_sovereign_os "$nation_info"; then
            deployed=true
            ((successful_deploys++))
        else
            ((retry_count++))
            if [ $retry_count -lt $max_retries ]; then
                log "Retrying $nation_name deployment (attempt $((retry_count + 1))/$max_retries)..."
                sleep 30
            fi
        fi
    done

    if [ "$deployed" = false ]; then
        log "Failed to deploy to $nation_name after $max_retries attempts"
        ((failed_deploys++))
    fi

    # Brief pause between deployments
    sleep 10
done

# Final status
log "Phase 2 expansion complete - $successful_deploys/$total_nations nations deployed"

if [ $successful_deploys -eq $total_nations ]; then
    success "Phase 2 expansion successful - All target nations joined the mesh"
    echo "🌍 CONTINENTAL EXPANSION: PHASE 2 COMPLETE"
    echo "54 nations now have Sovereign OS capability"
elif [ $successful_deploys -gt 2 ]; then
    success "Phase 2 expansion partially successful - $successful_deploys nations online"
    echo "⚠️ Partial expansion - $failed_deploys nations pending"
else
    error_exit "Phase 2 expansion failed - Only $successful_deploys/$total_nations nations deployed"
fi

log "BTNG Phase 2 Sovereign Expansion operation complete"