# Go Forum (MVP 后端项目)

一个使用 **Go + Gin + GORM + JWT** 实现的最小可用论坛后端 (MVP)。  
主要用于学习和实践 **后端服务的基础知识**，配合 Flutter 或其他前端进行调用。

---

## ✨ 项目简介
本项目是一个简易的论坛后端，实现了从 **用户注册 → 登录 → 发帖 → 评论** 的完整闭环。  
目标是帮助前端/客户端开发者快速上手后端开发，掌握常见功能和技术点。

---

## 🔧 技术栈
- **语言**: [Go](https://go.dev/)  
- **Web 框架**: [Gin](https://github.com/gin-gonic/gin)  
- **ORM**: [GORM](https://gorm.io/) + MySQL/MariaDB  
- **认证**: [JWT](https://github.com/golang-jwt/jwt)  

---

## 📂 项目结构
go-forum/
├── main.go # 入口文件，初始化服务和路由
├── models.go # 数据模型（User、Post、Comment）
├── handlers.go # 业务逻辑（注册、登录、发帖、评论）
├── middleware.go # 中间件（JWT 鉴权）
└── go.mod # 依赖管理


---

## 🚀 功能特性
1. **用户注册**
   - 用户名唯一
   - 简单密码存储（可优化为 bcrypt）

2. **用户登录**
   - 校验用户名和密码
   - 生成并返回 JWT token

3. **JWT 鉴权中间件**
   - 保护需要登录的接口
   - 自动解析 token 并注入用户信息

4. **发帖功能**（需登录）
   - 记录发帖用户 ID
   - 保存到数据库

5. **评论功能**（需登录）
   - 关联帖子和用户
   - 保存评论数据

---

## 📖 学习收获
通过本项目，你可以掌握：
- Go 项目基本结构与依赖管理 (`go mod`)
- Gin 框架的路由、请求处理与 JSON 响应
- GORM 的数据库建表、查询和插入操作
- JWT 的生成与校验
- 中间件的使用
- 最小可用后端（MVP）的搭建流程

---

## ⚙️ 快速开始

### 1. 克隆项目
```bash
git clone https://github.com/zeusghj/go-forum.git
cd go-forum
```

### 2. 安装依赖
```bash
go mod tidy
```

### 3. 配置数据库
在 models.go 中修改数据库配置：
```bash
dsn := "username:password@tcp(127.0.0.1:3306)/go_forum?charset=utf8mb4&parseTime=True&loc=Local"
```

### 4. 运行项目
在 models.go 中修改数据库配置：
```bash
go run .
```
默认启动在 http://localhost:8080

## 📡 API 示例

### 注册
在 models.go 中修改数据库配置：
```bash
go run .
```

### 登录
```bash
POST /api/register
{
  "username": "alice",
  "password": "123456"
}

```
返回
```bash
{ "token": "your.jwt.token" }
```

### 发帖
```bash
POST /api/posts
Authorization: Bearer <token>
{
  "title": "我的第一篇帖子",
  "content": "Hello, Go Forum!"
}
```

### 评论
```bash
POST /api/comments
Authorization: Bearer <token>
{
  "post_id": 1,
  "content": "不错的帖子！"
}
```

## 📌 后续优化方向

- 使用 bcrypt 对密码进行加密存储

- 增加帖子列表/评论列表接口

- 增加用户资料、头像等功能

- 完善错误处理和日志记录

- Docker 化部署，方便上线运行

- CI/CD 集成，实现持续更新

## 🤝 致谢

本项目主要作为学习用 MVP 后端，非常适合前端/客户端开发者入门后端开发。
欢迎 Fork、Star，一起学习与改进！