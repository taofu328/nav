# 第一阶段：构建前端
FROM --platform=$BUILDPLATFORM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# 复制前端 package.json 和 package-lock.json
COPY frontend/package*.json ./

# 安装前端依赖
RUN npm install

# 复制前端源代码
COPY frontend/ .

# 构建前端
RUN npm run build

# 第二阶段：构建后端
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

# 复制后端 go.mod 和 go.sum
COPY backend/go.mod backend/go.sum ./

# 下载后端依赖
RUN go mod download

# 复制后端源代码
COPY backend/ .

# 构建后端
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o nav-backend .

# 第三阶段：最终镜像
FROM --platform=linux/arm64 alpine:latest

WORKDIR /app

# 安装必要的依赖
RUN apk add --no-cache nginx supervisor

# 复制后端构建产物
COPY --from=backend-builder /app/backend/nav-backend .
# 复制前端构建产物
COPY --from=frontend-builder /app/frontend/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

# 创建数据目录
RUN mkdir -p /app/data/icons

# 创建 supervisor 配置
RUN echo '[supervisord]\nnodaemon=true\n\n[program:backend]\ncommand=/app/nav-backend\nautostart=true\nautorestart=true\nstdout_logfile=/dev/stdout\nstderr_logfile=/dev/stderr\n\n[program:nginx]\ncommand=nginx -g "daemon off;"\nautostart=true\nautorestart=true\nstdout_logfile=/dev/stdout\nstderr_logfile=/dev/stderr' > /etc/supervisord.conf

# 暴露端口
EXPOSE 80

# 运行
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
