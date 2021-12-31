package fastfiber

import (
	"flag"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"github.com/nerocho/fastfiber/interf"
	"github.com/nerocho/fastfiber/utils/orm"
	"github.com/nerocho/fastfiber/utils/redispool"
	"github.com/nerocho/fastfiber/utils/snowflake"
	"github.com/nerocho/fastfiber/utils/tracer"
)

var (
	// globals
	Logger    zerolog.Logger         //全局日志
	Conf      interf.ConfigInterface //全局配置
	Db        *gorm.DB               //数据库
	IdWorker  *snowflake.IdWorker    //id 生成器
	RedisPool *redis.Client          //全局redisPool

	configPath = flag.String("f", ".", "配置文件目录")
)

func Bootstrap() {
	flag.Parse()

	//绑定配置模块
	Conf = initConfig(*configPath)

	//绑定日志模块
	Logger = initZerolog()

	//初始化数据库
	if Conf.GetBool("Database.IsInit") {
		ops := &orm.DbOptions{
			SqlType:        Conf.GetString("Database.Type"),
			Dsn:            Conf.GetString("Database.Dsn.Write"),
			EnableReplicas: Conf.GetBool("Database.EnableReplicas"),
			Replicas:       Conf.GetStringSlice("Database.Dsn.Read"),
			MaxIdle:        Conf.GetInt("Database.MaxIdleConns"),
			MaxIdleTime:    Conf.GetDuration("Database.MaxIdleTime"),
			MaxLifeTime:    Conf.GetDuration("Database.MaxLifeTime "),
			MaxOpen:        Conf.GetInt("Database.MaxOpenConns"),
			SlowThreshold:  Conf.GetDuration("Database.SlowThreshold"),
		}

		if db, err := orm.GetSqlDriver(ops, Logger, Conf.GetBool("Tracer.Enable")); err != nil {
			log.Fatal(ErrorsDbInitFail + err.Error())
			return
		} else {
			Db = db
		}
	}

	// 初始化IdWorker
	if Conf.GetBool("IdWorker.IsInit") {
		if w, err := snowflake.NewIdWorker(Conf.GetInt64("IdWorker.WorkerId"), Conf.GetInt64("IdWorker.DataCenterId"), Conf.GetInt64("IdWorker.Twepoch")); err != nil {
			log.Fatal(ErrorsIdWorkerInitFail + err.Error())
		} else {
			IdWorker = w
		}
	}

	// RedisPool
	if Conf.GetBool("Redis.IsInit") {
		if redisPool, err := redispool.GetPool(Conf.GetString("Redis.Addr"), Conf.GetString("Redis.Password"), Conf.GetInt("Redis.MaxActive"), Conf.GetInt("Redis.MaxIdle"), Conf.GetInt("Redis.IdleTimeout"), Conf.GetInt("Redis.indexDb")); err != nil {
			log.Fatal(ErrorsRedisInitConnFail + err.Error())
		} else {
			RedisPool = redisPool
		}
	}

	if Conf.GetBool("Tracer.Enable") {
		_, _, err := tracer.NewJaegerTracer(Conf.GetString("System.AppName"), Conf.GetString("Tracer.HostPort"))
		if err != nil {
			log.Fatal(ErrorsTracerInitFail + err.Error())
		}
	}

}
