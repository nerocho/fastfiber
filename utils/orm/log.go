package orm

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	winner_logger "github.com/bfmTech/logger-go"
	gormLog "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var writer winner_logger.Logger

// 自定义日志格式, 对 gorm 自带日志进行拦截重写
func createCustomGormLog(slowThreshold time.Duration, winnerWriter winner_logger.Logger, options ...Options) gormLog.Interface {

	writer = winnerWriter

	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)
	logConf := gormLog.Config{
		SlowThreshold:             time.Millisecond * slowThreshold,
		LogLevel:                  gormLog.Warn,
		IgnoreRecordNotFoundError: true,
		Colorful:                  false,
	}
	log := &logger{
		Writer:       logOutPut{},
		Config:       logConf,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
	for _, val := range options {
		val.apply(log)
	}
	return log
}

type logOutPut struct{}

// 自定义格式化
func (l logOutPut) Printf(strFormat string, args ...interface{}) {
	logRes := fmt.Sprintf(strFormat, args...)
	if strings.HasPrefix(strFormat, "[info]") || strings.HasPrefix(strFormat, "[traceStr]") {
		writer.Info("[database] ", logRes)
	} else if strings.HasPrefix(strFormat, "[error]") || strings.HasPrefix(strFormat, "[traceErr]") {
		writer.Error(errors.New("[database] " + logRes))
	} else if strings.HasPrefix(strFormat, "[warn]") || strings.HasPrefix(strFormat, "[traceWarn]") {
		writer.Warn("database", logRes)
	}
}

// 尝试从外部重写内部相关的格式化变量
type Options interface {
	apply(*logger)
}
type OptionFunc func(log *logger)

func (f OptionFunc) apply(log *logger) {
	f(log)
}

// 设置 一般信息 日志格式
func SetInfoStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.infoStr = format
	})
}

// 设置 警告信息 日志格式
func SetWarnStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.warnStr = format
	})
}

// 设置 错误信息 日志格式
func SetErrStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.errStr = format
	})
}

// 设置 追踪信息 日志格式
func SetTraceStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.traceStr = format
	})
}

// 设置 追踪警告信息 日志格式
func SetTraceWarnStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.traceWarnStr = format
	})
}

// 设置 追踪错误信息 日志格式
func SetTraceErrStrFormat(format string) Options {
	return OptionFunc(func(log *logger) {
		log.traceErrStr = format
	})
}

//日志对象
type logger struct {
	gormLog.Writer
	gormLog.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// 日志模式
func (l *logger) LogMode(level gormLog.LogLevel) gormLog.Interface {
	nl := *l
	nl.LogLevel = level
	return &nl
}

// 一般信息
func (l logger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLog.Info {
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// 警告信息
func (l logger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLog.Warn {
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// 错误信息
func (l logger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLog.Error {
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// 打印sql信息
func (l logger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= gormLog.Error:
			sql, rows := fc()
			if rows == -1 {
				l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-1", sql)
			} else {
				l.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLog.Warn:
			sql, rows := fc()
			slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-1", sql)
			} else {
				l.Printf(l.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel >= gormLog.Info:
			sql, rows := fc()
			if rows == -1 {
				l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-1", sql)
			} else {
				l.Printf(l.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}
