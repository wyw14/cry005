# 多人游戏房间与匹配服务

房间生命周期、并发加入、房主转移、SSE 事件和超时解散。项目是可离线运行的 Go 1.24 演示基线，业务服务使用 Gin、GORM 仓储适配、Validator、Zap 与 Viper；测试默认使用线程安全内存仓储，部署配置对应 `MySQL 8`。

## 目录

- `cmd/server`：HTTP 服务入口与优雅停机
- `internal/domain`：状态、作用域、幂等和事务快照规则
- `internal/application`：用例编排和仓储端口
- `internal/repository`：内存实现与 GORM SQL 适配
- `internal/transport/http`：`/api/v1` 路由和稳定错误码
- `migrations`、`api/openapi`、`deploy`：迁移、接口与容器交付
- `web`：Vue 3 + TypeScript + Vite 演示控制台

## 本地验证

```bash
go build ./...
go test ./...
go test -race ./...
go vet ./...
cd web && npm ci && npm test && npm run build
```

## 启动

```bash
go run ./cmd/server
curl http://localhost:8080/healthz
```

配置通过环境变量加载，示例见 `.env.example`。仓库不包含真实密钥、运行数据库、依赖缓存或构建产物。
