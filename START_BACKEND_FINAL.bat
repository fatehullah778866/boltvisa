@echo off
REM ========================================
REM   BoltVisa Backend API - Final Startup
REM ========================================
title BoltVisa Backend API Server
color 0A
cls

echo.
echo ========================================
echo   BoltVisa Backend API Server
echo ========================================
echo.

REM Set Go path
set "GOROOT=C:\Program Files\Go"
set "PATH=%GOROOT%\bin;%PATH%"

REM Navigate to API directory
cd /d "%~dp0apps\api"
if not exist "main.go" (
    echo [ERROR] Cannot find main.go
    echo Current directory: %CD%
    pause
    exit /b 1
)

REM Check Go installation
"%GOROOT%\bin\go.exe" version >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Go is not installed or not in PATH
    echo Please install Go from: https://go.dev/dl/
    pause
    exit /b 1
)

echo [OK] Go found
"%GOROOT%\bin\go.exe" version
echo.

REM Check/create .env
if not exist ".env" (
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
echo [INFO] Keep this window open - closing it will stop the server
echo.
echo ========================================
echo.

REM Start server
"%GOROOT%\bin\go.exe" run main.go

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [ERROR] Server failed to start!
    echo Check the error messages above
    echo.
    pause
)

