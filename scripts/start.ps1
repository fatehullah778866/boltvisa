Write-Host "========================================" -ForegroundColor Cyan
Write-Host " BoltVisa" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Starting Next.js on http://localhost:3000" -ForegroundColor Green
Write-Host "API URL: http://localhost:8080" -ForegroundColor Green
Write-Host ""
Write-Host "Compiling... Please wait." -ForegroundColor Yellow

# 1) Ensure Corepack is available
if (-not (Get-Command corepack -ErrorAction SilentlyContinue)) {
  Write-Host "ERROR: Node.js 18+ is required (Corepack missing)." -ForegroundColor Red
  exit 1
}

# 2) Activate the pinned pnpm version
corepack enable | Out-Null
corepack prepare pnpm@9.12.0 --activate | Out-Null

# 3) Install deps on first run
if (-not (Test-Path "./node_modules")) {
  Write-Host "Installing dependencies with pnpm..." -ForegroundColor Yellow
  pnpm install
}

# 4) Start both backend and frontend via workspace script
pnpm run dev

