package orm

import (
	"strings"
	"time"

	winner_logger "github.com/bfmTech/logger-go"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/nerocho/fastfiber/utils/eventmanager"
)

const (
	ErrorsDbDriverNotExists   = "Database.Type 数据库驱动不被支持,请选择mysql|postgres:"
	ErrorsDialectorDbInitFail = "数据库驱动初始化失败:"
	EventDestroyPrefix        = "Destroy_Database"
)

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
// ops := &DbOptions
// wl winner_logger.Logger
// tracing 是否开启sql日志
func GetSqlDriver(options *DbOptions, wl winner_logger.Logger, tracing bool) (*gorm.DB, error) {
	var dbDialector gorm.Dialector
	if val, err := getDbDialector(options.SqlType, options.Dsn); err != nil {
		wl.Error(errors.WithMessage(err, ErrorsDialectorDbInitFail+options.Dsn))
	} else {
		dbDialector = val
	}

	gormDb, err := gorm.Open(dbDialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Silent), // 禁用日志
	})
	//gorm 数据库驱动初始化失败
	if err != nil {
		return nil, err
	}

	// 读写分离配置
	if options.EnableReplicas {
		resolverConf := getReplicas(options.SqlType, options.Replicas, wl)
		err = gormDb.Use(dbresolver.Register(*resolverConf).
			SetConnMaxIdleTime(time.Second * options.MaxIdleTime).
			SetConnMaxLifetime(time.Second * options.MaxLifeTime).
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

	// 是否开启SQL日志
	if tracing {
		gormDb.Use(NewMysqlTracingPlugin(wl))
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
func getReplicas(sqlType string, replicas []string, wl winner_logger.Logger) *dbresolver.Config {
	var dialectors []gorm.Dialector
	for i := 0; i < len(replicas); i++ {
		dsn := replicas[i]
		if val, err := getDbDialector(sqlType, dsn); err != nil {
			wl.Error(errors.WithMessage(err, ErrorsDialectorDbInitFail+dsn))
		} else {
			dialectors = append(dialectors, val)
		}
	}
	return &dbresolver.Config{
		Replicas: dialectors,
		Policy:   dbresolver.RandomPolicy{},
	}
}
