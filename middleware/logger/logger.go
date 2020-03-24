package logger

import (
	Logf "DIEM-API/tools/logfactory"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

var zlog zerolog.Logger

// init log writer
func initLogWriter(w io.Writer) {
	zlog = zerolog.New(w).
		With().Timestamp().Logger()
}

func init() {
	ginMode := gin.Mode()
	if ginMode == "release" {
		initLogWriter(Logf.GetWriter("access"))
		return
	}
	initLogWriter(Logf.GetWriter("std"))
}

// log every request
func Log(c *gin.Context) {
	c.Next()
	zlog.Debug().Int("| Status", c.Writer.Status()).
		Str("| Path", c.Request.URL.Path).
		Str("| Error", c.Errors.String()).
		Msg("SUCCESS")
}
