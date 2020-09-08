# 用于 Go 的 net 平滑重启/热重启（Graceful Restart）封装

让 Go 程序支持热更新

## 用法

提供两个基于 http server 的简单 Demo

### 热重启 http server

 [Demo http_reload](example/http_reload)
 
```shell script
cd example/http_reload
go run main.go
```

### 热更新到不同版本

[Demo http_switch](example/http_switch)

```shell script
cd example/http_switch
make run
```
