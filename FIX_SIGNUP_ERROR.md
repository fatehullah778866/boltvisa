# 🔧 Fix: "Failed to fetch" Error on Signup

## Problem
The signup page shows "Failed to fetch" because the **backend API server is not running**.

## Solution: Start the Backend Server

### Option 1: Quick Start (If Go is Installed)

1. **Open a new terminal/PowerShell window**

2. **Navigate to the API directory**:
   ```powershell
   cd C:\Users\dell\Desktop\boltvisa\apps\api
   ```

3. **Run the startup script**:
   ```powershell
   .\start-backend.ps1
   ```
   
   Or if using Command Prompt:
   ```cmd
   start-backend.bat
   ```

4. **Wait for the server to start** - You should see:
   ```
   ✅ Database connection established
   ✅ Database migrations completed
   🚀 Server starting on port 8080
   ```

5. **Keep this terminal open** - The server needs to keep running

### Option 2: Manual Start (If Go is Installed)

```powershell
cd C:\Users\dell\Desktop\boltvisa\apps\api
go mod download
go run main.go
```

### Option 3: Install Go First (If Go is NOT Installed)

1. **Download Go**: https://go.dev/dl/
   - Choose Windows installer (`.msi` file)
   - Download Go 1.21 or later

2. **Install Go**:
   - Run the installer
   - Follow the installation wizard
   - **Important**: Restart your terminal after installation

3. **Verify installation**:
   ```powershell
   go version
   ```
   Should show: `go version go1.21.x windows/amd64` (or similar)

4. **Then follow Option 1 or 2 above**

## ✅ Verify Backend is Running

Once started, you should be able to access:
- **Health Check**: http://localhost:8080/health
- **API Docs**: http://localhost:8080/openapi.json

Open these URLs in your browser - they should return JSON data.

## 🎯 Test Signup Again

1. **Backend is running** ✅ (check terminal)
2. **Frontend is running** ✅ (http://localhost:3000)
3. **Go to**: http://localhost:3000/signup
4. **Fill in the form** and click "Sign Up"
5. **Should work now!** ✅

## 🔍 Troubleshooting

### Error: "Go is not recognized"
- **Solution**: Install Go from https://go.dev/dl/
- **After installing**: Close and reopen your terminal

### Error: "Port 8080 already in use"
- **Solution**: Another process is using port 8080
- **Fix**: 
  ```powershell
  # Find what's using port 8080
  netstat -ano | findstr :8080
  # Kill the process (replace PID with actual process ID)
  taskkill /PID <PID> /F
  ```

### Error: "Cannot connect to database"
- **Solution**: The app uses SQLite by default (no setup needed)
- **Check**: `boltvisa.db` file should be created in `apps/api/` folder

### Still Getting "Failed to fetch"?
1. **Check backend is running**: Open http://localhost:8080/health
2. **Check browser console**: Press F12 → Console tab
3. **Check CORS**: Backend should allow requests from http://localhost:3000
4. **Verify API URL**: Frontend should use `http://localhost:8080`

## 📝 Quick Checklist

- [ ] Go is installed (`go version` works)
- [ ] Backend server is running (check terminal)
- [ ] Backend responds at http://localhost:8080/health
- [ ] Frontend is running at http://localhost:3000
- [ ] Try signup again

## 💡 Pro Tip

Keep **two terminal windows open**:
1. **Terminal 1**: Backend server (`cd apps/api && go run main.go`)
2. **Terminal 2**: Frontend server (`cd apps/web && npm run dev`)

Both need to run simultaneously!

---

**Need more help?** Check `START_PROJECT.md` for detailed setup instructions.

