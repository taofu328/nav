# 部署指南

## Docker部署（推荐）

### 前置要求
- Docker
- Docker Compose

### 部署步骤

1. 克隆项目
```bash
git clone <repository-url>
cd nav
```

2. 修改配置（可选）
编辑 `docker-compose.yml`，修改JWT_SECRET等环境变量

3. 启动服务
```bash
docker-compose up -d
```

4. 查看日志
```bash
docker-compose logs -f
```

5. 停止服务
```bash
docker-compose down
```

### 数据持久化

数据存储在Docker卷 `nav-data` 中，即使容器删除也不会丢失数据。

## 本地部署

### 后端部署

1. 安装Go 1.21+
2. 进入backend目录
```bash
cd backend
go mod download
go build -o nav-backend main.go
```

3. 运行
```bash
./nav-backend
```

### 前端部署

1. 安装Node.js 18+
2. 进入frontend目录
```bash
cd frontend
npm install
npm run build
```

3. 使用Nginx部署

将 `dist` 目录内容部署到Nginx的web根目录，配置反向代理到后端API。

## 云服务部署

### 部署到VPS

1. 安装Docker和Docker Compose
2. 上传项目文件
3. 运行 `docker-compose up -d`
4. 配置域名和SSL证书

### 使用Nginx反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;
    
    location / {
        proxy_pass http://localhost:80;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 备份与恢复

### 备份数据

```bash
docker cp nav-backend:/data/nav.db ./backup-$(date +%Y%m%d).db
```

### 恢复数据

```bash
docker cp ./backup.db nav-backend:/data/nav.db
docker-compose restart backend
```

## 监控与日志

查看容器状态：
```bash
docker-compose ps
```

查看日志：
```bash
docker-compose logs backend
docker-compose logs frontend
```

## 性能优化

1. 启用Gzip压缩
2. 配置CDN加速静态资源
3. 使用Redis缓存热点数据
4. 数据库索引优化

## 故障排查

### 前端无法访问后端
检查docker-compose.yml中的网络配置

### 数据库连接失败
检查DB_PATH环境变量和文件权限

### JWT验证失败
确认JWT_SECRET配置一致
