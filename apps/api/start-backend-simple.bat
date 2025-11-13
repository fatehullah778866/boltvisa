@echo off
title BoltVisa Backend Server
color 0A

echo ========================================
echo   BoltVisa Backend API Server
echo ========================================
echo.

REM Add Go to PATH
set PATH=C:\Program Files\Go\bin;%PATH%

REM Change to API directory
cd /d %~dp0

REM Check if .env exists
if not exist .env (
    echo Creating .env file...
    (
        echo DATABASE_URL=sqlite://boltvisa.db
        echo JWT_SECRET=dev-secret-key-change-in-production
        echo PORT=8080
        echo ENVIRONMENT=development
        echo FRONTEND_URL=http://localhost:3000
    ) > .env
    echo .env file created.
    echo.
)

echo Starting server on http://localhost:8080
echo.
echo Press Ctrl+C to stop the server
echo.

REM Pure Go SQLite driver - no CGO needed
REM Run the server
go run main.go

pause

