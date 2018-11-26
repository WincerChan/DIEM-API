package logger

import (
	"DIEM-API/config"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var zlog zerolog.Logger

func initLogWriter(w io.Writer) {
	zlog = zerolog.New(w).
		With().Timestamp().Logger()
}

func init() {
	ginMode := gin.Mode()
	initLogWriter(config.Log.GetWriter("std"))
	if ginMode == "release" {
		initLogWriter(config.Log.GetWriter("access"))
		return
	}
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
