@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

REM Navigation Backend Build Script
REM Supports building for Windows and Linux ARM platforms

set "PROJECT_ROOT=%~dp0"
set "BACKEND_DIR=%PROJECT_ROOT%backend"
set "FRONTEND_DIR=%PROJECT_ROOT%frontend"
set "OUTPUT_DIR=%PROJECT_ROOT%build"
set "BUILD_WINDOWS=true"
set "BUILD_LINUX_AMD64=true"
set "BUILD_LINUX_ARM=true"

REM Parse command line arguments
:parse_args
if "%~1"=="" goto :end_parse_args
if "%~1"=="-h" goto :help
if "%~1"=="--help" goto :help
if "%~1"=="-w" goto :windows_only
if "%~1"=="--windows" goto :windows_only
if "%~1"=="-l" goto :linux_only
if "%~1"=="--linux" goto :linux_only
if "%~1"=="-o" goto :set_output
if "%~1"=="--output" goto :set_output
echo Error: Unknown parameter '%~1'
goto :help

:help
echo Usage: %~n0 [options]
echo.
echo Options:
echo   -h, --help     Show this help message
echo   -w, --windows  Build only Windows platform
echo   -l, --linux    Build only Linux platform (AMD64 and ARM)
echo   -o, --output   Specify output directory (default: %OUTPUT_DIR%)
echo.
echo Examples:
echo   %~n0                    # Build all platforms
echo   %~n0 --windows          # Build only Windows platform
echo   %~n0 --linux            # Build only Linux platform (AMD64 and ARM)
echo   %~n0 --output .\dist    # Specify output directory as .\dist
exit /b 1

:windows_only
set "BUILD_WINDOWS=true"
set "BUILD_LINUX_AMD64=false"
set "BUILD_LINUX_ARM=false"
shift
goto :parse_args

:linux_only
set "BUILD_WINDOWS=false"
set "BUILD_LINUX_AMD64=true"
set "BUILD_LINUX_ARM=true"
shift
goto :parse_args

:set_output
if "%~2"=="" (
    echo Error: Missing output directory parameter
    goto :help
)
set "OUTPUT_DIR=%~2"
shift
shift
goto :parse_args

:end_parse_args

echo ========================================
echo Navigation Backend Build Script
echo ========================================
echo.

REM Initialize build environment
echo [1/4] Initializing build environment...
if not exist "%OUTPUT_DIR%" mkdir "%OUTPUT_DIR%"
cd "%BACKEND_DIR%"
if errorlevel 1 (
    echo Error: Cannot enter backend directory
    exit /b 1
)

echo [2/4] Building frontend...
cd "%FRONTEND_DIR%"
call npm run build
if errorlevel 1 (
    echo Error: Failed to build frontend
    exit /b 1
)

echo [3/4] Copying frontend dist to backend...
if not exist "%BACKEND_DIR%\dist" mkdir "%BACKEND_DIR%\dist"
xcopy /E /I /Y "%FRONTEND_DIR%\dist\*" "%BACKEND_DIR%\dist\"
if errorlevel 1 (
    echo Error: Failed to copy frontend files
    exit /b 1
)

cd "%BACKEND_DIR%"
echo [4/4] Downloading dependencies...
go mod tidy
if errorlevel 1 (
    echo Error: Failed to download dependencies
    exit /b 1
)

REM Build Windows platform
if "%BUILD_WINDOWS%"=="true" (
    echo [6/7] Building Windows 64-bit platform...
    set "CGO_ENABLED=0"
    set "GOOS=windows"
    set "GOARCH=amd64"
    go build -o "%OUTPUT_DIR%\nav-backend-windows-amd64.exe" .
    if errorlevel 1 (
        echo Error: Windows platform build failed
        set "WINDOWS_RESULT=1"
    ) else (
        if exist "%OUTPUT_DIR%\nav-backend-windows-amd64.exe" (
            echo Success: Windows platform built successfully
            set "WINDOWS_RESULT=0"
        ) else (
            echo Error: Windows build output not found
            set "WINDOWS_RESULT=1"
        )
    )
)

REM Build Linux AMD64 platform
if "%BUILD_LINUX_AMD64%"=="true" (
    echo [6/7] Building Linux AMD64 platform...
    set "CGO_ENABLED=0"
    set "GOOS=linux"
    set "GOARCH=amd64"
    go build -o "%OUTPUT_DIR%\nav-backend-linux-amd64" .
    if errorlevel 1 (
        echo Error: Linux AMD64 platform build failed
        set "LINUX_RESULT=1"
    ) else (
        if exist "%OUTPUT_DIR%\nav-backend-linux-amd64" (
            echo Success: Linux AMD64 platform built successfully
            set "LINUX_RESULT=0"
        ) else (
            echo Error: Linux AMD64 build output not found
            set "LINUX_RESULT=1"
        )
    )
)

REM Build Linux ARM platform
if "%BUILD_LINUX_ARM%"=="true" (
    echo [7/7] Building Linux ARM platform...
    
    echo Building Linux ARM64...
    set "CGO_ENABLED=0"
    set "GOOS=linux"
    set "GOARCH=arm64"
    go build -o "%OUTPUT_DIR%\nav-backend-linux-arm64" .
    if errorlevel 1 (
        echo Error: Linux ARM64 platform build failed
        set "LINUX_RESULT=1"
    ) else (
        if exist "%OUTPUT_DIR%\nav-backend-linux-arm64" (
            echo Success: Linux ARM64 platform built successfully
            set "LINUX_RESULT=0"
        ) else (
            echo Error: Linux ARM64 build output not found
            set "LINUX_RESULT=1"
        )
    )
    
    echo Building Linux ARMv7...
    set "CGO_ENABLED=0"
    set "GOOS=linux"
    set "GOARCH=arm"
    set "GOARM=7"
    go build -o "%OUTPUT_DIR%\nav-backend-linux-armv7" .
    if errorlevel 1 (
        echo Error: Linux ARMv7 platform build failed
        set "LINUX_RESULT=1"
    ) else (
        if exist "%OUTPUT_DIR%\nav-backend-linux-armv7" (
            echo Success: Linux ARMv7 platform built successfully
        ) else (
            echo Error: Linux ARMv7 build output not found
            set "LINUX_RESULT=1"
        )
    )
)

REM Show results
echo [7/7] Build completed!
echo.
echo ========================================
echo Build Summary
echo ========================================
echo Output directory: %OUTPUT_DIR%
echo.
echo Generated files:
dir "%OUTPUT_DIR%" /b

REM Check build results
if "%BUILD_WINDOWS%"=="true" (
    if "!WINDOWS_RESULT!"=="0" (
        echo.
        echo Build SUCCESS!
        exit /b 0
    )
)
if "%BUILD_LINUX_AMD64%"=="true" (
    if "!LINUX_RESULT!"=="0" (
        echo.
        echo Build SUCCESS!
        exit /b 0
    )
)
if "%BUILD_LINUX_ARM%"=="true" (
    if "!LINUX_RESULT!"=="0" (
        echo.
        echo Build SUCCESS!
        exit /b 0
    )
)

echo.
echo Build FAILED!
exit /b 1
