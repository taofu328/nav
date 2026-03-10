<#
导航网站后端构建脚本
支持编译Windows 64位和Linux ARM平台的可执行程序
#>

# 设置颜色
$GREEN = "`e[92m"
$YELLOW = "`e[93m"
$RED = "`e[91m"
$NC = "`e[0m"

# 项目根目录
$PROJECT_ROOT = Split-Path -Parent $MyInvocation.MyCommand.Definition
$BACKEND_DIR = Join-Path $PROJECT_ROOT "backend"
$OUTPUT_DIR = Join-Path $PROJECT_ROOT "build"

# 默认编译平台
$BUILD_WINDOWS = $true
$BUILD_LINUX_ARM = $true

# 显示帮助信息
function Show-Help {
    Write-Host "Usage: $($MyInvocation.MyCommand.Name) [options]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -h, --help     显示此帮助信息"
    Write-Host "  -w, --windows  只编译Windows平台"
    Write-Host "  -l, --linux    只编译Linux ARM平台"
    Write-Host "  -o, --output   指定输出目录 (默认: $OUTPUT_DIR)"
    Write-Host ""
    Write-Host "示例:"
    Write-Host "  $($MyInvocation.MyCommand.Name)                    # 编译两个平台"
    Write-Host "  $($MyInvocation.MyCommand.Name) --windows          # 只编译Windows平台"
    Write-Host "  $($MyInvocation.MyCommand.Name) --linux            # 只编译Linux ARM平台"
    Write-Host "  $($MyInvocation.MyCommand.Name) --output .\dist    # 指定输出目录为.\dist"
}

# 解析命令行参数
function Parse-Args {
    param(
        [string[]]$Args
    )
    
    $global:BUILD_WINDOWS = $true
    $global:BUILD_LINUX_ARM = $true
    $global:OUTPUT_DIR = Join-Path $PROJECT_ROOT "build"
    
    $i = 0
    while ($i -lt $Args.Length) {
        switch ($Args[$i]) {
            "-h" {
                Show-Help
                exit 0
            }
            "--help" {
                Show-Help
                exit 0
            }
            "-w" {
                $global:BUILD_WINDOWS = $true
                $global:BUILD_LINUX_ARM = $false
                $i++
            }
            "--windows" {
                $global:BUILD_WINDOWS = $true
                $global:BUILD_LINUX_ARM = $false
                $i++
            }
            "-l" {
                $global:BUILD_WINDOWS = $false
                $global:BUILD_LINUX_ARM = $true
                $i++
            }
            "--linux" {
                $global:BUILD_WINDOWS = $false
                $global:BUILD_LINUX_ARM = $true
                $i++
            }
            "-o" {
                if ($i + 1 -lt $Args.Length) {
                    $global:OUTPUT_DIR = $Args[$i + 1]
                    $i += 2
                } else {
                    Write-Host "$RED错误: 缺少输出目录参数$NC"
                    Show-Help
                    exit 1
                }
            }
            "--output" {
                if ($i + 1 -lt $Args.Length) {
                    $global:OUTPUT_DIR = $Args[$i + 1]
                    $i += 2
                } else {
                    Write-Host "$RED错误: 缺少输出目录参数$NC"
                    Show-Help
                    exit 1
                }
            }
            default {
                Write-Host "$RED错误: 未知参数 '$($Args[$i])'$NC"
                Show-Help
                exit 1
            }
        }
    }
}

# 初始化构建环境
function Init-Build {
    Write-Host "$GREEN初始化构建环境...$NC"
    
    # 创建输出目录
    if (-not (Test-Path $OUTPUT_DIR)) {
        New-Item -ItemType Directory -Path $OUTPUT_DIR -Force | Out-Null
    }
    
    # 进入后端目录
    Set-Location $BACKEND_DIR -ErrorAction Stop
    
    # 下载依赖
    Write-Host "$GREEN下载依赖...$NC"
    try {
        go mod tidy
    } catch {
        Write-Host "$RED错误: 依赖下载失败$NC"
        exit 1
    }
}

# 编译Windows平台
function Build-Windows {
    Write-Host "$GREEN编译Windows 64位平台...$NC"
    
    # 设置编译参数
    $env:CGO_ENABLED = "0"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    
    # 编译
    try {
        go build -o "$OUTPUT_DIR\nav-backend-windows-amd64.exe" .
    } catch {
        Write-Host "$RED错误: Windows平台编译失败$NC"
        return $false
    }
    
    # 验证输出文件
    if (Test-Path "$OUTPUT_DIR\nav-backend-windows-amd64.exe") {
        Write-Host "$GREENWindows平台编译成功: $OUTPUT_DIR\nav-backend-windows-amd64.exe$NC"
        return $true
    } else {
        Write-Host "$RED错误: Windows平台编译产物不存在$NC"
        return $false
    }
}

# 编译Linux ARM平台
function Build-Linux-ARM {
    Write-Host "$GREEN编译Linux ARM平台...$NC"
    
    # 编译ARM64
    Write-Host "$YELLOW编译Linux ARM64平台...$NC"
    $env:CGO_ENABLED = "0"
    $env:GOOS = "linux"
    $env:GOARCH = "arm64"
    
    # 编译
    try {
        go build -o "$OUTPUT_DIR\nav-backend-linux-arm64" .
    } catch {
        Write-Host "$RED错误: Linux ARM64平台编译失败$NC"
        return $false
    }
    
    # 验证输出文件
    if (Test-Path "$OUTPUT_DIR\nav-backend-linux-arm64") {
        Write-Host "$GREENLinux ARM64平台编译成功: $OUTPUT_DIR\nav-backend-linux-arm64$NC"
    } else {
        Write-Host "$RED错误: Linux ARM64平台编译产物不存在$NC"
        return $false
    }
    
    # 编译ARMv7
    Write-Host "$YELLOW编译Linux ARMv7平台...$NC"
    $env:CGO_ENABLED = "0"
    $env:GOOS = "linux"
    $env:GOARCH = "arm"
    $env:GOARM = "7"
    
    # 编译
    try {
        go build -o "$OUTPUT_DIR\nav-backend-linux-armv7" .
    } catch {
        Write-Host "$RED错误: Linux ARMv7平台编译失败$NC"
        return $false
    }
    
    # 验证输出文件
    if (Test-Path "$OUTPUT_DIR\nav-backend-linux-armv7") {
        Write-Host "$GREENLinux ARMv7平台编译成功: $OUTPUT_DIR\nav-backend-linux-armv7$NC"
        return $true
    } else {
        Write-Host "$RED错误: Linux ARMv7平台编译产物不存在$NC"
        return $false
    }
}

# 显示构建结果
function Show-Result {
    Write-Host "$GREEN
构建完成!$NC"
    Write-Host "$GREEN输出目录: $OUTPUT_DIR$NC"
    
    # 列出生成的文件
    Write-Host "$YELLOW
生成的文件:$NC"
    try {
        Get-ChildItem $OUTPUT_DIR | Format-Table -AutoSize
    } catch {
        Write-Host "$RED错误: 无法列出输出目录$NC"
    }
}

# 主函数
function Main {
    # 解析命令行参数
    Parse-Args -Args $args
    
    # 初始化构建环境
    Init-Build
    
    # 编译平台
    $windowsResult = $false
    $linuxResult = $false
    
    if ($BUILD_WINDOWS) {
        $windowsResult = Build-Windows
    }
    
    if ($BUILD_LINUX_ARM) {
        $linuxResult = Build-Linux-ARM
    }
    
    # 显示结果
    Show-Result
    
    # 检查构建结果
    if ($windowsResult -or $linuxResult) {
        Write-Host "$GREEN
构建成功!$NC"
        exit 0
    } else {
        Write-Host "$RED
构建失败!$NC"
        exit 1
    }
}

# 执行主函数
Main @args
