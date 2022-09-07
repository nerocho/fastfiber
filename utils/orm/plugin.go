package orm

import (
	"strconv"
	"time"

	winner_logger "github.com/bfmTech/logger-go"

	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

// 自定义sql日志插件
// 日志格式: tag(固定的标签), reqId(request_id), sql(sql文本), costSeconds(执行时间(ms)), affectedRows(受影响行数), stack(sql执行位置)

const (
	callBackBeforeName = "winnerTracing:before"
	callBackAfterName  = "winnerTracing:after"
	startTime          = "_winnerStartTime"
	tag                = "[sqlInfo]"
	requestId          = "request_id"
)

type WinnerTracingPlugin struct{}

var _logger winner_logger.Logger

// check
var _ gorm.Plugin = &WinnerTracingPlugin{}

func NewMysqlTracingPlugin(logger winner_logger.Logger) gorm.Plugin {
	_logger = logger
	return &WinnerTracingPlugin{}
}

func (op *WinnerTracingPlugin) Name() string {
	return "winnerTracingPlugin"
}

func (op *WinnerTracingPlugin) Initialize(db *gorm.DB) (err error) {
	// 开始前
	db.Callback().Create().Before("gorm:before_create").Register(callBackBeforeName, before)
	db.Callback().Query().Before("gorm:query").Register(callBackBeforeName, before)
	db.Callback().Delete().Before("gorm:before_delete").Register(callBackBeforeName, before)
	db.Callback().Update().Before("gorm:setup_reflect_value").Register(callBackBeforeName, before)
	db.Callback().Row().Before("gorm:row").Register(callBackBeforeName, before)
	db.Callback().Raw().Before("gorm:raw").Register(callBackBeforeName, before)

	// 结束后
	db.Callback().Create().After("gorm:after_create").Register(callBackAfterName, after)
	db.Callback().Query().After("gorm:after_query").Register(callBackAfterName, after)
	db.Callback().Delete().After("gorm:after_delete").Register(callBackAfterName, after)
	db.Callback().Update().After("gorm:after_update").Register(callBackAfterName, after)
	db.Callback().Row().After("gorm:row").Register(callBackAfterName, after)
	db.Callback().Raw().After("gorm:raw").Register(callBackAfterName, after)
	return
}

func before(db *gorm.DB) {
	db.InstanceSet(startTime, time.Now())
}

func after(db *gorm.DB) {
	_reqId := db.Statement.Context.Value(requestId)
	reqId, ok := _reqId.(string)
	if !ok {
		reqId = "#"
	}

	_st, isExist := db.InstanceGet(startTime)
	if !isExist {
		return
	}
	startTime, ok := _st.(time.Time)
	if !ok {
		return
	}

	cost := strconv.FormatInt(time.Since(startTime).Milliseconds(), 10)
	sql := db.Dialector.Explain(db.Statement.SQL.String(), db.Statement.Vars...)
	affectedRows := strconv.FormatInt(db.Statement.RowsAffected, 10)
	stack := utils.FileWithLineNum()

	_logger.Info(tag, "request_id:"+reqId, "sql:"+sql, "cost:"+cost, "affectedRows:"+affectedRows, "stack:"+stack)
}
