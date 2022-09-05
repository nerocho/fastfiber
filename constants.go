package fastfiber

const (
	// keys
	EventDestroyPrefix = "Destroy_"
	ProcessKilled      = "收到信号，进程被结束"

	// errors
	ErrInitConfigFile       = "初始化配置文件出错:"
	ErrInitLoggerFail       = "初始化日志出错:"
	ErrConfigFileNotExists  = "化配置文件不存在:"
	ErrorsDbInitFail        = "数据库初始化失败:"
	ErrorsIdWorkerInitFail  = "ID生成器 初始化失败:"
	ErrorsRedisInitConnFail = "Redis 初始化连接池失败:"
	ErrorsTracerInitFail    = "Tracer 初始化失败:"
)
