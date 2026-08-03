# Airis 项目现状分析

> 分析日期：2026-08-03
> 分析分支：`main`
> 基线提交：`63726d5 feat: load signing secret from environment`

## 1. 分析结论

Airis 当前是一个体量较小、调用链直接的 Go/Gin API 查询服务。项目可以正常编译，已有单元测试、竞态检测和静态检查均通过，但业务测试覆盖不足，安全、错误治理、并发安全和模块边界仍有明显缺口。

当前实现更适合作为可继续迭代的功能原型。在完成密钥管理、认证校验、日志脱敏、MongoDB 故障恢复和业务测试补齐之前，不建议直接作为高流量生产服务部署。

## 2. 项目概况

### 2.1 技术栈

| 类别 | 当前实现 |
| --- | --- |
| 开发语言 | Go 1.24 |
| HTTP 框架 | Gin 1.11.0 |
| 数据库 | MongoDB Driver 1.17.2 |
| 缓存 | go-redis 9.7.0，目前未接入核心业务 |
| 日志 | Logrus 1.8.1 |
| 配置 | godotenv + 环境变量 |
| RPC | gRPC/Protobuf，目前只有 HelloWorld 示例 |

`README.md` 中的“Go 1.18+”与 `go.mod` 实际声明的 Go 1.24、toolchain 1.24.10 不一致，应以 `go.mod` 为准。

### 2.2 项目规模

- 仓库约 81 个文件。
- 共有 26 个 Go 文件。
- 对外 HTTP 接口仅有 `GET /health` 和 `POST /loan`。
- 核心业务集中在贷款 Controller 和 MongoDB 查询模块。
- 仓库存在模型层和 Repository 接口，但尚未真正接入调用链。

### 2.3 核心调用链

```text
POST /loan
  → Gin 全局中间件
  → JSON 参数校验
  → API Key、时间戳和签名校验
  → MongoDB 按 phone 查询字段 t
  → gzip 解压字段 v
  → 解析逗号分隔数据
  → 映射为 var100001～var100006
  → 返回 JSON 响应
```

## 3. 代码结构

| 目录或文件 | 职责 | 当前状态 |
| --- | --- | --- |
| `main.go` | 配置加载、Gin 初始化、HTTP Server 和优雅关闭 | 已接入 |
| `bootstrap/` | 注册全局中间件及路由 | 已接入 |
| `routes/` | 健康检查和贷款路由 | 已接入 |
| `app/Http/controllers/loan/` | 请求校验、签名、Mongo 查询和响应构造 | 职责过重 |
| `app/middleware/` | CORS、日志、限流和认证 | 部分为占位实现 |
| `app/models/` | 贷款模型及通用响应 | 未进入核心调用链 |
| `app/repositories/` | Repository 接口 | 无实现、未使用 |
| `pkg/config/` | 应用配置 | 已接入；主程序执行必填校验，但配置入口仍重复 |
| `pkg/env/` | 重复读取 MongoDB 环境变量 | 已接入，和 config 职责重复 |
| `pkg/mongo/` | Mongo 连接池、查询、解压和结果映射 | 已接入 |
| `pkg/redis/` | Redis 客户端封装 | 未接入业务 |
| `pkg/rescode/` | 业务错误码 | 已接入 |
| `pkg/utils/` | SID、时间戳和签名工具 | 已接入 |
| `app/protos/helloworld/` | gRPC 示例 | 未接入实际服务 |

## 4. 验证基线

本次分析执行了以下命令：

```bash
go test ./...
go test -race ./...
go vet ./...
```

验证结果均为通过。

当前共有 10 个单元测试、1 个 Protobuf 测试和 1 个 benchmark，主要覆盖：

- SID 格式、长度和单线程唯一性；
- 时间戳与签名生成、校验；
- 错误码分类与消息；
- HelloWorld Protobuf 生成代码。

以下关键模块目前没有测试：

- `/loan` Controller 请求与响应行为；
- Gin 路由和中间件组合；
- MongoDB 查询、超时、断线和恢复；
- gzip 数据损坏及字段数量异常；
- 限流边界；
- 并发调用 SID 生成器。

需要特别说明：现有 `go test -race ./...` 虽然通过，但测试没有并发调用 SID 生成逻辑，因此不能证明该实现并发安全。

## 5. 主要风险

### 5.1 已整改：签名密钥配置化

签名密钥已改为通过 `.env` 的 `SECRET_KEY` 加载，并在应用启动时执行必填校验。生产环境仍应通过密钥管理系统注入强随机值，避免直接使用示例值；若 API Key 对应不同客户，后续还应改为按客户取得对应密钥。

### 5.2 P0：日志记录完整敏感数据

位置：`app/middleware/logger.go`

日志中间件会完整读取并记录请求体和响应体，其中可能包含：

- phone；
- apikey；
- sign；
- 查询结果数据。

同一成功请求的请求体还会分别在开始和完成日志中记录两次。

风险：

- 个人信息和凭据泄露；
- 大请求造成额外内存占用；
- 高 QPS 下形成严重日志和序列化开销；
- 响应数据量增长后可能放大磁盘和日志采集成本。

建议：默认只记录路径、耗时、HTTP 状态、业务码和 SID；对 phone、apikey、sign 做删除或掩码处理；增加 body 大小上限和可配置采样率。

### 5.3 P0：API Key 未进行真实性校验

当前 `/loan` 只检查 `apikey` 是否为空，没有校验其是否存在、启用或属于合法客户。`AuthRequired` 和 `APIKeyAuth` 中间件同样只是占位实现，且没有注册到 `/loan`。

建议：明确认证模型，在签名验证前查询 API Key 对应的主体、状态和密钥；认证失败统一返回认证错误，不进入 MongoDB 查询。

### 5.4 P1：SID 随机源并发不安全

位置：`pkg/utils/sid.go`

代码共享一个由 `rand.New` 创建的 `math/rand.Rand`。该对象不支持并发调用，多个 HTTP 请求同时生成 SID 时存在数据竞争风险。

建议：使用 `crypto/rand`、UUID 或 ULID；增加并发测试。如果 SID 用于安全、审计或幂等判断，不应继续使用 `math/rand`。

### 5.5 P1：MongoDB 首次初始化失败后无法恢复

位置：`pkg/mongo/mongo.go`

Mongo 客户端通过 `sync.Once` 初始化。若第一次连接或 Ping 失败，`sync.Once` 仍会被标记为已执行，后续请求永久返回同一个错误，只能通过重启进程恢复。

建议：

- 优先在应用启动阶段完成 MongoDB 初始化并快速失败；或
- 使用可重试、可替换的客户端生命周期管理；
- Ping 失败时及时断开已创建的客户端；
- 为连接和查询分别配置超时。

### 5.6 P1：错误分类不准确并泄露内部信息

Controller 当前将 Mongo 连接失败、查询超时、文档不存在、gzip 解压失败和数据解析失败统一返回为 `ErrDataNotFound`，HTTP 状态均为 200，并将底层错误文本直接返回客户端。

影响：

- 客户端无法判断是否应该重试；
- 监控会把基础设施故障误判为正常未查得；
- 内部数据库和数据格式信息可能泄露；
- HTTP 指标无法反映真实服务故障。

建议：定义稳定的领域错误，至少区分未查得、参数错误、认证失败、数据库不可用、超时和数据损坏；客户端返回安全文案，详细错误仅写内部日志。

### 5.7 部分整改：启动配置校验

主程序现已在加载配置后调用 `AppConfig.Validate()`，MongoDB DSN、数据库名或签名密钥缺失时会终止启动。`config.Load()` 当前仍永远返回 `nil`，配置加载和校验分为两个调用步骤。

建议后续进一步完成 MongoDB 连接初始化，确保关键依赖不可用时在启动阶段明确失败。

### 5.8 P2：Controller 职责过重

贷款 Controller 同时负责：

- 请求 DTO；
- 类型转换和参数校验；
- API Key 检查；
- 时间戳和签名验证；
- 环境配置读取；
- MongoDB 数据访问；
- 错误映射；
- HTTP 响应构造。

这使得业务逻辑难以独立测试，也使已有的模型和 Repository 层失去作用。

建议拆分为：

```text
HTTP Handler
  → Request Validator / Authenticator
  → Loan Service
  → Loan Repository
  → MongoDB
```

### 5.9 P2：配置入口重复

MongoDB 配置同时存在于：

- `pkg/config`；
- `pkg/env`；
- Controller 构造的 `mongo.Config`；
- `GetMongo` 内部再次读取环境变量。

传入 `GetMongo` 的 DSN 和 DB 实际会被内部环境变量覆盖，接口语义不清晰。

建议只保留统一的配置加载入口，并通过依赖注入将数据库或 Repository 传给 Service。

### 5.10 P2：未使用或半成品模块较多

- Redis 已封装并配置，但没有进入业务链路；
- `LoanRepository` 只有接口，没有实现；
- `LoanRequest` 模型未使用；
- gRPC 只有 HelloWorld 示例；
- `connectToMongoDB` 是未使用的旧连接实现；
- 认证中间件只有占位逻辑。

这些代码会增加理解成本并制造“功能已经具备”的错觉。应明确近期用途，否则删除或移入示例目录。

## 6. 性能与可观测性判断

仓库已有历史压测报告，但结果受测试机器、日志量、MongoDB 数据规模和连接配置影响，不能直接作为当前生产容量承诺。

当前最明显的性能风险不是 Gin 本身，而是：

- 每个请求完整复制请求体；
- 完整缓存并序列化响应体；
- 请求体重复记录；
- 高 QPS 下日志写入量过大；
- MongoDB 查询缺少明确的单请求超时；
- gzip 解压结果没有大小上限。

当前日志也缺少稳定的 request ID/SID 关联、错误类型、MongoDB 耗时和结果类型等字段。建议在压测前先调整日志策略，再重新建立容量基线。

## 7. 文档与仓库一致性

目前存在以下不一致：

- README 写 Go 1.18+，实际要求 Go 1.24；
- README 包含 Docker 构建和运行说明，仓库没有 Dockerfile；
- README 宣称 Redis 连接池属于功能特性，但当前核心流程没有初始化或使用 Redis；
- README 中的认证示例没有反映 `/loan` 当前真实注册方式；
- 项目根目录缺少 `AGENTS.md` 等工程规则文件；
- 仓库跟踪了约 27 MB 的 macOS arm64 编译产物 `main`，增加仓库体积且不具备跨平台复用价值；约 2.7 MB 的 `access.log` 已被忽略但会污染本地工作区。

### 7.1 AI Coding 文档体系

2026-08-03 盘点时，仓库不存在 `docs/`、`docs/specs/` 或 `docs/knowledges/`，只有既有的
扁平 `doc/` 目录。当时在根目录增加 `AGENTS.md` 和文档规则；后续文档已按
`doc/rules/`、`doc/specs/` 和 `doc/knowledge/` 迁移。当前入口为：

- `../README.md`：Rules、Specs、Knowledge 总索引和冲突处理；
- `../rules/project_rules.md`：长期项目规则；
- `../rules/spec_guidelines.md`：Spec 命名、模板和生命周期；
- `README.md`：现有知识文件索引和加载规则。

### 7.2 现状审计清单

| 优先级 | 问题与证据 | 影响 | 建议 |
| --- | --- | --- | --- |
| P0 | `app/middleware/logger.go` 完整读取并记录请求/响应，Loan body 含 phone、apikey、sign | 敏感信息泄露；无界内存和高 QPS 日志放大 | 先建安全 Spec，删除或掩码敏感字段并限制采集大小 |
| P0 | `app/Http/controllers/loan/index.go` 只检查 apikey 非空；`auth.go` 是未挂载占位代码 | 不能确认调用主体、状态和权限 | 确认 API Key 权威来源及客户密钥绑定后实现认证 |
| P1 | `pkg/utils/sid.go` 跨请求共享非并发安全的 `math/rand.Rand` | 并发数据竞争，SID 可靠性不可证明 | 使用并发安全随机源并增加并发测试 |
| P1 | `pkg/mongo/mongo.go` 使用 `sync.Once` 固化首次连接/Ping 结果 | 首次故障后不能进程内恢复 | 启动初始化或设计可重试、可替换生命周期 |
| P1 | Loan 把 Mongo、未查得、解压等错误统一映射为 HTTP 200 + `ErrDataNotFound` 并返回原始错误 | 误导调用方和监控，泄露内部信息 | 先确认现有客户契约，再区分稳定领域错误 |
| P1 | Controller、路由、中间件、Mongo 和配置无测试 | 核心行为、并发和故障回归无法发现 | 优先补 `httptest` 组件测试和隔离 Mongo 集成测试 |
| P1 | `pcode` 从 JSON `float64` 直接转 `int` | 小数可能被截断，输入校验不严格 | 以契约测试固定整数要求后修正 DTO/解析 |
| P1 | `cmd/init_mongo.go` 对 16 个同名集合执行 `Drop` | 误用会造成数据丢失 | 仅授权测试环境运行，增加目标确认与回滚方案 |
| P2 | `pkg/config`、`pkg/env`、Controller 和 `GetMongo` 重复读取/覆盖配置 | 语义混乱，测试困难，运行配置可能不一致 | 在独立重构 Spec 中统一配置边界 |
| P2 | Redis、Repository/Model、gRPC 示例和旧 Mongo 连接函数未进入主链路 | 增加维护成本并造成能力误判 | 由负责人确认规划后再保留、隔离或删除 |
| P2 | `README.md` 的 Go 1.18+、JWT/API Key、Redis、Docker 和限流描述与代码不一致 | 新成员和 AI 使用错误事实 | 单独审查 README；本次只记录，不擅自改产品说明 |
| P2 | 无 CI、项目级 Lint、漏洞扫描、覆盖率门槛或部署清单 | 质量门禁依赖人工执行 | 待团队确认平台和门槛后通过 Spec 引入 |
| P2 | 跟踪约 27 MB 的本机 Mach-O `main` 二进制；`.gitignore` 未覆盖根构建产物 | 仓库膨胀、跨平台无效 | 确认发布流程后单独清理历史与忽略规则 |
| P2 | `/health` 只返回 `ok`，无 readiness/liveness、Metrics、Tracing 或 Request ID | 依赖故障不可见，排障关联不足 | 先定义运维语义和平台，再补可观测能力 |
| P2 | 历史 `tests/bench/bench.sh` 请求体只有 phone，与当前 Loan 必填字段不兼容 | 脚本当前不能有效压测成功业务路径 | 在性能 Spec 中修订测试数据生成与安全注入 |

## 8. 建议实施路线

### 阶段一：安全与正确性

1. 为不同环境配置独立的强随机签名密钥，并建立轮换机制。
2. 实现 API Key 真实性校验及其和签名密钥的绑定关系。
3. 日志删除敏感字段，限制请求和响应采集大小。
4. SID 改用并发安全实现。
5. 区分业务未查得和系统异常，停止向客户端返回底层错误。

### 阶段二：稳定性

1. 启动时校验 MongoDB 配置并完成连接初始化。
2. 重构 Mongo 客户端生命周期，解决首次失败后无法恢复的问题。
3. 为 MongoDB 查询增加独立超时。
4. 增加 gzip 解压大小及数据字段完整性校验。
5. 明确 HTTP 状态码与业务错误码规范。

### 阶段三：结构治理

1. 抽取明确的 LoanRequest DTO。
2. 引入 LoanService，迁移认证后的业务编排逻辑。
3. 实现并使用 LoanRepository。
4. 合并 `pkg/config` 和 `pkg/env` 的重复职责。
5. 清理或明确 Redis、gRPC、旧 Mongo 连接代码和占位认证模块。

### 阶段四：测试与交付

1. 为 Controller 增加表驱动参数校验测试。
2. 使用 `httptest` 覆盖 `/health`、`/loan` 和中间件组合。
3. 增加 MongoDB 集成测试及异常数据测试。
4. 增加 SID 并发测试和限流边界测试。
5. 生成覆盖率报告并设定最低门槛。
6. 补齐 Dockerfile、健康检查及部署说明。
7. 同步 README，并增加 `AGENTS.md` 或等效开发规则文件。

## 9. 建议验收标准

完成上述整改后，建议至少满足：

- 仓库及日志中不存在明文签名密钥、完整 API Key 和完整手机号；
- 非法或未知 API Key 无法通过认证；
- `go test ./...`、`go test -race ./...`、`go vet ./...` 全部通过；
- `/loan` 的参数、认证、签名、未查得、数据库异常和成功路径均有自动化测试；
- MongoDB 启动失败和运行期故障具有明确行为及可观测指标；
- Controller 不直接读取环境变量或操作 MongoDB；
- README、Go 版本、Docker 部署说明和实际代码保持一致；
- 重新执行压测并记录测试环境、数据规模、日志配置、P95/P99 和错误率。

## 10. 总体评价

项目的优点是规模小、核心流程清晰、依赖数量可控，且已有基础测试和优雅关闭机制，重构成本相对较低。

当前主要短板集中在安全边界、异常语义、并发安全、配置生命周期和业务测试覆盖。建议先解决 P0/P1 风险，再进行分层重构和性能优化，避免在不可靠的基线上继续增加功能。
