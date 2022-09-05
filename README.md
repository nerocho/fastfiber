# fastfiber

> [fastfiber](https://github.com/nerocho/fastfiber)是一个基于[fiber](https://github.com/gofiber/fiber)的快速开始`goweb`开发的脚手架。目的是为了让开发者更快的进入开发，而无需关注`项目组织`,`系统配置`,`驱动集成`,`日志`等基础模块。

## QuickStart

详细示例 [fastfiber-demo](https://github.com/nerocho/fastfiber-demo)

```bash
go get -u github.com/nerocho/fastfiber@latest
```

## Features

- 很棒的性能，参见[fasthttp](https://github.com/valyala/fasthttp)
- 低内存占用
- 集成 `viper` 的配置文件管理
- 集成 `zerolog` 的高性能日志输出
- 集成 `gorm` 的数据库操作
- 集成 `go-redis` 的缓存操作
- 自动的参数校验及全局错误处理
- `Graceful Shutdown`
- 开箱即用的中间件&工具包
  - 访问日志,限流,跨域,JWT,接口级缓存...
  - 邮件,tracing,response,加解密...

## 配置文件 config.env.yml

配置文件依赖于环境变量`APP_ENV`,如`APP_ENV=test`,则配置文件会读取`config.test.yml`

```yaml
# 系统配置
System:
  AppName: 'go-express' # 应用名称
  Port: 3000 # 启动端口
  LogResponseBody: false # 是否打印response日志，默认关闭，开启比较影响性能和存储，比如返回给前端数据很大时，记录日志消耗会比较大

# 数据库配置
Database:
  IsInit: true # 是否初始化到fastfiber对象上
  Type: 'mysql' # 仅支持mysql 和 postgre
  SlowThreshold: 100 # 慢日志，单位毫秒，执行时间大于SlowThreshold的sql会被记录到日志中
  MaxIdleConns: 10 # 最大空闲连接数即一直持有
  MaxIdleTime: 1800 # 默认30分钟，最大空闲时间秒
  MaxOpenConns: 128 # 连接池大小
  EnableReplicas: false # 是否开启读写分离
  Dsn:
    Write: 'root:123456@tcp(127.0.0.1:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local' # mysql
    # Write: 'host=127.0.0.1 port=5432 dbname=dev user=root password=123456 sslmode=disable TimeZone=Asia/Shanghai' # postgre
    Read:
      - 'root:123456@tcp(127.0.0.2:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local'
      - 'root:123456@tcp(127.0.0.3:3306)/dev?charset=utf8mb4&parseTime=True&loc=Local'

# redis 配置
Redis:
  IsInit: true # 是否初始化到fastfiber对象上
  Addr: '127.0.0.1:6379'
  Password: ''
  MaxIdle: 10 #最大空闲连接数
  MaxActive: 1000 # 连接池大小
  IdleTimeout: 60 #空闲超时时间
  IndexDb: 0 #数据库

# Id生成器
IdWorker:
  IsInit: true # 是否初始化到全局对象
  WorkerId: 0 # 为0 则使用默认值 建议按照节点数量自行设置
  DataCenterId: 0 # 为0 则使用默认值 建议按照业务线自行设置
  Twepoch: 0 # 为0 则使用默认值

# 限流配置
Limiter:
  Enable: true # true开启 false 关闭
  IpWhiteList: # 白名单
    - '127.0.0.1'
    - '172.16.1.1'

# 接口缓存配置
ApiCache:
  Enable: true #true开启 false 关闭

# Tracing
Tracer:
  Enable: false
  HostPort: 127.0.0.1:6831 # jaeger server
```

## 数据库操作

```bash
// 安装模型类生成工具
go install github.com/xxjwxc/gormt@latest

// 项目目录下执行，默认会把所有的dev库里面所有的表生成到models文件夹下
gormt -H=127.0.0.1 -d=dev -p=123456 -u=root --port=3306 -F=true -o=models

// 具体代码可以查看相关文件夹
```
## Stargazers

如果您觉得本项目对您有所帮助，请不要吝啬一个⭐哦！

[![Stargazers over time](https://starchart.cc/nerocho/fastfiber.svg)](https://starchart.cc/nerocho/fastfiber)
