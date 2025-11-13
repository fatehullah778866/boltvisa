# Quick start script - one-liner setup
# Usage: .\scripts\start-simple.ps1

Write-Host "Quick Setup: Enabling Corepack and starting services..." -ForegroundColor Cyan

corepack enable
corepack prepare pnpm@9.12.0 --activate
pnpm install
pnpm run dev

