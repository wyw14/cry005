# Bug reproduction

## Bug

不同房间的数据会在成员列表中串在一起，邀请码查询也会越界。

## Trigger

在项目根目录执行：

```bash
go test ./internal/application -run ^TestRoomMembersStayWithinRoom$ -count=1
```

## Error

`application_test.go:37: scope leak`：同一次 alpha 房间成员查询结果同时包含 `Scope:beta` 与 `Scope:alpha` 的记录。
