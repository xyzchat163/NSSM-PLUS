$ErrorActionPreference = "Stop"

Write-Host "=== NSSM Plus - Install Dependencies ===" -ForegroundColor Cyan

if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "[ERROR] Go is not installed. Install from https://go.dev/dl/" -ForegroundColor Red
    exit 1
}

if (!(Get-Command node -ErrorAction SilentlyContinue)) {
    Write-Host "[ERROR] Node.js is not installed. Install from https://nodejs.org/" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "[1/3] Installing Go dependencies..." -ForegroundColor Yellow
go mod download
if ($LASTEXITCODE -ne 0) { Write-Host "[ERROR] Go mod download failed" -ForegroundColor Red; exit 1 }

Write-Host ""
Write-Host "[2/3] Installing frontend dependencies..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\frontend"
npm install
if ($LASTEXITCODE -ne 0) { Write-Host "[ERROR] npm install failed" -ForegroundColor Red; exit 1 }

Set-Location $PSScriptRoot

Write-Host ""
Write-Host "[3/3] Verifying Go build..." -ForegroundColor Yellow
go build ./...
if ($LASTEXITCODE -ne 0) { Write-Host "[ERROR] Go build failed" -ForegroundColor Red; exit 1 }

Write-Host ""
Write-Host "=== Install complete! ===" -ForegroundColor Green
Write-Host "Run '.\Run.ps1' for dev mode, or '.\Build.ps1' for production build." -ForegroundColor White
