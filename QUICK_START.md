# Quick Start Guide

## One-Command Setup (After First Time)

```bash
npm run dev
```

This will automatically:
1. Check if Corepack/pnpm is available
2. Install dependencies if needed
3. Start both backend and frontend

## First-Time Setup

### Windows (PowerShell)

```powershell
# Option 1: Automated script (recommended)
.\scripts\start.ps1

# Option 2: Manual setup
npm run setup
npm run bootstrap
npm run dev
```

**Note:** If `corepack enable` fails with permission error, you can:
- Run PowerShell as Administrator, OR
- Skip Corepack and install pnpm directly: `npm install -g pnpm`

### macOS / Linux

```bash
npm run setup
npm run bootstrap
npm run dev
```

## Troubleshooting

### "pnpm: command not found"

**Quick Fix:**
```bash
npm install -g pnpm
```

Then continue with:
```bash
pnpm install
pnpm run dev
```

### "Unsupported URL Type workspace:"

**Cause:** Used `npm install` instead of `pnpm install`

**Fix:**
```bash
# Remove and reinstall with pnpm
rm -rf node_modules apps/*/node_modules packages/*/node_modules
pnpm install
```

### Corepack Permission Error (Windows)

**Solution 1:** Run as Administrator
```powershell
# Right-click PowerShell -> Run as Administrator
npm run setup
```

**Solution 2:** Skip Corepack
```bash
npm install -g pnpm
pnpm install
pnpm run dev
```

## Verify Installation

After setup, check:

1. **Backend**: http://localhost:8080/health
2. **Frontend**: http://localhost:3000
3. **API Docs**: http://localhost:8080/openapi.json

## What's Running?

- **Backend (Go)**: Port 8080
- **Frontend (Next.js)**: Port 3000

Both start automatically with `npm run dev` or `pnpm run dev`.

