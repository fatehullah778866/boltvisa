# Test script to verify backend can start
Write-Host "Testing backend startup..." -ForegroundColor Cyan

$env:Path = "C:\Program Files\Go\bin;$env:Path"
cd C:\Users\dell\Desktop\boltvisa\apps\api

Write-Host "Checking Go..." -ForegroundColor Yellow
& "C:\Program Files\Go\bin\go.exe" version

Write-Host "Checking .env..." -ForegroundColor Yellow
if (Test-Path .env) {
    Write-Host ".env exists" -ForegroundColor Green
} else {
    Write-Host ".env missing - creating..." -ForegroundColor Yellow
    @"
DATABASE_URL=sqlite://boltvisa.db
JWT_SECRET=dev-secret-key-change-in-production
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000
"@ | Out-File -FilePath ".env" -Encoding utf8
}

Write-Host "Building..." -ForegroundColor Yellow
& "C:\Program Files\Go\bin\go.exe" build -o test-build.exe .

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    Remove-Item test-build.exe -ErrorAction SilentlyContinue
    Write-Host ""
    Write-Host "Starting server..." -ForegroundColor Cyan
    Write-Host "Server will run in this window. Press Ctrl+C to stop." -ForegroundColor Yellow
    Write-Host ""
    & "C:\Program Files\Go\bin\go.exe" run main.go
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    pause
}

