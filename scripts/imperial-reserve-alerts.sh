#!/bin/bash

# BTNG Imperial Reserve Alert System
# Automated monitoring and alerts for dashboard metrics

set -e

# Configuration
DASHBOARD_API="http://74.118.126.72:64799/api/dashboard/imperial-reserve"
ALERT_LOG="imperial-reserve-alerts.log"
NOTIFICATION_WEBHOOK="https://discord.com/api/webhooks/YOUR_WEBHOOK_URL"

# Alert thresholds
GOLD_BACKING_MIN=1.00
MERCHANT_MIN=1000
TRADE_VOLUME_MIN=500000
UPTIME_MIN=99.5

# Previous values for change detection
PREV_BACKING=""
PREV_MERCHANTS=""
PREV_VOLUME=""

log() {
    echo "$(date '+%Y-%m-%d %H:%M:%S') - $1" | tee -a "$ALERT_LOG"
}

alert() {
    local message="$1"
    local priority="${2:-normal}"

    log "ALERT [$priority]: $message"

    # Send to webhook if configured
    if [ -n "$NOTIFICATION_WEBHOOK" ] && [ "$NOTIFICATION_WEBHOOK" != "https://discord.com/api/webhooks/YOUR_WEBHOOK_URL" ]; then
        local emoji="📊"
        if [ "$priority" = "critical" ]; then
            emoji="🚨"
        elif [ "$priority" = "warning" ]; then
            emoji="⚠️"
        fi

        curl -X POST "$NOTIFICATION_WEBHOOK" \
            -H "Content-Type: application/json" \
            -d "{\"content\": \"$emoji BTNG Imperial Reserve: $message\"}" \
            --silent --output /dev/null
    fi
}

check_dashboard_metrics() {
    log "Checking Imperial Reserve dashboard metrics..."

    # Fetch current metrics
    local response
    response=$(curl -s -X GET "$DASHBOARD_API" \
        -H "Accept: application/json" \
        --max-time 30)

    if [ $? -ne 0 ] || [ -z "$response" ]; then
        alert "Failed to fetch dashboard metrics" "critical"
        return 1
    fi

    # Parse JSON response (using jq if available, otherwise basic parsing)
    local total_nations active_nations gold_backing merchant_count trade_volume

    if command -v jq >/dev/null 2>&1; then
        total_nations=$(echo "$response" | jq -r '.totalNations // 0')
        active_nations=$(echo "$response" | jq -r '.activeNations // 0')
        gold_backing=$(echo "$response" | jq -r '.goldBackingRatio // 0')
        merchant_count=$(echo "$response" | jq -r '.merchantCount // 0')
        trade_volume=$(echo "$response" | jq -r '.tradeVolume24h // 0')
    else
        # Basic parsing without jq
        total_nations=$(echo "$response" | grep -o '"totalNations":[0-9]*' | cut -d':' -f2)
        active_nations=$(echo "$response" | grep -o '"activeNations":[0-9]*' | cut -d':' -f2)
        gold_backing=$(echo "$response" | grep -o '"goldBackingRatio":[0-9.]*' | cut -d':' -f2)
        merchant_count=$(echo "$response" | grep -o '"merchantCount":[0-9]*' | cut -d':' -f2)
        trade_volume=$(echo "$response" | grep -o '"tradeVolume24h":[0-9.]*' | cut -d':' -f2)
    fi

    # Validate we got numbers
    if ! [[ "$gold_backing" =~ ^[0-9]+(\.[0-9]+)?$ ]]; then
        alert "Invalid gold backing ratio received: $gold_backing" "warning"
        return 1
    fi

    log "Metrics: Nations=$active_nations/$total_nations, Gold=$gold_backing, Merchants=$merchant_count, Volume=$$trade_volume"

    # Check critical thresholds
    if (( $(echo "$gold_backing < $GOLD_BACKING_MIN" | bc -l 2>/dev/null || echo "1") )); then
        alert "CRITICAL: Gold backing ratio below $GOLD_BACKING_MIN ($gold_backing)" "critical"
    fi

    if [ "$merchant_count" -lt "$MERCHANT_MIN" ]; then
        alert "WARNING: Merchant count below $MERCHANT_MIN ($merchant_count)" "warning"
    fi

    if (( $(echo "$trade_volume < $TRADE_VOLUME_MIN" | bc -l 2>/dev/null || echo "1") )); then
        alert "WARNING: Trade volume below $TRADE_VOLUME_MIN ($$trade_volume)" "warning"
    fi

    # Check for significant changes
    if [ -n "$PREV_BACKING" ]; then
        local backing_change
        backing_change=$(echo "scale=4; ($gold_backing - $PREV_BACKING) / $PREV_BACKING * 100" | bc -l 2>/dev/null)

        if (( $(echo "${backing_change#-} > 1.0" | bc -l 2>/dev/null || echo "0") )); then
            local direction="increased"
            if (( $(echo "$backing_change < 0" | bc -l 2>/dev/null || echo "0") )); then
                direction="decreased"
            fi
            alert "Gold backing ratio $direction by ${backing_change}% (now $gold_backing)" "normal"
        fi
    fi

    # Update previous values
    PREV_BACKING="$gold_backing"
    PREV_MERCHANTS="$merchant_count"
    PREV_VOLUME="$trade_volume"

    # Success metrics
    if [ "$active_nations" -eq "$total_nations" ]; then
        alert "✅ All $total_nations nations operational - Full continental coverage maintained" "normal"
    fi

    return 0
}

# Main alert loop
log "Starting BTNG Imperial Reserve Alert System"

while true; do
    if check_dashboard_metrics; then
        log "Dashboard check completed successfully"
    else
        log "Dashboard check failed"
    fi

    # Wait 5 minutes before next check
    sleep 300
done