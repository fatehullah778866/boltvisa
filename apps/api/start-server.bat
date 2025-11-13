@echo off
title BoltVisa Backend API Server
color 0A
echo.
echo ========================================
echo   BoltVisa Backend API Server
echo ========================================
echo.

REM Add Go to PATH
set PATH=C:\Program Files\Go\bin;%PATH%

REM Navigate to API directory
cd /d "%~dp0"

REM Check if Go is available
where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [WARNING] Go is not in PATH. Using full path...
    set GOCMD="C:\Program Files\Go\bin\go.exe"
) else (
    set GOCMD=go
)

REM Check if .env exists
if not exist .env (
    echo [INFO] Creating .env file...
    (
        echo DATABASE_URL=sqlite://boltvisa.db
        echo JWT_SECRET=dev-secret-key-change-in-production
        echo PORT=8080
        echo ENVIRONMENT=development
        echo FRONTEND_URL=http://localhost:3000
    ) > .env
    echo [OK] .env file created
    echo.
)

echo [INFO] Starting server on http://localhost:8080
echo [INFO] Press Ctrl+C to stop the server
echo.
echo ========================================
echo.

REM Start the server
%GOCMD% run main.go

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Server failed to start!
    echo.
    pause
)

