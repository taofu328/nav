@echo off
setlocal

echo ========================================
echo Navigation Backend Build Script
echo ========================================
echo.

REM Set paths
set "BACKEND_DIR=%~dp0backend"
set "OUTPUT_DIR=%~dp0build"

REM Create output directory
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"

REM Go to backend directory
cd "%BACKEND_DIR%"

REM Download dependencies
echo [1/4] Downloading dependencies...
go mod tidy
if errorlevel 1 (
    echo Error: Failed to download dependencies
    pause
    exit /b 1
)

REM Build Windows platform
echo [2/4] Building Windows 64-bit platform...
set "CGO_ENABLED=0"
set "GOOS=windows"
set "GOARCH=amd64"
go build -o "%OUTPUT_DIR%\nav-backend-windows-amd64.exe" .
if errorlevel 1 (
    echo Error: Windows platform build failed
    pause
    exit /b 1
)

REM Build Linux ARM64 platform
echo [3/4] Building Linux ARM64 platform...
set "CGO_ENABLED=0"
set "GOOS=linux"
set "GOARCH=arm64"
go build -o "%OUTPUT_DIR%\nav-backend-linux-arm64" .
if errorlevel 1 (
    echo Error: Linux ARM64 platform build failed
    pause
    exit /b 1
)

REM Build Linux ARMv7 platform
echo [4/4] Building Linux ARMv7 platform...
set "CGO_ENABLED=0"
set "GOOS=linux"
set "GOARCH=arm"
set "GOARM=7"
go build -o "%OUTPUT_DIR%\nav-backend-linux-armv7" .
if errorlevel 1 (
    echo Error: Linux ARMv7 platform build failed
    pause
    exit /b 1
)

echo.
echo ========================================
echo Build Summary
echo ========================================
echo Output directory: %OUTPUT_DIR%
echo.
echo Generated files:
dir "%OUTPUT_DIR%" /b

echo.
echo Build SUCCESS!
pause
