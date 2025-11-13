# 🚀 Quick Start Guide

Get the Visa Help Center up and running in minutes!

## Prerequisites

- **Node.js** 18+ and npm
- **Go** 1.21+
- **PostgreSQL** 14+ (or Docker)

## Step 1: Clone and Install

```bash
# Install root dependencies
npm install

# Install Go dependencies (if needed)
cd apps/api
go mod download
go mod tidy
cd ../..
```

## Step 2: Setup Database

### Option A: Local PostgreSQL

```bash
createdb boltvisa
```

### Option B: Docker

```bash
docker run --name boltvisa-postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=boltvisa \
  -p 5432:5432 \
  -d postgres:14
```

## Step 3: Configure Environment

Create `.env` files:

**Root `.env`:**
```env
DATABASE_URL=postgres://user:password@localhost:5432/boltvisa?sslmode=disable
JWT_SECRET=dev-secret-key-change-in-production
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000
NEXT_PUBLIC_API_URL=http://localhost:8080
```

**apps/api/.env:**
```env
DATABASE_URL=postgres://user:password@localhost:5432/boltvisa?sslmode=disable
JWT_SECRET=dev-secret-key-change-in-production
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000
```

**apps/web/.env:**
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## Step 4: Start Backend

```bash
cd apps/api
go run main.go
```

Backend will start on `http://localhost:8080`

## Step 5: Start Frontend

In a new terminal:

```bash
npm run dev --workspace=apps/web
```

Frontend will start on `http://localhost:3000`

## Step 6: Test the Application

1. Open `http://localhost:3000` in your browser
2. Click "Get Started"
3. Create an account or login
4. Explore the dashboard!

## API Testing

Test the API directly:

```bash
# Health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

## Project Structure

```
boltvisa/
├── apps/
│   ├── web/          # Next.js frontend (port 3000)
│   └── api/          # Golang backend (port 8080)
├── packages/
│   ├── types/        # Shared TypeScript types
│   ├── utils/        # Shared utilities
│   └── ui/           # Shared UI components
└── docs/             # Documentation
```

## Common Commands

```bash
# Development
npm run dev                    # Start all apps in dev mode
npm run dev --workspace=apps/web   # Start only frontend
cd apps/api && go run main.go      # Start only backend

# Building
npm run build                 # Build all apps
npm run build --workspace=apps/web # Build frontend only

# Linting
npm run lint                  # Lint all packages

# Testing
npm run test                  # Run all tests
```

## Troubleshooting

### Database Connection Error

- Verify PostgreSQL is running: `pg_isready`
- Check `DATABASE_URL` in `.env` files
- Ensure database exists: `createdb boltvisa`

### Port Already in Use

- Backend: Change `PORT` in `apps/api/.env`
- Frontend: Change port in `apps/web/package.json` scripts

### Module Not Found (Go)

```bash
cd apps/api
go mod download
go mod tidy
```

### Module Not Found (Node)

```bash
rm -rf node_modules package-lock.json
npm install
```

## Next Steps

- Read [ARCHITECTURE.md](./docs/ARCHITECTURE.md) for system design
- Check [API.md](./docs/API.md) for API documentation
- See [DEPLOYMENT.md](./DEPLOYMENT.md) for production deployment

## Need Help?

- Check the documentation in `docs/`
- Review error logs in the terminal
- Verify all environment variables are set correctly

