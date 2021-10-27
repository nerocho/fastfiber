package fastfiber

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
)

func initZerolog() zerolog.Logger {
	size := Conf.GetInt("System.LogBufferSize")
	appName := Conf.GetString("System.AppName")
	logTimeFormat := Conf.GetString("System.LogTimeFormat")

	// 绑定日志
	if logTimeFormat == "Unix" {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	}
	wr := diode.NewWriter(os.Stdout, size, 10*time.Millisecond, nil)
	return zerolog.New(wr).With().Timestamp().Str("app", appName).Logger()
}
