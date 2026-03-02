# 网址导航 - 个人收藏管理系统

一个功能完整的个人网址收藏管理网站，支持用户认证、网址分类管理、搜索筛选和数据导出导入功能。

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
- SQLite数据库

## 功能特性

- 用户注册和登录
- 网址增删改查
- 分类管理
- 搜索和筛选
- 访问统计
- 数据导出/导入
- 响应式设计

## 快速开始

### 本地开发

#### 后端启动

```bash
cd backend
go mod download
go run main.go
```

后端服务将在 http://localhost:8080 启动

#### 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 http://localhost:3000 启动

### Docker部署

```bash
docker-compose up -d
```

服务将在以下端口启动：
- 前端: http://localhost:80
- 后端: http://localhost:8080

## API文档

### 认证接口

- POST `/api/auth/register` - 用户注册
- POST `/api/auth/login` - 用户登录

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

- GET `/api/export` - 导出数据
- POST `/api/import` - 导入数据

## 数据库结构

### 用户表 (users)
- id, username, email, password, created_at, updated_at

### 分类表 (categories)
- id, user_id, name, description, icon, sort_order, created_at, updated_at

### 网址表 (bookmarks)
- id, user_id, category_id, title, url, description, icon, sort_order, visit_count, created_at, updated_at

## 环境变量

后端支持以下环境变量：

- `PORT` - 服务端口 (默认: 8080)
- `DB_PATH` - 数据库文件路径 (默认: ./nav.db)
- `JWT_SECRET` - JWT密钥 (生产环境请修改)

## 项目结构

```
nav/
├── backend/           # Go后端
│   ├── main.go
│   ├── database/      # 数据库配置
│   ├── handlers/      # API处理器
│   ├── middleware/    # 中间件
│   ├── models/        # 数据模型
│   └── utils/         # 工具函数
├── frontend/          # Vue前端
│   ├── src/
│   │   ├── views/     # 页面组件
│   │   ├── stores/    # Pinia状态
│   │   ├── router/    # 路由配置
│   │   └── utils/     # 工具函数
│   └── index.html
└── docker-compose.yml # Docker编排
```

## 安全建议

1. 生产环境请修改JWT_SECRET
2. 使用HTTPS部署
3. 定期备份数据库文件
4. 限制API访问频率

## 许可证

MIT
