$ErrorActionPreference = "Stop"

Write-Host "=== NSSM Plus - Dev Mode ===" -ForegroundColor Cyan

$wailsCmd = Get-Command wails -ErrorAction SilentlyContinue
if (!$wailsCmd) {
    Write-Host "[ERROR] Wails CLI is not installed." -ForegroundColor Red
    Write-Host "Install with: go install github.com/wailsapp/wails/v2/cmd/wails@latest" -ForegroundColor Yellow
    exit 1
}

Write-Host "Starting Wails dev server (hot reload)..." -ForegroundColor Yellow
Write-Host "NOTE: Run this as Administrator for service management operations." -ForegroundColor Yellow
Write-Host ""

Set-Location $PSScriptRoot
wails dev
