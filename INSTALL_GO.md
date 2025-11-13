# 🚀 Install Go to Run the Backend

## ❌ Current Status

**Go is NOT installed** on your system. The backend requires Go 1.21+ to run.

## 📥 Quick Installation Guide

### Step 1: Download Go

1. **Visit**: https://go.dev/dl/
2. **Download**: Click on the Windows installer (`.msi` file)
   - Choose the latest version (Go 1.21 or later)
   - File will be something like: `go1.21.x.windows-amd64.msi`

### Step 2: Install Go

1. **Run the installer** (`go1.21.x.windows-amd64.msi`)
2. **Follow the installation wizard**:
   - Click "Next" through the prompts
   - Accept the license agreement
   - Choose installation location (default is fine: `C:\Program Files\Go`)
   - Click "Install"
   - Wait for installation to complete
   - Click "Finish"

### Step 3: Verify Installation

1. **Close ALL terminal/PowerShell windows** (important!)
2. **Open a NEW terminal/PowerShell window**
3. **Verify Go is installed**:
   ```powershell
   go version
   ```
   
   Should show: `go version go1.21.x windows/amd64` (or similar)

### Step 4: Start the Backend

Once Go is installed:

```powershell
cd C:\Users\dell\Desktop\boltvisa\apps\api
go mod download
go run main.go
```

## 🎯 Alternative: Use Docker (If You Have Docker)

If you have Docker installed, you can run the backend without installing Go:

```powershell
cd C:\Users\dell\Desktop\boltvisa\apps\api
docker build -t boltvisa-api .
docker run -p 8080:8080 boltvisa-api
```

## ✅ After Installing Go

Once Go is installed and verified:

1. **Navigate to API directory**:
   ```powershell
   cd C:\Users\dell\Desktop\boltvisa\apps\api
   ```

2. **Download dependencies**:
   ```powershell
   go mod download
   ```

3. **Start the server**:
   ```powershell
   go run main.go
   ```

4. **You should see**:
   ```
   Database connection established
   Database migrations completed
   Server starting on port 8080
   ```

5. **Keep the terminal open** - server must keep running

6. **Test**: Open http://localhost:8080/health in browser

## 🔍 Troubleshooting

### "go: command not found" after installation
- **Solution**: Close and reopen your terminal
- **Check**: Restart your computer if needed
- **Verify**: Run `go version` in a new terminal

### Installation fails
- **Check**: You have administrator rights
- **Try**: Right-click installer → "Run as administrator"
- **Check**: Windows Defender isn't blocking the installer

### Still having issues?
- **Check Go installation**: `where.exe go` (should show path)
- **Check PATH**: Go should be in your system PATH
- **Reinstall**: Uninstall and reinstall Go

## 📝 Quick Checklist

- [ ] Downloaded Go installer from https://go.dev/dl/
- [ ] Installed Go successfully
- [ ] Closed and reopened terminal
- [ ] Verified: `go version` works
- [ ] Navigated to `apps/api` directory
- [ ] Ran `go mod download`
- [ ] Started server: `go run main.go`
- [ ] Backend running on http://localhost:8080

---

**Once Go is installed, come back and I'll help you start the backend!**

