# WMS 仓库管理系统 - 全功能库存管理平台

一个基于 Go 语言开发的现代化仓库管理系统后端服务，提供完整的库存管理、盘点、查询和报告功能。系统采用分层架构设计，具有高性能、高可用性和可扩展性。

## 项目简介

本项目采用分层架构设计，使用 Go 语言实现了一个完整的 WMS 库存管理系统后端。系统提供 REST API 接口用于库存上传、盘点、查询和报告，通过事务确保业务逻辑的一致性，并使用 Gorm 框架将数据持久化到 PostgreSQL 数据库。

### 核心特性

- ✅ **分层架构设计** - 清晰的职责分离:数据层、服务层、API 层
- ✅ **结构化日志** - 使用 Zap 提供高性能的结构化日志记录
- ✅ **事务一致性** - 基于 Gorm 的事务管理确保数据完整性
- ✅ **输入验证** - 使用 Gin 框架的验证机制确保数据安全性
- ✅ **优雅停机** - 支持信号处理和优雅停机机制
- ✅ **数据库连接池** - 高效的数据库连接管理
- ✅ **自动迁移** - 数据库模型自动创建和更新
- ✅ **批量处理** - 支持批量库存操作
- ✅ **实时数据同步** - 库存数据与盘点记录保持一致
- ✅ **健康检查** - 支持容器化部署和健康监控

## 技术栈

- **Web 框架**: [Gin](https://github.com/gin-gonic/gin) v1.11.0
- **ORM 框架**: [Gorm](https://gorm.io/) v1.25.10
- **数据库**: PostgreSQL
- **日志系统**: [Zap](https://github.com/uber-go/zap) v1.26.0
- **配置管理**: [godotenv](https://github.com/joho/godotenv) v1.5.1
- **Go 版本**: 1.23.0

## 系统架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   API 层        │    │   服务层        │    │   数据访问层     │
│  (handlers)     │───▶│ (service)       │───▶│ (repository)    │
│                 │    │                 │    │                 │
│  - Gin 路由     │    │  - 业务逻辑      │    │  - 数据库交互    │
│  - 请求验证     │    │  - 事务管理      │    │  - 查询构建     │
│  - 响应格式化   │    │  - 数据验证      │    │  - 连接池管理   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   数据模型层     │
                       │   (models)      │
                       │                 │
                       │  - Stock        │
                       │  - InventoryCheck│
                       │    Record       │
                       └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   数据库        │
                       │   (PostgreSQL)  │
                       └─────────────────┘
```

## 项目结构

```
wms/
├── cmd/
│   └── server/              # 应用程序入口
│       └── main.go          # 主函数,依赖注入和优雅停机
├── internal/
│   ├── api/                 # API 层
│   │   ├── dto/             # 数据传输对象(DTO)
│   │   ├── handlers/        # HTTP 请求处理器
│   │   └── routes/          # 路由配置
│   ├── model/               # 数据模型层
│   │   ├── inventory_check.go    # 库存盘点记录模型
│   │   └── stock.go              # 库存模型
│   ├── repository/          # 数据访问层
│   │   └── inventory_check_repository.go  # 库存盘点仓储
│   └── service/             # 业务逻辑层
│       └── inventory_service.go           # 库存服务
├── pkg/
│   ├── config/              # 配置管理
│   │   └── config.go
│   └── logger/              # 日志系统
│       └── logger.go
├── scripts/                 # 脚本文件
│   └── init-db.sql          # 数据库初始化脚本
├── .env.example             # 环境变量示例文件
├── docker-compose.yml       # Docker 容器配置
├── Dockerfile               # Docker 镜像构建文件
├── go.mod                   # Go 模块依赖
└── go.sum                   # 依赖版本锁定
```

## 快速开始

### 前置要求

- Go 1.23.0 或更高版本
- PostgreSQL 数据库
- Git
- Docker (可选，用于容器化部署)

### 安装步骤

1. **克隆项目**

```bash
git clone 
cd wms
```

2. **安装依赖**

```bash
make deps
```

3. **配置环境变量**

运行 `make setup` 自动复制 `.env.example` 为 `.env`:

```bash
make setup
```

编辑 `.env` 文件,配置数据库连接信息:

```env
# 服务器配置
SERVER_ADDRESS=0.0.0.0
SERVER_PORT=8080

# 数据库配置
DATABASE_DSN=host=localhost user=wms_user password=wms_password dbname=wms_db port=5432 sslmode=disable TimeZone=Asia/Shanghai

# 运行环境
ENVIRONMENT=development
```

4. **创建数据库**

连接到 PostgreSQL 并创建数据库:

```sql
CREATE DATABASE wms_db;
CREATE USER wms_user WITH PASSWORD 'wms_password';
GRANT ALL PRIVILEGES ON DATABASE wms_db TO wms_user;
```

5. **运行应用**

```bash
make dev
```

> `make dev` Requires Air (auto-installed on first run) 并提供热重载体验。

或者编译后运行:

```bash
make build
./bin/wms-server
```

需要直接启动服务并自动拉起依赖容器时，可执行:

```bash
make run
```

### Docker 部署

1. **使用 Docker Compose 一键部署（推荐）**

```bash
# 构建并启动所有服务（Postgres + WMS Server）
docker-compose up -d --build

# 查看日志
docker-compose logs -f wms-server

# 停止服务
docker-compose down

# 停止并删除数据卷（清空数据库）
docker-compose down -v
```

2. **单独构建 Docker 镜像**

```bash
# 构建镜像
docker build -t wms-server:latest .

# 运行容器（需要先启动 Postgres）
docker run -d \
  --name wms-server \
  -p 8080:8080 \
  -e DATABASE_DSN="host=postgres user=wms_user password=wms_password dbname=wms_db port=5432 sslmode=disable TimeZone=Asia/Shanghai" \
  wms-server:latest
```

**重要提示**：
- Docker 环境中，数据库主机名使用 `postgres`（服务名）而不是 `localhost`
- 生产环境配置参考 `.env.production` 文件
- 健康检查端点：`http://localhost:8080/health`

服务将启动在 `http://localhost:8080`

#### 快速启动命令

```bash
make setup && make docker-up && make dev
```

## Makefile Targets

Makefile 覆盖了开发、测试、容器及常用维护流程，默认执行 `make` 将运行 `all` 目标(构建 + 测试)。

### Build & Development

| Target | Description | Usage |
|--------|-------------|-------|
| `all` | 构建二进制并运行全部测试 | `make` |
| `build` | 编译应用并输出到 `bin/wms-server` | `make build` |
| `test` | 执行 `go test ./... -v` | `make test` |
| `run` | 构建、确保 `.env` 存在与 Docker 服务运行后启动二进制 | `make run` |
| `dev` | 使用 Air 热重载启动开发服务器(首次自动安装 Air) | `make dev` |
| `clean` | 删除 `bin/` 构建产物 | `make clean` |

### Docker Orchestration

| Target | Description | Usage |
|--------|-------------|-------|
| `docker-up` | 启动 docker-compose 服务 | `make docker-up` |
| `docker-down` | 停止并移除 docker-compose 服务 | `make docker-down` |
| `docker-logs` | 持续跟踪 docker-compose 日志 | `make docker-logs` |
| `docker-restart` | 重启所有 docker-compose 服务 | `make docker-restart` |

### Utility & Maintenance

| Target | Description | Usage |
|--------|-------------|-------|
| `deps` | 下载 Go 模块依赖 | `make deps` |
| `tidy` | 同步 `go.mod` 与 `go.sum` | `make tidy` |
| `fmt` | 对所有 Go 文件执行 `go fmt` | `make fmt` |
| `vet` | 运行 `go vet` 静态分析 | `make vet` |
| `setup` | 从 `.env.example` 创建 `.env` | `make setup` |
| `help` | 列出所有可用目标与说明 | `make help` |

## API 文档

###  上传库存盘点数据

**接口**: `POST /api/wms/inventory/check/upload`

**描述**: 处理单次库存盘点操作，更新库存并记录差异

**请求头**:
```
Content-Type: application/json
```

**请求体**:
```json
{
  "checker_id": "CHECKER001",
  "location_code": "A-01-01",
  "material_code": "MAT001",
  "actual_quantity": 100
}
```

**响应示例**:

成功 (200 OK):
```json
{
  "code": 0,
  "message": "success"
}
```

验证失败 (400 Bad Request):
```json
{
  "code": -1,
  "message": "Invalid request: checker_id is required"
}
```

服务器错误 (500 Internal Server Error):
```json
{
  "code": -1,
  "message": "Failed to process inventory check: database connection timeout"
}
```



## 数据模型

### Stock (库存表)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | uint | PRIMARY KEY, AUTO_INCREMENT | 主键 |
| material_code | varchar(100) | NOT NULL, UNIQUE INDEX | 物料代码 |
| location_code | varchar(100) | NOT NULL, UNIQUE INDEX | 库位代码 |
| quantity | int | NOT NULL, DEFAULT: 0 | 库存数量 |
| created_at | datetime | AUTO_CREATE | 创建时间 |
| updated_at | datetime | AUTO_UPDATE | 更新时间 |

### InventoryCheckRecord (库存盘点记录表)

| 字段名 | 类型 | 约束 | 说明 |
|--------|------|------|------|
| id | uint | PRIMARY KEY, AUTO_INCREMENT | 主键 |
| checker_id | varchar(100) | NOT NULL, INDEX | 盘点人员ID |
| location_code | varchar(100) | NOT NULL, INDEX | 库位代码 |
| material_code | varchar(100) | NOT NULL, INDEX | 物料代码 |
| actual_quantity | int | NOT NULL | 实盘数量 |
| stock_quantity | int | NOT NULL | 系统库存数量 |
| difference | int | NOT NULL | 差异数量 (实际-系统) |
| check_time | timestamp | NOT NULL, INDEX | 盘点时间 |
| is_processed | boolean | DEFAULT: false, NOT NULL | 是否已处理 |
| created_at | datetime | AUTO_CREATE | 创建时间 |
| updated_at | datetime | AUTO_UPDATE | 更新时间 |

## 部署

### Docker Compose 部署 (推荐)

创建 `docker-compose.yml` 文件:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: wms-postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: wms_user
      POSTGRES_PASSWORD: wms_password
      POSTGRES_DB: wms_db
      TZ: Asia/Shanghai
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init-db.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U wms_user -d wms_db"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - wms-network

  wms-server:
    build: .
    container_name: wms-server
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      DATABASE_DSN: host=wms-postgres user=wms_user password=wms_password dbname=wms_db port=5432 sslmode=disable TimeZone=Asia/Shanghai
      ENVIRONMENT: production
    depends_on:
      postgres:
        condition: service_healthy
    networks:
      - wms-network

volumes:
  postgres_data:
    driver: local

networks:
  wms-network:
    driver: bridge
```

启动服务:

```bash
make docker-up
```

### 直接部署

```bash
# 构建二进制文件
make build

# 运行
./bin/wms-server
```


## 核心功能

### 1. 库存盘点处理流程

系统通过以下步骤处理库存盘点数据:

1. **接收请求** - API 层接收并验证盘点数据
2. **业务处理** - 服务层开启数据库事务
3. **查询库存** - 查询当前库存数据
4. **计算差异** - 计算盘点数量与库存数量的差异
5. **记录盘点** - 创建盘点记录
6. **更新库存** - 更新库存数量
7. **提交事务** - 提交事务或回滚(发生错误时)

### 2. 事务管理

所有涉及数据修改的操作都使用事务确保:
- **原子性** - 要么全部成功,要么全部回滚
- **一致性** - 库存数据与盘点记录保持一致
- **隔离性** - 避免并发操作导致的数据不一致
- **持久性** - 提交后数据永久保存

### 3. 日志系统

使用 Zap 提供结构化日志,支持:
- 不同日志级别 (DEBUG, INFO, WARN, ERROR)
- JSON 格式输出
- 性能优化的零内存分配
- 开发/生产环境自动切换

### 4. 优雅停机

服务支持优雅停机机制:
- 监听 `SIGINT` (Ctrl+C) 和 `SIGTERM` 信号
- 5 秒超时等待正在处理的请求
- 自动关闭数据库连接
- 完整的停机日志记录



## 故障排查

### 常见问题及解决方案

1. **数据库连接失败**
   
   **问题**: `failed to connect to database`
   
   **解决方案**:
   - 检查数据库服务是否启动
   - 验证 `DATABASE_DSN` 配置是否正确
   - 确认数据库用户权限

2. **端口被占用**
   
   **问题**: `bind: address already in use`
   
   **解决方案**:
   - 更改 `SERVER_PORT` 环境变量
   - 使用 `lsof -i :8080` 查看占用端口的进程

3. **API 调用失败**
   
   **问题**: `400 Bad Request`
   
   **解决方案**:
   - 检查请求体格式和必填字段
   - 确认 Content-Type 设置为 `application/json`

4. **事务回滚异常**
   
   **问题**: `failed to commit transaction`
   
   **解决方案**:
   - 检查数据库锁情况
   - 确认网络连接稳定性

### 日志分析

系统使用 Zap 记录结构化日志，常见的日志类型包括：

- **INFO**: 业务操作记录，如库存更新、盘点完成
- **WARN**: 潜在问题，如无效输入、重试操作
- **ERROR**: 系统错误，如数据库连接失败、业务异常
- **DEBUG**: 详细调试信息，仅在开发环境启用

### 性能监控

1. **API 响应时间**: 监控各接口的响应时间，识别性能瓶颈
2. **数据库连接**: 监控数据库连接数和查询性能
3. **内存使用**: 监控应用内存使用情况，防止内存泄漏
4. **并发请求**: 监控并发请求数，确保系统稳定性

## 开发指南

### 添加新的 API 接口

1. 在 `internal/api/dto/` 定义请求/响应 DTO
2. 在 `internal/service/` 实现业务逻辑
3. 在 `internal/api/handlers/` 创建处理器
4. 在 `internal/api/routes/` 注册路由

### 数据库迁移

系统使用 Gorm 自动迁移功能。新增模型后，只需在 `cmd/server/main.go` 的 `AutoMigrate` 中添加:

```go
db.AutoMigrate(&model.NewModel{})
```

### 日志记录

在代码中使用日志:

```go
log.Info("Operation successful",
    zap.String("operation", "inventory_check"),
    zap.Int("count", itemCount),
)

log.Error("Operation failed",
    zap.Error(err),
    zap.String("context", "additional info"),
)
```

## 配置说明

### 环境变量

| 变量名 | 说明 | 默认值 | 必填 |
|--------|------|--------|------|
| `SERVER_ADDRESS` | 服务器监听地址 | `0.0.0.0` | 否 |
| `SERVER_PORT` | 服务器监听端口 | `8080` | 否 |
| `DATABASE_DSN` | PostgreSQL 连接字符串 | - | 是 |
| `ENVIRONMENT` | 运行环境 (`development`/`production`) | `development` | 否 |

### 数据库连接池配置

默认配置(在 `cmd/server/main.go` 中):
- 最大空闲连接数: 10
- 最大打开连接数: 100
- 连接最大生命周期: 1 小时



## 版本历史
### v1.0.1 (2025-12-03)
- 增加makefile配置
- 增加docker配置
- 支持热重载

### v1.0.0 (2025-12-03)
- 初始版本，包含基础库存管理功能
- 支持库存盘点和更新
- 实现事务一致性保证

## 支持与联系

**版本**: 1.0.0

**最后更新**: 2025-12-03

**作者**: 小王同学
