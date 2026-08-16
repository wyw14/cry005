# Bug reproduction

## Bug

房间事件断线续传会漏掉 closed 终态，客户端一直显示游戏中。

## Trigger

在项目根目录执行：

```bash
go test ./internal/application -run ^TestEventReplayPreservesOrder$ -count=1
```

## Error

`application_test.go:91: bad replay: []`
