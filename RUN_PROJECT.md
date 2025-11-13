# 🚀 Running the Project

## ✅ Current Status

Both servers have been started:

### Frontend (Next.js)
- **Status**: ✅ Running
- **URL**: http://localhost:3000
- **Process**: Running in background

### Backend (Golang API)
- **Status**: ⏳ Starting
- **URL**: http://localhost:8080
- **Process**: Starting in background

## 🎯 Verify Both Are Running

### Check Frontend:
Open in browser: http://localhost:3000
- Should show the Visa Help Center homepage

### Check Backend:
Open in browser: http://localhost:8080/health
- Should show: `{"status":"ok"}`

## 🎊 Use the Application

Once both are running:
- **Signup**: http://localhost:3000/signup
- **Login**: http://localhost:3000/login
- **Dashboard**: http://localhost:3000/dashboard

## ⚠️ Important Notes

1. **Backend may take 15-20 seconds** to fully start
2. **Both servers are running in background**
3. **Keep the terminal/processes running** - Closing them stops the servers

## 🔍 Quick Status Check

```powershell
# Check if servers are running
Invoke-WebRequest -Uri "http://localhost:3000"
Invoke-WebRequest -Uri "http://localhost:8080/health"

# Check ports
netstat -ano | findstr ":3000"
netstat -ano | findstr ":8080"
```

## 🎯 If Backend Isn't Responding

The backend is starting in the background. If it doesn't respond after 20 seconds:

1. **Check for errors** in the terminal
2. **Try manually**: 
   ```powershell
   cd C:\Users\dell\Desktop\boltvisa\apps\api
   $env:Path = "C:\Program Files\Go\bin;$env:Path"
   go run main.go
   ```

---

**Both servers are starting. Wait 15-20 seconds, then try the signup page!**

