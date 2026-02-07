# Message Pocket

一个基于 PocketBase 的消息转发服务，用于接收 Webhook 事件并将消息转发到不同的目的地（如 QQ 群）。

## 功能特性

- **Webhook 接收**：接收 EdgeOne 等服务的 Webhook 事件
- **消息转发**：将接收到的消息转发到配置的目的地（目前支持 QQ 群）
- **消息存储**：所有消息都会保存到数据库，便于追溯和审计
- **统一消息处理**：通过 MessageBoxService 统一处理所有消息发送逻辑
- **Trace 追踪**：每个请求都有唯一的 trace_id，便于日志追踪
- **Token 验证**：支持 Bearer Token 验证，确保接口安全

## 技术栈

- **后端框架**：[PocketBase](https://pocketbase.io/) - Go 编写的开源后端框架
- **数据库**：SQLite（PocketBase 内置）
- **消息推送**：NapCat API（QQ 机器人）
- **日志**：slog（结构化日志）
- **配置管理**：Viper + YAML

## 项目结构

```
message-pocket/
├── internal/
│   ├── config/          # 配置管理
│   ├── constants/       # 常量定义
│   │   └── message_box_enum/  # 消息相关枚举
│   ├── controllers/     # 控制器层
│   ├── define/          # 数据定义
│   │   ├── dtos/        # 数据传输对象
│   │   └── model/       # 数据模型
│   ├── middlewares/     # 中间件
│   ├── repo/           # 数据访问层
│   ├── services/       # 业务服务层
│   │   └── logic/      # 业务逻辑
│   └── utils/          # 工具函数
├── migrations/         # 数据库迁移
├── config.yaml        # 配置文件
├── main.go           # 应用入口
├── go.mod            # Go 模块定义
└── README.md         # 项目说明
```

## 核心组件

### 1. MessageBoxService
统一的消息处理服务，负责：
- 保存消息到数据库
- 根据目的地类型发送消息
- 处理发送失败和重试逻辑

### 2. EOService
EdgeOne Webhook 事件处理服务：
- 解析 EdgeOne 事件
- 构建格式化消息
- 调用 MessageBoxService 保存并发送消息

### 3. 中间件
- **TraceMiddleware**：生成 trace_id 并存入 context，便于请求追踪
- **TokenAuthMiddleware**：验证请求的 Bearer Token

### 4. 数据模型
- **MessageBoxModel**：消息存储模型，包含消息内容、来源、目的地等信息

## 快速开始

### 1. 环境要求
- Go 1.21+
- PocketBase

### 2. 配置
复制 `config.yaml.example` 为 `config.yaml` 并修改配置：

```yaml
napcat:
  url: "http://your-napcat-server:port"
  token: "your-napcat-token"
  group_id: "your-qq-group-id"
```

### 3. 运行
```bash
# 安装依赖
go mod download

# 运行应用
go run main.go
```

### 4. API 接口

#### EdgeOne Webhook
```
POST /api/eo/webhook
Authorization: Bearer <your-token>
Content-Type: application/json

{
  "eventType": "deployment.succeeded",
  "appId": "your-app-id",
  "projectId": "your-project-id",
  "deploymentId": "your-deployment-id",
  "projectName": "Your Project",
  "repoBranch": "main",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## 开发规范

项目遵循严格的开发规范，详见 [SKILL.md](SKILL.md)。主要规范包括：

### 1. 数据访问层
- **先创建模型，后赋值给 SQL**：插入数据前先创建完整的模型对象
- **返回完整模型**：数据访问方法返回完整的模型对象，而非仅 ID

### 2. 消息发送
- **统一通过 MessageBoxService**：所有消息发送必须通过 MessageBoxService
- **参数封装**：多个非 context 参数应封装到结构体中

### 3. 错误处理
- 使用 `fmt.Errorf` 包装错误，提供上下文信息
- 在服务层记录适当的日志

### 4. 日志记录
- 使用 slog 进行结构化日志记录
- 包含 trace_id 等上下文信息

### 5. 规范维护
- **即时更新规范**：每次修改代码后都要即时更新规范文档

## 配置说明

### NapCat 配置
```yaml
napcat:
  url: "NapCat API 地址"
  token: "NapCat 认证 Token"
  group_id: "QQ 群号"
```

### 服务器配置
```yaml
server:
  open_token: "API 访问 Token"
```

## 消息格式

EdgeOne 事件会被格式化为以下消息：

```
🚀 EdgeOne 部署事件
📋 事件类型: 部署成功
📁 项目名称: Your Project
🌿 代码分支: main
🆔 项目ID: your-project-id
🆔 部署ID: your-deployment-id
⏰ 时间: 2024-01-01T00:00:00Z
```

## 事件类型支持

目前支持以下 EdgeOne 事件类型：
- `deployment.created` - 开始部署
- `deployment.succeeded` - 部署成功
- `deployment.failed` - 部署失败
- `deployment.cancelled` - 部署取消
- `deployment.rollback` - 部署回滚
- `deployment.in_progress` - 部署进行中
- `build.started` - 构建开始
- `build.succeeded` - 构建成功
- `build.failed` - 构建失败
- `project.created` - 项目创建
- `project.updated` - 项目更新
- `project.deleted` - 项目删除

## 扩展开发

### 添加新的消息来源
1. 在 `message_box_enum/source_type.go` 中添加新的 SourceType
2. 创建对应的 Service 处理新来源的消息
3. 在控制器中添加对应的路由

### 添加新的消息目的地
1. 在 `message_box_enum/destiantion_type.go` 中添加新的 DestinationType
2. 在 MessageBoxService 的 `SendMessage` 方法中添加对应的发送逻辑
3. 实现具体的发送方法

## 许可证

MIT License