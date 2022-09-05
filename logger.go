package fastfiber

import (
	"log"

	winner_logger "github.com/bfmTech/logger-go"
)

// 阿里云sls 服务
// Console loggerMethod = "console" // 控制台输出日志
// File    loggerMethod = "file"    // 文件记录日志
// Http    loggerMethod = "http"    // http同步日志
func initSlSLogger(appName, logType string) winner_logger.Logger {
	method := winner_logger.Console
	if logType == string(winner_logger.File) {
		method = winner_logger.File
	} else if logType == string(winner_logger.Http) {
		method = winner_logger.Http

		//check log env
		if len(GetEnv("LOGGER_ALIYUN_ACCESSKEYID")) == 0 ||
			len(GetEnv("LOGGER_ALIYUN_ACCESSKEYSECRET")) == 0 {
			log.Fatal(ErrInitLoggerFail + " Http模式需要配置环境变量LOGGER_ALIYUN_ACCESSKEYID,LOGGER_ALIYUN_ACCESSKEYSECRET")
		}
	}

	logger, err := winner_logger.NewLogger(appName, method)

	if err != nil {
		log.Fatal(ErrInitLoggerFail + err.Error())
	}
	return logger
}
