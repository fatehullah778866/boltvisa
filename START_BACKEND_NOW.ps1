# Quick script to start backend (run this AFTER Go is installed)
# This script will use the full path to go.exe if PATH isn't updated yet

Write-Host "🚀 Starting BoltVisa Backend API..." -ForegroundColor Cyan
Write-Host ""

$goPath = "C:\Program Files\Go\bin\go.exe"

# Check if Go is installed
if (-not (Test-Path $goPath)) {
    Write-Host "❌ Go is not installed!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install Go first:" -ForegroundColor Yellow
    Write-Host "1. Run install-go-admin.ps1 as Administrator" -ForegroundColor Yellow
    Write-Host "   OR" -ForegroundColor Yellow
    Write-Host "2. Download from https://go.dev/dl/ and install manually" -ForegroundColor Yellow
    Write-Host ""
    pause
    exit 1
}

# Navigate to API directory
$apiDir = "C:\Users\dell\Desktop\boltvisa\apps\api"
if (-not (Test-Path $apiDir)) {
    Write-Host "❌ API directory not found: $apiDir" -ForegroundColor Red
    pause
    exit 1
}

Set-Location $apiDir

# Check/create .env file
if (-not (Test-Path ".env")) {
    Write-Host "📝 Creating .env file..." -ForegroundColor Cyan
    @"
DATABASE_URL=sqlite://boltvisa.db
JWT_SECRET=dev-secret-key-change-in-production-$(Get-Random)
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000
"@ | Out-File -FilePath ".env" -Encoding utf8
    Write-Host "✅ .env file created" -ForegroundColor Green
}

# Download dependencies
Write-Host ""
Write-Host "📦 Downloading Go dependencies..." -ForegroundColor Cyan
& $goPath mod download
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Failed to download dependencies" -ForegroundColor Red
    pause
    exit 1
}
Write-Host "✅ Dependencies downloaded" -ForegroundColor Green

# Start server
Write-Host ""
Write-Host "🚀 Starting server on http://localhost:8080..." -ForegroundColor Cyan
Write-Host "   Press Ctrl+C to stop the server" -ForegroundColor Gray
Write-Host ""

& $goPath run main.go

