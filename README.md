# hackathon-68 WhisperTrail

## 项目简介
WhisperTrail 是一个图片分享平台，使用交互式地图与世界各地的人分享美好

## 运行说明
```bash
go mod tidy
go run main.go
```
使用 nginx 作为反向代理服务器将 localhost:8080 的后端 和 localhost:5173 的前端程序在 localhost:80 监听
具体配置见 `nginx.conf`

数据库使用 MySQL, 配置文件见 `config.yaml`
