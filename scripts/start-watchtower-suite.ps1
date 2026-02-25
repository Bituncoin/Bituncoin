param(
    [ValidateSet('dev', 'prod')]
    [string]$ApiMode = 'dev',
    [int]$Port = 3001,
    [switch]$DisablePortAutoFallback
)

$ErrorActionPreference = 'Stop'

function Import-BtngEnvFile {
    param(
        [Parameter(Mandatory = $true)]
        [string]$Path
    )

    if (-not (Test-Path $Path)) {
        return
    }

    Get-Content -Path $Path | ForEach-Object {
        $line = $_.Trim()
        if (-not $line -or $line.StartsWith('#')) {
            return
        }

        $parts = $line -split '=', 2
        if ($parts.Count -ne 2) {
            return
        }

        $name = $parts[0].Trim()
        $value = $parts[1].Trim().Trim('"')

        [Environment]::SetEnvironmentVariable($name, $value, 'Process')
        Set-Item -Path "Env:$name" -Value $value
    }
}

function Test-LocalPortInUse {
    param(
        [Parameter(Mandatory = $true)]
        [int]$PortToTest
    )

    $client = New-Object System.Net.Sockets.TcpClient
    try {
        $asyncResult = $client.BeginConnect('127.0.0.1', $PortToTest, $null, $null)
        $connected = $asyncResult.AsyncWaitHandle.WaitOne(250)
        if (-not $connected) {
            return $false
        }

        $client.EndConnect($asyncResult)
        return $true
    }
    catch {
        return $false
    }
    finally {
        $client.Dispose()
    }
}

function Get-AvailablePort {
    param(
        [Parameter(Mandatory = $true)]
        [int]$StartPort,

        [int]$MaxAttempts = 30
    )

    for ($offset = 0; $offset -lt $MaxAttempts; $offset++) {
        $candidatePort = $StartPort + $offset
        if (-not (Test-LocalPortInUse -PortToTest $candidatePort)) {
            return $candidatePort
        }
    }

    throw "No available port found from $StartPort to $($StartPort + $MaxAttempts - 1)."
}

$repoRoot = Split-Path -Path $PSScriptRoot -Parent
$heartbeatTemplate = Join-Path $PSScriptRoot 'heartbeat.env.example'
$apiEnvTemplate = Join-Path $repoRoot 'btng-api\.env.example'
$apiEnvLocal = Join-Path $repoRoot 'btng-api\.env.local'

$selectedPort = $Port
if (-not $DisablePortAutoFallback -and (Test-LocalPortInUse -PortToTest $Port)) {
    $selectedPort = Get-AvailablePort -StartPort ($Port + 1)
    Write-Host "WARNING: Port $Port is in use. Falling back to port $selectedPort."
}

[Environment]::SetEnvironmentVariable('BTNG_WATCHTOWER_ACTIVE_PORT', "$selectedPort", 'Process')
Set-Item -Path 'Env:BTNG_WATCHTOWER_ACTIVE_PORT' -Value "$selectedPort"

Write-Host 'BTNG Watchtower Suite Launcher'
Write-Host '================================'
Write-Host "Mode: $ApiMode"
Write-Host "Port: $selectedPort"
Write-Host "Watchtower URL: http://localhost:$selectedPort/watchtower"
Write-Host "Merchant Package URL: http://localhost:$selectedPort/api/merchant/deployment-package"
Write-Host "Nodes API URL: http://localhost:$selectedPort/api/watchtower/nodes"

Import-BtngEnvFile -Path $heartbeatTemplate
if (Test-Path $apiEnvLocal) {
    Import-BtngEnvFile -Path $apiEnvLocal
    Write-Host "Loaded API env from: $apiEnvLocal"
}
else {
    Import-BtngEnvFile -Path $apiEnvTemplate
    Write-Host "WARNING: .env.local not found, loaded template vars from: $apiEnvTemplate"
}

Write-Host 'Starting heartbeat process...'
$heartbeatProcess = Start-Process -FilePath 'node' -ArgumentList 'scripts/heartbeat.js' -WorkingDirectory $repoRoot -PassThru -NoNewWindow
Write-Host "Heartbeat PID: $($heartbeatProcess.Id)"

try {
    if ($ApiMode -eq 'prod') {
        Write-Host 'Starting Watchtower app (btng-api production server)...'
    }
    else {
        Write-Host 'Starting Watchtower app (btng-api dev server)...'
    }

    Set-Location $repoRoot
    if ($ApiMode -eq 'prod') {
        npx next start btng-api -p $selectedPort
    }
    else {
        npx next dev btng-api -p $selectedPort
    }
}
finally {
    if ($heartbeatProcess -and -not $heartbeatProcess.HasExited) {
        Write-Host 'Stopping heartbeat process...'
        Stop-Process -Id $heartbeatProcess.Id -Force
    }
}
