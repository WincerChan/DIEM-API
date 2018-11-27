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
	initLogWriter(config.GetWriter("std"))
	if ginMode == "release" {
		initLogWriter(config.GetWriter("access"))
		return
	}
}

func Log(c *gin.Context) {
	c.Next()
	zlog.Debug().Int("| Status", c.Writer.Status()).
		Str("| Path", c.Request.URL.Path).
		Str("| Error", c.Errors.String()).
		Msg("SUCCESS")
}
