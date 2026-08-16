# Bug reproduction

## Bug

开始游戏写入成员状态失败后房间已进入 playing，成员却仍未准备。

## Trigger

在项目根目录执行：

```bash
go test ./internal/application -run ^TestStartGameFailureRollsBack$ -count=1
```

## Error

`application_test.go:119: partial commit: before={ID:workflow State:queued Secondary:ready Version:1} after={ID:workflow State:completed Secondary:ready Version:2}`
