# BoltVisa Backend Startup Script
# This script ensures the backend starts properly

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  BoltVisa Backend API Server" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Set Go path
$env:Path = "C:\Program Files\Go\bin;$env:Path"

# Change to API directory
Set-Location $PSScriptRoot

# Check if Go is installed
try {
    $goVersion = & go version 2>&1
    Write-Host "✅ Go found: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Please install Go from https://go.dev/dl/" -ForegroundColor Yellow
    pause
    exit 1
}

# Download dependencies
Write-Host ""
Write-Host "Downloading dependencies..." -ForegroundColor Yellow
& go mod download 2>&1 | Out-Null
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to download dependencies" -ForegroundColor Red
    pause
    exit 1
}
Write-Host "✅ Dependencies downloaded" -ForegroundColor Green

# Set environment variables
$env:DATABASE_URL = "sqlite://boltvisa.db"
$env:JWT_SECRET = if ($env:JWT_SECRET) { $env:JWT_SECRET } else { "dev-secret-key-change-in-production" }
$env:PORT = "8080"

Write-Host ""
Write-Host "Starting server on http://localhost:8080" -ForegroundColor Cyan
Write-Host "Database: SQLite (boltvisa.db)" -ForegroundColor Gray
Write-Host ""
Write-Host "All optimizations applied:" -ForegroundColor Green
Write-Host "  • Async audit logging" -ForegroundColor Gray
Write-Host "  • Async notifications" -ForegroundColor Gray
Write-Host "  • No blocking operations" -ForegroundColor Gray
Write-Host ""
Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Yellow
Write-Host ""

# Run the server
& go run main.go

