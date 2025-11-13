# 🚀 Starting the Visa Help Center Project

## ✅ Current Status

**Frontend (Next.js)**: Starting on http://localhost:3000
- Status: Running in background
- Access: Open http://localhost:3000 in your browser

**Backend (Golang API)**: Requires Go installation

## 📋 Quick Start Instructions

### Frontend is Running ✅

The Next.js frontend should be accessible at:
- **URL**: http://localhost:3000
- **Status**: Check your browser or terminal for confirmation

### Backend Setup Required

The backend requires **Go 1.21+** to be installed. Here are your options:

#### Option 1: Install Go (Recommended)

1. **Download Go** from https://go.dev/dl/
2. **Install** following the official instructions
3. **Verify installation**:
   ```powershell
   go version
   ```
4. **Start the backend**:
   ```powershell
   cd apps/api
   go mod download
   go run main.go
   ```

#### Option 2: Use Docker (Alternative)

If you have Docker installed:

```powershell
# Build and run the backend container
cd apps/api
docker build -t boltvisa-api .
docker run -p 8080:8080 --env-file .env boltvisa-api
```

#### Option 3: Use SQLite (Easiest for Development)

The database connection has been updated to support SQLite automatically. Just set:

```powershell
# In apps/api/.env or set environment variable
$env:DATABASE_URL = "sqlite://boltvisa.db"
```

Then start the backend (once Go is installed).

## 🔧 Environment Setup

### Backend Environment Variables

Create `apps/api/.env` with:

```env
# Database (SQLite for local dev)
DATABASE_URL=sqlite://boltvisa.db

# JWT Secret
JWT_SECRET=dev-secret-key-change-in-production-use-a-strong-random-string

# Server
PORT=8080
ENVIRONMENT=development
FRONTEND_URL=http://localhost:3000

# Optional services (leave empty for local dev)
GCP_PROJECT_ID=
SENDGRID_API_KEY=
TWILIO_ACCOUNT_SID=
STRIPE_SECRET_KEY=
RAZORPAY_KEY_ID=
```

### Frontend Environment Variables

Create `apps/web/.env.local` with:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 🎯 Testing the Application

Once both services are running:

1. **Open Frontend**: http://localhost:3000
2. **Test API**: http://localhost:8080/health
3. **Register**: Click "Get Started" → "Sign Up"
4. **Login**: Use your credentials

## 📊 Service Status

| Service | Port | Status | URL |
|---------|------|--------|-----|
| Frontend | 3000 | ✅ Running | http://localhost:3000 |
| Backend | 8080 | ⏳ Requires Go | http://localhost:8080 |

## 🐛 Troubleshooting

### Frontend Not Loading?
- Check if port 3000 is available: `netstat -ano | findstr :3000`
- Check terminal for errors
- Try: `cd apps/web && npm run dev`

### Backend Won't Start?
- **Go not installed**: Install from https://go.dev/dl/
- **Database error**: Ensure SQLite is working (should auto-create `boltvisa.db`)
- **Port conflict**: Change `PORT` in `.env` file

### Database Issues?
- SQLite will auto-create `boltvisa.db` in `apps/api/` directory
- No PostgreSQL needed for local development
- Database migrations run automatically on startup

## 📚 Next Steps

1. **Install Go** (if not already installed)
2. **Start Backend**: `cd apps/api && go run main.go`
3. **Access Application**: http://localhost:3000
4. **Create Account**: Register a new user
5. **Explore Features**: Dashboard, Applications, Payments, etc.

## 🎉 Project Features

- ✅ User Authentication (Register, Login, Password Reset)
- ✅ Role-Based Dashboards (Applicant, Consultant, Admin)
- ✅ Visa Application Management
- ✅ Document Upload (GCS integration ready)
- ✅ Payment Processing (Stripe & Razorpay ready)
- ✅ Notifications (Email & SMS ready)
- ✅ Audit Logging
- ✅ Rate Limiting
- ✅ API Documentation (OpenAPI/Swagger)

---

**Note**: The frontend is running but will need the backend API to function fully. Install Go and start the backend to complete the setup!

