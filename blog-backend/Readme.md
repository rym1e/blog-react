# 博客系统后端开发文档 (Go + MySQL)

## 1. 项目概述

本项目是一个博客系统后端API，使用Go语言开发，MySQL作为数据库，为前端React应用提供数据支持。提供用户管理、文章管理、评论管理等功能。

## 2. 技术栈

- **后端语言**: Go (Golang)
- **Web框架**: Gin
- **数据库**: MySQL
- **ORM**: GORM
- **认证**: JWT (JSON Web Tokens)
- **密码加密**: bcrypt
- **环境配置**: godotenv

## 3. 项目结构

```
blog-backend/
├── cmd/
│   └── main.go              # 应用入口
├── config/                  # 配置文件
│   └── config.go
├── internal/
│   ├── controllers/         # 控制器
│   ├── middleware/          # 中间件
│   ├── models/              # 数据模型
│   ├── routes/              # 路由
│   └── utils/               # 工具函数
├── pkg/                     # 第三方包封装
├── .env                     # 环境变量
├── go.mod                   # 依赖管理
├── go.sum                   # 依赖校验
└── README.md
```

## 4. 数据库设计

### 4.1 用户表 (users)

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### 4.2 文章表 (articles)

```sql
CREATE TABLE articles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    author_id INT NOT NULL,
    views INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);
```

### 4.3 评论表 (comments)

```sql
CREATE TABLE comments (
    id INT AUTO_INCREMENT PRIMARY KEY,
    content TEXT NOT NULL,
    article_id INT NOT NULL,
    author_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);
```

## 5. 数据模型 (GORM)

### 5.1 用户模型 (User)

```go
// internal/models/user.go
package models

import (
    "time"
)

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Username  string    `gorm:"size:50;not null;unique" json:"username"`
    Email     string    `gorm:"size:100;not null;unique" json:"email"`
    Password  string    `gorm:"size:255;not null" json:"password"`
    Avatar    string    `gorm:"size:255" json:"avatar"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Articles  []Article `gorm:"foreignKey:AuthorID" json:"articles"`
    Comments  []Comment `gorm:"foreignKey:AuthorID" json:"comments"`
}
```

### 5.2 文章模型 (Article)

```go
// internal/models/article.go
package models

import (
    "time"
)

type Article struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"size:200;not null" json:"title"`
    Content   string    `gorm:"type:text;not null" json:"content"`
    AuthorID  uint      `gorm:"not null" json:"author_id"`
    Author    User      `gorm:"foreignKey:AuthorID" json:"author"`
    Views     int       `gorm:"default:0" json:"views"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Comments  []Comment `gorm:"foreignKey:ArticleID" json:"comments"`
}
```

### 5.3 评论模型 (Comment)

```go
// internal/models/comment.go
package models

import (
    "time"
)

type Comment struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Content   string    `gorm:"type:text;not null" json:"content"`
    ArticleID uint      `gorm:"not null" json:"article_id"`
    Article   Article   `gorm:"foreignKey:ArticleID" json:"article"`
    AuthorID  uint      `gorm:"not null" json:"author_id"`
    Author    User      `gorm:"foreignKey:AuthorID" json:"author"`
    CreatedAt time.Time `json:"created_at"`
}
```

## 6. API 接口设计

### 6.1 认证相关接口

#### 用户注册
- **URL**: `/api/v1/auth/register`
- **Method**: `POST`
- **请求参数**:
```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```
- **响应**:
```json
{
  "success": true,
  "message": "注册成功",
  "data": {
    "token": "jwt_token"
  }
}
```

#### 用户登录
- **URL**: `/api/v1/auth/login`
- **Method**: `POST`
- **请求参数**:
```json
{
  "email": "string",
  "password": "string"
}
```
- **响应**:
```json
{
  "success": true,
  "message": "登录成功",
  "data": {
    "token": "jwt_token",
    "user": {
      "id": 1,
      "username": "string",
      "email": "string",
      "avatar": "string"
    }
  }
}
```

### 6.2 用户相关接口

#### 获取当前用户信息
- **URL**: `/api/v1/users/me`
- **Method**: `GET`
- **Headers**: `Authorization: Bearer <token>`
- **响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "string",
    "email": "string",
    "avatar": "string",
    "created_at": "2023-07-01T12:00:00Z"
  }
}
```

#### 更新用户信息
- **URL**: `/api/v1/users/me`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer <token>`
- **请求参数**:
```json
{
  "username": "string",
  "avatar": "string"
}
```
- **响应**:
```json
{
  "success": true,
  "message": "更新成功",
  "data": {
    "id": 1,
    "username": "string",
    "email": "string",
    "avatar": "string"
  }
}
```

### 6.3 文章相关接口

#### 获取文章列表
- **URL**: `/api/v1/articles`
- **Method**: `GET`
- **查询参数**:
    - `page`: 页码 (默认: 1)
    - `limit`: 每页数量 (默认: 10)
- **响应**:
```json
{
  "success": true,
  "data": {
    "articles": [
      {
        "id": 1,
        "title": "string",
        "content": "string",
        "author": {
          "id": 1,
          "username": "string"
        },
        "views": 128,
        "created_at": "2023-07-01T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 100,
      "total_pages": 10
    }
  }
}
```

#### 获取文章详情
- **URL**: `/api/v1/articles/:id`
- **Method**: `GET`
- **响应**:
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "string",
    "content": "string",
    "author": {
      "id": 1,
      "username": "string"
    },
    "views": 128,
    "created_at": "2023-07-01T12:00:00Z",
    "updated_at": "2023-07-01T12:00:00Z"
  }
}
```

#### 创建文章
- **URL**: `/api/v1/articles`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **请求参数**:
```json
{
  "title": "string",
  "content": "string"
}
```
- **响应**:
```json
{
  "success": true,
  "message": "创建成功",
  "data": {
    "id": 1,
    "title": "string",
    "content": "string",
    "author_id": 1,
    "views": 0,
    "created_at": "2023-07-01T12:00:00Z"
  }
}
```

#### 更新文章
- **URL**: `/api/v1/articles/:id`
- **Method**: `PUT`
- **Headers**: `Authorization: Bearer <token>`
- **请求参数**:
```json
{
  "title": "string",
  "content": "string"
}
```
- **响应**:
```json
{
  "success": true,
  "message": "更新成功",
  "data": {
    "id": 1,
    "title": "string",
    "content": "string",
    "author_id": 1,
    "views": 0,
    "updated_at": "2023-07-01T12:00:00Z"
  }
}
```

#### 删除文章
- **URL**: `/api/v1/articles/:id`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer <token>`
- **响应**:
```json
{
  "success": true,
  "message": "删除成功"
}
```

### 6.4 评论相关接口

#### 获取文章评论列表
- **URL**: `/api/v1/articles/:id/comments`
- **Method**: `GET`
- **查询参数**:
    - `page`: 页码 (默认: 1)
    - `limit`: 每页数量 (默认: 10)
- **响应**:
```json
{
  "success": true,
  "data": {
    "comments": [
      {
        "id": 1,
        "content": "string",
        "author": {
          "id": 1,
          "username": "string"
        },
        "created_at": "2023-07-01T12:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 100,
      "total_pages": 10
    }
  }
}
```

#### 发表评论
- **URL**: `/api/v1/articles/:id/comments`
- **Method**: `POST`
- **Headers**: `Authorization: Bearer <token>`
- **请求参数**:
```json
{
  "content": "string"
}
```
- **响应**:
```json
{
  "success": true,
  "message": "评论发表成功",
  "data": {
    "id": 1,
    "content": "string",
    "article_id": 1,
    "author_id": 1,
    "created_at": "2023-07-01T12:00:00Z"
  }
}
```

#### 删除评论
- **URL**: `/api/v1/comments/:id`
- **Method**: `DELETE`
- **Headers**: `Authorization: Bearer <token>`
- **响应**:
```json
{
  "success": true,
  "message": "删除成功"
}
```

## 7. 错误响应格式

所有错误响应遵循统一格式:

```json
{
  "success": false,
  "message": "错误信息",
  "error_code": "ERROR_CODE"
}
```

常见错误码:
- `INVALID_INPUT`: 输入参数无效
- `UNAUTHORIZED`: 未授权访问
- `FORBIDDEN`: 权限不足
- `NOT_FOUND`: 资源不存在
- `INTERNAL_ERROR`: 服务器内部错误

## 8. 环境变量配置

在 `.env` 文件中配置以下环境变量:

```env
# 服务器配置
PORT=8080
HOST=localhost

# 数据库配置
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=blog

# JWT配置
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRE=24h
```

## 9. 启动项目

1. 安装依赖:
```bash
go mod tidy
```

2. 设置数据库:
```sql
CREATE DATABASE blog;
```

3. 启动服务:
```bash
go run cmd/main.go
```

服务将在 `http://localhost:8080` 启动。

## 10. 部署说明

1. 构建二进制文件:
```bash
go build -o blog-backend cmd/main.go
```

2. 设置生产环境变量

3. 运行服务:
```bash
./blog-backend
```

## 11. 前端集成说明

前端应用需要在请求头中添加认证信息:

```
Authorization: Bearer <jwt_token>
```

所有API响应都包含 `success` 字段，前端应根据此字段判断请求是否成功。