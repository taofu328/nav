#!/bin/bash

# 导航网站后端构建脚本
# 支持编译Windows 64位和Linux ARM平台的可执行程序

# 颜色定义
GREEN="\033[0;32m"
YELLOW="\033[1;33m"
RED="\033[0;31m"
NC="\033[0m" # No Color

# 项目根目录
PROJECT_ROOT="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"
OUTPUT_DIR="$PROJECT_ROOT/build"

# 默认编译平台
BUILD_WINDOWS=true
BUILD_LINUX_AMD64=true
BUILD_LINUX_ARM=true

# 显示帮助信息
show_help() {
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -h, --help     显示此帮助信息"
    echo "  -w, --windows  只编译Windows平台"
    echo "  -l, --linux    只编译Linux平台 (AMD64和ARM)"
    echo "  -o, --output   指定输出目录 (默认: $OUTPUT_DIR)"
    echo ""
    echo "示例:"
    echo "  $0                    # 编译所有平台"
    echo "  $0 --windows          # 只编译Windows平台"
    echo "  $0 --linux            # 只编译Linux平台 (AMD64和ARM)"
    echo "  $0 --output ./dist    # 指定输出目录为./dist"
}

# 解析命令行参数
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -w|--windows)
                BUILD_WINDOWS=true
                BUILD_LINUX_AMD64=false
                BUILD_LINUX_ARM=false
                shift
                ;;
            -l|--linux)
                BUILD_WINDOWS=false
                BUILD_LINUX_AMD64=true
                BUILD_LINUX_ARM=true
                shift
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            *)
                echo -e "${RED}错误: 未知参数 '$1'${NC}"
                show_help
                exit 1
                ;;
        esac
    done
}

# 初始化构建环境
init_build() {
    echo -e "${GREEN}[1/7] 初始化构建环境...${NC}"
    
    # 创建输出目录
    mkdir -p "$OUTPUT_DIR"
    
    # 构建前端
    echo -e "${GREEN}[2/7] 构建前端...${NC}"
    cd "$FRONTEND_DIR" || {
        echo -e "${RED}错误: 无法进入前端目录${NC}"
        exit 1
    }
    npm run build || {
        echo -e "${RED}错误: 前端构建失败${NC}"
        exit 1
    }
    
    # 复制前端dist到backend
    echo -e "${GREEN}[3/7] 复制前端dist到backend...${NC}"
    mkdir -p "$BACKEND_DIR/dist"
    cp -r "$FRONTEND_DIR/dist/"* "$BACKEND_DIR/dist/" || {
        echo -e "${RED}错误: 复制前端文件失败${NC}"
        exit 1
    }
    
    # 进入后端目录
    cd "$BACKEND_DIR" || {
        echo -e "${RED}错误: 无法进入后端目录${NC}"
        exit 1
    }
    
    # 下载依赖
    echo -e "${GREEN}[4/7] 下载依赖...${NC}"
    go mod tidy || {
        echo -e "${RED}错误: 依赖下载失败${NC}"
        exit 1
    }
}

# 编译Windows平台
build_windows() {
    echo -e "${GREEN}[5/7] 编译Windows 64位平台...${NC}"
    
    # 设置编译参数
    export CGO_ENABLED=0
    export GOOS=windows
    export GOARCH=amd64
    
    # 编译
    go build -o "$OUTPUT_DIR/nav-backend-windows-amd64.exe" . || {
        echo -e "${RED}错误: Windows平台编译失败${NC}"
        return 1
    }
    
    # 验证输出文件
    if [[ -f "$OUTPUT_DIR/nav-backend-windows-amd64.exe" ]]; then
        echo -e "${GREEN}Windows平台编译成功: $OUTPUT_DIR/nav-backend-windows-amd64.exe${NC}"
        return 0
    else
        echo -e "${RED}错误: Windows平台编译产物不存在${NC}"
        return 1
    fi
}

# 编译Linux AMD64平台
build_linux_amd64() {
    echo -e "${GREEN}[6/7] 编译Linux AMD64平台...${NC}"
    
    # 设置编译参数
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=amd64
    
    # 编译
    go build -o "$OUTPUT_DIR/nav-backend-linux-amd64" . || {
        echo -e "${RED}错误: Linux AMD64平台编译失败${NC}"
        return 1
    }
    
    # 验证输出文件
    if [[ -f "$OUTPUT_DIR/nav-backend-linux-amd64" ]]; then
        echo -e "${GREEN}Linux AMD64平台编译成功: $OUTPUT_DIR/nav-backend-linux-amd64${NC}"
        return 0
    else
        echo -e "${RED}错误: Linux AMD64平台编译产物不存在${NC}"
        return 1
    fi
}

# 编译Linux ARM平台
build_linux_arm() {
    echo -e "${GREEN}[7/7] 编译Linux ARM平台...${NC}"
    
    # 编译ARM64
    echo -e "${YELLOW}编译Linux ARM64平台...${NC}"
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=arm64
    
    # 编译
    go build -o "$OUTPUT_DIR/nav-backend-linux-arm64" . || {
        echo -e "${RED}错误: Linux ARM64平台编译失败${NC}"
        return 1
    }
    
    # 验证输出文件
    if [[ -f "$OUTPUT_DIR/nav-backend-linux-arm64" ]]; then
        echo -e "${GREEN}Linux ARM64平台编译成功: $OUTPUT_DIR/nav-backend-linux-arm64${NC}"
    else
        echo -e "${RED}错误: Linux ARM64平台编译产物不存在${NC}"
        return 1
    fi
    
    # 编译ARMv7
    echo -e "${YELLOW}编译Linux ARMv7平台...${NC}"
    export CGO_ENABLED=0
    export GOOS=linux
    export GOARCH=arm
    export GOARM=7
    
    # 编译
    go build -o "$OUTPUT_DIR/nav-backend-linux-armv7" . || {
        echo -e "${RED}错误: Linux ARMv7平台编译失败${NC}"
        return 1
    }
    
    # 验证输出文件
    if [[ -f "$OUTPUT_DIR/nav-backend-linux-armv7" ]]; then
        echo -e "${GREEN}Linux ARMv7平台编译成功: $OUTPUT_DIR/nav-backend-linux-armv7${NC}"
        return 0
    else
        echo -e "${RED}错误: Linux ARMv7平台编译产物不存在${NC}"
        return 1
    fi
}

# 显示构建结果
show_result() {
    echo -e "${GREEN}\n[7/7] 构建完成!${NC}"
    echo -e "${GREEN}输出目录: ${OUTPUT_DIR}${NC}"
    
    # 列出生成的文件
    echo -e "${YELLOW}\n生成的文件:${NC}"
    ls -la "$OUTPUT_DIR" || {
        echo -e "${RED}错误: 无法列出输出目录${NC}"
    }
}

# 主函数
main() {
    # 解析命令行参数
    parse_args "$@"
    
    # 初始化构建环境
    init_build
    
    # 编译平台
    if [[ "$BUILD_WINDOWS" == true ]]; then
        build_windows
        WINDOWS_RESULT=$?
    fi
    
    if [[ "$BUILD_LINUX_AMD64" == true ]]; then
        build_linux_amd64
        LINUX_AMD64_RESULT=$?
    fi
    
    if [[ "$BUILD_LINUX_ARM" == true ]]; then
        build_linux_arm
        LINUX_ARM_RESULT=$?
    fi
    
    # 显示结果
    show_result
    
    # 检查构建结果
    if [[ "$WINDOWS_RESULT" == 0 || "$LINUX_AMD64_RESULT" == 0 || "$LINUX_ARM_RESULT" == 0 ]]; then
        echo -e "${GREEN}\n构建成功!${NC}"
        exit 0
    else
        echo -e "${RED}\n构建失败!${NC}"
        exit 1
    fi
}

# 执行主函数
main "$@"
