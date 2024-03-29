package fastfiber

import (
	"flag"
	"log"
	"net/url"
	"strings"

	winner_logger "github.com/bfmTech/logger-go"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"github.com/kos-v/dsnparser"
	"github.com/nerocho/fastfiber/interf"
	"github.com/nerocho/fastfiber/utils/orm"
	"github.com/nerocho/fastfiber/utils/redispool"
	"github.com/nerocho/fastfiber/utils/snowflake"
)

var (
	// globals
	Logger    winner_logger.Logger   //全局日志
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

	appName := Conf.GetString("System.AppName")

	//绑定日志模块
	Logger = initSlSLogger(appName, Conf.GetString("System.LogType"))

	//初始化数据库
	if Conf.GetBool("Database.IsInit") {

		ops := &orm.DbOptions{
			SqlType:        Conf.GetString("Database.Type"),
			EnableReplicas: Conf.GetBool("Database.EnableReplicas"),
			MaxIdle:        Conf.GetInt("Database.MaxIdleConns"),
			MaxIdleTime:    Conf.GetDuration("Database.MaxIdleTime"),
			MaxLifeTime:    Conf.GetDuration("Database.MaxLifeTime"),
			MaxOpen:        Conf.GetInt("Database.MaxOpenConns"),
			SlowThreshold:  Conf.GetDuration("Database.SlowThreshold"),
		}

		// 配置写链接
		ops.Dsn = orm.ParseDsn(ops.SqlType, GetEnv(Conf.GetString("Database.Dsn.Write")))

		// 配置读库
		if ops.EnableReplicas {
			replicas := strings.Split(GetEnv(Conf.GetString("Database.Dsn.Read")), ",")
			for i := range replicas {
				replicas[i] = orm.ParseDsn(ops.SqlType, replicas[i])
			}
			ops.Replicas = replicas
		}

		if db, err := orm.GetSqlDriver(ops, Logger, Conf.GetBool("Database.EnableSQLLog")); err != nil {
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
		s := GetEnv(Conf.GetString("Redis.Addr"))
		before := dsnparser.Parse(s)
		addr := before.GetHost() + ":" + before.GetPort()

		pwd, _ := url.QueryUnescape(before.GetPassword())

		if redisPool, err := redispool.GetPool(addr, pwd, Conf.GetInt("Redis.MaxActive"), Conf.GetInt("Redis.MaxIdle"), Conf.GetInt("Redis.IdleTimeout"), Conf.GetInt("Redis.indexDb"), Logger, Conf.GetBool("Redis.EnableTraceLog")); err != nil {
			log.Fatal(ErrorsRedisInitConnFail + err.Error())
		} else {
			RedisPool = redisPool
		}
	}

}
