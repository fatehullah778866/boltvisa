# 🏗️ Visa Help Center - Architecture Documentation

## Overview

The Visa Help Center is a monorepo-based application built with Next.js (frontend) and Golang (backend), designed to manage visa applications with role-based access control.

## System Architecture

```
┌─────────────────┐
│   Next.js App   │  (Frontend - Port 3000)
│   (React/TS)    │
└────────┬────────┘
         │ HTTP/REST
         │
┌────────▼────────┐
│  Golang API     │  (Backend - Port 8080)
│  (Gin/GORM)     │
└────────┬────────┘
         │
    ┌────┴────┬──────────┬──────────────┐
    │         │          │              │
┌───▼───┐ ┌──▼───┐ ┌────▼────┐ ┌──────▼──────┐
│PostgreSQL│ │ Redis │ │   GCS   │ │  Pub/Sub   │
│(Cloud SQL)│ │(Cache)│ │(Storage)│ │(Messages) │
└──────────┘ └───────┘ └─────────┘ └────────────┘
```

## Monorepo Structure

```
boltvisa/
├── apps/
│   ├── web/              # Next.js frontend application
│   │   ├── src/
│   │   │   ├── app/      # Next.js 14 App Router
│   │   │   ├── components/
│   │   │   └── lib/
│   │   └── package.json
│   │
│   └── api/              # Golang backend API
│       ├── internal/
│       │   ├── config/   # Configuration management
│       │   ├── database/ # Database connection & migrations
│       │   ├── handlers/ # HTTP request handlers
│       │   ├── middleware/# Middleware (auth, CORS, etc.)
│       │   ├── models/   # GORM models
│       │   ├── router/   # Route definitions
│       │   └── utils/    # Utility functions
│       ├── main.go
│       └── go.mod
│
├── packages/
│   ├── types/            # Shared TypeScript types
│   ├── utils/            # Shared utility functions
│   ├── ui/               # Shared React components
│   └── config/           # Shared configuration
│
├── .github/workflows/    # CI/CD pipelines
├── docs/                 # Documentation
└── turbo.json            # Turborepo configuration
```

## Database Schema

### Users Table
- `id` (PK)
- `email` (unique)
- `password` (hashed)
- `first_name`, `last_name`
- `role` (admin, consultant, applicant)
- `active` (boolean)
- `created_at`, `updated_at`, `deleted_at`

### Visa Categories Table
- `id` (PK)
- `name`, `description`
- `country`
- `duration`, `price`
- `active` (boolean)
- `created_at`, `updated_at`

### Visa Applications Table
- `id` (PK)
- `user_id` (FK → users)
- `consultant_id` (FK → users, nullable)
- `category_id` (FK → visa_categories)
- `status` (draft, submitted, in_review, approved, rejected, cancelled)
- `passport_number`, `date_of_birth`, `nationality`, `travel_date`
- `notes` (text)
- `created_at`, `updated_at`, `deleted_at`

### Documents Table
- `id` (PK)
- `application_id` (FK → visa_applications)
- `type` (passport, photo, bank_statement, invitation, other)
- `name`, `file_name`
- `gcs_url` (Google Cloud Storage URL)
- `size`, `mime_type`
- `uploaded_by` (FK → users)
- `created_at`, `updated_at`, `deleted_at`

### Notifications Table
- `id` (PK)
- `user_id` (FK → users)
- `type` (application_update, document_request, payment, system)
- `title`, `message` (text)
- `read` (boolean)
- `read_at` (timestamp, nullable)
- `created_at`, `updated_at`

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user
- `POST /api/v1/auth/refresh` - Refresh JWT token

### Users
- `GET /api/v1/users/me` - Get current user (protected)
- `PUT /api/v1/users/me` - Update current user (protected)
- `GET /api/v1/admin/users` - List all users (admin only)
- `PUT /api/v1/admin/users/:id` - Update user (admin only)

### Visa Categories
- `GET /api/v1/visa-categories` - List all active categories
- `GET /api/v1/visa-categories/:id` - Get category details
- `POST /api/v1/admin/visa-categories` - Create category (admin only)
- `PUT /api/v1/admin/visa-categories/:id` - Update category (admin only)

### Visa Applications
- `GET /api/v1/applications` - List applications (role-based filtering)
- `POST /api/v1/applications` - Create new application
- `GET /api/v1/applications/:id` - Get application details
- `PUT /api/v1/applications/:id` - Update application
- `DELETE /api/v1/applications/:id` - Delete application (draft only)

### Documents
- `POST /api/v1/applications/:id/documents` - Upload document
- `GET /api/v1/applications/:id/documents` - List documents
- `DELETE /api/v1/documents/:id` - Delete document

### Notifications
- `GET /api/v1/notifications` - List user notifications
- `PUT /api/v1/notifications/:id/read` - Mark notification as read

### Consultant Routes
- `GET /api/v1/consultant/clients` - List consultant's clients
- `GET /api/v1/consultant/applications` - List consultant's applications

## Authentication & Authorization

### JWT Token Structure
```json
{
  "user_id": 123,
  "email": "user@example.com",
  "role": "applicant",
  "exp": 1234567890,
  "iat": 1234567890,
  "iss": "boltvisa-api"
}
```

### Role-Based Access Control
- **Applicant**: Can manage own applications and documents
- **Consultant**: Can view and manage assigned client applications
- **Admin**: Full system access, can manage users and visa categories

## Technology Stack

### Frontend
- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State Management**: React Context / Zustand (future)
- **HTTP Client**: Fetch API (via utils package)

### Backend
- **Language**: Golang 1.21+
- **Framework**: Gin
- **ORM**: GORM
- **Database**: PostgreSQL (Cloud SQL)
- **Authentication**: JWT (golang-jwt/jwt/v5)
- **Password Hashing**: bcrypt

### Infrastructure
- **Monorepo**: Turborepo
- **CI/CD**: GitHub Actions
- **Container Registry**: Google Container Registry
- **Deployment**: Cloud Run / GKE
- **Storage**: Google Cloud Storage
- **Messaging**: Cloud Pub/Sub
- **Monitoring**: Stackdriver / Prometheus + Grafana

## Development Workflow

1. **Local Development**
   ```bash
   npm install
   npm run dev  # Starts both frontend and backend
   ```

2. **Database Setup**
   ```bash
   # Create PostgreSQL database
   createdb boltvisa
   
   # Set DATABASE_URL in .env
   # Migrations run automatically on API startup
   ```

3. **Testing**
   ```bash
   npm run test        # Run all tests
   npm run lint        # Lint all packages
   npm run type-check  # TypeScript type checking
   ```

## Deployment

### Environment Variables
See `.env.example` for required environment variables.

### Cloud Run Deployment
```bash
# Build and deploy backend
gcloud run deploy boltvisa-api --source apps/api

# Build and deploy frontend
gcloud run deploy boltvisa-web --source apps/web
```

### GKE Deployment
Kubernetes manifests should be created in `k8s/` directory (future).

## Security Considerations

1. **Password Security**: Passwords are hashed using bcrypt
2. **JWT Tokens**: Short-lived tokens (24h) with refresh capability
3. **CORS**: Configured for production domains
4. **Input Validation**: All inputs validated using Gin's binding
5. **SQL Injection**: Prevented by GORM's parameterized queries
6. **File Uploads**: Validated file types and sizes before GCS upload

## Future Enhancements

- [ ] Redis caching for frequently accessed data
- [ ] Rate limiting middleware
- [ ] WebSocket support for real-time updates
- [ ] GraphQL API option
- [ ] Comprehensive test coverage
- [ ] API documentation with Swagger/OpenAPI
- [ ] Multi-language support (i18n)
- [ ] Advanced analytics dashboard

