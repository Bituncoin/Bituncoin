@echo off
REM BTNG Sovereign Platform - Complete Verification Suite (Windows)
REM Runs all verification checks for Sepolia deployment readiness

echo 🇰🇪 BTNG SOVEREIGN PLATFORM - VERIFICATION SUITE
echo ===============================================
echo Testing complete BTNG sovereign gold standard system
echo.

REM Colors for Windows (limited support)
set "GREEN=[92m"
set "RED=[91m"
set "YELLOW=[93m"
set "NC=[0m"

REM Track results
set TOTAL_TESTS=0
set PASSED_TESTS=0

:run_test
set "test_name=%~1"
set "test_script=%~2"

echo.
echo %YELLOW%▶ Running: %test_name%%NC%
echo ----------------------------------------

set /a TOTAL_TESTS+=1

if exist "%test_script%" (
    REM For Node.js scripts, run directly
    if "%test_script:~-3%"==".js" (
        node "%test_script%"
        if %errorlevel% equ 0 (
            echo %GREEN%✅ PASSED: %test_name%%NC%
            set /a PASSED_TESTS+=1
        ) else (
            echo %RED%❌ FAILED: %test_name%%NC%
        )
    ) else (
        REM For bash scripts, try to run with bash if available
        where bash >nul 2>nul
        if %errorlevel% equ 0 (
            bash "%test_script%"
            if %errorlevel% equ 0 (
                echo %GREEN%✅ PASSED: %test_name%%NC%
                set /a PASSED_TESTS+=1
            ) else (
                echo %RED%❌ FAILED: %test_name%%NC%
            )
        ) else (
            echo %YELLOW%⚠️ SKIPPED: %test_name% (bash not available on Windows)%NC%
        )
    )
) else (
    echo %RED%❌ MISSING: %test_script%%NC%
)
goto :eof

REM Pre-flight checks
echo 🔍 Pre-flight checks...
echo ----------------------

if not exist "package.json" (
    echo %RED%❌ Not in BTNG project directory (missing package.json)%NC%
    exit /b 1
)

if not exist "contracts" (
    echo %RED%❌ Contracts directory not found%NC%
    exit /b 1
)

if not exist "k8s" (
    echo %RED%❌ Kubernetes policies directory not found%NC%
    exit /b 1
)

echo %GREEN%✅ Project structure verified%NC%

REM Check Node.js
where node >nul 2>nul
if %errorlevel% equ 0 (
    for /f "tokens=*" %%i in ('node --version') do set NODE_VERSION=%%i
    echo %GREEN%✅ Node.js %NODE_VERSION%%NC%
) else (
    echo %RED%❌ Node.js not found%NC%
    exit /b 1
)

REM Check if dependencies are installed
if exist "node_modules" (
    echo %GREEN%✅ Dependencies installed%NC%
) else (
    echo %YELLOW%⚠️ Installing dependencies...%NC%
    npm install
)

echo.
echo 🚀 Starting BTNG Verification Suite...
echo ======================================

REM Run individual tests
call :run_test "Environment Configuration" "scripts\check-testnet-env.js"
call :run_test "Sepolia Smart Contract Deployment" "scripts\verify-sepolia-deployment.sh"
call :run_test "Sovereign Gatekeeper Policies" "scripts\verify-gatekeeper-policies.sh"
call :run_test "Zero-Trust Admission Control (JWT/TLS)" "scripts\verify-admission-control.sh"
call :run_test "Cryptographic Signature Verification" "scripts\verify-cryptography.sh"

REM Summary
echo.
echo 📊 VERIFICATION SUITE RESULTS
echo ============================
echo Total Tests: %TOTAL_TESTS%
echo Passed: %PASSED_TESTS%
set /a FAILED_TESTS=%TOTAL_TESTS%-%PASSED_TESTS%
echo Failed: %FAILED_TESTS%

if %PASSED_TESTS% equ %TOTAL_TESTS% (
    echo.
    echo %GREEN%🎉 ALL TESTS PASSED! 🇰🇪%NC%
    echo.
    echo 🇰🇪 Your BTNG Sovereign Gold Standard is READY for Sepolia deployment!
    echo.
    echo 📋 Final Deployment Steps:
    echo 1. Update .env with your real credentials
    echo 2. Run: npm run deploy:testnet
    echo 3. Verify contracts on Etherscan
    echo 4. Deploy Gatekeeper policies: .\scripts\deploy-gatekeeper-policies.sh
    echo.
    echo 🌟 BTNG will bring sovereign prosperity to African nations! 🌟
    exit /b 0
) else (
    echo.
    echo %RED%❌ SOME TESTS FAILED%NC%
    echo.
    echo 🔧 Please fix the failed tests before deploying to production.
    echo    Check the error messages above for guidance.
    exit /b 1
)