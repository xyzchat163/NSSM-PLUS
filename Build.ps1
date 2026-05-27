$ErrorActionPreference = "Stop"

Write-Host "=== NSSM Plus - Production Build ===" -ForegroundColor Cyan

$wailsCmd = Get-Command wails -ErrorAction SilentlyContinue
if (!$wailsCmd) {
    Write-Host "[ERROR] Wails CLI is not installed." -ForegroundColor Red
    Write-Host "Install with: go install github.com/wailsapp/wails/v2/cmd/wails@latest" -ForegroundColor Yellow
    exit 1
}

Set-Location $PSScriptRoot

Write-Host ""
Write-Host "[1/2] Building frontend..." -ForegroundColor Yellow
Set-Location "$PSScriptRoot\frontend"
npm run build
if ($LASTEXITCODE -ne 0) { Write-Host "[ERROR] Frontend build failed" -ForegroundColor Red; exit 1 }

Set-Location $PSScriptRoot

Write-Host ""
Write-Host "[2/2] Building Wails application..." -ForegroundColor Yellow
wails build
if ($LASTEXITCODE -ne 0) { Write-Host "[ERROR] Wails build failed" -ForegroundColor Red; exit 1 }

Write-Host ""
Write-Host "=== Build complete! ===" -ForegroundColor Green
Write-Host "Output: $PSScriptRoot\build\bin\nssm-plus.exe" -ForegroundColor White
