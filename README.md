# PisaList 后端服务

这是一个基于 Go 语言开发的待办事项和心愿清单管理系统的后端服务。

## 技术栈

- Go 1.21
- Gin 框架
- GORM
- MySQL
- Docker
- JWT 认证

## 功能特性

### 用户管理
- 用户注册
- 用户登录（JWT认证）

### 待办事项管理
- 创建任务
- 删除任务
- 更新任务
- 完成任务
- 设置任务优先级
- 查看任务时间线（近7天）
- 获取今日任务列表

### 心愿管理
- 创建心愿
- 删除心愿
- 更新心愿
- 分享心愿到社区
- 查看个人心愿列表
- 查看心愿社区
- 随机获取心愿

## 项目结构

```
.
├── api
│   └── v1              # API 处理器
├── config              # 配置文件
├── internal
│   ├── middleware      # 中间件
│   ├── model          # 数据模型
│   └── service        # 业务逻辑
├── pkg
│   ├── database       # 数据库工具
│   ├── jwt           # JWT 工具
│   └── util          # 通用工具
├── Dockerfile         # Docker 构建文件
├── docker-compose.yml # Docker 编排文件
├── go.mod            # Go 模块文件
└── main.go           # 主程序入口
```

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/yourusername/PisaListBE.git
cd PisaListBE
```

2. 修改配置文件
```bash
cp config/config.example.yaml config/config.yaml
# 编辑 config.yaml 文件，设置数据库连接信息和 JWT 密钥
```

3. 使用 Docker Compose 启动服务
```bash
docker-compose up -d
```

服务将在 http://localhost:8080 启动

## API 文档

### 认证相关
- POST /api/v1/register - 用户注册
- POST /api/v1/login - 用户登录

### 任务相关
- POST /api/v1/tasks - 创建任务
- DELETE /api/v1/tasks/:id - 删除任务
- PUT /api/v1/tasks/:id - 更新任务
- PUT /api/v1/tasks/:id/complete - 完成任务
- PUT /api/v1/tasks/:id/importance - 更新任务优先级
- GET /api/v1/tasks/timeline - 获取任务时间线
- GET /api/v1/tasks/today - 获取今日任务

### 心愿相关
- POST /api/v1/wishes - 创建心愿
- DELETE /api/v1/wishes/:id - 删除心愿
- PUT /api/v1/wishes/:id - 更新心愿
- POST /api/v1/wishes/:id/share - 分享心愿
- GET /api/v1/wishes - 获取用户心愿列表
- GET /api/v1/wishes/community - 获取心愿社区列表
- GET /api/v1/wishes/random - 获取随机心愿

## 性能优化

1. 数据库优化
   - 使用适当的索引
   - 连接池配置优化
   - 查询优化

2. API 性能
   - 使用适当的缓存策略
   - 请求限流
   - 响应压缩

3. 代码优化
   - 使用 goroutine 处理异步任务
   - 合理的错误处理
   - 内存使用优化

## 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 许可证

MIT License 