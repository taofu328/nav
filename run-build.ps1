<#
Navigation Backend Build Script Wrapper
Runs the build.bat script using CMD to avoid PowerShell encoding issues
#>

$scriptPath = Join-Path $PSScriptRoot "build.bat"
$workingDir = $PSScriptRoot

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Navigation Backend Build Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Run the batch file using CMD
$process = Start-Process -FilePath "cmd.exe" -ArgumentList "/c `"$scriptPath`"" -WorkingDirectory $workingDir -NoNewWindow -Wait -PassThru

# Check the exit code
if ($process.ExitCode -eq 0) {
    Write-Host ""
    Write-Host "Build completed successfully!" -ForegroundColor Green
    exit 0
} else {
    Write-Host ""
    Write-Host "Build failed with exit code: $($process.ExitCode)" -ForegroundColor Red
    exit $process.ExitCode
}
