# Go-DDD：快速构建领域驱动的 Go 服务

[![CI](https://github.com/sklinkert/go-ddd/actions/workflows/go.yml/badge.svg)](https://github.com/sklinkert/go-ddd/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/sklinkert/go-ddd/branch/main/graph/badge.svg)](https://codecov.io/gh/sklinkert/go-ddd)
[![Go Reference](https://pkg.go.dev/badge/github.com/sklinkert/go-ddd.svg)](https://pkg.go.dev/github.com/sklinkert/go-ddd)
[![Go Version](https://img.shields.io/github/go-mod/go-version/sklinkert/go-ddd)](go.mod)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

[English](README.md) | [简体中文](README.zh-CN.md)

`go-ddd` 帮助你快速搭建生产级的 Go 后端服务，将业务规则、基础设施与交付层代码清晰分离。开箱即用：一套明确取舍的 DDD 构建块、CQRS 命令与查询流、幂等的写入路径、基于事务性发件箱（Transactional Outbox）的领域事件，以及让数据库模式与代码保持同步的工具链。

> 📚 **第一次接触 DDD？[从教程开始 →](https://sklinkert.github.io/go-ddd/)**
> 九章教程从零讲解领域驱动设计，全程以本仓库代码为实例。（教程为英文）

## 为什么选择这个模板

- **模型优先的默认设计** —— 洋葱架构保持领域层纯净，应用服务负责编排基础设施相关的事务。
- **经过实战检验的模式** —— 命令、查询、仓储、值对象、领域事件、软删除，均对应真实企业应用中的常用模式。
- **幂等的写入管道** —— 无竞态的幂等键（原子性预留，杜绝先查后写），让每个命令都可以安全重试。
- **领域事件 + 发件箱** —— 状态变更与其事件在同一事务中持久化，由中继进程以至少一次（at-least-once）语义发布。
- **严格的迁移纪律** —— SQL 迁移、`migrate.go` 与 `sqlc` 让数据库模式演进显式且可复现。

## 你将得到什么

- 一个市场（Marketplace）示例：演示聚合（Seller、Product）、`Money` 值对象、跨模块交互与校验规则。
- `internal/` 下按层组织的模块：`domain`、`application`、`infrastructure`、`interface`，以及可复用测试夹具的 `testhelpers`。
- 可执行入口 `cmd/marketplace/main.go`，随时接入你选择的适配器或框架。
- `migrations/` 与 `sql/` 中的数据库资产，以及由 `sqlc` 生成的数据访问代码。
- [OpenAPI 规范](api/openapi.yaml)、`/healthz` 与 `/readyz` 探针，以及一条命令即可启动的 Docker Compose 环境。

## 技术栈

- **Go 1.26**，惯用的 Go 写法，测试基于 testify。
- **Echo v4** 提供 REST 接口。
- **pgx/v5** 与 `sqlc` 实现类型安全的 PostgreSQL 访问。
- **golang-migrate** 管理 SQL 模式迁移。
- **Testcontainers** 在测试中启动一次性的 Postgres 实例。
- **google/uuid** 在领域层内生成确定性的 Id。
- **golangci-lint** 做静态分析，另有 `Makefile`、`Dockerfile` 与 `docker-compose.yml` 支持一条命令启动本地环境。

## 设计原则实战

领域驱动设计让实现与持续演进的模型保持一致。`go-ddd` 通过一个简单的市场领域来展示这一点：`Seller`（卖家）管理 `Product`（商品），其中涉及聚合、值对象与校验流程。

一次写请求的全过程：

```mermaid
sequenceDiagram
    participant C as 客户端
    participant R as REST 控制器
    participant S as 应用服务
    participant D as 领域层
    participant P as Postgres

    C->>R: POST /api/v1/products (idempotency_key)
    R->>S: CreateProductCommand
    S->>P: 原子性预留幂等键
    S->>D: NewProduct(name, Money, ValidatedSeller)
    D->>D: validate() → ValidatedProduct + ProductCreated 事件
    S->>P: 同一事务写入商品 + 发件箱事件
    P-->>S: 已持久化的行
    S->>P: 缓存幂等键对应的响应
    S-->>R: CommandResult
    R-->>C: 201 Created (JSON)
    Note over P: 发件箱中继异步发布<br/>ProductCreated 事件
```

## 文档

📚 **[DDD 零基础教程](https://sklinkert.github.io/go-ddd/)** —— 九个章节：为什么需要 DDD、实体、值对象、聚合、仓储、CQRS、发件箱、幂等性与测试，全部结合本仓库代码讲解。（英文）

📖 **[DDD 与 CQRS 原则完整指南](DDD_CQRS_PRINCIPLES.md)** —— 学习如何将这些模式应用到任意业务领域。（英文）

## 仓库结构

![ddd-diagram-onion.png](ddd-diagram-onion.png)

- `domain`：软件的核心，承载业务逻辑与规则。
    - `entities`：系统中的基础对象，如 `Product` 和 `Seller`，包含基本的校验逻辑。
- `application`：面向具体用例的操作，与领域层交互。
- `infrastructure`：为上层提供数据库访问等技术能力。
    - `db`：数据库访问与模型。
    - `repositories`：存储需求的具体实现。
- `interface`：与外部世界交互的最外层，如 API 端点。
    - `api/rest`：处理 HTTP 请求与响应的控制器。

## 分层原则

- 领域层（Domain）
  - 不得依赖其他任何层。
  - 为基础设施层提供接口，但不得访问基础设施。
  - 实现业务逻辑与规则。
  - 对实体执行校验；只有通过校验的实体才会传递给基础设施层。
  - 实体的默认值（如 Id 的 uuid、创建时间戳）由领域层设置。不要在基础设施层甚至数据库中设置默认值！
  - 不要把领域对象泄漏到外部。
- 应用层（Application）
  - 领域层与基础设施层之间的粘合代码。
- 基础设施层（Infrastructure）
  - 仓储负责领域实体与数据库模型之间的转换与读写，这里不执行任何业务逻辑。
  - 实现由领域层定义的接口。
  - 实现持久化逻辑，例如访问 Postgres 或 MySQL。
  - 写入存储后先读回再返回，确保数据被正确写入。

## 最佳实践

- 仓储的读取方法不要返回已校验实体（ValidatedX），而应直接返回领域实体类型。
  - 校验规则会随时间演变。你不会想为此迁移数据库中的全部数据；应保证无论校验逻辑如何演进，历史数据总是可以加载。
  - 否则，用旧校验逻辑写入的数据将无法读出，只能在运行时处理错误。
  - 把校验放在写入侧——创建（NewX）与更新方法——那里本来就必须维护不变量。
- 不要把默认值（如当前时间戳、Id）放在数据库里，而应在领域层（工厂方法）中设置：
  - 两个事实来源（source of truth）非常危险。
  - 领域层更容易测试。
  - 数据库可能被替换，你不会想连带修改所有默认值。
- 基础设施层写入后总是读回实体。
  - 确保数据被正确写入，并且我们永远不会基于过期数据继续操作。
- `find` 与 `get` 的语义区分：
  - `find` 方法可以返回空值或空列表。
  - `get` 方法必须返回值；找不到时应返回错误。
- 删除：始终使用软删除。在数据库中创建 `deleted_at` 列，删除实体时将其设为当前时间戳。这样在需要时总能恢复数据。

## CQRS 与幂等性

### 命令查询职责分离（CQRS）

CQRS 将应用中的读操作（查询）与写操作（命令）分开。在本代码库中：

- **命令**修改状态（CreateSellerCommand、CreateProductCommand、UpdateSellerCommand）
- **查询**读取数据、无副作用（FindAllSellers、FindSellerById）

这种分离带来不同的优化空间：

- **写优化**：命令可以使用规范化模式、ACID 事务与强一致性。
- **读优化**：查询可以使用反规范化视图、缓存、只读副本，甚至不同的数据库（例如写用 PostgreSQL、读用 Elasticsearch）。
- **可扩展性**：读写负载可以按实际使用模式独立扩容。
- **性能**：复杂查询不影响写入性能，写锁也不会阻塞读取。

### 幂等键

幂等性保证多次相同的请求与单次请求产生同样的效果，这对处理分布式系统中的网络故障与重试至关重要。实现方式：

- 每个命令可在请求中携带可选的 `idempotency_key`。
- 幂等键通过**原子性预留**（`INSERT ... ON CONFLICT DO NOTHING`）占用，两个并发的同键请求绝不会同时执行——不存在先查后写的竞态。
- 已完成的请求返回缓存的响应；仍在执行中的请求返回"处理中"错误，客户端稍后重试即可。
- 用**不同的请求体**复用同一个键会被拒绝，而不是悄悄返回错误的缓存响应。
- 命令失败时预留会被释放，客户端可以重试；因进程崩溃而遗留的预留会在 TTL 后过期。

这可以防止客户端重试失败请求时创建出重复的实体。

### 领域事件与事务性发件箱

聚合在业务相关的事情发生时记录事件（例如 `ProductCreated`）。事件不会直接发布到消息代理——那样在数据库提交与发布之间进程崩溃时会丢失事件——而是先存入 `outbox_events` 表。中继进程轮询发件箱，以至少一次（at-least-once）语义发布未发布的事件。参见 `internal/domain/events/` 与 `internal/infrastructure/outbox/`。

## 数据库迁移

本项目使用 [golang-migrate](https://github.com/golang-migrate/migrate) 管理数据库模式。迁移文件存放在 `migrations/` 目录下，按序号版本化。

### 迁移文件结构

```
migrations/
├── 000001_initial_schema.up.sql    # 创建初始表
├── 000001_initial_schema.down.sql  # 初始模式的回滚
├── 000002_price_as_money.up.sql    # 金额改为整数分 + 货币
├── 000003_outbox.up.sql            # 事务性发件箱表
└── ...
```

### 执行迁移

#### 使用内置工具：

```bash
# 应用全部待执行的迁移
go run migrate.go -database-url "postgres://user:pass@localhost/db?sslmode=disable" -command up

# 回滚最近一次迁移
go run migrate.go -database-url "postgres://user:pass@localhost/db?sslmode=disable" -command down -steps 1

# 查看当前版本
go run migrate.go -database-url "postgres://user:pass@localhost/db?sslmode=disable" -command version

# 强制设置到指定版本（谨慎使用）
go run migrate.go -database-url "postgres://user:pass@localhost/db?sslmode=disable" -command force -version 1
```

#### 直接使用 CLI 工具：

```bash
# 安装 CLI 工具
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# 应用全部待执行的迁移
migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" up

# 回滚最近一次迁移
migrate -path migrations -database "postgres://user:pass@localhost/db?sslmode=disable" down 1
```

### 创建新迁移

```bash
migrate create -ext sql -dir migrations -seq add_user_email_column
```

将生成两个文件：

- `0000NN_add_user_email_column.up.sql` —— 正向迁移
- `0000NN_add_user_email_column.down.sql` —— 回滚迁移

### 迁移最佳实践

- 始终同时编写 `up` 和 `down` 迁移
- 在生产数据的副本上测试迁移
- 保持迁移小而聚焦
- 已在生产环境应用过的迁移文件绝不修改
- 使用描述性的迁移文件名

## 快速开始

> 需要 **Go 1.26+**。运行 `make help` 查看全部可用目标。

### 用 Docker Compose 一键启动

一条命令启动 Postgres、应用迁移并运行 API：

```bash
make docker-up        # docker compose up --build
# API 已在 http://localhost:8080 就绪
make docker-down      # 全部销毁
```

### 30 秒体验 API

创建一个卖家：

```bash
curl -s -X POST http://localhost:8080/api/v1/sellers \
  -H 'Content-Type: application/json' \
  -d '{"name": "Acme Corp", "idempotency_key": "create-acme-1"}'
```

```json
{"id":"0197a3c2-...","name":"Acme Corp","created_at":"2026-07-14T09:00:00Z","updated_at":"2026-07-14T09:00:00Z"}
```

为该卖家创建商品（价格为整数分——绝不使用浮点数）：

```bash
curl -s -X POST http://localhost:8080/api/v1/products \
  -H 'Content-Type: application/json' \
  -d '{"name": "Wooden Chair", "price_cents": 4999, "currency": "EUR", "seller_id": "<上一步返回的 seller-id>"}'
```

```json
{"id":"0197a3c3-...","name":"Wooden Chair","price_cents":4999,"currency":"EUR","seller_id":"0197a3c2-...","created_at":"...","updated_at":"..."}
```

用相同的 `idempotency_key` 重放请求——返回缓存的响应，而不是创建重复卖家：

```bash
curl -s -X POST http://localhost:8080/api/v1/sellers \
  -H 'Content-Type: application/json' \
  -d '{"name": "Acme Corp", "idempotency_key": "create-acme-1"}'
# → 响应完全一致，不会插入第二行
```

列出商品并检查服务健康状态：

```bash
curl -s http://localhost:8080/api/v1/products
curl -s http://localhost:8080/readyz
```

完整的 API 描述见 [OpenAPI 规范](api/openapi.yaml)。

### 本地开发

1. 克隆仓库：
```bash
git clone https://github.com/sklinkert/go-ddd.git
cd go-ddd
go mod download
```

2. 安装 sqlc（开发用）：
```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

3. 生成数据库代码（修改 SQL 查询后执行）：
```bash
sqlc generate
```

4. 配置 PostgreSQL 并运行迁移：
```bash
# 设置数据库连接 URL
export DATABASE_URL="postgres://user:password@localhost/dbname?sslmode=disable"

# 使用内置工具运行迁移
go run migrate.go -command up

# 或直接使用 CLI 工具
migrate -path migrations -database $DATABASE_URL up
```

5. 运行应用：
```bash
make run            # 或：go run ./cmd/marketplace
# 通过 DATABASE_URL 环境变量覆盖数据库配置（libpq DSN 或 postgres:// URL）。
```

### 常用 Make 目标

```bash
make build        # 构建服务端二进制到 ./bin
make test         # 运行全部测试，含竞态检测（需要 Docker）
make test-unit    # 仅运行不需要 Docker 的测试
make lint         # 运行 golangci-lint
make fmt          # 格式化代码（gofmt + goimports）
make migrate-up   # 对 $DATABASE_URL 应用迁移
make sqlc         # 重新生成 sqlc 代码
```

### 参与贡献

欢迎贡献代码、提交 Issue 与功能建议！请查看 Issues 页面。

### 使用此模板

在 GitHub 上点击 **"Use this template"**，以此结构为基础创建你自己的服务；或者 fork 后把市场领域替换成你的业务领域。[DDD 与 CQRS 指南](DDD_CQRS_PRINCIPLES.md)会带你把这些模式适配到任意业务领域。

如果这个模板对你有帮助，**请点一个 ⭐** —— 这能帮助更多人发现它。

[![Star History Chart](https://api.star-history.com/svg?repos=sklinkert/go-ddd&type=Date)](https://star-history.com/#sklinkert/go-ddd&Date)

### 许可证

基于 MIT 许可证分发。详情见 LICENSE。
