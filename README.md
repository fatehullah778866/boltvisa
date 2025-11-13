# 🧭 Visa Help Center

A comprehensive visa application management system built with Next.js, Golang, and Google Cloud Platform.

## Quickstart

1. Install Node.js 18+  

2. In a terminal:

   ```bash
   npm run setup     # activates pnpm via Corepack
   npm run bootstrap # installs dependencies
   npm run dev       # runs backend and frontend
   ```

**Do not install pnpm globally; Corepack manages it for you.**

**Notes for Windows:** Run commands in PowerShell or VS Code terminal.

### Safety Fallback (if Corepack doesn't work)

If you must run without Corepack:

```bash
npm i -g pnpm@9.12.0
pnpm install
pnpm run dev
```

---

## 🏗️ Architecture

This is a monorepo containing:

- **Frontend**: Next.js application (`apps/web`)
- **Backend**: Golang API server (`apps/api`)
- **Shared Packages**: Common utilities, types, and UI components (`packages/*`)

### Prerequisites

- **Node.js ≥ 18** (includes Corepack)
- Go 1.21+
- PostgreSQL 14+ (or SQLite for local dev)
- Docker (optional)

### Setup Details

The project uses **pnpm** with **Corepack** for automatic package manager management. The `setup` script enables Corepack and activates pnpm@9.12.0 automatically.

#### Windows PowerShell Script

```powershell
# Run the automated setup script:
.\scripts\start.ps1
```

#### Safety Fallback (if Corepack doesn't work)

If you must run without Corepack (not recommended):

```bash
npm i -g pnpm@9.12.0
pnpm install
pnpm run dev
```

**Note:** This is an escape hatch. Prefer using Corepack via `npm run setup`.

### Development

Once setup is complete:

- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **API Docs**: http://localhost:8080/openapi.json

### Troubleshooting

**"pnpm: command not found"**
- Run `npm run setup` first to enable Corepack
- Or install pnpm globally: `npm install -g pnpm`

**"Unsupported URL Type workspace:"**
- This means npm was used instead of pnpm
- Use `pnpm install` instead of `npm install`

**Windows Issues**
- Use PowerShell (not CMD)
- Run as Administrator if Corepack fails
- Ensure Node.js 18+ is installed

### Development

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## 📁 Project Structure

```
boltvisa/
├── apps/
│   ├── web/          # Next.js frontend
│   └── api/          # Golang backend
├── packages/
│   ├── ui/           # Shared UI components
│   ├── utils/        # Shared utilities
│   ├── types/        # TypeScript types
│   └── config/       # Shared configuration
├── turbo.json        # Turborepo configuration
└── package.json      # Root package.json
```

## 🛠️ Tech Stack

- **Frontend**: Next.js 14, React, TypeScript, Tailwind CSS
- **Backend**: Golang, Gin, GORM
- **Database**: PostgreSQL (Cloud SQL)
- **Storage**: Google Cloud Storage
- **Auth**: JWT with password reset
- **Notifications**: Pub/Sub, SendGrid (Email), Twilio (SMS)
- **Payments**: Stripe, Razorpay
- **CI/CD**: GitHub Actions, Cloud Build
- **Deployment**: Cloud Run / GKE
- **Monitoring**: Custom metrics, logging, audit trails
- **Testing**: Unit tests, E2E (Playwright), Load testing (k6)

## 📚 Documentation

- [Architecture](./docs/ARCHITECTURE.md) - System design and architecture
- [API Documentation](./docs/API.md) - API endpoints and usage
- [Security Guide](./docs/SECURITY.md) - Security features and best practices
- [Testing Guide](./docs/TESTING.md) - Testing instructions
- [Monitoring Guide](./docs/MONITORING.md) - Monitoring and observability
- [Deployment Guide](./DEPLOYMENT.md) - Deployment instructions
- [Production Guide](./README_PRODUCTION.md) - Production deployment checklist

## 🔐 Security Features

- ✅ Rate limiting (100 req/15min anonymous, 200 req/15min authenticated)
- ✅ JWT authentication with refresh tokens
- ✅ Password reset flow with secure tokens
- ✅ Audit logging for sensitive operations
- ✅ Role-based access control (RBAC)
- ✅ Input validation and SQL injection prevention
- ✅ Error boundaries and graceful error handling

## 🚀 Quick Links

- **API Spec**: `/openapi.json` or `/swagger.json`
- **Health Check**: `/health`
- **Metrics**: `/metrics`
- **Load Testing**: `k6 run load-test.js`

## 📝 License

MIT

