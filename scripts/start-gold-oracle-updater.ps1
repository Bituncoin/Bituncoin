param(
    [int]$IntervalMs = 60000,
    [switch]$Once
)

$ErrorActionPreference = 'Stop'

function Assert-EnvVar {
    param([string]$Name)
    $value = (Get-Item -Path ("Env:" + $Name) -ErrorAction SilentlyContinue).Value
    if ([string]::IsNullOrWhiteSpace($value)) {
        throw "Missing required environment variable: $Name"
    }
}

try {
    Write-Output "Verifying BTNG documentation identity..."
    & node scripts/verify-doc-identity.js
    if ($LASTEXITCODE -ne 0) {
        throw "Documentation identity verification failed. Refusing to start oracle updater."
    }

    Assert-EnvVar -Name 'BTNG_ORACLE_RPC_URL'
    Assert-EnvVar -Name 'BTNG_ORACLE_ADDRESS'
    Assert-EnvVar -Name 'BTNG_ORACLE_ADMIN_PRIVATE_KEY'

    $apiUrl = (Get-Item -Path 'Env:BTNG_GOLD_PRICE_API_URL' -ErrorAction SilentlyContinue).Value
    if ([string]::IsNullOrWhiteSpace($apiUrl)) {
        $env:BTNG_GOLD_PRICE_API_URL = 'http://localhost:64799/api/btng/gold/price/latest'
    }

    if ($IntervalMs -lt 1000) {
        throw 'IntervalMs must be at least 1000 ms.'
    }

    $args = @('scripts/push-gold-oracle-price.js')
    if ($Once.IsPresent) {
        $args += '--once'
    }
    else {
        $args += '--intervalMs'
        $args += "$IntervalMs"
    }

    Write-Output "Starting BTNG oracle updater..."
    Write-Output "RPC URL: $($env:BTNG_ORACLE_RPC_URL)"
    Write-Output "Oracle Address: $($env:BTNG_ORACLE_ADDRESS)"
    Write-Output "Gold API URL: $($env:BTNG_GOLD_PRICE_API_URL)"
    if ($Once.IsPresent) {
        Write-Output "Mode: once"
    }
    else {
        Write-Output "Mode: continuous"
        Write-Output "IntervalMs: $IntervalMs"
    }

    & node @args
    $exitCode = $LASTEXITCODE
    if ($exitCode -ne 0) {
        throw "Oracle updater exited with code $exitCode"
    }
}
catch {
    Write-Error $_
    exit 1
}
