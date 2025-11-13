# 🚀 Quick Start Backend - Simple Instructions

## ✅ Backend is NOT Running

The backend server needs to be started manually. Here's the easiest way:

## Method 1: Double-Click the Batch File (Easiest!)

1. **Navigate to**: `C:\Users\dell\Desktop\boltvisa\apps\api`
2. **Double-click**: `start-server.bat`
3. **Wait**: A command window will open and the server will start
4. **Look for**: "Server starting on port 8080"
5. **Keep the window open** - Don't close it!

## Method 2: Manual Start (If batch file doesn't work)

1. **Open PowerShell or Command Prompt**
2. **Run these commands**:
   ```powershell
   cd C:\Users\dell\Desktop\boltvisa\apps\api
   $env:Path = "C:\Program Files\Go\bin;$env:Path"
   go run main.go
   ```

3. **Wait for**: "Server starting on port 8080"
4. **Keep the terminal open**

## ✅ Verify Backend is Running

Open in your browser:
- **Health Check**: http://localhost:8080/health
- **Should show**: `{"status":"ok"}`

## 🎯 Then Try Signup Again

Once backend is running:
- **Signup Page**: http://localhost:3000/signup
- **Should work now!** ✅

## ⚠️ Important Notes

- **Keep the backend window open** - Closing it stops the server
- **Frontend is already running** at http://localhost:3000
- **Both need to run** at the same time

## 🔍 Troubleshooting

### Port 8080 already in use?
```powershell
# Find what's using port 8080
netstat -ano | findstr :8080
# Kill the process (replace PID with actual number)
taskkill /PID <PID> /F
```

### Go command not found?
- Make sure Go is installed: `go version`
- If not, it was installed at: `C:\Program Files\Go\bin\go.exe`
- Use full path: `& "C:\Program Files\Go\bin\go.exe" run main.go`

### Database errors?
- SQLite database will be created automatically
- File location: `apps/api/boltvisa.db`

---

**Quick Start**: Just double-click `apps/api/start-server.bat` and keep the window open!

