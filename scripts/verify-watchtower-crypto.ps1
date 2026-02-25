param(
    [int]$PreferredPort = 3001,
    [int]$StartupTimeoutSeconds = 45
)

$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Path $PSScriptRoot -Parent
$privateKeyPath = Join-Path $repoRoot 'btng-es256-private-key.pem'
$publicKeyPath = Join-Path $repoRoot 'btng-es256-public-key.pem'
$envLocalPath = Join-Path $repoRoot 'btng-api\.env.local'

function Get-ActiveWatchtowerBaseUrl {
    param(
        [int]$StartPort = 3001,
        [int]$EndPort = 3100
    )

    for ($port = $StartPort; $port -le $EndPort; $port++) {
        $url = "http://localhost:$port/api/watchtower/nodes"
        try {
            $response = Invoke-WebRequest -Uri $url -UseBasicParsing -TimeoutSec 5
            if ($response.StatusCode -eq 200) {
                return "http://localhost:$port"
            }
        }
        catch {
            # continue scan
        }
    }

    return $null
}

function Set-EnvLine {
    param(
        [Parameter(Mandatory = $true)]
        [string]$FilePath,
        [Parameter(Mandatory = $true)]
        [string]$Key,
        [Parameter(Mandatory = $true)]
        [string]$Value
    )

    if (-not (Test-Path $FilePath)) {
        New-Item -Path $FilePath -ItemType File -Force | Out-Null
    }

    $lines = Get-Content -Path $FilePath
    $pattern = "^$([Regex]::Escape($Key))="
    $replacement = "$Key=$Value"

    if ($lines -match $pattern) {
        $updatedLines = $lines | ForEach-Object {
            if ($_ -match $pattern) { $replacement } else { $_ }
        }
        Set-Content -Path $FilePath -Value $updatedLines
    }
    else {
        Add-Content -Path $FilePath -Value $replacement
    }
}

function Initialize-Es256Keys {
    if ((Test-Path $privateKeyPath) -and (Test-Path $publicKeyPath)) {
        Write-Host 'ES256 keypair already present.'
        return
    }

    Write-Host 'Generating ES256 keypair...'
    Set-Location $repoRoot
    node -e "const {generateKeyPairSync}=require('crypto');const fs=require('fs');const {privateKey,publicKey}=generateKeyPairSync('ec',{namedCurve:'prime256v1'});fs.writeFileSync('btng-es256-private-key.pem',privateKey.export({type:'pkcs8',format:'pem'}));fs.writeFileSync('btng-es256-public-key.pem',publicKey.export({type:'spki',format:'pem'}));"
}

Initialize-Es256Keys
Set-EnvLine -FilePath $envLocalPath -Key 'BTNG_ES256_PUBLIC_KEY_PATH' -Value './btng-es256-public-key.pem'
Set-EnvLine -FilePath $envLocalPath -Key 'BTNG_ES256_PRIVATE_KEY_PATH' -Value './btng-es256-private-key.pem'
Set-EnvLine -FilePath $envLocalPath -Key 'BTNG_WATCHTOWER_MAX_TIMESTAMP_SKEW_SECONDS' -Value '300'

$existingBaseUrl = Get-ActiveWatchtowerBaseUrl
$launchedByScript = $false
$suiteProcess = $null

$logDir = Join-Path $repoRoot 'cache'
if (-not (Test-Path $logDir)) {
    New-Item -Path $logDir -ItemType Directory -Force | Out-Null
}

$suiteOutLog = Join-Path $logDir 'verify-watchtower-suite.out.log'
$suiteErrLog = Join-Path $logDir 'verify-watchtower-suite.err.log'

try {
    if ($existingBaseUrl) {
        Write-Host "Using existing Watchtower instance at $existingBaseUrl"
    }
    else {
        Write-Host 'No active Watchtower instance detected. Launching suite...'
        $launchedByScript = $true
        $suiteProcess = Start-Process -FilePath 'powershell' -ArgumentList @('-NoProfile', '-ExecutionPolicy', 'Bypass', '-File', 'scripts/start-watchtower-suite.ps1', '-ApiMode', 'dev', '-Port', "$PreferredPort") -WorkingDirectory $repoRoot -PassThru -RedirectStandardOutput $suiteOutLog -RedirectStandardError $suiteErrLog

        $deadline = (Get-Date).AddSeconds($StartupTimeoutSeconds)
        do {
            Start-Sleep -Milliseconds 750
            $existingBaseUrl = Get-ActiveWatchtowerBaseUrl
            if ($existingBaseUrl) {
                break
            }

            if ($suiteProcess.HasExited) {
                break
            }
        } while ((Get-Date) -lt $deadline)

        if (-not $existingBaseUrl) {
            $outTail = ''
            $errTail = ''
            if (Test-Path $suiteOutLog) {
                $outTail = (Get-Content -Path $suiteOutLog -Tail 30) -join "`n"
            }
            if (Test-Path $suiteErrLog) {
                $errTail = (Get-Content -Path $suiteErrLog -Tail 30) -join "`n"
            }

            throw "Watchtower did not become reachable within $StartupTimeoutSeconds seconds.`n--- suite stdout ---`n$outTail`n--- suite stderr ---`n$errTail"
        }

        Write-Host "Watchtower became reachable at $existingBaseUrl"
    }

    Set-Location $repoRoot
    node scripts/test-signed-node-registration.js
    if ($LASTEXITCODE -ne 0) {
        throw "Signed registration test failed with exit code $LASTEXITCODE"
    }

    Write-Host 'verify-watchtower-crypto: PASS'
}
catch {
    Write-Error "verify-watchtower-crypto failed: $($_.Exception.Message)"
    exit 1
}
finally {
    if ($launchedByScript -and $suiteProcess -and -not $suiteProcess.HasExited) {
        Write-Host 'Stopping temporary Watchtower suite process...'
        Stop-Process -Id $suiteProcess.Id -Force
    }
}
