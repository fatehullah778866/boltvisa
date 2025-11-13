# Run this script as Administrator to install Go automatically
# Right-click PowerShell -> "Run as Administrator" -> Run this script

Write-Host "🚀 Installing Go Programming Language..." -ForegroundColor Cyan
Write-Host ""

# Check if running as admin
$isAdmin = ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)

if (-not $isAdmin) {
    Write-Host "❌ This script requires Administrator privileges!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please:" -ForegroundColor Yellow
    Write-Host "1. Right-click PowerShell" -ForegroundColor Yellow
    Write-Host "2. Select 'Run as Administrator'" -ForegroundColor Yellow
    Write-Host "3. Run this script again" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "Or install Go manually from: https://go.dev/dl/" -ForegroundColor Yellow
    pause
    exit 1
}

# Download Go installer
$goVersion = "1.21.6"
$goInstaller = "$env:TEMP\go-installer.msi"
$goUrl = "https://go.dev/dl/go${goVersion}.windows-amd64.msi"

Write-Host "📥 Downloading Go ${goVersion}..." -ForegroundColor Cyan
try {
    $ProgressPreference = 'SilentlyContinue'
    Invoke-WebRequest -Uri $goUrl -OutFile $goInstaller -UseBasicParsing
    Write-Host "✅ Download complete" -ForegroundColor Green
} catch {
    Write-Host "❌ Download failed: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host "Please download manually from: https://go.dev/dl/" -ForegroundColor Yellow
    pause
    exit 1
}

# Install Go silently
Write-Host ""
Write-Host "📦 Installing Go..." -ForegroundColor Cyan
try {
    $process = Start-Process msiexec.exe -ArgumentList "/i `"$goInstaller`" /quiet /norestart" -Wait -PassThru -NoNewWindow
    
    if ($process.ExitCode -eq 0) {
        Write-Host "✅ Go installed successfully!" -ForegroundColor Green
    } else {
        Write-Host "⚠️ Installation completed with exit code: $($process.ExitCode)" -ForegroundColor Yellow
    }
} catch {
    Write-Host "❌ Installation failed: $($_.Exception.Message)" -ForegroundColor Red
    pause
    exit 1
}

# Update PATH for current session
$goPath = "C:\Program Files\Go\bin"
if (Test-Path $goPath) {
    $env:Path = "$goPath;$env:Path"
    Write-Host ""
    Write-Host "✅ Go added to PATH for this session" -ForegroundColor Green
    
    # Verify installation
    Write-Host ""
    Write-Host "🔍 Verifying installation..." -ForegroundColor Cyan
    $goVersionOutput = & "$goPath\go.exe" version
    Write-Host "✅ $goVersionOutput" -ForegroundColor Green
    
    Write-Host ""
    Write-Host "🎉 Go installation complete!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Next steps:" -ForegroundColor Cyan
    Write-Host "1. Close and reopen your terminal (to refresh PATH)" -ForegroundColor Yellow
    Write-Host "2. Navigate to: cd C:\Users\dell\Desktop\boltvisa\apps\api" -ForegroundColor Yellow
    Write-Host "3. Run: go mod download" -ForegroundColor Yellow
    Write-Host "4. Run: go run main.go" -ForegroundColor Yellow
    Write-Host ""
} else {
    Write-Host "⚠️ Go may need a system restart to be available" -ForegroundColor Yellow
    Write-Host "Please restart your computer or install manually from: https://go.dev/dl/" -ForegroundColor Yellow
}

# Cleanup
Remove-Item $goInstaller -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "Press any key to exit..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")

