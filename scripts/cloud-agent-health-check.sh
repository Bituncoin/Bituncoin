#!/bin/bash
# BTNG Cloud Agent Health Check Script
# This script monitors cloud agent status and helps diagnose timeout issues

echo "🔍 BTNG Cloud Agent Health Check"
echo "================================="

# Check current timestamp
echo "Current time: $(date)"
echo ""

# Check GitLab CI pipeline status
echo "📊 Checking GitLab CI Pipeline Status..."
if command -v gitlab-ci-pipelines-exporter &> /dev/null; then
    echo "GitLab CI monitoring available"
else
    echo "⚠️  GitLab CI monitoring not available locally"
fi
echo ""

# Check system resources
echo "💻 System Resource Check:"
echo "CPU Usage: $(top -bn1 | grep "Cpu(s)" | sed "s/.*, *\([0-9.]*\)%* id.*/\1/" | awk '{print 100 - $1"%"}')"
echo "Memory Usage: $(free | grep Mem | awk '{printf "%.2f%%", $3/$2 * 100.0}')"
echo "Disk Usage: $(df / | tail -1 | awk '{print $5}')"
echo ""

# Check network connectivity
echo "🌐 Network Connectivity Check:"
echo "Testing connection to GitLab..."
if ping -c 1 gitlab.com &> /dev/null; then
    echo "✅ GitLab.com reachable"
else
    echo "❌ GitLab.com unreachable"
fi

echo "Testing connection to GitHub..."
if ping -c 1 github.com &> /dev/null; then
    echo "✅ GitHub.com reachable"
else
    echo "❌ GitHub.com unreachable"
fi
echo ""

# Check BTNG services
echo "🔧 BTNG Service Status:"
services=("btng-api" "gold-oracle" "fabric-network" "imperial-reserve")

for service in "${services[@]}"; do
    if pgrep -f "$service" > /dev/null; then
        echo "✅ $service: Running"
    else
        echo "❌ $service: Not running"
    fi
done
echo ""

# Check recent logs for timeout patterns
echo "📝 Recent Timeout Patterns in Logs:"
if [ -f "imperial-reserve-alerts.log" ]; then
    timeout_count=$(grep -c "TIMEOUT\|timeout\|failed" imperial-reserve-alerts.log | tail -10)
    if [ "$timeout_count" -gt 0 ]; then
        echo "⚠️  Found $timeout_count timeout/failure entries in recent logs"
        echo "Last 3 timeout entries:"
        grep -i "timeout\|failed" imperial-reserve-alerts.log | tail -3
    else
        echo "✅ No recent timeout patterns found"
    fi
else
    echo "ℹ️  No imperial-reserve-alerts.log found"
fi
echo ""

# Recommendations
echo "💡 Recommendations:"
echo "1. Ensure GitLab Runner has sufficient resources (CPU, Memory, Disk)"
echo "2. Check network connectivity to GitLab.com"
echo "3. Verify SSH keys and VPS credentials are valid"
echo "4. Monitor resource usage during pipeline execution"
echo "5. Consider increasing job timeouts if operations are resource-intensive"
echo ""

echo "🏁 Health check complete. If issues persist, check issue #13 for updates."