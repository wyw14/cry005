# Bug reproduction

## Bug

超时解散任务取消后仍会关闭刚刚恢复活跃的房间。

## Trigger

在项目根目录执行：

```bash
go test ./internal/application -run ^TestCanceledExpiryJobDoesNotCloseRoom$ -count=1
```

## Error

`application_test.go:102: error=<nil>, want context.Canceled`
