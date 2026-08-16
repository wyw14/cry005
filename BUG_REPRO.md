# Bug reproduction

## Bug

同一加入请求并发重试时会重复占位并返回多个成员记录。

## Trigger

在项目根目录执行：

```bash
go test ./internal/application -run ^TestConcurrentJoinUsesOneMembership$ -count=1
```

## Error

`application_test.go:72: idempotent requests returned different IDs`：同一幂等加入请求返回了两个不同的成员 ID。
