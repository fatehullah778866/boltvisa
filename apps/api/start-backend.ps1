# PowerShell script to start the backend server
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  BoltVisa Backend API Server" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Add Go to PATH
$env:Path = "C:\Program Files\Go\bin;$env:Path"

# Navigate to API directory
$scriptPath = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $scriptPath

# Check if Go is available
$goCmd = "C:\Program Files\Go\bin\go.exe"
if (-not (Test-Path $goCmd)) {
    Write-Host "[ERROR] Go is not installed at: $goCmd" -ForegroundColor Red
    Write-Host "Please install Go from: https://go.dev/dl/" -ForegroundColor Yellow
    pause
    exit 1
}

# Check if .env exists
if (-not (Test-Path ".env")) {
    Write-Host "[INFO] Creating .env file..." -ForegroundColor Yellow
    @"
DATABASE_URL=sqlite://boltvisa.db
JWT_SECRET=dev-secret-key-change-in-production
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000
"@ | Out-File -FilePath ".env" -Encoding utf8
    Write-Host "[OK] .env file created" -ForegroundColor Green
    Write-Host ""
}

Write-Host "[INFO] Starting server on http://localhost:8080" -ForegroundColor Cyan
Write-Host "[INFO] Press Ctrl+C to stop the server" -ForegroundColor Gray
Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Start the server
& $goCmd run main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "[ERROR] Server failed to start!" -ForegroundColor Red
    Write-Host ""
    pause
}
