# Setup Guide

This guide helps you set up the BoltVisa project on Windows, macOS, or Linux.

## Prerequisites

- **Node.js ≥ 18** (includes Corepack)
- **Go 1.21+** (for backend)
- **Git**

## Quick Setup

### Windows (PowerShell)

```powershell
# Option 1: Automated script
.\scripts\start.ps1

# Option 2: Manual commands
npm run setup
npm run bootstrap
npm run dev
```

### macOS / Linux

```bash
# Setup Corepack and pnpm
npm run setup

# Install dependencies
npm run bootstrap

# Start services
npm run dev
```

## What Each Command Does

### `npm run setup`
- Enables Corepack (Node.js package manager manager)
- Activates pnpm@9.12.0 automatically
- No manual installation needed!

### `npm run bootstrap`
- Installs all dependencies for the monorepo
- Links workspace packages
- Sets up both frontend and backend dependencies

### `npm run dev`
- Starts both backend (Go) and frontend (Next.js)
- Backend runs on port 8080
- Frontend runs on port 3000

## Individual Service Commands

If you want to run services separately:

```bash
# Backend only
pnpm run dev:api

# Frontend only
pnpm run dev:web
```

## Troubleshooting

### Issue: "pnpm: command not found"

**Solution:**
```bash
npm run setup
```

This enables Corepack and activates pnpm automatically.

### Issue: "Unsupported URL Type workspace:"

**Cause:** You used `npm install` instead of `pnpm install`

**Solution:**
```bash
# Remove node_modules and reinstall with pnpm
rm -rf node_modules apps/*/node_modules packages/*/node_modules
pnpm install
```

### Issue: Corepack not available

**Cause:** Node.js version is too old (< 16.13)

**Solution:**
1. Upgrade to Node.js 18+ from https://nodejs.org/
2. Or install pnpm manually: `npm install -g pnpm`

### Issue: Port already in use

**Solution:**
```bash
# Kill processes on ports 3000 and 8080
# Windows:
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# macOS/Linux:
lsof -ti:3000 | xargs kill
lsof -ti:8080 | xargs kill
```

## Environment Variables

Create `.env` files if needed:

- `apps/api/.env` - Backend configuration
- `apps/web/.env.local` - Frontend configuration

See `.env.example` files for reference.

## Verification

After setup, verify everything works:

1. **Backend Health**: http://localhost:8080/health
2. **Frontend**: http://localhost:3000
3. **API Docs**: http://localhost:8080/openapi.json

## Next Steps

- Read [README.md](./README.md) for project overview
- Check [docs/ARCHITECTURE.md](./docs/ARCHITECTURE.md) for system design
- See [docs/API.md](./docs/API.md) for API documentation

