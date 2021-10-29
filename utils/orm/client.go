package orm

import (
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	gormopentracing "gorm.io/plugin/opentracing"

	"github.com/nerocho/fastfiber/utils/eventmanager"
)

const (
	ErrorsDbDriverNotExists   = "Database.Type 数据库驱动不被支持,请选择mysql|postgre:"
	ErrorsDialectorDbInitFail = "数据库驱动初始化失败:"
	EventDestroyPrefix        = "Destroy_Database"
)

// 根据默认配置获取数据库实例，需要自行创建数据库实例，请使用GetSqlDriver
// func GetDefaultDb(tracing bool) (*gorm.DB, error) {
// 	ops := &DbOptions{
// 		SqlType:        "Database.Type",
// 		Dsn:            "Database.Dsn",
// 		EnableReplicas: "Database.EnableReplicas",
// 		MaxIdle:        "Database.MaxIdleConns",
// 		MaxIdleTime:    "Database.MaxIdleTime",
// 		MaxOpen:        "Database.MaxOpenConns",
// 	}

// 	return GetSqlDriver(ops, tracing)
// }

type DbOptions struct {
	SqlType        string
	Dsn            string
	Replicas       []string
	EnableReplicas bool
	MaxIdle        int
	MaxIdleTime    time.Duration
	MaxLifeTime    time.Duration
	MaxOpen        int
	SlowThreshold  time.Duration
}

// 获取数据库驱动
// ops := &DbOptions{
// 	SqlType:        "Database.Type",
// 	Dsn:            "Database.Dsn",
// 	EnableReplicas: "Database.EnableReplicas",
// 	MaxIdle:        "Database.MaxIdleConns",
// 	MaxIdleTime:    "Database.MaxIdleTime",
// 	MaxOpen:        "Database.MaxOpenConns",
// }
func GetSqlDriver(options *DbOptions, logger zerolog.Logger, tracing bool) (*gorm.DB, error) {
	var dbDialector gorm.Dialector
	if val, err := getDbDialector(options.SqlType, options.Dsn); err != nil {
		logger.Error().Err(err).Msg(ErrorsDialectorDbInitFail + options.Dsn)
	} else {
		dbDialector = val
	}

	gormDb, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 redefineLog(options.SlowThreshold, logger), //拦截、接管 gorm v2 自带日志
	})
	//gorm 数据库驱动初始化失败
	if err != nil {
		return nil, err
	}

	// 读写分离配置
	if options.EnableReplicas {
		resolverConf := getReplicas(options.SqlType, options.Replicas, logger)

		err = gormDb.Use(dbresolver.Register(*resolverConf).
			SetConnMaxIdleTime(time.Second * options.MaxIdleTime).
			SetConnMaxLifetime(options.MaxLifeTime * time.Hour).
			SetMaxIdleConns(options.MaxIdle).
			SetMaxOpenConns(options.MaxOpen))
		if err != nil {
			return nil, err
		}
	}

	// 查询没有数据，屏蔽 gorm v2 包中会爆出的错误
	_ = gormDb.Callback().Query().Before("gorm:query").Register("disable_raise_record_not_found", func(d *gorm.DB) {
		d.Statement.RaiseErrorOnNotFound = false
	})

	// tracing
	if tracing {
		gormDb.Use(gormopentracing.New())
	}

	// 为主连接设置连接池
	if rawDb, err := gormDb.DB(); err != nil {
		return nil, err
	} else {
		rawDb.SetMaxIdleConns(options.MaxIdle)
		rawDb.SetConnMaxLifetime(time.Hour * options.MaxLifeTime)
		rawDb.SetConnMaxIdleTime(time.Second * options.MaxIdleTime)
		rawDb.SetMaxOpenConns(options.MaxOpen)

		eventmanager.CreateEventManageFactory().Set(EventDestroyPrefix, func(args ...interface{}) {
			_ = rawDb.Close()
		})

		return gormDb, nil
	}
}

// 获取驱动
func getDbDialector(sqlType, dsn string) (gorm.Dialector, error) {
	var dbDialector gorm.Dialector
	switch strings.ToLower(sqlType) {
	case "mysql":
		dbDialector = mysql.Open(dsn)
	case "postgres", "postgresql", "postgre":
		dbDialector = postgres.Open(dsn)
	default:
		return nil, errors.New(ErrorsDbDriverNotExists + sqlType)
	}
	return dbDialector, nil
}

// 获取读节点，遍历时会忽略连接不上的节点
func getReplicas(sqlType string, replicas []string, logger zerolog.Logger) *dbresolver.Config {
	var dialectors []gorm.Dialector
	for i := 0; i < len(replicas); i++ {
		dsn := replicas[i]
		if val, err := getDbDialector(sqlType, dsn); err != nil {
			logger.Error().Err(err).Msg(ErrorsDialectorDbInitFail + dsn)
		} else {
			dialectors = append(dialectors, val)
		}
	}
	return &dbresolver.Config{
		Replicas: dialectors,
		Policy:   dbresolver.RandomPolicy{},
	}

}

// 创建自定义日志模块，对 gorm 日志进行拦截、
func redefineLog(slowThreshold time.Duration, zerowriter zerolog.Logger) gormLog.Interface {
	return createCustomGormLog(slowThreshold, zerowriter,
		SetInfoStrFormat("[info] %s\n"),
		SetWarnStrFormat("[warn] %s\n"),
		SetErrStrFormat("[error] %s\n"),
		SetTraceStrFormat("[traceStr] %s [%.3fms] [rows:%v] %s\n"),
		SetTraceWarnStrFormat("[traceWarn] %s %s [%.3fms] [rows:%v] %s\n"),
		SetTraceErrStrFormat("[traceErr] %s %s [%.3fms] [rows:%v] %s\n"))
}
