# 网址导航 - 个人收藏管理系统

一个功能完整的个人网址收藏管理网站，支持网址分类管理、搜索筛选、数据导出导入功能和单用户管理模式。

## 技术栈

### 前端
- Vue 3 + Vite
- Element Plus UI组件库
- Pinia状态管理
- Vue Router路由
- Axios HTTP客户端
- Tailwind CSS响应式设计

### 后端
- Go 1.21
- Gin Web框架
- GORM ORM
- JWT认证
- bcrypt密码加密
- SQLite数据库 (modernc.org/sqlite纯Go驱动)

## 功能特性

- 单用户管理模式（默认管理员账户）
- 网址增删改查
- 分类管理
- 排序功能（默认值99）
- 搜索和筛选
- 访问统计
- 数据导出/导入（JSON格式）
- 响应式设计
- 支持Windows和Linux ARM平台部署

## 快速开始

### 本地开发

#### 后端启动

```bash
cd backend
go mod tidy
go run main.go
```

后端服务将在 http://localhost:8081 启动（默认端口）

#### 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 http://localhost:3000 启动

### 构建可执行文件

项目提供了自动化构建脚本，支持编译多个平台的可执行文件：

#### Windows平台

```bash
# 使用简化的批处理脚本（推荐）
.\build-simple.bat

# 或使用完整功能脚本
.\build.bat
```

#### Linux/Mac平台

```bash
# 使用bash脚本
./build.sh
```

#### 手动构建

```bash
cd backend

# 编译Windows 64位
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o nav-backend-windows-amd64.exe .

# 编译Linux ARM64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o nav-backend-linux-arm64 .

# 编译Linux ARMv7
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -o nav-backend-linux-armv7 .
```

编译后的文件将存储在 `build/` 目录中。

### Docker部署

#### 构建Docker镜像

```bash
docker build -t nav-app .
```

#### 使用Docker Compose部署

```bash
docker-compose up -d
```

服务将在以下端口启动：
- 前端: http://localhost:80
- 后端: http://localhost:8081

#### 数据持久化

Docker容器使用volume映射实现数据持久化：
- `./data` 目录映射到容器的 `/app/data` 目录
- 数据库文件：`./data/nav.db`
- 图标文件：`./data/icons/`

## API文档

### 认证接口

- POST `/api/auth/login` - 用户登录
- POST `/api/admin/login` - 管理员登录

### 公开接口

- GET `/api/public/categories` - 获取公开分类列表
- GET `/api/public/bookmarks` - 获取公开网址列表
- GET `/api/icons/:filename` - 获取图标文件

### 分类管理

- GET `/api/categories` - 获取分类列表
- POST `/api/categories` - 创建分类
- PUT `/api/categories/:id` - 更新分类
- DELETE `/api/categories/:id` - 删除分类

### 网址管理

- GET `/api/bookmarks` - 获取网址列表
- POST `/api/bookmarks` - 创建网址
- PUT `/api/bookmarks/:id` - 更新网址
- DELETE `/api/bookmarks/:id` - 删除网址
- POST `/api/bookmarks/:id/visit` - 增加访问次数

### 数据管理

- GET `/api/export` - 导出数据（JSON格式）
- POST `/api/import` - 导入数据（JSON格式）

### 管理员接口

- GET `/api/admin/settings` - 获取站点设置
- PUT `/api/admin/settings` - 更新站点设置
- PUT `/api/admin/user` - 更新用户信息
- DELETE `/api/admin/clear-all` - 清空所有数据
- POST `/api/admin/update-icons` - 更新所有图标

### 图标管理

- GET `/api/icons` - 获取图标列表
- POST `/api/icons/upload` - 上传图标
- DELETE `/api/icons` - 删除图标

## 数据库结构

### 用户表 (users)
- id, username, email, password, created_at, updated_at

### 分类表 (categories)
- id, name, description, icon, sort_order, created_at, updated_at

### 网址表 (bookmarks)
- id, category_id, title, url, description, icon, sort_order, visit_count, created_at, updated_at

### 站点设置表 (site_settings)
- id, site_name, site_description, site_keywords, created_at, updated_at

## 环境变量

后端支持以下环境变量：

- `PORT` - 服务端口 (默认: 8081)
- `DB_PATH` - 数据库文件路径 (默认: ./data/nav.db)
- `JWT_SECRET` - JWT密钥 (生产环境请修改)

## 命令行参数

后端支持以下命令行参数：

- `-port` - 指定Web服务端口（默认：8081）

示例：
```bash
# 使用自定义端口启动
./nav-backend -port 9000
```

## 项目结构

```
nav/
├── backend/           # Go后端
│   ├── main.go
│   ├── config/       # 配置文件
│   ├── database/     # 数据库配置
│   ├── handlers/     # API处理器
│   ├── middleware/   # 中间件
│   ├── models/       # 数据模型
│   └── utils/        # 工具函数
├── frontend/        # Vue前端
│   ├── src/
│   │   ├── views/   # 页面组件
│   │   ├── stores/  # Pinia状态
│   │   ├── router/  # 路由配置
│   │   └── utils/   # 工具函数
│   └── index.html
├── data/           # 数据目录
│   ├── nav.db      # SQLite数据库文件
│   └── icons/     # 图标文件
├── build/          # 构建输出目录
├── Dockerfile       # Docker镜像构建文件
├── docker-compose.yml # Docker编排文件
├── build.bat       # Windows构建脚本
├── build.sh        # Linux/Mac构建脚本
└── build-simple.bat # 简化Windows构建脚本
```

## 默认管理员账户

系统首次启动时会自动创建默认管理员账户：

- 用户名：`admin`
- 密码：`admin`

**重要：** 生产环境请立即修改默认密码！

## 构建脚本说明

### build.bat (Windows)
- 支持命令行参数
- 可选择编译特定平台
- 可指定输出目录

参数：
- `-h, --help` - 显示帮助信息
- `-w, --windows` - 只编译Windows平台
- `-l, --linux` - 只编译Linux ARM平台
- `-o, --output` - 指定输出目录

### build-simple.bat (Windows)
- 简化版本，适合快速构建
- 编译所有平台
- 无需参数

### build.sh (Linux/Mac)
- 支持命令行参数
- 功能与build.bat相同

## 安全建议

1. 生产环境请修改默认管理员密码
2. 修改JWT_SECRET环境变量
3. 使用HTTPS部署
4. 定期备份 `data/nav.db` 数据库文件
5. 限制API访问频率
6. 使用防火墙限制端口访问

## 部署注意事项

### 单用户模式
- 系统采用单用户模式，无需注册功能
- 所有数据归属于默认管理员账户
- 移除了基于user_id的权限控制

### 数据持久化
- 数据库文件和图标文件存储在 `data/` 目录
- Docker部署时请确保volume映射正确
- 建议定期备份 `data/` 目录

### 跨平台部署
- 使用纯Go SQLite驱动，无需CGO
- 支持Windows、Linux ARM64、Linux ARMv7平台
- 编译后的可执行文件可直接运行，无需额外依赖

## 常见问题

### 如何修改默认端口？
在启动后端时使用 `-port` 参数：
```bash
./nav-backend -port 9000
```

### 如何备份数据？
直接复制 `data/nav.db` 文件即可。

### 如何迁移数据？
1. 导出数据：使用管理界面的导出功能
2. 在新系统中导入数据：使用导入功能
3. 或直接复制 `data/` 目录到新系统

### Docker构建失败怎么办？
确保Docker版本支持多阶段构建，并且网络连接正常。如果遇到SQLite编译错误，请确认使用的是纯Go驱动。

## 许可证

MIT
