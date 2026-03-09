FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o nav-backend .

# 最终镜像
FROM --platform=linux/arm64 alpine:latest

WORKDIR /app

# 复制构建产物
COPY --from=builder /app/nav-backend .
# 复制数据目录
COPY --from=builder /app/data ./data

# 暴露端口
EXPOSE 8081

# 运行
CMD ["./nav-backend"]
