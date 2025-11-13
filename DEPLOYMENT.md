# 🚀 Deployment Guide

## Prerequisites

- Google Cloud Platform account
- `gcloud` CLI installed and configured
- Docker installed (for container builds)
- PostgreSQL database (Cloud SQL recommended)

## Local Development Setup

### 1. Install Dependencies

```bash
npm install
```

### 2. Setup Database

```bash
# Create PostgreSQL database
createdb boltvisa

# Or using Docker
docker run --name boltvisa-postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=boltvisa \
  -p 5432:5432 \
  -d postgres:14
```

### 3. Configure Environment Variables

```bash
# Copy example files
cp .env.example .env
cp apps/api/.env.example apps/api/.env
cp apps/web/.env.example apps/web/.env

# Edit .env files with your configuration
```

### 4. Run Backend

```bash
cd apps/api
go mod download
go mod tidy
go run main.go
```

Backend will be available at `http://localhost:8080`

### 5. Run Frontend

```bash
npm run dev --workspace=apps/web
```

Frontend will be available at `http://localhost:3000`

## Google Cloud Platform Deployment

### 1. Setup GCP Project

```bash
# Set your project ID
export GCP_PROJECT_ID=your-project-id
gcloud config set project $GCP_PROJECT_ID

# Enable required APIs
gcloud services enable \
  cloudbuild.googleapis.com \
  run.googleapis.com \
  sqladmin.googleapis.com \
  storage-api.googleapis.com \
  pubsub.googleapis.com
```

### 2. Create Cloud SQL Database

```bash
# Create Cloud SQL instance
gcloud sql instances create boltvisa-db \
  --database-version=POSTGRES_14 \
  --tier=db-f1-micro \
  --region=us-central1

# Create database
gcloud sql databases create boltvisa --instance=boltvisa-db

# Create user
gcloud sql users create boltvisa-user \
  --instance=boltvisa-db \
  --password=YOUR_SECURE_PASSWORD
```

### 3. Create GCS Bucket

```bash
gsutil mb -p $GCP_PROJECT_ID -l us-central1 gs://boltvisa-documents
```

### 4. Create Pub/Sub Topic

```bash
gcloud pubsub topics create visa-notifications
```

### 5. Deploy Backend to Cloud Run

```bash
cd apps/api

# Build and deploy
gcloud run deploy boltvisa-api \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars DATABASE_URL="postgres://user:pass@/boltvisa?host=/cloudsql/PROJECT_ID:REGION:INSTANCE_NAME" \
  --set-env-vars JWT_SECRET="your-secret-key" \
  --set-env-vars GCP_PROJECT_ID="$GCP_PROJECT_ID" \
  --set-env-vars GCS_BUCKET_NAME="boltvisa-documents" \
  --add-cloudsql-instances PROJECT_ID:REGION:INSTANCE_NAME
```

### 6. Deploy Frontend to Cloud Run

```bash
cd apps/web

# Build
npm run build

# Deploy (using Cloud Run or Vercel)
# For Cloud Run:
gcloud run deploy boltvisa-web \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars NEXT_PUBLIC_API_URL="https://boltvisa-api-xxx.run.app"
```

## Docker Deployment

### Build Images

```bash
# Build backend
docker build -t boltvisa-api:latest -f apps/api/Dockerfile apps/api

# Build frontend
docker build -t boltvisa-web:latest -f apps/web/Dockerfile .
```

### Run with Docker Compose

Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      POSTGRES_DB: boltvisa
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  api:
    build:
      context: .
      dockerfile: apps/api/Dockerfile
    ports:
      - "8080:8080"
    environment:
      DATABASE_URL: postgres://user:password@postgres:5432/boltvisa?sslmode=disable
      JWT_SECRET: your-secret-key
      PORT: 8080
    depends_on:
      - postgres

  web:
    build:
      context: .
      dockerfile: apps/web/Dockerfile
    ports:
      - "3000:3000"
    environment:
      NEXT_PUBLIC_API_URL: http://api:8080
    depends_on:
      - api

volumes:
  postgres_data:
```

Run:
```bash
docker-compose up -d
```

## CI/CD with GitHub Actions

The CI/CD pipeline is configured in `.github/workflows/ci.yml`. It will:

1. Run linting and type checking
2. Build frontend and backend
3. Run tests (when implemented)
4. Deploy to preview environment on PRs
5. Deploy to production on merge to main

## Monitoring & Logging

### Stackdriver Logging

Logs are automatically sent to Stackdriver when deployed to Cloud Run.

View logs:
```bash
gcloud logging read "resource.type=cloud_run_revision" --limit 50
```

### Health Checks

Backend health endpoint:
```bash
curl https://your-api-url.run.app/health
```

## Troubleshooting

### Database Connection Issues

- Verify Cloud SQL instance is running
- Check connection string format
- Ensure Cloud Run service account has Cloud SQL Client role

### CORS Issues

- Verify `FRONTEND_URL` environment variable matches your frontend domain
- Check CORS middleware configuration

### Build Failures

- Ensure all dependencies are installed: `npm ci`
- Check Go module dependencies: `cd apps/api && go mod tidy`
- Verify Node.js and Go versions match requirements

