package logger

import (
	"DIEM-API/config"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var zlog zerolog.Logger

func init() {
	ginMode := gin.Mode()
	if ginMode == "release" {
		return
	}
	zlog = zerolog.New(config.LogWriter.ErrFile).
		With().Timestamp().Logger()
}

func logWithLevel(e *zerolog.Event, c *gin.Context) {
	e.Int("| Status", c.Writer.Status()).
		Str("| Path", c.Request.URL.Path).
		Str("| Error", c.Errors.String()).
		Msg("")
}

func Log(c *gin.Context) {
	c.Next()
	if c.Errors != nil {
		logWithLevel(zlog.Error(), c)
		return
	}
	logWithLevel(zlog.Info(), c)
}
