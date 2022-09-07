package redispool

import (
	"context"
	"net"
	"sync"
	"time"

	winner_logger "github.com/bfmTech/logger-go"
	redis "github.com/go-redis/redis/v8"

	"github.com/nerocho/fastfiber/utils/eventmanager"
)

//RedisClient 连接客户端
var redisClientPool *redis.Client
var EventDestroyPrefix = "Destroy_Redis"
var once sync.Once

// 初始化redis 连接池
func GetPool(addr, password string, poolSize, maxIdle, idleTimeout, indexDb int, wl winner_logger.Logger, tracing bool) (*redis.Client, error) {

	once.Do(func() {
		redisClientPool = redis.NewClient(&redis.Options{
			//连接信息
			Network:  "tcp",    //网络类型, tcp 或者 unix, 默认tcp
			Addr:     addr,     //ip:port
			Password: password, //密码,
			DB:       indexDb,  //连接后选中的redis数据库index

			//命令执行失败时的重试策略
			MaxRetries:      3,                      //命令执行失败时最大重试次数，默认3次重试。
			MinRetryBackoff: 8 * time.Millisecond,   //每次重试最小间隔时间，默认8ms，-1表示取消间隔
			MaxRetryBackoff: 512 * time.Millisecond, //每次重试最大时间间隔，默认512ms，-1表示取消间隔

			//超时
			DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒
			ReadTimeout:  3 * time.Second, //读超时，默认3秒，-1表示取消读超时
			WriteTimeout: 3 * time.Second, //写超时，默认与读超时相等

			//连接池容量、闲置连接数量、闲置连接检查
			PoolSize:           poolSize,                                 //连接池最大Socket连接数，默认为10倍CPU数量，10 * runtime.NumCPU()
			MinIdleConns:       maxIdle,                                  //启动阶段创建指定数量的Idle连接，并长期维持Idle状态的连接数不少于指定数量。
			MaxConnAge:         0 * time.Second,                          //连接存活时长，超过指定时长则关闭连接。默认为0，不关闭旧连接。
			PoolTimeout:        4 * time.Second,                          //当所有连接都处于繁忙状态时，客户端等待可用连接的最大等待时长。默认为读超时+1秒
			IdleTimeout:        time.Duration(idleTimeout) * time.Minute, //关闭闲置连接时间，默认5分钟，-1表示取消闲置超时检查
			IdleCheckFrequency: 1 * time.Minute,                          //闲置连接检查周期，默认为1分钟；-1表示不做检查，只在客户端获取连接时对闲置连接进行处理。

			//自定义连接函数
			Dialer: func(ctx context.Context, network string, addr string) (net.Conn, error) {
				netDialer := net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 5 * time.Minute,
				}
				return netDialer.Dial(network, addr)
			},

			//钩子函数，建立新连接时调用
			// OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			// 	return nil
			// },
		})
		// 将redis的关闭事件，注册在全局事件统一管理器，由程序退出时统一销毁
		eventmanager.CreateEventManageFactory().Set(EventDestroyPrefix, func(args ...interface{}) {
			_ = redisClientPool.Close()
		})

		// 开启日志
		if tracing {
			redisClientPool.AddHook(NewRedisTracingHook(wl))
		}
	})

	//ping
	_, err := redisClientPool.Ping(context.Background()).Result()

	return redisClientPool, err
}
