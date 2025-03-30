# PS - 简单的订阅发布服务

## 编译

```bash
# 普通编译
go build -o psc ./cmd/psc

# 指定默认服务器地址的编译方式
go build -ldflags "-X 'main.DefaultAddr=localhost:8080'" -o psc ./cmd/psc
```

## 用法

### 命令行工具

1. 运行服务器：

```bash
psc -s -k my-token
```

2. 订阅消息：

```bash
psc -k my-token -n my-namespace -t my-topic
```

3. 发布消息：

```bash
psc -k my-token -n my-namespace -t my-topic -d "Hello World"
```

### 使用 curl 操作

1. 发布消息(pub)：

```bash
curl -X POST "http://localhost:8080/p/namespace/topic" \
     -H "Authorization: your-token" \
     -H "Content-Type: application/json" \
     -d '{"ID": 1648538490000, "Data": "your message", "Timestamp": 1648538490}'
```

2. 订阅消息(sub)：

```bash
curl -N "http://localhost:8080/s/namespace/topic?lastID=0" \
     -H "Authorization: your-token"
```

## 参数说明

```
  -a, --addr string        服务器地址 (默认为 "localhost:8080")
  -d, --data string        要发布的数据
  -h, --help               帮助信息
  -k, --key string         认证令牌 (必填)
  -s, --server             以服务器模式运行
  -n, --namespace string   命名空间
  -t, --topic string       主题
```
