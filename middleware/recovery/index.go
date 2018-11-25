package recovery

import (
	"DIEM-API/config"
	"errors"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"github.com/rs/zerolog"
)

var zlog zerolog.Logger

func init() {
	zlog = zerolog.New(config.LogWriter.ErrFile).
		With().Timestamp().Logger()
}

func storeError() {
	zlog.Error().Msg(string(debug.Stack()))
}

func Recover(c *gin.Context) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}

		storeError()

		e := r.(error)
		c.Error(errors.New(e.Error()))

		c.String(500, "Internal Error")

	}()
	c.Next()
}
