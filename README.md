# CYI-Note 智能笔记应用

CYI-Note 是一个现代化的笔记应用，支持 Markdown 编辑、标签管理、附件上传和 AI 辅助功能。

## 功能特点

- 用户注册和登录
- 笔记创建、编辑和删除
- Markdown 编辑支持
- 标签分类管理
- 附件上传和管理
- 全文搜索
- AI 自动标签推荐
- AI 内容摘要生成

## 技术栈

### 前端

- Vue.js
- Tailwind CSS
- Vue Router
- Axios

### 后端

- Go
- Gin Web 框架
- GORM ORM 框架
- JWT 认证

### 数据库

- 默认: MySQL
- 也支持: PostgreSQL, SQLite

## 快速开始

### 后端

1. 安装 Go 和依赖项

```bash
cd backend
go mod download
```

2. 复制配置文件并编辑

```bash
cp .env.example .env
# 编辑 .env 文件设置数据库和其他配置
```

3. 运行后端服务

```bash
go run main.go
```

后端服务将在 http://localhost:8080 运行。

### 前端

1. 安装依赖项

```bash
cd fronted
npm install
```

2. 运行开发服务器

```bash
npm run serve
```

前端开发服务器将在 http://localhost:3000 运行。

## API 文档

### 认证 API

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/auth/user` - 获取当前用户信息

### 笔记 API

- `GET /api/notes` - 获取笔记列表
- `POST /api/notes` - 创建笔记
- `GET /api/notes/:id` - 获取笔记详情
- `PUT /api/notes/:id` - 更新笔记
- `DELETE /api/notes/:id` - 删除笔记
- `GET /api/notes/search` - 搜索笔记

### 标签 API

- `GET /api/tags` - 获取标签列表
- `POST /api/tags` - 创建标签
- `DELETE /api/tags/:id` - 删除标签
- `POST /api/tags/:id/notes/:noteId` - 给笔记添加标签
- `DELETE /api/tags/:id/notes/:noteId` - 从笔记中移除标签

### 附件 API

- `POST /api/attachments` - 上传附件
- `GET /api/attachments/:id` - 获取附件
- `DELETE /api/attachments/:id` - 删除附件

### AI API

- `POST /api/ai/tags` - 生成标签推荐
- `POST /api/ai/summary` - 生成内容摘要

## 部署

### 服务器部署

1. 编译后端

```bash
cd backend
go build -o cyi-note
```

2. 编译前端

```bash
cd fronted
npm run build
```

3. 配置 Nginx 或其他 Web 服务器

4. 使用 systemd 或其他工具管理后端服务

### Docker 部署

```bash
docker-compose up -d
```

## 许可证

MIT 